package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/invopop/gobl/dsig"
)

type keygenOpts struct {
	overwrite bool
}

func keygen() *keygenOpts {
	return &keygenOpts{}
}

func (k *keygenOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "keygen [flags] [outfile]",
		Args: cobra.MaximumNArgs(1),
		RunE: k.runE,
	}

	f := cmd.Flags()

	f.BoolVarP(&k.overwrite, "force", "f", false, "force writing output file, even if it exists")

	return cmd
}

func outputKeyfile(args []string) string {
	if len(args) == 0 {
		return "~/.gobl/id_es256"
	}
	return args[0]
}

func (k *keygenOpts) runE(cmd *cobra.Command, args []string) error {
	key := dsig.NewES256Key()
	priv, err := json.Marshal(key)
	if err != nil {
		return err
	}
	pub, err := json.Marshal(key.Public())
	if err != nil {
		return err
	}
	outfile := outputKeyfile(args)
	if outfile == "-" {
		fmt.Fprintln(cmd.OutOrStdout(), string(priv))
		return nil
	}
	if err = writeKey(outfile, priv, 0o600, k.overwrite); err != nil {
		return err
	}
	if err = writeKey(outfile+".pub", pub, 0o666, k.overwrite); err != nil {
		return err
	}
	return nil
}

func writeKey(filename string, key []byte, mode os.FileMode, force bool) error {
	dir, base := filepath.Dir(filename), filepath.Base(filename)
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

func safeRename(old, new string, force bool) error {
	if force {
		return os.Rename(old, new)
	}
	err := os.Link(old, new)
	if err != nil {
		return fmt.Errorf("target %q exists", new)
	}
	return os.Remove(old)
}
