package packages

// contains libmount (for glib)

import (
	"fmt"
	"strings"
	"github.com/alexandrecarlton/gogurt"
)

type UtilLinux struct{}

func (utilLinux UtilLinux) Name() string {
	return "utilLinux"
}

func (utilLinux UtilLinux) URL(version string) string {
	majorVersion := strings.Join(strings.Split(version, ".")[:2], ".")
	return fmt.Sprintf("https://www.kernel.org/pub/linux/utils/util-linux/v%s/util-linux-%s.tar.xz", majorVersion, version)
}

func (utilLinux UtilLinux) Build(config gogurt.Config) error {
	// configure := gogurt.ConfigureCmd{
	// 	Prefix: config.InstallDir(utilLinux),
	// 	Args: []string{
	// 		"--disable-shared",
	// 		"--enable-static",
	// 	},
	// }.Cmd()
	// if err := configure.Run(); err != nil {
	// 	return err
	// }
	// make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	// return make.Run()
	return nil
}

func (utilLinux UtilLinux) Install(config gogurt.Config) error {
	// makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	// return makeInstall.Run()
	return nil
}

func (utilLinux UtilLinux) Dependencies() []gogurt.Package {
	return []gogurt.Package{

	}
}
