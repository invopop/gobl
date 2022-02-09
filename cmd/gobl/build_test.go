package main

import (
	"bytes"
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

func Test_build_args(t *testing.T) {
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
		{
			name: "set values",
			args: []string{"--set", "foo=bar", "--set", "bar=baz", "--set", "foo=qux"},
		},
		{
			name: "set files",
			args: []string{"--set-file", "foo=foo.json"},
		},
		{
			name: "set string values",
			args: []string{"--set-string", "foo=foo", "--set-string", "bar=1234"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			opts := build()

			cmd := opts.cmd()
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

func Test_build_preRun(t *testing.T) {
	tests := []struct {
		name string
		opts *buildOpts
		err  string
	}{
		{
			name: "invalid yaml value on command line",
			opts: &buildOpts{
				set: map[string]string{"foo": ":"},
			},
			err: `yaml: did not find expected key`,
		},
		{
			name: "valid yaml on command line",
			opts: &buildOpts{
				set: map[string]string{
					"string":  "one two three",
					"number":  "123",
					"boolean": "false",
					"array":   "[one,two,three]",
				},
			},
		},
		{
			name: "valid string",
			opts: &buildOpts{
				setStrings: map[string]string{
					"string":  "one two three",
					"number":  "123",
					"boolean": "false",
					"array":   "[one,two,three]",
				},
			},
		},
		{
			name: "missing file",
			opts: &buildOpts{
				setFiles: map[string]string{
					"foo": "missing.yaml",
				},
			},
			err: `open missing.yaml: no such file or directory`,
		},
		{
			name: "valid file",
			opts: &buildOpts{
				setFiles: map[string]string{
					"foo": "testdata/supplier.yaml",
				},
			},
		},
		{
			name: "nested key",
			opts: &buildOpts{
				set: map[string]string{
					"one.two.three": "123",
				},
			},
		},
		{
			name: "merge",
			opts: &buildOpts{
				set: map[string]string{
					"one.two.three": "123",
					"one.two.four":  "124",
				},
			},
		},
		{
			name: "root key",
			opts: &buildOpts{
				setFiles: map[string]string{".": "testdata/exists.json"},
			},
		},
		{
			name: "literal period",
			opts: &buildOpts{
				setFiles: map[string]string{`\.`: "testdata/exists.json"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.preRunE(nil, nil)
			if tt.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
			}
			if err != nil {
				return
			}
			if d := testy.DiffInterface(testy.Snapshot(t), tt.opts); d != nil {
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
			name: "invalid yaml value on command line",
			opts: &buildOpts{
				set: map[string]string{"foo": ":"},
			},
			err: `yaml: did not find expected key`,
		},
		{
			name: "invalid stdin",
			in:   strings.NewReader("this isn't JSON"),
			err:  "code=400, message=yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `this is...` into map[string]interface {}",
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
			err: "code=422, message=no document included",
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
			err: "code=422, message=unrecognized document type: duck",
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
			err: "code=422, message=json: cannot unmarshal string into Go value of type bill.Invoice",
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
			err: "code=422, message=validation: code: cannot be blank; currency: cannot be blank; issue_date: required; lines: cannot be blank; supplier: cannot be blank.",
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
		{
			name: "merge values",
			opts: &buildOpts{
				setValues: map[string]interface{}{
					"doc": map[string]interface{}{
						"currency": "MXN",
					},
				},
			},
			args: []string{"testdata/success.json"},
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
			err := opts.runE(c, tt.args)
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
