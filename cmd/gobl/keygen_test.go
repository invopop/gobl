package main

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/spf13/cobra"
	"gitlab.com/flimzy/testy"
)

func Test_keygen(t *testing.T) {
	type tt struct {
		opts *keygenOpts
		args []string
		err  string
	}

	tests := testy.NewTable()
	tests.Add("stdout", tt{
		args: []string{"-"},
	})
	tests.Add("target does not exist", tt{
		args: []string{"/some/path/that/does/not/exist"},
		err:  "open /some/path/that/does/not/.exist-.*: no such file or directory",
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		t.Parallel()

		c := &cobra.Command{}
		buf := &bytes.Buffer{}
		c.SetOut(buf)
		opts := tt.opts
		if opts == nil {
			opts = &keygenOpts{}
		}
		err := opts.runE(c, tt.args)
		if !testy.ErrorMatchesRE(tt.err, err) {
			t.Errorf("Unexpected error: %s", err)
		}

		res := []testy.Replacement{
			{
				Regexp:      regexp.MustCompile(`"kid":"[^"]*"`),
				Replacement: `"kid":"..."`,
			},
			{
				Regexp:      regexp.MustCompile(`"x":"[^"]*"`),
				Replacement: `"x":"..."`,
			},
			{
				Regexp:      regexp.MustCompile(`"y":"[^"]*"`),
				Replacement: `"y":"..."`,
			},
			{
				Regexp:      regexp.MustCompile(`"d":"[^"]*"`),
				Replacement: `"d":"..."`,
			},
		}

		if d := testy.DiffText(testy.Snapshot(t), buf.String(), res...); d != nil {
			t.Error(d)
		}
	})
}
