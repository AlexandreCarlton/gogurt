package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type TexInfo struct{}

func (texinfo TexInfo) Name() string {
	return "texinfo"
}

func (texinfo TexInfo) URL(version string) string {
	return fmt.Sprintf("https://ftp.gnu.org/gnu/texinfo/texinfo-%s.tar.gz", version)
}

func (texinfo TexInfo) Build(config gogurt.Config) error {

	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(texinfo),
		Args: []string{
			"--disable-perl-api-texi-build",
		},
		Paths: []string{
			config.BinDir(AutoMake{}),
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}

	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
		Paths: []string{
			config.BinDir(AutoMake{}),
		},
	}.Cmd()
	return make.Run()
}

func (texinfo TexInfo) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (texinfo TexInfo) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		AutoMake{},
	}
}
