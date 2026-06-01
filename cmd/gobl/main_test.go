package main

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/flimzy/testy"

	"github.com/invopop/gobl"
)

func Test_root(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		stdin io.Reader
		err   string
	}{
		{
			name: "unsupported command",
			args: []string{"foo"},
			err:  `unknown command "foo" for "gobl"`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cmd := root().cmd()
			cmd.SetArgs(tt.args)
			var err error
			stdout, stderr := testy.RedirIO(tt.stdin, func() {
				err = cmd.Execute()
			})
			if d := testy.DiffText(testy.Snapshot(t, "_stdout"), stdout); d != nil {
				t.Errorf("STDOUT: %s", d)
			}
			if d := testy.DiffText(testy.Snapshot(t, "_stderr"), stderr); d != nil {
				t.Errorf("STDERR: %s", d)
			}
			if tt.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
			}
		})
	}
}

// withLogger pins slog.Default() to a logger that writes text-handler
// output into buf for the duration of the test.
func withLogger(t *testing.T) *bytes.Buffer {
	t.Helper()
	prev := slog.Default()
	buf := new(bytes.Buffer)
	slog.SetDefault(slog.New(slog.NewTextHandler(buf, nil)))
	t.Cleanup(func() { slog.SetDefault(prev) })
	return buf
}

func TestPrintErrorPlain(t *testing.T) {
	buf := withLogger(t)
	printError(errors.New("boom"))
	out := buf.String()
	assert.Contains(t, out, "key=internal")
	assert.Contains(t, out, "message=boom")
}

func TestPrintErrorGoblError(t *testing.T) {
	buf := withLogger(t)
	printError(gobl.ErrInput.WithReason("nope"))
	out := buf.String()
	assert.Contains(t, out, "key=input")
	assert.Contains(t, out, "message=nope")
}

func TestPrintErrorJSON(t *testing.T) {
	prev := slog.Default()
	buf := new(bytes.Buffer)
	slog.SetDefault(slog.New(slog.NewJSONHandler(buf, nil)))
	t.Cleanup(func() { slog.SetDefault(prev) })

	printError(gobl.ErrInput.WithReason("missing field"))
	out := buf.String()
	assert.Contains(t, out, `"key":"input"`)
	assert.Contains(t, out, `"message":"missing field"`)
}

func TestInputFilename(t *testing.T) {
	assert.Equal(t, "", inputFilename(nil))
	assert.Equal(t, "", inputFilename([]string{"-"}))
	assert.Equal(t, "foo.json", inputFilename([]string{"foo.json"}))
}

func TestRunExecutesRootCommand(_ *testing.T) {
	// Drive `run()` via an invocation that doesn't actually need any
	// stdin. Cobra prints help when no args are given and returns nil;
	// we only need to confirm it doesn't panic.
	stdout, _ := testy.RedirIO(nil, func() {
		_ = run()
	})
	_, _ = io.ReadAll(stdout)
}

func Test_version(t *testing.T) {
	cmd := versionCmd()
	stdout, stderr := testy.RedirIO(nil, func() {
		err := cmd.Execute()
		if err != nil {
			t.Fatal(err)
		}
	})
	wantOut := `"version": "dev"` // just check it's there somewhere!
	wantErr := ""
	if sout, _ := io.ReadAll(stdout); !strings.Contains(string(sout), wantOut) {
		t.Errorf("Unexpected STDOUT: %s", sout)
	}
	if serr, _ := io.ReadAll(stderr); !strings.Contains(string(serr), wantErr) {
		t.Errorf("Unexpected STDERR: %s", serr)
	}
}
