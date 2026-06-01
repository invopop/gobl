package main

import (
	"errors"
	"io"
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

func TestPrintErrorPlain(t *testing.T) {
	// Plain error -> wrapped as gobl.ErrInternal so the JSON body has
	// "key" and "message".
	_, stderr := testy.RedirIO(nil, func() {
		printError(errors.New("boom"))
	})
	b, _ := io.ReadAll(stderr)
	out := string(b)
	assert.Contains(t, out, "internal")
}

func TestPrintErrorGoblError(t *testing.T) {
	// Already a gobl.Error -> passes through unchanged.
	_, stderr := testy.RedirIO(nil, func() {
		printError(gobl.ErrInput.WithReason("nope"))
	})
	b, _ := io.ReadAll(stderr)
	out := string(b)
	assert.Contains(t, out, "input")
	assert.Contains(t, out, "nope")
}

func TestInputFilename(t *testing.T) {
	assert.Equal(t, "", inputFilename(nil))
	assert.Equal(t, "", inputFilename([]string{"-"}))
	assert.Equal(t, "foo.json", inputFilename([]string{"foo.json"}))
}

func TestRunExecutesRootCommand(t *testing.T) {
	// Drive `run()` via an invocation that doesn't actually need any
	// stdin. The default --help-like flow exits successfully.
	stdout, _ := testy.RedirIO(nil, func() {
		// run() routes through root().cmd().ExecuteContext, which
		// reads os.Args. We don't expose a clean override, so we just
		// confirm calling run with no args doesn't panic.
		// (Without args, cobra prints help and returns nil.)
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
