package gobl

import "strings"

// VERSION identifies which version of GOBL is in use.
const VERSION Version = "gobl.org/0.10.0"

// Version string
type Version string

// Domain provides the version domain
func (v Version) Domain() string {
	res := v.Split()
	return res[0]
}

// Semver extracts the semversion component of the version string
func (v Version) Semver() string {
	res := v.Split()
	return res[1]
}

// Split divides version into a host and semver
func (v Version) Split() []string {
	return strings.SplitN(string(v), "/", 2)
}
