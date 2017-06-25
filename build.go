package gogurt

// Contains structures to ease running of common operations.
// TODO: Have RunCmds(exec.Command...) and it just runs through all in sequence, stopping if one errors.

import (
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type ConfigureCmd struct {
	Prefix string

	Args []string

	CFlags []string

	CppFlags []string

	LdFlags []string

	Libs []string

	Paths []string

	PkgConfigPaths []string
}

func (configure ConfigureCmd) Cmd() *exec.Cmd {
	args := append(configure.Args, "--prefix=" + configure.Prefix)
	cmd := exec.Command("./configure", args...)
	cmd.Env = []string{
		"CFLAGS=" + strings.Join(configure.CFlags, " "),
		"CPPFLAGS=" + strings.Join(configure.CppFlags, " "),
		"LDFLAGS=" + strings.Join(configure.LdFlags, " "),
		"LIBS=" + strings.Join(configure.Libs, " "),
		"PATH=" + strings.Join(configure.Paths, ":") + ":" + os.Getenv("PATH"),
		"PKG_CONFIG_PATHS=" + strings.Join(configure.PkgConfigPaths, ":"),
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}


type MakeCmd struct {
	Jobs uint

	Args []string
}

func (makeCmd MakeCmd) Cmd() *exec.Cmd {
	jobs := int(math.Max(1, float64(makeCmd.Jobs)))
	args := append(makeCmd.Args, "--jobs=" + strconv.Itoa(jobs))
	cmd := exec.Command("make", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
