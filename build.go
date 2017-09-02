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

	Dir string
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
		"PKG_CONFIG_PATH=" + strings.Join(configure.PkgConfigPaths, ":"),
	}
	if len(configure.Dir) > 0 {
		cmd.Dir = configure.Dir
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

	// Path to cmake binary
	// If empty, we will look up cmake on $PATH
	Path string

	// TODO: Rename to InstallPrefix
	Prefix string

	// Used when searching for include files and libraries
	// TODO: Rename to PrefixPath
	PathPrefix []string

	SourceDir string

	BuildDir string

	CacheEntries map[string]string

	// Generator string // TODO: Add once we get Ninja

	Paths[] string

	CFlags []string

	PkgConfigPaths []string
}

func (cmakeCmd CMakeCmd) Cmd() *exec.Cmd {
	cacheEntries := make([]string, 1)
	if len(cmakeCmd.Prefix) > 0 {
		cacheEntries = append(cacheEntries, "-DCMAKE_INSTALL_PREFIX=" + cmakeCmd.Prefix)
	}
	for key, value := range cmakeCmd.CacheEntries {
		cacheEntries = append(cacheEntries, fmt.Sprintf("-D%s=%s", key, value))
	}
	// Hacky way to use our own CMake if provided.
	var cmd *exec.Cmd
	if len(cmakeCmd.Path) > 0 {
		cmd = exec.Command(cmakeCmd.Path, cacheEntries...)
	} else {
		cmd = exec.Command("cmake", cacheEntries...)
	}
	cmd.Args = append(cmd.Args, cmakeCmd.SourceDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = cmakeCmd.BuildDir
	cmd.Env = os.Environ()
	if len(cmakeCmd.PathPrefix) > 0 {
		cmd.Env = append(cmd.Env, "CMAKE_PREFIX_PATH=" + strings.Join(cmakeCmd.PathPrefix, ":"))
	}
	if len(cmakeCmd.PkgConfigPaths) > 0 {
		cmd.Env = append(cmd.Env, "PKG_CONFIG_PATH=" + strings.Join(cmakeCmd.PkgConfigPaths, ":"))
	}
	if len(cmakeCmd.Paths) > 0 {
		cmd.Env = append(cmd.Env, "PATH=" + strings.Join(cmakeCmd.Paths, ":") + ":" + os.Getenv("PATH"))
	}
	fmt.Println(cmd)
	return cmd
}

