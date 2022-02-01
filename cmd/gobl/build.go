package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/invopop/gobl"
)

type buildOpts struct {
	overwriteOutputFile bool
	inPlace             bool
	set                 map[string]string
	setFiles            map[string]string
	setStrings          map[string]string
	// setValues contains the parsed values from `set`, `setFiles`, and
	// `setStrings`, ready to be merged into the GOBL document in RunE.
	setValues map[string]interface{}
}

func build() *buildOpts {
	return &buildOpts{}
}

func (b *buildOpts) preRunE(*cobra.Command, []string) error {
	b.setValues = make(map[string]interface{}, len(b.set)+len(b.setFiles)+len(b.setStrings))
	for k, v := range b.setStrings {
		if err := b.setValue(k, v); err != nil {
			return err
		}
	}
	for k, v := range b.set {
		var val interface{}
		if err := yaml.Unmarshal([]byte(v), &val); err != nil {
			return err
		}
		if err := b.setValue(k, val); err != nil {
			return err
		}
	}
	for k, v := range b.setFiles {
		content, err := ioutil.ReadFile(v)
		if err != nil {
			return err
		}
		var val interface{}
		if err := yaml.Unmarshal(content, &val); err != nil {
			return err
		}
		if err := b.setValue(k, val); err != nil {
			return err
		}
	}
	return nil
}

func (b *buildOpts) setValue(key string, value interface{}) error {
	for {
		i := strings.LastIndex(key, ".")
		if i == -1 {
			break
		}
		value = map[string]interface{}{
			key[i+1:]: value,
		}
		key = key[:i]
	}
	return mergo.Merge(&b.setValues, map[string]interface{}{
		key: value,
	}, mergo.WithOverride)
}

func (b *buildOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "build [infile] [outfile]",
		Args:    cobra.MaximumNArgs(2),
		PreRunE: b.preRunE,
		RunE:    b.runE,
	}

	f := cmd.Flags()

	f.BoolVarP(&b.overwriteOutputFile, "force", "f", false, "force writing output file, even if it exists (only outputs JSON)")
	f.BoolVarP(&b.inPlace, "in-place", "w", false, "overwrite the input file in place")
	f.StringToStringVar(&b.set, "set", nil, "set value from the command line")
	f.StringToStringVar(&b.setFiles, "set-file", nil, "set value from the specified YAML or JSON file")
	f.StringToStringVar(&b.setStrings, "set-string", nil, "set STRING value from the command line")

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

type ctxReader struct {
	context.Context
	io.Reader
}

func (r *ctxReader) Read(p []byte) (int, error) {
	var c int
	var err error
	wait := make(chan struct{}, 1)
	go func() {
		c, err = r.Reader.Read(p)
		close(wait)
	}()
	select {
	case <-r.Context.Done():
		return 0, r.Context.Err()
	case <-wait:
		return c, err
	}
}

func cmdContext(cmd *cobra.Command) context.Context {
	if ctx := cmd.Context(); ctx != nil {
		return ctx
	}
	return context.Background()
}

func (b *buildOpts) runE(cmd *cobra.Command, args []string) error {
	input, err := openInput(cmd, args)
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
	defer input.Close() // nolint:errcheck

	in := &ctxReader{
		Context: cmdContext(cmd),
		Reader:  input,
	}
	var intermediate map[string]interface{}
	if err := yaml.NewDecoder(in).Decode(&intermediate); err != nil {
		return err
	}
	if err := mergo.Merge(&intermediate, b.setValues, mergo.WithOverride); err != nil {
		return err
	}
	encoded, err := json.Marshal(intermediate)
	if err != nil {
		return err
	}
	env := new(gobl.Envelope)
	if err := json.Unmarshal(encoded, &env); err != nil {
		return err
	}

	if err := reInsertDoc(env); err != nil {
		return err
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "\t")
	return enc.Encode(env)
}

func reInsertDoc(env *gobl.Envelope) error {
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
	return nil
}
