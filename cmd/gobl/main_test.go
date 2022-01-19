package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"gitlab.com/flimzy/testy"
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
			cmd := root()
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

func Test_verify(t *testing.T) {
	tests := []struct {
		name string
		in   io.Reader
		args []string
		err  string
	}{
		{
			name: "invalid stdin",
			in:   strings.NewReader("this isn't JSON"),
			err:  "error unmarshaling JSON: json: cannot unmarshal string into Go value of type gobl.Envelope",
		},
		{
			name: "read error",
			in:   testy.ErrorReader(`{"foo":`, errors.New("read error")),
			err:  "read error",
		},
		{
			name: "empty envelope",
			in:   strings.NewReader(`{}`),
			err:  "doc: cannot be blank; head: cannot be blank.",
		},
		{
			name: "success",
			in: func() io.Reader {
				f, err := os.Open("testdata/success.json")
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					_ = f.Close()
				})
				return f
			}(),
		},
		{
			name: "digest mismatch",
			in: func() io.Reader {
				f, err := os.Open("testdata/digest-mismatch.json")
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					_ = f.Close()
				})
				return f
			}(),
			err: "digest mismatch",
		},
		{
			name: "read from file",
			args: []string{"testdata/digest-mismatch.json"},
			err:  "digest mismatch",
		},
		{
			name: "file missing",
			args: []string{"asdf"},
			err:  "open asdf: no such file or directory",
		},
		{
			name: "explicit stdin",
			args: []string{"-"},
			in: func() io.Reader {
				f, err := os.Open("testdata/digest-mismatch.json")
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					_ = f.Close()
				})
				return f
			}(),
			err: "digest mismatch",
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
			err := verify(c, tt.args)
			if tt.err != "" {
				assert.EqualError(t, err, tt.err)
			} else {
				assert.Nil(t, err)
			}
			if d := testy.DiffText(testy.Snapshot(t), buf.String()); d != nil {
				t.Error(d)
			}
		})
	}
}

func Test_build(t *testing.T) {
	tmpdir := testy.CopyTempDir(t, "testdata", 0)
	t.Cleanup(func() {
		_ = os.RemoveAll(tmpdir)
	})

	tests := []struct {
		name   string
		opts   *buildOpts
		in     io.Reader
		args   []string
		err    string
		target string
	}{
		{
			name: "invalid stdin",
			in:   strings.NewReader("this isn't JSON"),
			err:  "error unmarshaling JSON: json: cannot unmarshal string into Go value of type gobl.Envelope",
		},
		{
			name: "success",
			in: func() io.Reader {
				f, err := os.Open("testdata/nototals.json")
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					_ = f.Close()
				})
				return f
			}(),
		},
		{
			name: "no document",
			in: strings.NewReader(`{
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
			err: "no document included",
		},
		{
			name: "invalid type",
			in: strings.NewReader(`{
				"head": {
					"uuid": "9d8eafd5-77be-11ec-b485-5405db9a3e49",
					"typ": "duck",
					"rgn": "ES",
					"dig": {
						"alg": "sha256",
						"val": "dce3bc3c8bf28f3d209f783917b3082ddc0339a66e9ba3aa63849e4357db1422"
					}
				},
				doc: {
					"walk": "like a duck",
					"talk": "like a duck",
					"look": "like a duck"
				}
			}`),
			err: "unrecognized document type: duck",
		},
		{
			name: "invalid doc",
			in: strings.NewReader(`{
				"head": {
					"uuid": "9d8eafd5-77be-11ec-b485-5405db9a3e49",
					"typ": "bill.Invoice",
					"rgn": "ES",
					"dig": {
						"alg": "sha256",
						"val": "dce3bc3c8bf28f3d209f783917b3082ddc0339a66e9ba3aa63849e4357db1422"
					}
				},
				doc: "foo bar baz"
			}`),
			err: "json: cannot unmarshal string into Go value of type bill.Invoice",
		},
		{
			name: "incomplete",
			in: strings.NewReader(`{
				"head": {
					"uuid": "9d8eafd5-77be-11ec-b485-5405db9a3e49",
					"typ": "bill.Invoice",
					"rgn": "ES",
					"dig": {
						"alg": "sha256",
						"val": "dce3bc3c8bf28f3d209f783917b3082ddc0339a66e9ba3aa63849e4357db1422"
					}
				},
				doc: {}
			}`),
			err: "calculation: issue date cannot be empty",
		},
		{
			name: "input file",
			args: []string{"testdata/success.json"},
		},
		{
			name: "recalculate",
			args: []string{"testdata/nototals.json"},
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
			opts: &buildOpts{
				overwriteOutputFile: true,
			},
			args:   []string{"testdata/success.json", filepath.Join(tmpdir, "overwrite.json")},
			target: filepath.Join(tmpdir, "overwrite.json"),
		},
		{
			name: "overwrite input file",
			opts: &buildOpts{
				inPlace: true,
			},
			args:   []string{filepath.Join(tmpdir, "input.json")},
			target: filepath.Join(tmpdir, "input.json"),
		},
		{
			name: "overwrite stdin",
			opts: &buildOpts{
				inPlace: true,
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
				opts = &buildOpts{}
			}
			err := opts.RunE(c, tt.args)
			if tt.err != "" {
				assert.EqualError(t, err, tt.err)
			} else {
				assert.Nil(t, err)
			}
			if d := testy.DiffText(testy.Snapshot(t), buf.String()); d != nil {
				t.Error(d)
			}
			if tt.target != "" {
				result, err := ioutil.ReadFile(tt.target)
				if err != nil {
					t.Fatal(err)
				}
				if d := testy.DiffText(testy.Snapshot(t, "outfile"), result); d != nil {
					t.Error(d)
				}
			}
		})
	}
}
