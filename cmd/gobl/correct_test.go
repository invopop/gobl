package main

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"gitlab.com/flimzy/testy"
)

func Test_correct_args(t *testing.T) {
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
			name: "options",
			args: []string{"--options"},
		},
		{
			name: "credit",
			args: []string{"--credit"},
		},
		{
			name: "debit",
			args: []string{"--debit"},
		},
		{
			name: "data",
			args: []string{"--data", `{"type":"credit"}`},
		},
		{
			name: "data short",
			args: []string{"-d", `{"type":"credit"}`},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			root := root()
			rootCmd := &cobra.Command{}
			root.setFlags(rootCmd)

			opts := correct(root)

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

func Test_correct(t *testing.T) {
	tests := []struct {
		name    string
		opts    *correctOpts
		in      io.Reader
		args    []string
		replace []testy.Replacement
		err     string
	}{
		{
			name: "invalid stdin",
			in:   strings.NewReader("this isn't JSON"),
			err:  "code=400, message=yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `this is...` into map[string]interface {}",
		},
		{
			name: "options schema",
			args: []string{"testdata/success.json"},
			opts: &correctOpts{
				options: true,
			},
		},
		{
			name: "credit note",
			args: []string{"testdata/success.json"},
			opts: &correctOpts{
				credit: true,
			},
			replace: []testy.Replacement{
				{
					Regexp:      regexp.MustCompile(`"uuid":\s*"[^"]*"`),
					Replacement: `"uuid": "00000000-0000-0000-0000-000000000000"`,
				},
				{
					Regexp:      regexp.MustCompile(`"val":\s*"[a-f0-9]+"`),
					Replacement: `"val": "0000000000000000000000000000000000000000000000000000000000000000"`,
				},
				{
					Regexp:      regexp.MustCompile(`"issue_date":\s*"[^"]*"`),
					Replacement: `"issue_date": "2024-01-01"`,
				},
			},
		},
		{
			name: "file missing",
			args: []string{"asdf"},
			err:  "open asdf: no such file or directory",
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
				opts = &correctOpts{}
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

			if d := testy.DiffText(testy.Snapshot(t), buf.String(), tt.replace...); d != nil {
				t.Error(d)
			}
		})
	}
}
