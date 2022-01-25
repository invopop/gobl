package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/divideandconquer/go-merge/merge"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type buildOpts struct {
	overwriteOutputFile bool
	inPlace             bool
	set                 map[string]string
	setFiles            map[string]string
	setStrings          map[string]string
	setValues           map[string]interface{}
}

func build() *buildOpts {
	return &buildOpts{}
}

func (b *buildOpts) preRunE(*cobra.Command, []string) error {
	b.setValues = make(map[string]interface{}, len(b.set)+len(b.setFiles)+len(b.setStrings))
	for k, v := range b.setStrings {
		b.setValue(k, v)
	}
	for k, v := range b.set {
		var val interface{}
		if err := yaml.Unmarshal([]byte(v), &val); err != nil {
			return err
		}
		b.setValue(k, val)
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
		b.setValue(k, val)
	}
	return nil
}

func (b *buildOpts) setValue(key string, value interface{}) {
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
	b.setValues = merge.Merge(b.setValues, map[string]interface{}{
		key: value,
	}).(map[string]interface{})
}

func (b *buildOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "build [infile] [outfile]",
		Args:    cobra.MaximumNArgs(2),
		PreRunE: b.preRunE,
		RunE:    b.runE,
	}

	f := cmd.Flags()

	f.BoolVarP(&b.overwriteOutputFile, "force", "f", false, "force writing output file, even if it exists")
	f.BoolVarP(&b.inPlace, "in-place", "w", false, "overwrite the input file in place")
	f.StringToStringVar(&b.set, "set", nil, "set value from the command line")
	f.StringToStringVar(&b.setFiles, "set-file", nil, "set value from the specified YAML or JSON file")
	f.StringToStringVar(&b.setStrings, "set-string", nil, "set STRING value from the command line")

	return cmd
}

func buildCmd() *cobra.Command {
	return build().cmd()
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

func (b *buildOpts) runE(cmd *cobra.Command, args []string) error {
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
