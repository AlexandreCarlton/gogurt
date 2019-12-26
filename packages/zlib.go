package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type Zlib struct{}

func (zlib Zlib) Name() string {
	return "zlib"
}

func (zlib Zlib) URL(version string) string {
	return fmt.Sprintf("http://zlib.net/zlib-%s.tar.gz", version)
}

func (zlib Zlib) Build(config gogurt.Config) error {
	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
		Args: []string{
			"CFLAGS=" + config.IncludeDir(Zlib{}),
			"LDFLAGS=" + config.LibDir(Zlib{}),
		},
	}.Cmd()
	return make.Run()
}

func (zlib Zlib) Install(config gogurt.Config) error {

	make := gogurt.MakeCmd{
		Args: []string{"install"},
		Jobs: config.NumCores,
	}.Cmd()
	return make.Run()
}

func (zlib Zlib) Dependencies() []gogurt.Package {
	return []gogurt.Package{}
}
