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

	// TODO: Viper this.
	config := gogurt.Config{
		Prefix: "/home/alexandre/.local",
		CacheFolder: "/home/alexandre/.cache/gogurt",
		BuildFolder: "/home/alexandre/.gogurt/build",
		InstallFolder: "/home/alexandre/.gogurt/install",
		NumCores: 3,
		// TODO We'll have a set of default versions, and our config will override the defaults.
		PackageVersions: map[string]string {
			"autoconf": "2.69",
			"automake": "1.15",
			"bzip2": "1.0.6",
			"curl": "7.54.1",
			"editorconfig-core-c": "0.12.1",
			"fish": "2.6.0",
			"libevent": "2.1.8-stable",
			"ncurses": "6.0",
			"openssl": "1.0.2k",
			"pcre": "8.40",
			"python2": "2.7.5",
			"texinfo": "6.3",
			"tmux": "2.5",
			"zlib": "1.2.11",
			"libffi": "3.2.1",
			"vim": "8.0.0045",
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
		"gettext" : packages.GetText{},
		"libevent": packages.Libevent{},
		"libffi": packages.LibFFI{},
		"openssl": packages.OpenSSL{},
		"pcre": packages.Pcre{},
		"python2": packages.Python2{},
		"ncurses": packages.Ncurses{},
		"tmux": packages.Tmux{},
		"texinfo": packages.TexInfo{},
		"vim": packages.Vim{},
		"zlib": packages.Zlib{},
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
	defer response.Body.Close()

	_, err = io.Copy(destination, response.Body)
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
		default:
			log.Println("Header is ", header)
			log.Println("Typeflag is ", header.Typeflag)
			// '103' is g = TypeXGlobalHeader
			log.Printf("No idea what '%s' is (original: '%s').\n", headerName, header.Name)
		}
	}
	return nil
}
