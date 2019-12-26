package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type HTop struct{}

func (htop HTop) Name() string {
	return "htop"
}

func (htop HTop) URL(version string) string {
	return fmt.Sprintf("http://hisham.hm/htop/releases/%s/htop-%s.tar.gz", version, version)
}

func (htop HTop) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(htop),
		Args: []string{
			"--disable-shared",
			"--enable-static",
			"--enable-proc",
			"--enable-cgroup",
			"--enable-openvz",
			"--enable-taskstats",
			"--enable-unicode",
			"--enable-linux-affinity",
		},
		CFlags: []string{
			"-I" + config.IncludeDir(Ncurses{}),
		},
		CppFlags: []string{
			"-I" + config.IncludeDir(Ncurses{}),
		},
		LdFlags: []string{
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

func (htop HTop) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (htop HTop) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Ncurses{},
	}
}
