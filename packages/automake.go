package packages

import (
	"fmt"
	"os/exec"
	"github.com/alexandrecarlton/gogurt"
)

type AutoMake struct{}

func (automake AutoMake) Name() string {
	return "automake"
}

func (automake AutoMake) URL(version string) string {
	return fmt.Sprintf("https://ftp.gnu.org/gnu/automake/automake-%s.tar.gz", version)
}

func (automake AutoMake) Build(config gogurt.Config) error {
	// This *shouldn't* be necessary but errors are generated without it.
	bootstrap := exec.Command("./bootstrap.sh")
	if err := bootstrap.Run(); err != nil {
		return err
	}

	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(automake),
		Args: []string{
			// Disable documentation since we need texinfo, which depends on AutoMake.
			"--disable-perl-api-texi-build",
			"--disable-pod-simple-texinfo-tests",
		},
		Paths: []string{
			config.BinDir(AutoConf{}),
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}

	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (automake AutoMake) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (automake AutoMake) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		AutoConf{}, // doesn't really need it on centos 7
		//TexInfo{}, // for documentation... But Texinfo Requires Automake with recent version.
	}
}

