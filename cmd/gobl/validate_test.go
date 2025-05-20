package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"gitlab.com/flimzy/testy"
)

func Test_validate_args(t *testing.T) {
	tests := []struct {
		name string
		args []string
		err  string
	}{
		{
			name: "no args",
		},
		{
			name: "invalid flag",
			args: []string{"--foo"},
			err:  `unknown flag: --foo`,
		},
		{
			name: "force long",
			args: []string{"--force"},
		},
		{
			name: "force short",
			args: []string{"-f"},
		},
		{
			name: "in-place long",
			args: []string{"--in-place"},
		},
		{
			name: "in-place short",
			args: []string{"-w"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			root := root()
			rootCmd := &cobra.Command{}
			root.setFlags(rootCmd)

			opts := validate(root)

			cmd := opts.cmd()
			rootCmd.AddCommand(cmd)
			err := cmd.ParseFlags(tt.args)
			if tt.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
			}
			if err != nil {
				return
			}
			if d := testy.DiffInterface(testy.Snapshot(t), opts); d != nil {
				t.Error(d)
			}
		})
	}
}

func Test_validate(t *testing.T) {
	tmpdir := testy.CopyTempDir(t, "testdata", 0)
	t.Cleanup(func() {
		_ = os.RemoveAll(tmpdir)
	})

	tests := []struct {
		name   string
		opts   *validateOpts
		in     io.Reader
		args   []string
		err    string
		target string
	}{
		{
			name: "invalid stdin",
			in:   strings.NewReader("this isn't JSON"),
			err:  "code=400, message=yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `this is...` into map[string]interface {}",
		},
		{
			name: "no document",
			in: strings.NewReader(`{
				"$schema": "https://gobl.org/draft-0/envelope",
				"head": {
					"uuid": "9d8eafd5-77be-11ec-b485-5405db9a3e49",
					"typ": "duck",
					"rgn": "ES",
					"dig": {
						"alg": "sha256",
						"val": "dce3bc3c8bf28f3d209f783917b3082ddc0339a66e9ba3aa63849e4357db1422"
					}
				},
			}`),
			err: "code=422, message=doc: cannot be blank.",
		},
		{
			name: "invalid doc",
			in: strings.NewReader(`{
				"$schema": "https://gobl.org/draft-0/envelope",
				"head": {
					"uuid": "9d8eafd5-77be-11ec-b485-5405db9a3e49",
					"dig": {
						"alg": "sha256",
						"val": "dce3bc3c8bf28f3d209f783917b3082ddc0339a66e9ba3aa63849e4357db1422"
					}
				},
				doc: "foo bar baz"
			}`),
			err: "code=400, message=json: cannot unmarshal string into Go struct field Envelope.doc of type schema.document",
		},
		{
			name: "incomplete",
			in: strings.NewReader(`{
				"$schema": "https://gobl.org/draft-0/envelope",
				"head": {
					"uuid": "9d8eafd5-77be-11ec-b485-5405db9a3e49",
					"dig": {
						"alg": "sha256",
						"val": "dce3bc3c8bf28f3d209f783917b3082ddc0339a66e9ba3aa63849e4357db1422"
					}
				},
				doc: {}
			}`),
			err: "code=422, message=doc: ($schema: cannot be blank.).",
		},
		{
			name: "input file",
			args: []string{"testdata/success.json"},
		},
		{
			name: "without totals",
			args: []string{"testdata/nototals.json"},
			err:  "code=422, message=doc: (totals: cannot be blank.).",
		},
		{
			name:   "output file",
			args:   []string{"testdata/success.json", filepath.Join(tmpdir, "output-file.json")},
			target: filepath.Join(tmpdir, "output-file.json"),
		},
		{
			name: "explicit stdout",
			args: []string{"testdata/success.json", "-"},
		},
		{
			name: "output file exists",
			args: []string{"testdata/success.json", filepath.Join(tmpdir, "exists.json")},
			err:  "open " + tmpdir + "/exists.json: file exists",
		},
		{
			name: "overwrite output file",
			opts: &validateOpts{
				rootOpts: &rootOpts{
					overwriteOutputFile: true,
				},
			},
			args:   []string{"testdata/success.json", filepath.Join(tmpdir, "overwrite.json")},
			target: filepath.Join(tmpdir, "overwrite.json"),
		},
		{
			name: "overwrite input file",
			opts: &validateOpts{
				rootOpts: &rootOpts{
					inPlace: true,
				},
			},
			args:   []string{filepath.Join(tmpdir, "success.json")},
			target: filepath.Join(tmpdir, "success.json"),
		},
		{
			name: "overwrite stdin",
			opts: &validateOpts{
				rootOpts: &rootOpts{
					inPlace: true,
				},
			},
			err: "cannot overwrite STDIN",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &cobra.Command{}
			if tt.in != nil {
				c.SetIn(tt.in)
			}
			buf := &bytes.Buffer{}
			c.SetOut(buf)
			opts := tt.opts
			if opts == nil {
				opts = &validateOpts{}
			}
			if opts.rootOpts == nil {
				opts.rootOpts = &rootOpts{}
			}
			opts.indent = true
			err := opts.runE(c, tt.args)
			if tt.err != "" {
				assert.EqualError(t, err, tt.err)
			} else if err != nil {
				t.Errorf("Unexpected error: %q", err)
			}

			if d := testy.DiffText(testy.Snapshot(t), buf.String()); d != nil {
				t.Error(d)
			}
			if tt.target != "" {
				if d := testy.DiffText(testy.Snapshot(t, "outfile"), &testy.File{Path: tt.target}); d != nil {
					t.Errorf("outfile:\n%s", d)
				}
			}
		})
	}
}
