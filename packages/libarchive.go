package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type LibArchive struct{}

func (libarchive LibArchive) Name() string {
	return "libarchive"
}

func (libarchive LibArchive) URL(version string) string {
	return fmt.Sprintf("https://www.libarchive.org/downloads/libarchive-%s.tar.gz", version)
}

func (libarchive LibArchive) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(libarchive),
		Args: []string{
			"--disable-shared",
			"--enable-static",
			"--enable-bsdtar",
			"--enable-bsdcat",
			"--enable-bsdcpio",
			"--with-bz2lib",
			"--with-lzma",
			"--with-zlib",
			"--with-openssl",
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (libarchive LibArchive) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (libarchive LibArchive) Dependencies() []gogurt.Package {
	return []gogurt.Package{

	}
}
