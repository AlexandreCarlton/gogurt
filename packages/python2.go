package packages

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"github.com/alexandrecarlton/gogurt"
)

type Python2 struct{}

func (python2 Python2) Name() string {
	return "python2"
}

func (python2 Python2) URL(version string) string {
	return fmt.Sprintf("https://www.python.org/ftp/python/%s/Python-%s.tar.xz", version, version)
}

func (python2 Python2) Build(config gogurt.Config) error {
	openssl := OpenSSL{}
	zlib := Zlib{}

	gogurt.ReplaceInFile("Modules/Setup.dist", "^#SSL=.*", "SSL=" + strings.Replace(config.InstallDir(openssl), "/", "\\/", -1))
	uncommentModule("_ssl")
	uncommentModule("\t-DUSE_SSL")
	uncommentModule("\t-L\\$\\(SSL\\)") // TODO: Redo this section, these aren't modules.
	gogurt.ReplaceInFile("Modules/Setup.dist", "-lcrypto", "-lcrypto -ldl")

	addModule("future_builtins", []string{"future_builtins.c"})
	addModule("_multiprocessing", []string{
		"_multiprocessing/multiprocessing.c",
		"_multiprocessing/semaphore.c",
		"_multiprocessing/socket_connection.c",
		"-I./Modules/_multiprocessing", // TODO: Fix this.
		"-IModules/_multiprocessing",
	})
	uncommentModule("_bisect")
	uncommentModule("_collections")
	uncommentModule("_elementtree")
	uncommentModule("_functools")
	uncommentModule("_heapq")
	uncommentModule("_io")
	uncommentModule("_md5")
	uncommentModule("_random")
	uncommentModule("_sha")
	uncommentModule("_sha256")
	uncommentModule("_sha512")
	uncommentModule("_socket")
	uncommentModule("_struct")

	uncommentModule("array")
	uncommentModule("binascii")
	uncommentModule("cmath")
	uncommentModule("cPickle")
	uncommentModule("cStringIO")
	uncommentModule("datetime")
	uncommentModule("fcntl")
	uncommentModule("itertools")
	uncommentModule("math")
	uncommentModule("operator")
	uncommentModule("pyexpat")
	uncommentModule("select")
	uncommentModule("time")
	uncommentModule("unicodedata")
	uncommentModule("zlib")

	gogurt.ReplaceInFile("Modules/Setup.dist", "#\\*shared\\*", "\\*static\\*")


	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(python2),
		Args: []string{
			"--disable-shared",
			"--enable-unicode=ucs4",
			"--with-system-expat",
		},
		CFlags: []string{
			"-I" + config.IncludeDir(Expat{}),
			"-I" + config.IncludeDir(OpenSSL{}),
		},
		LdFlags: []string{
			"-L" + config.LibDir(Expat{}),
			"-L" + config.LibDir(OpenSSL{}),
		},
		Libs: []string{
			"-lexpat",
			"-lssl",
			"-lcrypto",
			"-lz",
		},
	}.Cmd();
	configure.Env = append(
		configure.Env,
		"LINKFORSHARED= ",
		"DYNLOADFILE=dynload_stub.o",
		"CFLAGSFORSHARED=-fPIC")
	if err := configure.Run(); err != nil {
		return err
	}

	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
	}.Cmd()

	return make.Run();
}

// TODO: Optimise - Taking too slow
// Perhaps only open it once, and modify it then.

func uncommentModule(module string) error {
	return gogurt.ReplaceInFile(
		"Modules/Setup.dist",
		"^#(" + module + ".*)",
		"\\1",
	)
}

func addModule(module string, sources []string) error {
	setupDistContents, _ := ioutil.ReadFile("Modules/Setup.dist")
	if bytes.Contains(setupDistContents, []byte(module)) {
		log.Printf("Modules/Setup.dist already contains %s.\n", module)
	} else {
		log.Printf("Adding module %s to Modules/Setup.dist.\n", module)
		setupDist, err := os.OpenFile("Modules/Setup.dist", os.O_WRONLY | os.O_APPEND, 0644)
		if err != nil {
			fmt.Println(err)
			return err
		}
		writer := bufio.NewWriter(setupDist)
		fmt.Fprintf(writer, "%s %s\n", module, strings.Join(sources, " "))
		writer.Flush()
		setupDist.Close()
	}
	return nil
}

func (python2 Python2) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{
		Args: []string{"install"},
	}.Cmd()
	return makeInstall.Run();
}

func (python2 Python2) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Expat{},
		LibFFI{}, // On second thought, is libffi really necessary?
		Zlib{},
		OpenSSL{},
	}
}

