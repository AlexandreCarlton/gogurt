package main

import (
	"archive/tar"
	"compress/gzip"

	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"


	"github.com/alexandrecarlton/gogurt"
	"github.com/alexandrecarlton/gogurt/packages"
	"github.com/ulikunitz/xz"
)

// have table of latest versions?
// and checksums...


func main() {

	// TODO: Create a default config object, merge that in.
	// TODO: look at mitchellh's homedir instead of using os.GetEnv(HOME)
	// TODO: Viper this.
	config := gogurt.Config{
		Prefix: filepath.Join(os.Getenv("HOME"), ".local"),
		CacheFolder: filepath.Join(os.Getenv("HOME"), ".cache/gogurt"),
		BuildFolder: filepath.Join(os.Getenv("HOME"),  ".gogurt/build"),
		InstallFolder: filepath.Join(os.Getenv("HOME"), ".gogurt/install"),
		NumCores: 3,
		// TODO We'll have a set of default versions, and our config will override the defaults.
		PackageVersions: map[string]string {
			"autoconf": "2.69",
			"automake": "1.15",
			"bzip2": "1.0.6",
			"curl": "7.54.1",
			"editorconfig-core-c": "0.12.1",
			"fish": "2.6.0",
			"gcc": "7.1.0",
			"git": "2.13.2",
			"gmp": "6.1.2",
			"go": "1.8.3",
			"help2man": "1.47.4",
			"libevent": "2.1.8-stable",
			"libtool": "2.4.6",
			"lua": "5.3.4",
			"mpc": "1.0.3",
			"mpfr": "3.1.5",
			"ncurses": "6.0",
			"neovim": "0.2.0",
			"nodejs": "8.1.3",
			"openssl": "1.0.2k",
			"pcre": "8.40",
			"python2": "2.7.5",
			"readline": "7.0",
			"stow": "2.2.2",
			"texinfo": "6.3",
			"the_silver_searcher": "2.0.0",
			"tmux": "2.5",
			"universal-ctags": "b5c9b76",
			"libffi": "3.2.1",
			"xz": "5.2.3",
			"vim": "8.0.0045",
			"zsh": "5.3.1",
			"zlib": "1.2.11",
		},
	}

	name := os.Args[1]
	mappings := map[string]gogurt.Package{
		"autoconf": packages.AutoConf{},
		"automake": packages.AutoMake{},
		"bzip2": packages.Bzip2{},
		"curl": packages.Curl{},
		"editorconfig-core-c": packages.EditorConfigCoreC{},
		"fish": packages.Fish{},
		"gcc": packages.GCC{},
		"gettext" : packages.GetText{},
		"git": packages.Git{},
		"gmp": packages.GMP{},
		"go": packages.Go{},
		"help2man": packages.Help2Man{},
		"libevent": packages.Libevent{},
		"libffi": packages.LibFFI{},
		"libtool": packages.LibTool{},
		"lua": packages.Lua{},
		"mpc": packages.MPC{},
		"mpfr": packages.MPFR{},
		"neovim": packages.NeoVim{},
		"openssl": packages.OpenSSL{},
		"pcre": packages.Pcre{},
		"python2": packages.Python2{},
		"ncurses": packages.Ncurses{},
		"nodejs": packages.NodeJS{},
		"readline": packages.ReadLine{},
		"stow": packages.Stow{},
		"texinfo": packages.TexInfo{},
		"the_silver_searcher": packages.TheSilverSearcher{},
		"tmux": packages.Tmux{},
		"universal-ctags": packages.UniversalCTags{},
		"vim": packages.Vim{},
		"xz": packages.XZ{},
		"zlib": packages.Zlib{},
		"zsh": packages.Zsh{},
	}
	pac := mappings[name]

	installPackage(pac, config)
}

func installPackage(pac gogurt.Package, config gogurt.Config) {

	if _, err := os.Stat(config.InstallDir(pac)); err == nil {
		log.Printf("Package '%s' already installed, skipping.", pac.Name())
		return
	}

	for _, dependency := range pac.Dependencies() {
		installPackage(dependency, config)
	}
	version := config.PackageVersions[pac.Name()]

	// Download tarball
	url := pac.URL(version)
	// TODO: This assumes we have a tarball. Should account for this.
	extension := filepath.Ext(url)
	cacheFilename := filepath.Join(config.CacheFolder, pac.Name(), pac.Name() + "-" + version  + ".tar" + extension)
	if err := os.MkdirAll(filepath.Dir(cacheFilename), os.ModePerm); err != nil {
		log.Fatalf("Error creating cache directory '%s': %s", config.CacheFolder, err.Error())
	}
	if _, err := os.Stat(cacheFilename); err == nil {
		log.Printf("File '%s' already exists, not downloading a new copy.", cacheFilename)
	} else if err := downloadFile(url, cacheFilename); err != nil {
		log.Fatalf("Could not download url '%s' to file '%s': %s\n", url, cacheFilename, err.Error())
	}

	buildDirname := config.BuildDir(pac)
	extractCompressedTar(cacheFilename, buildDirname)

	if err := os.Chdir(buildDirname); err != nil {
		log.Fatalf("Error changing to directory '%s': %s", buildDirname, err.Error())
	}
	if err := pac.Build(config); err != nil {
		log.Fatalf("Error building package '%s': %s", pac.Name(), err.Error())
	}
	if err := pac.Install(config); err != nil {
		log.Fatalf("Error installing package '%s': %s", pac.Name(), err.Error())
	}
}

func downloadFile(url string, destinationFilename string) error {
	log.Printf("Downloading url '%s' to file '%s'.\n", url, destinationFilename)

	destination, err := os.Create(destinationFilename)
	if err != nil {
		return err
	}
	defer destination.Close()

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	// Follow redirects. TODO - clean up.
	finalUrl := response.Request.URL.String()
	redirectedResponse, err := http.Get(finalUrl)
	if err != nil {
		return err
	}
	defer redirectedResponse.Body.Close()

	_, err = io.Copy(destination, redirectedResponse.Body)
	if err != nil {
		return err
	}

	return nil
}

func extractCompressedTar(filename string, dir string) error {

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("Error creating build directory '%s': %s", dir, err.Error())
	}


	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	ext := filepath.Ext(filename)
	switch ext {
	case ".gz":
		compressedFile, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer compressedFile.Close()
		return extractTar(compressedFile, dir)
	case ".xz":
		compressedFile, err := xz.NewReader(file)
		if err != nil {
			return err
		}
		return extractTar(compressedFile, dir)
	default:
		log.Fatalf("Unknown compression format for file '%s'.", filename)
		return nil
	}
}

func extractTar(file io.Reader, dir string) error {

	tarFile := tar.NewReader(file)

	for {
		header, err := tarFile.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		// Hack to strip out the leading component.
		headerName := strings.Join(strings.Split(header.Name, "/")[1:], "/")
		newFilename := filepath.Join(dir, headerName)

		switch header.Typeflag {
		case tar.TypeReg: fallthrough
		case tar.TypeRegA:
			func() {
				newFile, _ := os.Create(newFilename)
				defer newFile.Close()
				io.Copy(newFile, tarFile)
				os.Chmod(newFilename, header.FileInfo().Mode())
			}()
		case tar.TypeDir:
			os.MkdirAll(newFilename, os.ModePerm)
		case tar.TypeSymlink:
			source := filepath.Join(dir, strings.Join(strings.Split(header.Linkname, "/")[1:], "/"))
			os.Symlink(source, newFilename)
		default:
			log.Println("Header is ", header)
			log.Println("Typeflag is ", header.Typeflag)
			// '103' is g = TypeXGlobalHeader
			log.Printf("No idea what '%s' is (original: '%s').\n", headerName, header.Name)
		}
	}
	return nil
}
