package packages

// TODO: Complains about missing std::index_sequence,
// despite being in the C++14 standard and us using a
// C++14-compliant compiler.

// TODO: Install jemalloc.

import (
	"fmt"
	"path/filepath"
	"github.com/alexandrecarlton/gogurt"
)

type LNav struct{}

func (lnav LNav) Name() string {
	return "lnav"
}

func (lnav LNav) URL(version string) string {
	return fmt.Sprintf("https://github.com/tstack/lnav/releases/download/v%s/lnav-%s.tar.gz", version, version)
}

func (lnav LNav) Build(config gogurt.Config) error {

	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(lnav),
		Args: []string{
			"--enable-static",
			"--with-pcre=" + config.InstallDir(Pcre{}),
			"--with-readline=" + config.InstallDir(ReadLine{}),
			"--with-libcurl=" + config.InstallDir(Curl{}),
		},
		// A C++14 compiler is required.
		// clang is acting up, may need gcc instead.
		CXX: filepath.Join(config.BinDir(LLVM{}), "clang++"),
		CC: filepath.Join(config.BinDir(LLVM{}), "clang"),
		CPP: filepath.Join(config.BinDir(LLVM{}), "clang-cpp"),
		CxxFlags: []string{
			"-std=c++1z",
			"-I" + config.IncludeDir(Bzip2{}),
			"-I" + config.IncludeDir(Ncurses{}),
			"-I" + config.IncludeDir(OpenSSL{}),
			"-I" + config.IncludeDir(Pcre{}),
			"-I" + config.IncludeDir(ReadLine{}),
			"-I" + config.IncludeDir(SQLite{}),
			"-I" + config.IncludeDir(Zlib{}),
		},
		CFlags: []string{
			"-I" + config.IncludeDir(Bzip2{}),
			"-I" + config.IncludeDir(Ncurses{}),
			"-I" + config.IncludeDir(OpenSSL{}),
			"-I" + config.IncludeDir(Pcre{}),
			"-I" + config.IncludeDir(ReadLine{}),
			"-I" + config.IncludeDir(SQLite{}),
			"-I" + config.IncludeDir(Zlib{}),
		},
		CppFlags: []string{
			"-I" + config.IncludeDir(Bzip2{}),
			"-I" + config.IncludeDir(Ncurses{}),
			"-I" + config.IncludeDir(Zlib{}),
		},
		LdFlags: []string{
			"-L" + config.LibDir(Bzip2{}),
			"-L" + config.LibDir(Curl{}),
			"-L" + config.LibDir(OpenSSL{}),
			"-L" + config.LibDir(Ncurses{}),
			"-L" + config.LibDir(Pcre{}),
			"-L" + config.LibDir(ReadLine{}),
			"-L" + config.LibDir(SQLite{}),
			"-L" + config.LibDir(Zlib{}),
		},
	}.Cmd()
	configure.Env = append(configure.Env, "CURSES_LIB=-lncursesw -ltinfow")
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (lnav LNav) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (lnav LNav) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Bzip2{},
		Curl{},
		Ncurses{},
		Pcre{},
		ReadLine{},
		SQLite{},
		Zlib{},

		OpenSSL{}, // Due to curl
		LLVM{}, // for a C++14 compatible compiler.
	}
}
