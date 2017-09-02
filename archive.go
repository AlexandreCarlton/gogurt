package gogurt

import (
	"archive/tar"
	"compress/bzip2"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"github.com/ulikunitz/xz"
)

func DecompressSourceArchive(filename string, dir string) error {

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
	case ".tgz":
		fallthrough
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
	case ".bz2":
		compressedFile := bzip2.NewReader(file)
		return extractTar(compressedFile, dir)
	default:
		log.Fatalf("Unknown compression format for file '%s'.", filename)
		return nil
	}
}

// Strips out the leading folder.
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
			dir := filepath.Dir(newFilename)
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				os.MkdirAll(dir, 0755)
			}
			newFile, _ := os.Create(newFilename)
			defer newFile.Close()
			io.Copy(newFile, tarFile)
			os.Chmod(newFilename, header.FileInfo().Mode())
			// TODO: Make a PR for this in either archiver or go-getter
			// to set times, and strip components.
			os.Chtimes(newFilename, header.FileInfo().ModTime(), header.FileInfo().ModTime())
		case tar.TypeDir:
			//os.MkdirAll(newFilename, os.ModePerm)
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
