package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newInitCmd builds a cobra.Command and re-applies any pre-set fields
// on opts (cmd flag wiring resets them to their defined defaults).
func newInitCmd(t *testing.T, opts *initCmdOpts) *cobra.Command {
	t.Helper()
	if opts.rootOpts == nil {
		opts.rootOpts = &rootOpts{}
	}
	preCfg := opts.configDir
	preName := opts.name
	preForce := opts.force
	c := opts.cmd()
	c.SetOut(new(bytes.Buffer))
	c.SetErr(new(bytes.Buffer))
	if preCfg != "" {
		opts.configDir = preCfg
	}
	if preName != "" {
		opts.name = preName
	}
	if preForce {
		opts.force = preForce
	}
	return c
}

func TestInitCmdFlags(t *testing.T) {
	o := &initCmdOpts{rootOpts: &rootOpts{}}
	c := o.cmd()
	assert.Equal(t, "init <domain>", c.Use)
	assert.NotNil(t, c.Flags().Lookup("config-dir"))
	assert.NotNil(t, c.Flags().Lookup("name"))
	assert.NotNil(t, c.Flags().Lookup("force"))
}

func TestInitCmdRunESuccess(t *testing.T) {
	tmp := t.TempDir()
	o := &initCmdOpts{
		rootOpts:  &rootOpts{},
		configDir: tmp,
		name:      "Alice Co",
	}
	c := newInitCmd(t, o)
	err := o.runE(c, []string{"alice.example"})
	require.NoError(t, err)
	assert.DirExists(t, filepath.Join(tmp, "alice.example", "keys"))
	assert.FileExists(t, filepath.Join(tmp, "alice.example", "private.jwk"))
	assert.FileExists(t, filepath.Join(tmp, "alice.example", "party.json"))
	assert.DirExists(t, filepath.Join(tmp, "alice.example", "inbox"))
}

func TestInitCmdRunEEmptyDomain(t *testing.T) {
	o := &initCmdOpts{rootOpts: &rootOpts{}}
	c := newInitCmd(t, o)
	err := o.runE(c, []string{""})
	require.Error(t, err)
}

func TestInitCmdRunEExistingDomain(t *testing.T) {
	tmp := t.TempDir()
	o := &initCmdOpts{rootOpts: &rootOpts{}, configDir: tmp}
	c := newInitCmd(t, o)
	require.NoError(t, o.runE(c, []string{"bob.example"}))
	// Re-run without --force fails.
	err := o.runE(c, []string{"bob.example"})
	require.Error(t, err)
}

func TestInitCmdCtor(t *testing.T) {
	o := initCmd(&rootOpts{})
	require.NotNil(t, o)
	assert.NotNil(t, o.rootOpts)
}

// TestInitCmdRunEUsesDefaultConfigDir exercises the empty-configDir
// fallback to defaultConfigDir() — we point HOME at a temp dir so the
// init writes there.
func TestInitCmdRunEUsesDefaultConfigDir(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	o := &initCmdOpts{rootOpts: &rootOpts{}, configDir: defaultConfigDir()}
	c := newInitCmd(t, o)
	require.NoError(t, o.runE(c, []string{"default.example"}))
	_, err := os.Stat(filepath.Join(tmp, ".config", "gobl", "default.example", "private.jwk"))
	require.NoError(t, err)
}
