//go:build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

const (
	name       = "gobl"
	mainBranch = "main"
)

// Schema generates the JSON Schema from the base models
func Schema() error {
	return errors.New("please now run `go generate .` instead")
}

// Regimes generates JSON version of each regimes's data.
func Regimes() error {
	return errors.New("please now run `go generate .` instead")
}

// Build the binary
func Build() error {
	changed, err := target.Dir("./"+name, ".")
	if os.IsNotExist(err) || (err == nil && changed) {
		return build("build")
	}
	return nil
}

// Install the binary into your go bin path
func Install() error {
	return build("install")
}

func build(action string) error {
	args := []string{action}
	flags, err := buildFlags()
	if err != nil {
		return err
	}
	args = append(args, flags...)
	args = append(args, "./cmd/"+name)
	return sh.RunV("go", args...)
}

func buildFlags() ([]string, error) {
	ldflags := []string{
		fmt.Sprintf("-X 'main.date=%s'", date()),
	}
	if v, err := version(); err != nil {
		return nil, err
	} else if v != "" {
		ldflags = append(ldflags, fmt.Sprintf("-X 'main.version=%s'", v))
	}

	out := []string{}
	if len(ldflags) > 0 {
		out = append(out, "-ldflags="+strings.Join(ldflags, " "))
	}
	return out, nil
}

func version() (string, error) {
	vt, err := versionTag()
	if err != nil {
		return "", err
	}
	v := []string{vt}
	if b, err := branch(); err != nil {
		return "", err
	} else if b != mainBranch {
		v = append(v, b)
	}
	if uncommittedChanges() {
		v = append(v, "uncommitted")
	}
	return strings.Join(v, "-"), nil
}

func versionTag() (string, error) {
	return trimOutput("git", "describe", "--tags") // no "--exact-match"
}

func branch() (string, error) {
	return trimOutput("git", "rev-parse", "--abbrev-ref", "HEAD")
}

func uncommittedChanges() bool {
	err := sh.Run("git", "diff-index", "--quiet", "HEAD", "--")
	return err != nil
}

func date() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func trimOutput(cmd string, args ...string) (string, error) {
	txt, err := sh.Output(cmd, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(txt), nil
}
