//go:build mage
// +build mage

package main

import "errors"

// Schema generates the JSON Schema from the base models
func Schema() error {
	return errors.New("please now run `go generate .` instead")
}

// Regimes generates JSON version of each regimes's data.
func Regimes() error {
	return errors.New("please now run `go generate .` instead")
}
