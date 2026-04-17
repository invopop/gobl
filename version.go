package gobl

import (
	"github.com/Masterminds/semver/v3"
)

// Version defines the semver for this version of GOBL
type Version string

// VERSION is the current version of the GOBL library
// 2026-04-15
const VERSION Version = "v0.401.0"

// Semver parses and returns semver
func (v Version) Semver() *semver.Version {
	return semver.MustParse(string(v))
}
