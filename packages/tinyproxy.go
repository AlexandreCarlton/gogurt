package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)


// requires asciidoc.

type TinyProxy struct{}

func (tinyproxy TinyProxy) Name() string {
	return "tinyproxy"
}

func (tinyproxy TinyProxy) URL(version string) string {
	return fmt.Sprintf("https://github.com/tinyproxy/tinyproxy/releases/download/%s/tinyproxy-%s.tar.xz", version, version)
}

func (tinyproxy TinyProxy) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(tinyproxy),
		Paths: []string{
			config.BinDir(AsciiDoc{}),
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
		// a2x needs to be able to find asciidoc and xmllint.
		Paths: []string{
			config.BinDir(AsciiDoc{}),
			config.BinDir(LibXML2{}), // for xmllint, required by asciidoc
			config.BinDir(LibXSLT{}), // for xsltproc, required by asciidoc
		},
	}.Cmd()
	return make.Run()
}

func (tinyproxy TinyProxy) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (tinyproxy TinyProxy) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		AsciiDoc{},
		LibXML2{},
		LibXSLT{},
	}
}
