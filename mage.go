//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

const (
	name       = "gobl"
	mainBranch = "main"
)

// Lint runs golangci-lint
func Lint() error {
	return runQuiet("✓ Lint passed", "golangci-lint", "run")
}

// Fix runs golangci-lint with auto-fix
func Fix() error {
	return runQuiet("✓ Fix complete", "golangci-lint", "run", "--fix")
}

// Test runs all tests
func Test() error {
	return runQuiet("✓ Tests passed", "go", "test", "./...")
}

// TestRace runs all tests with the race detector
func TestRace() error {
	return runQuiet("✓ Tests passed (race)", "go", "test", "-race", "./...")
}

// Generate runs go generate
func Generate() error {
	return runQuiet("✓ Generate complete", "go", "generate", ".")
}

// Check runs the full pipeline: lint, generate, test, and verify no uncommitted changes
func Check() error {
	if err := Lint(); err != nil {
		return err
	}
	if err := Generate(); err != nil {
		return err
	}
	if err := Test(); err != nil {
		return err
	}
	// Verify generate didn't produce uncommitted changes
	if err := runQuiet("No uncommitted changes", "git", "diff", "--exit-code"); err != nil {
		return err
	}
	fmt.Println("✓ All checks passed")
	return nil
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
	return runQuiet("✓ "+action+" complete", "go", args...)
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

// runQuiet buffers output and only shows it on failure, printing msg on success.
// Use mage -v to stream everything.
func runQuiet(msg, cmd string, args ...string) error {
	if mg.Verbose() {
		return sh.RunV(cmd, args...)
	}
	c := exec.Command(cmd, args...)
	out, err := c.CombinedOutput()
	if err != nil {
		os.Stderr.Write(out)
		return err
	}
	fmt.Println(msg)
	return nil
}

func trimOutput(cmd string, args ...string) (string, error) {
	txt, err := sh.Output(cmd, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(txt), nil
}
