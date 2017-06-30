package packages

import (
	"fmt"
	"path/filepath"
	"github.com/alexandrecarlton/gogurt"
)

type Tmux struct{}

func (tmux Tmux) Name() string {
	return "tmux"
}

func (tmux Tmux) URL(version string) string {
	return fmt.Sprintf("https://github.com/tmux/tmux/releases/download/%s/tmux-%s.tar.gz", version, version)
}

func (tmux Tmux) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.Prefix,
	 	Args: []string{
			"--enable-static",
		},
		CFlags: []string{
			"-I" + config.IncludeDir("libevent"),
			"-I" + config.IncludeDir("ncurses"),
		},
		LdFlags: []string{
			"-L" + config.LibDir("libevent"),
			"-L" + config.LibDir("ncurses"),
		},
		PkgConfigPaths: []string{
			filepath.Join(config.InstallDir("ncurses"), "share", "pkgconfig"),
		},
	}.Cmd()
	configure.Env = append(
		configure.Env,
		"LIBNCURSES_CFLAGS=-D_GNU_SOURCE",
		"LIBNCURSES_LIBS=-lncursesw -ltinfow",
	)

	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (tmux Tmux) Install(config gogurt.Config) error {
	make := gogurt.MakeCmd{
		Args: []string{
			"install",
			"prefix=" + config.InstallDir("tmux"),
		},
	}.Cmd()
	return make.Run()
}

func (tmux Tmux) Dependencies() []string {
	return []string{"libevent", "ncurses"}
}
