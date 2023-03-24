package gobl_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var skipExamplePaths = []string{
	"build/",
	".out.",
	".github",
}

// TestConvertExamplesToJSON finds all of the `.json` and `.yaml` files in the
// package and attempts to convert the to JSON Envelopes.
func TestConvertExamplesToJSON(t *testing.T) {
	// Find all .yaml files in subdirectories
	var files []string
	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		switch filepath.Ext(path) {
		case ".yaml", ".json":
			for _, skip := range skipExamplePaths {
				if strings.Contains(path, skip) {
					return nil
				}
			}
			files = append(files, path)
		}
		return nil
	})
	require.NoError(t, err)

	for _, path := range files {
		assert.NoError(t, processFile(t, path))
	}
}

func processFile(t *testing.T, path string) error {
	t.Helper()
	t.Logf("processing file: %v", path)

	// attempt to load and convert
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	var env *gobl.Envelope
	if strings.Contains(path, ".env.") {
		// Handle Envelopes
		env = new(gobl.Envelope)
		if err := yaml.Unmarshal(data, env); err != nil {
			return fmt.Errorf("invalid contents: %w", err)
		}
		if err := env.Calculate(); err != nil {
			return fmt.Errorf("failed to complete: %w", err)
		}
	} else {
		// Handle documents
		doc := new(gobl.Document)
		if err := yaml.Unmarshal(data, doc); err != nil {
			return fmt.Errorf("invalid contents: %w", err)
		}
		env, err = gobl.Envelop(doc)
		if err != nil {
			return fmt.Errorf("failed to envelop: %w", err)
		}
	}

	if err := env.Sign(testKey); err != nil {
		return fmt.Errorf("failed to sign the doc: %w", err)
	}

	// Output to the filesystem (.out.json is defined in .gitignore)
	np := strings.TrimSuffix(path, filepath.Ext(path)) + ".out.json"
	out, err := json.MarshalIndent(env, "", "	")
	if err != nil {
		return fmt.Errorf("marshalling output: %w", err)
	}
	if err := ioutil.WriteFile(np, out, 0644); err != nil {
		return fmt.Errorf("saving file data: %w", err)
	}

	t.Logf("wrote file: %v", np)
	return nil
}
