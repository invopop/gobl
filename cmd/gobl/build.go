package main

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

type buildOpts struct {
	overwriteOutputFile bool
	inPlace             bool
}

func build() *cobra.Command {
	opts := &buildOpts{}
	cmd := &cobra.Command{
		Use:  "build [infile] [outfile]",
		Args: cobra.MaximumNArgs(2),
		RunE: opts.RunE,
	}

	f := cmd.Flags()

	f.BoolVarP(&opts.overwriteOutputFile, "force", "f", false, "force writing output file, even if it exists")
	f.BoolVarP(&opts.inPlace, "in-place", "w", false, "overwrite the input file in place")

	return cmd
}

func (b *buildOpts) outputFilename(args []string) string {
	if b.inPlace {
		return inputFilename(args)
	}
	if len(args) >= 2 && args[1] != "-" {
		return args[1]
	}
	return ""
}

func (b *buildOpts) RunE(cmd *cobra.Command, args []string) error {
	env, err := readEnv(cmd, args)
	if err != nil {
		return err
	}
	out := cmd.OutOrStdout()
	if outFile := b.outputFilename(args); outFile != "" {
		flags := os.O_CREATE | os.O_WRONLY
		if !b.overwriteOutputFile && !b.inPlace {
			flags |= os.O_EXCL
		}
		f, err := os.OpenFile(outFile, flags, os.ModePerm)
		if err != nil {
			return err
		}
		defer f.Close() // nolint:errcheck
		out = f
	} else if b.inPlace {
		return errors.New("cannot overwrite STDIN")
	}
	if env.Document == nil {
		return errors.New("no document included")
	}
	doc, err := extractDoc(env)
	if err != nil {
		return err
	}
	if err := env.Insert(doc); err != nil {
		return err
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "\t")
	return enc.Encode(env)
}
