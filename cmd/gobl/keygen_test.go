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

var jwkREs = []testy.Replacement{
	{
		Regexp:      regexp.MustCompile(`"kid": ?"[^"]*"`),
		Replacement: `"kid":"..."`,
	},
	{
		Regexp:      regexp.MustCompile(`"x": ?"[^"]*"`),
		Replacement: `"x":"..."`,
	},
	{
		Regexp:      regexp.MustCompile(`"y": ?"[^"]*"`),
		Replacement: `"y":"..."`,
	},
	{
		Regexp:      regexp.MustCompile(`"d": ?"[^"]*"`),
		Replacement: `"d":"..."`,
	},
}

func TestExpandHome(t *testing.T) {
	// Without the leading ~/, returned verbatim.
	got, err := expandHome("/abs/path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/abs/path" {
		t.Errorf("expected pass-through, got %q", got)
	}

	// With ~/, resolved against $HOME.
	t.Setenv("HOME", "/tmp/home")
	got, err = expandHome("~/sub/key.jwk")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/tmp/home/sub/key.jwk" {
		t.Errorf("expected expanded path, got %q", got)
	}
}

func TestHomedir(t *testing.T) {
	t.Setenv("HOME", "/tmp/home")
	got, err := homedir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/tmp/home" {
		t.Errorf("expected /tmp/home, got %q", got)
	}
}

func TestHomedirNoHOME(t *testing.T) {
	// HOME unset -> falls back to user.Current(). Whichever path the
	// underlying machine takes, the function should return a value
	// (or an error). Either way the fallback branch is exercised.
	t.Setenv("HOME", "")
	_, _ = homedir()
}

func TestKeygenRunEPubFileExists(t *testing.T) {
	// Pre-create the .pub.jwk file so writeKey on the public side
	// errors out, exercising the second writeKey error branch.
	dir := t.TempDir()
	priv := filepath.Join(dir, "id_test")
	pub := pubfileFromPriv(priv)
	if err := os.WriteFile(pub, []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}
	o := &keygenOpts{rootOpts: &rootOpts{}}
	c := &cobra.Command{}
	c.SetOut(&bytes.Buffer{})
	err := o.runE(c, []string{priv})
	if err == nil {
		t.Fatal("expected error when pub file pre-exists")
	}
}

func TestWriteKeyMkdirFailure(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("write-permission tests do not apply when running as root")
	}
	// Make a non-writable parent so MkdirAll fails inside writeKey.
	dir := t.TempDir()
	ro := filepath.Join(dir, "ro")
	if err := os.MkdirAll(ro, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(ro, 0o755) })

	err := writeKey(filepath.Join(ro, "sub", "id.jwk"), []byte("data"), 0o600, false)
	if err == nil {
		t.Fatal("expected mkdir failure")
	}
}

func TestWriteKeyDefaultDirMkdirFailure(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("write-permission tests do not apply when running as root")
	}
	// Point HOME at a non-writable parent so the default keyfile's
	// containing dir cannot be created, triggering the L119 MkdirAll
	// branch inside writeKey.
	parent := t.TempDir()
	if err := os.Chmod(parent, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(parent, 0o755) })
	t.Setenv("HOME", parent)

	def, err := defaultKeyfile()
	if err != nil {
		t.Fatal(err)
	}
	if err := writeKey(def, []byte("data"), 0o600, false); err == nil {
		t.Fatal("expected mkdir failure for default keyfile dir")
	}
}

func TestPubfileFromPriv(t *testing.T) {
	got := pubfileFromPriv("id.jwk")
	if got != "id.pub.jwk" {
		t.Errorf("got %q", got)
	}
}

func TestOutputKeyfile(t *testing.T) {
	got, err := outputKeyfile([]string{"a.jwk"})
	if err != nil {
		t.Fatal(err)
	}
	if got != "a.jwk" {
		t.Errorf("got %q", got)
	}
}

func Test_keygen(t *testing.T) {
	type tt struct {
		env  map[string]string
		opts *keygenOpts
		args []string
		err  string
	}

	tests := testy.NewTable()
	tests.Add("stdout", tt{
		args: []string{"-"},
	})
	tests.Add("indent", tt{
		opts: &keygenOpts{
			rootOpts: &rootOpts{
				indent: true,
			},
		},
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
	tests.Add("Create .gobl dir", func(t *testing.T) interface{} {
		tmp := t.TempDir()

		return tt{
			env: map[string]string{
				"HOME": tmp,
			},
		}
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		if len(tt.env) > 0 {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
		}
		c := &cobra.Command{}
		buf := &bytes.Buffer{}
		c.SetOut(buf)
		opts := tt.opts
		if opts == nil {
			opts = &keygenOpts{}
		}
		if opts.rootOpts == nil {
			opts.rootOpts = &rootOpts{}
		}
		err := opts.runE(c, tt.args)
		if !testy.ErrorMatchesRE(tt.err, err) {
			t.Errorf("Unexpected error: %s", err)
		}
		if err != nil {
			return
		}

		if d := testy.DiffText(testy.Snapshot(t), buf.String(), jwkREs...); d != nil {
			t.Error(d)
		}

		outfile, err := outputKeyfile(tt.args)
		if err != nil {
			t.Fatal(err)
		}

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

		pub, err := os.Open(pubfileFromPriv(outfile))
		if err != nil {
			t.Fatal(err)
		}
		_, err = pub.Stat()
		if err != nil {
			t.Fatal(err)
		}
	})
}
