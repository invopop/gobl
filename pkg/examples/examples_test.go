package examples_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/invopop/gobl/pkg/examples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// messageDoc is a minimal GOBL document (a note.Message) that calculates and
// validates without any regime or addon, so the conversion helpers can be
// exercised without heavy fixtures.
const messageDoc = `$schema: "https://gobl.org/draft-0/note/message"
uuid: "cfefe56a-9812-4b97-9b86-23eb6e9fb15a"
title: "Test Message"
content: "Hello"
`

func TestIsEnvelope(t *testing.T) {
	assert.True(t, examples.IsEnvelope("foo.env.yaml"))
	assert.True(t, examples.IsEnvelope("dir/bar.env.json"))
	assert.False(t, examples.IsEnvelope("foo.yaml"))
	assert.False(t, examples.IsEnvelope("dir/baz.json"))
}

func TestGoldenPath(t *testing.T) {
	assert.Equal(t, filepath.Join("examples", "out", "invoice.json"),
		examples.GoldenPath(filepath.Join("examples", "invoice.yaml")))
	// Only the final extension is replaced.
	assert.Equal(t, filepath.Join("out", "doc.env.json"),
		examples.GoldenPath("doc.env.yaml"))
}

func TestSources(t *testing.T) {
	root := t.TempDir()
	mk := func(rel, body string) {
		p := filepath.Join(root, filepath.FromSlash(rel))
		require.NoError(t, os.MkdirAll(filepath.Dir(p), 0o755))
		require.NoError(t, os.WriteFile(p, []byte(body), 0o644))
	}
	mk("a.yaml", "x")
	mk("b.json", "x")
	mk("sub/c.yaml", "x")
	mk("sub/out/c.json", "x") // generated golden -> skipped by default
	mk("readme.txt", "x")     // wrong extension -> skipped
	mk("vendor/d.yaml", "x")  // skipped via extra fragment

	got, err := examples.Sources(root, "vendor/")
	require.NoError(t, err)

	var rel []string
	for _, p := range got {
		r, err := filepath.Rel(root, p)
		require.NoError(t, err)
		rel = append(rel, r)
	}
	assert.ElementsMatch(t, []string{
		"a.yaml",
		"b.json",
		filepath.Join("sub", "c.yaml"),
	}, rel)
}

func TestSourcesError(t *testing.T) {
	_, err := examples.Sources(filepath.Join(t.TempDir(), "does-not-exist"))
	assert.Error(t, err)
}

func TestConvert(t *testing.T) {
	out, err := examples.Convert([]byte(messageDoc), false)
	require.NoError(t, err)
	s := string(out)

	// The envelope head UUID is forced to TestUUID for deterministic output.
	assert.Contains(t, s, examples.TestUUID.String())
	assert.Contains(t, s, "Test Message")

	// Deterministic: converting the same input again yields identical bytes.
	out2, err := examples.Convert([]byte(messageDoc), false)
	require.NoError(t, err)
	assert.Equal(t, out, out2)
}

const messageEnv = `$schema: "https://gobl.org/draft-0/envelope"
head:
  uuid: "cfefe56a-9812-4b97-9b86-23eb6e9fb15a"
doc:
  $schema: "https://gobl.org/draft-0/note/message"
  title: "Env Message"
  content: "Hi"
`

func TestConvertEnvelope(t *testing.T) {
	out, err := examples.Convert([]byte(messageEnv), true)
	require.NoError(t, err)
	s := string(out)
	assert.Contains(t, s, examples.TestUUID.String())
	assert.Contains(t, s, "Env Message")
}

func TestConvertInvalid(t *testing.T) {
	_, err := examples.Convert([]byte("\tnot valid yaml: ["), false)
	assert.Error(t, err)
}

func TestRun(t *testing.T) {
	root := t.TempDir()
	exDir := filepath.Join(root, "examples")
	require.NoError(t, os.MkdirAll(exDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(exDir, "message.yaml"), []byte(messageDoc), 0o644))

	golden := filepath.Join(exDir, "out", "message.json")

	// First pass with update=true creates the golden.
	examples.Run(t, exDir, true)
	require.FileExists(t, golden)

	// Second pass with update=false compares against it and must pass.
	examples.Run(t, exDir, false)
}
