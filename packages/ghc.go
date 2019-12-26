package packages

// Needs an existing ghc. 7.10 is needed for 8.2.
// We'll need to download an archive (like LLVM does) and use that instead.

import (
	"fmt"
	// "os/exec"
	"github.com/alexandrecarlton/gogurt"
)

type GHC struct{}

func (ghc GHC) Name() string {
	return "ghc"
}

func (ghc GHC) URL(version string) string {
	return fmt.Sprintf("https://downloads.haskell.org/~ghc/%s/ghc-%s-src.tar.xz", version, version)
}

func (ghc GHC) Build(config gogurt.Config) error {
	// boot := exec.Command("./boot")
	// boot.Dir = config.BuildDir(ghc)
	// if err := boot.Run(); err != nil {
	// 	return err
	// }
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(ghc),
		// Args: []string{},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (ghc GHC) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (ghc GHC) Dependencies() []gogurt.Package {
	return []gogurt.Package{

	}
}
