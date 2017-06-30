package packages

// TODO: Requires aclocal-1.15
import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type Pcre struct{}

func (pcre Pcre) Name() string {
	return "pcre"
}

func (pcre Pcre) URL(version string) string {
	return fmt.Sprintf("http://ftp.csx.cam.ac.uk/pub/software/programming/pcre/pcre-%s.tar.gz", version)
}

func (pcre Pcre) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.Prefix,
		Args: []string{
			"--enable-unicode-properties",
			"--enable-static",
			"--disable-shared",
			"--enable-pcre16",
			"--enable-pcre32",
			"--enable-pcregrep-libz",
			"--enable-jit",
		},
		CFlags: []string{
			"-I" + config.IncludeDir("zlib"),
			"-I" + config.IncludeDir("bzip2"),
		},
		LdFlags: []string{
			"-L" + config.LibDir("zlib"),
			"-L" + config.LibDir("bzip2"),
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (pcre Pcre) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (pcre Pcre) Dependencies() []string {
	return []string{
		"bzip2",
		"zlib",
	}
}

