package gogurt

import (
	"path/filepath"
)

// Shit. We need to put in the latest versions here as well.
// Have a like a subsection, e.g. 'bzip2': '1.0.6'

type Config struct {
	// where we assume applications will finally be installed.
	Prefix string

	// Where we store our tarballs.
	CacheFolder string

	// Where we extract and build our applications
	BuildFolder string

	// Where we install applications to keep them separate.
	InstallFolder string

	// How many cores are allocated to building applications.
	NumCores uint

  // We load from defaults, overwriting with the entries found in the config file.
	PackageVersions map[string]string

}

func (config Config) BuildDir(packageName string) string {
  return filepath.Join(config.BuildFolder, packageName, config.PackageVersions[packageName])
}

func (config Config) InstallDir(packageName string) string {
  return filepath.Join(config.InstallFolder, packageName, config.PackageVersions[packageName])
}

func (config Config) BinDir(packageName string) string {
	return filepath.Join(config.InstallDir(packageName), "bin")
}

func (config Config) IncludeDir(packageName string) string {
	return filepath.Join(config.InstallDir(packageName), "include")
}

func (config Config) LibDir(packageName string) string {
	return filepath.Join(config.InstallDir(packageName), "lib")
}
