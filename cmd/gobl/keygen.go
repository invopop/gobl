package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/invopop/gobl/dsig"
)

const defaultKeyFilename = "~/.gobl/id_es256.jwk"

type keygenOpts struct {
	*rootOpts
	overwrite bool
}

func keygen(root *rootOpts) *keygenOpts {
	return &keygenOpts{
		rootOpts: root,
	}
}

func (k *keygenOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keygen [flags] [outfile]",
		Short: "Generate a keypair",
		Args:  cobra.MaximumNArgs(1),
		RunE:  k.runE,
	}

	f := cmd.Flags()

	f.BoolVarP(&k.overwrite, "force", "f", false, "force writing output file, even if it exists")

	return cmd
}

func expandHome(in string) (string, error) {
	if !strings.HasPrefix(in, "~/") {
		return in, nil
	}
	home, err := homedir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, strings.TrimPrefix(in, "~/")), nil
}

func homedir() (string, error) {
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir, nil
}

func defaultKeyfile() (string, error) {
	return expandHome(defaultKeyFilename)
}

func outputKeyfile(args []string) (string, error) {
	if len(args) == 0 {
		return defaultKeyfile()
	}
	return args[0], nil
}

func pubfileFromPriv(priv string) string {
	return strings.TrimSuffix(priv, ".jwk") + ".pub.jwk"
}

func (k *keygenOpts) runE(cmd *cobra.Command, args []string) error {
	key := dsig.NewES256Key()
	marshal := json.Marshal
	if k.indent {
		marshal = func(i interface{}) ([]byte, error) {
			return json.MarshalIndent(i, "", "\t")
		}
	}
	priv, err := marshal(key)
	if err != nil {
		return err
	}
	pub, err := marshal(key.Public())
	if err != nil {
		return err
	}
	outfile, err := outputKeyfile(args)
	if err != nil {
		return err
	}
	if outfile == "-" {
		_, err := fmt.Fprintln(cmd.OutOrStdout(), string(priv))
		return err
	}
	if err = writeKey(outfile, priv, 0o600, k.overwrite); err != nil {
		return err
	}
	if err = writeKey(pubfileFromPriv(outfile), pub, 0o666, k.overwrite); err != nil {
		return err
	}
	return nil
}

func writeKey(filename string, key []byte, mode os.FileMode, force bool) error {
	dir, base := filepath.Dir(filename), filepath.Base(filename)
	def, _ := defaultKeyfile()
	if dir == filepath.Dir(def) {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return err
		}
	}
	tmp, err := os.CreateTemp(dir, "."+base+"-*")
	if err != nil {
		return err
	}
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
	}()

	if err = tmp.Chmod(mode); err != nil {
		return err
	}
	if _, err = tmp.Write(key); err != nil {
		return err
	}
	return safeRename(tmp.Name(), filename, force)
}

func safeRename(old, in string, force bool) error {
	if force {
		return os.Rename(old, in)
	}
	err := os.Link(old, in)
	if err != nil {
		return fmt.Errorf("target %q exists", in)
	}
	return os.Remove(old)
}
