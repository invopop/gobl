package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"gitlab.com/flimzy/testy"
)

func Test_bulk(t *testing.T) {
	type tt struct {
		in     io.Reader
		opts   *bulkOpts
		args   []string
		err    string
		target string
	}

	tests := testy.NewTable()
	tests.Add("no input", tt{
		in: strings.NewReader(""),
	})
	tests.Add("non json", tt{
		in: strings.NewReader("not json"),
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		t.Parallel()

		c := &cobra.Command{}
		if tt.in != nil {
			c.SetIn(tt.in)
		}
		buf := &bytes.Buffer{}
		c.SetOut(buf)
		opts := tt.opts
		if opts == nil {
			opts = &bulkOpts{}
		}
		if opts.rootOpts == nil {
			opts.rootOpts = &rootOpts{}
		}
		opts.rootOpts.indent = true
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
			if d := testy.DiffText(testy.Snapshot(t, "outfile"), &testy.File{Path: tt.target}); d != nil {
				t.Errorf("outfile:\n%s", d)
			}
		}
	})
}
