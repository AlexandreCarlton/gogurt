package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type AsciiDoc struct{}

func (asciidoc AsciiDoc) Name() string {
	return "asciidoc"
}

func (asciidoc AsciiDoc) URL(version string) string {
	return fmt.Sprintf("https://downloads.sourceforge.net/asciidoc/asciidoc-%s.tar.gz", version)
}

func (asciidoc AsciiDoc) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(asciidoc),
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (asciidoc AsciiDoc) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (asciidoc AsciiDoc) Dependencies() []gogurt.Package {
	// Note that this is a Python application, so some form of Python is needed to actually use it.
	return []gogurt.Package{}
}
