package packages

// TODO: This requires an extra rct thing.
// Probably will do this later.

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/alexandrecarlton/gogurt"
)

type RTags struct{}
type RCT struct{}

func (rtags RTags) Name() string {
	return "rtags"
}

func (rct RCT) Name() string{
	return "rct"
}


func (rtags RTags) URL(version string) string {
	return fmt.Sprintf("https://github.com/Andersbakken/rtags/archive/v%s.tar.gz", version)
}

func (rtags RTags) Build(config gogurt.Config) error {
	build := filepath.Join(config.BuildDir(rtags), "build")
	os.Mkdir(build, 0755)

	cmake := gogurt.CMakeCmd{
		Prefix: config.InstallDir(rtags),
		SourceDir: config.BuildDir(rtags),
		BuildDir: build,
		CacheEntries: map[string]string{
			"BUILD_SHARED_LIBS": "OFF",
			"RTAGS_NO_ELISP_FILES": "ON", // TODO: Install emacs?
		},
		PathPrefix: []string{
			config.InstallDir(LLVM{}),
		},
	}.Cmd()
	if err := cmake.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
		Dir: build,
	}.Cmd()
	return make.Run()
}

func (rtags RTags) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{
		Args: []string{"install"},
		Dir: filepath.Join(config.BuildDir(rtags), "build"),
	}.Cmd()
	return makeInstall.Run()
}

func (rtags RTags) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		LLVM{}, // It's clang based, after all.
		// Lua, zlib, openssl...
	}
}
