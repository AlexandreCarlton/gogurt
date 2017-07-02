package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type Fish struct{}

func (fish Fish) Name() string {
	return "fish"
}

func (fish Fish) URL(version string) string {
	return fmt.Sprintf("https://fishshell.com/files/%s/fish-%s.tar.gz", version, version)
}

func (fish Fish) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(fish),
		Args: []string{
			"--without-doxygen",
		},
		CxxFlags: []string{
			"-I" + config.IncludeDir(Ncurses{}),
		},
		CppFlags: []string{
			"-I" + config.IncludeDir(Ncurses{}),
		},
		LdFlags: []string{
			"-static",
			"-L" + config.LibDir(Ncurses{}),
		},
		Libs: []string{
			"-ltinfow",
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (fish Fish) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (fish Fish) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Ncurses{},
	}
}
