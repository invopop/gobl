package main

import (
	"errors"
	"io"
	"os"

	"github.com/spf13/cobra"
)

type rootOpts struct {
	indent              bool // when true, indent output, mainly for testing
	overwriteOutputFile bool
	inPlace             bool
}

func root() *rootOpts {
	return &rootOpts{}
}

func (o *rootOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "gobl",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	o.setFlags(cmd)

	cmd.AddCommand(verify().cmd())
	cmd.AddCommand(validate(o).cmd())
	cmd.AddCommand(build(o).cmd())
	cmd.AddCommand(sign(o).cmd())
	cmd.AddCommand(correct(o).cmd())
	cmd.AddCommand(replicate(o).cmd())
	cmd.AddCommand(versionCmd())
	cmd.AddCommand(serve().cmd())
	cmd.AddCommand(keygen(o).cmd())
	return cmd
}

func (o *rootOpts) setFlags(cmd *cobra.Command) {
	f := cmd.PersistentFlags()
	f.BoolVarP(&o.indent, "indent", "i", false, "format JSON output with indentation")
	f.BoolVarP(&o.overwriteOutputFile, "force", "f", false, "force writing output file, even if it exists")
	f.BoolVarP(&o.inPlace, "in-place", "w", false, "overwrite the input file in place  (only outputs JSON)")
}

func (o *rootOpts) outputFilename(args []string) string {
	if o.inPlace {
		return inputFilename(args)
	}
	if len(args) >= 2 && args[1] != "-" {
		return args[1]
	}
	return ""
}

func openInput(cmd *cobra.Command, args []string) (io.ReadCloser, error) {
	if inFile := inputFilename(args); inFile != "" {
		return os.Open(inFile)
	}
	return io.NopCloser(cmd.InOrStdin()), nil
}

func (o *rootOpts) openOutput(cmd *cobra.Command, args []string) (io.WriteCloser, error) {
	if outFile := o.outputFilename(args); outFile != "" {
		flags := os.O_CREATE | os.O_WRONLY
		if !o.overwriteOutputFile && !o.inPlace {
			flags |= os.O_EXCL
		}
		return os.OpenFile(outFile, flags, os.ModePerm)
	}
	if o.inPlace {
		return nil, errors.New("cannot overwrite STDIN")
	}
	return writeCloser{cmd.OutOrStdout()}, nil
}

type writeCloser struct {
	io.Writer
}

func (writeCloser) Close() error { return nil }
