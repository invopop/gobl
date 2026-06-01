package main

import (
	"bytes"
	"context"
	"encoding/json"
	stdnet "net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/internal/ops"
	"github.com/invopop/gobl/net"
	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/uuid"
)

// initDomainForCLI scaffolds a domain at <configDir>/<domain>/ using
// the internal ops layer so cmd/gobl tests can run without invoking
// the full init command.
func initDomainForCLI(t *testing.T, configDir, domain string) {
	t.Helper()
	require.NoError(t, ops.InitDomain(&ops.InitOptions{
		ConfigDir: configDir,
		Domain:    domain,
		Name:      domain,
		Out:       new(bytes.Buffer),
	}))
}

func TestNetCmdSubcommands(t *testing.T) {
	n := netCmd(&rootOpts{})
	c := n.cmd()
	assert.Equal(t, "net", c.Use)
	have := map[string]bool{}
	for _, sub := range c.Commands() {
		have[sub.Name()] = true
	}
	assert.True(t, have["serve"])
	assert.True(t, have["send"])
	assert.True(t, have["who"])
}

// ---------- net send -----------

func signedNoteBody(t *testing.T) []byte {
	t.Helper()
	priv := dsig.NewES256Key()
	msg := &note.Message{Content: "hi"}
	msg.SetUUID(uuid.V7())
	env, err := gobl.Envelop(msg)
	require.NoError(t, err)
	require.NoError(t, env.Sign(priv,
		net.Address("peer.example").URI(),
		net.Address("acme.example").URI()))
	out, err := json.Marshal(env)
	require.NoError(t, err)
	return out
}

func TestNetSendCmdMissingTo(t *testing.T) {
	o := netSend(&rootOpts{})
	c := o.cmd()
	c.SetArgs([]string{"-"})
	c.SetOut(new(bytes.Buffer))
	c.SetErr(new(bytes.Buffer))
	err := c.Execute()
	require.Error(t, err)
}

func TestNetSendCmdSuccess(t *testing.T) {
	body := signedNoteBody(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	tmp := t.TempDir()
	infile := filepath.Join(tmp, "env.json")
	require.NoError(t, os.WriteFile(infile, body, 0o644))

	o := netSend(&rootOpts{})
	c := o.cmd()
	c.SetOut(new(bytes.Buffer))
	c.SetErr(new(bytes.Buffer))
	u := strings.TrimPrefix(srv.URL, "http://")
	c.SetArgs([]string{"--to", u, "--insecure", infile})
	require.NoError(t, c.Execute())
}

func TestNetSendCmdBadInput(t *testing.T) {
	o := netSend(&rootOpts{})
	c := o.cmd()
	c.SetOut(new(bytes.Buffer))
	c.SetErr(new(bytes.Buffer))
	c.SetArgs([]string{"--to", "acme.example", "--insecure", "/no/such/file.json"})
	err := c.Execute()
	require.Error(t, err)
}

// ---------- net who -----------

func TestNetWhoCmdMissingFromFlag(t *testing.T) {
	o := netWho(&rootOpts{})
	c := o.cmd()
	c.SetOut(new(bytes.Buffer))
	c.SetErr(new(bytes.Buffer))
	c.SetArgs([]string{"acme.example"})
	err := c.Execute()
	require.Error(t, err)
}

func TestNetWhoCmdMissingKey(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	o := netWho(&rootOpts{})
	c := o.cmd()
	c.SetOut(new(bytes.Buffer))
	c.SetErr(new(bytes.Buffer))
	c.SetArgs([]string{"--from", "missing.example", "acme.example"})
	err := c.Execute()
	require.Error(t, err)
}

func TestNetWhoCmdMissingParty(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	initDomainForCLI(t, filepath.Join(tmp, ".config", "gobl"), "from.example")
	require.NoError(t, os.Remove(filepath.Join(tmp, ".config", "gobl", "from.example", "party.json")))

	o := netWho(&rootOpts{})
	c := o.cmd()
	c.SetOut(new(bytes.Buffer))
	c.SetErr(new(bytes.Buffer))
	c.SetArgs([]string{"--from", "from.example", "target.example"})
	err := c.Execute()
	require.Error(t, err)
}

func TestNetWhoCmdBadPartyJSON(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	initDomainForCLI(t, filepath.Join(tmp, ".config", "gobl"), "from.example")
	pj := filepath.Join(tmp, ".config", "gobl", "from.example", "party.json")
	require.NoError(t, os.WriteFile(pj, []byte("not json"), 0o644))

	o := netWho(&rootOpts{})
	c := o.cmd()
	c.SetOut(new(bytes.Buffer))
	c.SetErr(new(bytes.Buffer))
	c.SetArgs([]string{"--from", "from.example", "target.example"})
	err := c.Execute()
	require.Error(t, err)
}

func TestNetWhoCmdRunENilCase(t *testing.T) {
	// Direct runE call with empty --from -> short-circuits with an
	// explicit error message.
	o := &netWhoOpts{rootOpts: &rootOpts{}}
	err := o.runE(&cobra.Command{}, []string{"acme.example"})
	require.Error(t, err)
}

// ---------- net serve -----------

func TestNetServeCmdValidateMutualACME(t *testing.T) {
	o := &netServeOpts{rootOpts: &rootOpts{}, acmeLive: true, acmeTest: true}
	require.Error(t, o.validate())
}

func TestNetServeCmdValidateMutualACMEAndTLS(t *testing.T) {
	o := &netServeOpts{rootOpts: &rootOpts{}, acmeLive: true, tlsCert: "x.pem"}
	require.Error(t, o.validate())
}

func TestNetServeCmdValidatePartialTLS(t *testing.T) {
	o := &netServeOpts{rootOpts: &rootOpts{}, tlsCert: "x.pem"}
	require.Error(t, o.validate())
}

func TestNetServeCmdValidateOK(t *testing.T) {
	o := &netServeOpts{rootOpts: &rootOpts{}}
	require.NoError(t, o.validate())
}

func TestNetServeCmdNoDomainsErrors(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	o := netServe(&rootOpts{})
	c := o.cmd()
	c.SetOut(new(bytes.Buffer))
	c.SetErr(new(bytes.Buffer))
	c.SetArgs([]string{"--config-dir", tmp, "--http-port", strconv.Itoa(freeCLIPort(t))})
	err := c.Execute()
	require.Error(t, err)
}

func TestNetServeCmdManualMode(t *testing.T) {
	// Manual mode wires --party + --keys-dir + --private-key explicitly.
	tmp := t.TempDir()
	priv := dsig.NewES256Key()
	keysDir := filepath.Join(tmp, "keys")
	require.NoError(t, os.MkdirAll(keysDir, 0o755))
	pub, err := json.Marshal(priv.Public())
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(keysDir, priv.ID()+".json"), pub, 0o644))
	privBytes, err := json.MarshalIndent(priv, "", "  ")
	require.NoError(t, err)
	privFile := filepath.Join(tmp, "private.jwk")
	require.NoError(t, os.WriteFile(privFile, privBytes, 0o600))
	partyFile := filepath.Join(tmp, "party.json")
	partyBytes, err := json.Marshal(&org.Party{Name: "Solo"})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(partyFile, partyBytes, 0o644))

	o := netServe(&rootOpts{})
	c := o.cmd()
	c.SetOut(new(bytes.Buffer))
	c.SetErr(new(bytes.Buffer))
	port := freeCLIPort(t)
	c.SetArgs([]string{
		"--keys-dir", keysDir,
		"--party", partyFile,
		"--private-key", privFile,
		"--inbox", filepath.Join(tmp, "inbox"),
		"--http-port", strconv.Itoa(port),
	})

	ctx, cancel := context.WithCancel(context.Background())
	c.SetContext(ctx)
	done := make(chan error, 1)
	go func() { done <- c.Execute() }()
	time.Sleep(50 * time.Millisecond)
	cancel()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("netServe did not return")
	}
}

func freeCLIPort(t *testing.T) int {
	t.Helper()
	ln, err := stdnet.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close() //nolint:errcheck
	return ln.Addr().(*stdnet.TCPAddr).Port
}
