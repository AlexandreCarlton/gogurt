package packages

// NOTE:
// There are several things that make Vim
// Vim hardcodes the prefix used to look for its runtime files when configuring (in src/auto/pathdef.c).
// These files are in PREFIX/share/vim/vim<ver>. which is problematic if your prefix is in /home.
// This is normally not a problem for personal use, but if your prefix is in /home then it is problematic to distribute this to your friends.
// They'll need to set the environment variable VIMRUNTIME=<path/to>/share/vim/vim<ver>. Then everything should just work.
//

import (
	"fmt"
	"path/filepath"
	// "strings"
	"github.com/alexandrecarlton/gogurt"
)

type Vim struct{}

func (vim Vim) Name() string {
	return "vim"
}

func (vim Vim) URL(version string) string {
	return fmt.Sprintf("https://github.com/vim/vim/archive/v%s.tar.gz", version)
}

func (vim Vim) Build(config gogurt.Config) error {

	ncurses := Ncurses{}
	openssl := OpenSSL{}
	python2 := Python2{}
	zlib := Zlib{}
	//pythonMajorVersion := strings.Join(strings.Split(config.PackageVersions["python2"], ".")[:2], ".")

	configure := gogurt.ConfigureCmd{
		// Note: Vim cares where it thinks it's prefix is, as it uses this to find its
		// accompanying *.vim files
		Prefix: config.Prefix,
		Args: []string{
			"--with-features=huge",
			"--enable-multibyte",
			"--enable-cscope",
			"--enable-rubyinterp", // TODO: Need to install ruby.
			"--enable-pythoninterp",
			"--with-python-config-dir=" + filepath.Join(config.LibDir(python2.Name()), "python2.7", "config"), // TODO: Get major the version out of the config.
			"--enable-luainterp",
			"--disable-darwin",
			"--disable-gui",
			"--disable-netbeans",
			"--without-x",
			"--with-tlib=tinfow",
		},
		CFlags: []string{
			"-Os", // Vim tries to compile with _FORTIFY_SOURCE, which requires -O
			"-I" + config.IncludeDir("openssl"),
			"-I" + config.IncludeDir("ncurses"),
			"-I" + config.IncludeDir("python2"),
			"-I" + filepath.Join(config.IncludeDir(python2.Name()), "python2.7"),
		},
		CppFlags: []string{
			"-I" + config.IncludeDir(ncurses.Name()),
		},
		LdFlags: []string{
			"-static",
			"-L" + config.LibDir(ncurses.Name()),
			"-L" + config.LibDir(python2.Name()),
			"-L" + config.LibDir(openssl.Name()),
			"-L" + config.LibDir(zlib.Name()),
		},
		Libs: []string{
			"-lncursesw",
			"-ltinfow",
			"-lpython2.7",
			"-lssl",
			"-lcrypto",
			"-ldl",
			"-lz",
		},
		Paths: []string{
			config.BinDir(python2.Name()), // We want our static Python2 to be picked up.
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
	}.Cmd()
	return make.Run()
}

func (vim Vim) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{
		Args: []string{
			"install",
			"prefix=" + config.InstallDir(vim.Name()),
		},
	}.Cmd()
	return makeInstall.Run()
}

func (vim Vim) Dependencies() []string {
	return []string{
		"ncurses",
		"python2",
	}
}

