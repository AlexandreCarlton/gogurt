package packages

// Downloads a binary release of fzf.
// TODO: Won't work - no folder :/

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/alexandrecarlton/gogurt"
)

type FZF struct{}

func (fzf FZF) Name() string {
	return "fzf"
}

func (fzf FZF) URL(version string) string {
	return fmt.Sprintf("https://github.com/junegunn/fzf-bin/releases/download/%s/fzf-%s-linux_amd64.tgz", version, version)
}

func (fzf FZF) Build(config gogurt.Config) error {
	// It's prebuilt!
	return nil
}

func (fzf FZF) Install(config gogurt.Config) error {
	err := os.MkdirAll(config.BinDir(fzf), 0755)
	if err != nil {
		return err
	}

	return gogurt.CopyFile(
		filepath.Join(config.BuildDir(fzf), "fzf"),
		filepath.Join(config.BinDir(fzf), "fzf"))
}

func (fzf FZF) Dependencies() []gogurt.Package {
	return []gogurt.Package{}
}


