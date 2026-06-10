//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
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

// Generate runs go generate (regenerates schemas, definitions, rules data)
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
