package gobl

import "strings"

// Version contains a domain and semver for this version of GOBL.
type Version string

// VERSION is the version of the GOBL library.
const VERSION Version = "gobl.org/v0.10.0"

// Domain returns the domain portion of the version.
func (v Version) Domain() string {
	parts := strings.SplitN(string(v), "/", 2)
	return parts[0]
}

// Semver returns the semver portion of the version.
func (v Version) Semver() string {
	parts := strings.SplitN(string(v), "/", 2)
	return parts[1]
}
