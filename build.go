package gogurt

// Contains structures to ease running of common operations.
// TODO: Have RunCmds(exec.Command...) and it just runs through all in sequence, stopping if one errors.

import (
	"fmt"
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

	CxxFlags []string

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
		"CXXFLAGS=" + strings.Join(configure.CxxFlags, " "),
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

	Dir string

	Paths []string
}

func (makeCmd MakeCmd) Cmd() *exec.Cmd {
	jobs := int(math.Max(1, float64(makeCmd.Jobs)))
	args := append(makeCmd.Args, "--jobs=" + strconv.Itoa(jobs))
	cmd := exec.Command("make", args...)
	cmd.Env = []string{
		"PATH=" + strings.Join(makeCmd.Paths, ":") + ":" + os.Getenv("PATH"),
	}
	if len(makeCmd.Dir) > 0 {
		cmd.Dir = makeCmd.Dir
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

type CMakeCmd struct {
	Prefix string

	SourceDir string

	BuildDir string

	CacheEntries map[string]string

	// Generator string // TODO: Add once we get Ninja

	CFlags []string
}

func (cmakeCmd CMakeCmd) Cmd() *exec.Cmd {
	cacheEntries := make([]string, 1)
	if len(cmakeCmd.Prefix) > 0 {
		cacheEntries = append(cacheEntries, "-DCMAKE_INSTALL_PREFIX=" + cmakeCmd.Prefix,)
	}
	for key, value := range cmakeCmd.CacheEntries {
		cacheEntries = append(cacheEntries, fmt.Sprintf("-D%s=%s", key, value))
	}

	cmd := exec.Command("cmake", cacheEntries...)
	cmd.Args = append(cmd.Args, cmakeCmd.SourceDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = cmakeCmd.BuildDir
	fmt.Println(cmd)
	return cmd
}

