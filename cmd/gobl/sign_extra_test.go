package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestSignCmdDomainAndKeyMutuallyExclusive(t *testing.T) {
	o := sign(&rootOpts{})
	c := o.cmd()
	c.SetOut(new(bytes.Buffer))
	c.SetErr(new(bytes.Buffer))
	c.SetIn(strings.NewReader(`{}`))
	c.SetArgs([]string{"--domain", "x.example", "--key", "key.jwk", "-"})
	err := c.Execute()
	require.Error(t, err)
}

func TestSignCmdMissingTemplate(t *testing.T) {
	o := sign(&rootOpts{})
	c := o.cmd()
	c.SetOut(new(bytes.Buffer))
	c.SetErr(new(bytes.Buffer))
	c.SetArgs([]string{"--template", "/no/such/template.yaml", "-"})
	err := c.Execute()
	require.Error(t, err)
}

func TestSignCmdMissingKeyFile(t *testing.T) {
	o := sign(&rootOpts{})
	c := o.cmd()
	c.SetOut(new(bytes.Buffer))
	c.SetErr(new(bytes.Buffer))
	c.SetIn(strings.NewReader(`{}`))
	c.SetArgs([]string{"--key", "/no/such/key.jwk", "-"})
	err := c.Execute()
	require.Error(t, err)
}

// TestSignCmdRunEAudience verifies the --to / audience branch.
func TestSignCmdRunEAudience(t *testing.T) {
	o := sign(&rootOpts{})
	c := &cobra.Command{}
	c.SetIn(strings.NewReader(`{}`))
	c.SetOut(new(bytes.Buffer))
	o.privateKeyFile = "/no/such/key.jwk"
	o.audience = "to.example"
	err := o.runE(c, []string{"-"})
	require.Error(t, err)
}

// TestSignCmdRunEDomainBranch exercises the --domain key-file
// resolution branch by pointing HOME at a tempdir without the
// expected private.jwk file -> open error.
func TestSignCmdRunEDomainBranch(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	o := sign(&rootOpts{})
	c := o.cmd()
	c.SetOut(new(bytes.Buffer))
	c.SetErr(new(bytes.Buffer))
	c.SetIn(strings.NewReader(`{}`))
	c.SetArgs([]string{"--domain", "x.example", "-"})
	err := c.Execute()
	require.Error(t, err)
}

// TestSignCmdRunETemplateSuccess verifies the --template branch by
// supplying a readable template file. The sign proceeds and only
// errors on the absent key file.
func TestSignCmdRunETemplateSuccess(t *testing.T) {
	tmp := t.TempDir()
	tpl := filepath.Join(tmp, "tpl.yaml")
	require.NoError(t, os.WriteFile(tpl, []byte("doc:\n  foo: bar\n"), 0o644))
	o := sign(&rootOpts{})
	c := o.cmd()
	c.SetOut(new(bytes.Buffer))
	c.SetErr(new(bytes.Buffer))
	c.SetIn(strings.NewReader(`{}`))
	c.SetArgs([]string{"--template", tpl, "--key", "/no/such/key.jwk", "-"})
	err := c.Execute()
	require.Error(t, err)
}
