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

func Test_replicate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *replicateOpts
		in      io.Reader
		args    []string
		replace []testy.Replacement
		err     string
	}{
		{
			name: "invalid stdin",
			in:   strings.NewReader("this isn't JSON"),
			err:  "input: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `this is...` into map[string]interface {}",
		},
		{
			name: "success",
			args: []string{"testdata/success.json"},
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
				opts = &replicateOpts{}
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
