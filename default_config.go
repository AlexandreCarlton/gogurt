package gogurt

import (
	"os"
	"path/filepath"
)

// TODO: use mitchellh's homedir
var DefaultConfig Config = Config{
	Prefix: filepath.Join(os.Getenv("HOME"), ".local"),

	CacheFolder: filepath.Join(os.Getenv("HOME"), ".cache/gogurt"),

	BuildFolder: filepath.Join(os.Getenv("HOME"),  ".gogurt/build"),

	InstallFolder: filepath.Join(os.Getenv("HOME"), ".gogurt/install"),

	NumCores: 1,

	PackageVersions: map[string]string {
		"autoconf": "2.69",
		"automake": "1.15",
		"bzip2": "1.0.6",
		"curl": "7.54.1",
		"editorconfig-core-c": "0.12.1",
		"fish": "2.6.0",
		"gettext": "0.19.8",
		"git": "2.13.2",
		"libevent": "2.1.8-stable",
		"ncurses": "6.0",
		"neovim": "0.2.0",
		"openssl": "1.0.2k",
		"pcre": "8.40",
		"python2": "2.7.5",
		// "texinfo": "6.3",
		"tmux": "2.5",
		"zlib": "1.2.11",
		"libffi": "3.2.1",
		"vim": "8.0.0045",
	},
}
