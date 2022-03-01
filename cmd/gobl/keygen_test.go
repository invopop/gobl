package main

import (
	"bytes"
	"os"
	"path/filepath"
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
	tests.Add("success", func(t *testing.T) interface{} {
		dir := t.TempDir()

		return tt{
			args: []string{filepath.Join(dir, "id_test")},
		}
	})
	tests.Add("target exists", tt{
		args: []string{"testdata/id_es256"},
		err:  `target "testdata/id_es256" exists`,
	})
	tests.Add("force overwrite", func(t *testing.T) interface{} {
		tmp := t.TempDir()
		f, err := os.Create(filepath.Join(tmp, "id_test"))
		if err != nil {
			t.Fatal(err)
		}
		if err = f.Close(); err != nil {
			t.Fatal(err)
		}

		return tt{
			opts: &keygenOpts{overwrite: true},
			args: []string{f.Name()},
		}
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
		if err != nil {
			return
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

		outfile := outputKeyfile(tt.args)

		if outfile == "-" {
			return
		}

		priv, err := os.Open(outfile)
		if err != nil {
			t.Fatal(err)
		}
		stat, err := priv.Stat()
		if err != nil {
			t.Fatal(err)
		}
		if stat.Mode() != 0o600 {
			t.Errorf("Unexpected file mode on private key file: %v", stat.Mode())
		}

		pub, err := os.Open(outfile + ".pub")
		if err != nil {
			t.Fatal(err)
		}
		_, err = pub.Stat()
		if err != nil {
			t.Fatal(err)
		}
	})
}
