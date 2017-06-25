package packages

import (
	"github.com/alexandrecarlton/gogurt"
)

// Package ...
type Package interface {
	URL(version string) string

	// These execute in the checked out directory.
	// Should these take in a config object? Would make things much easier and flexible.
	Build(config gogurt.Config) error
	Install(config gogurt.Config) error
	// Clean() int

	Dependencies() []string
}
