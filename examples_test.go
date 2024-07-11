package gobl_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var skipExamplePaths = []string{
	"build/",
	".out.",
	"/out/",
	"data/",
	".github",
}

var updateExamples = flag.Bool("update", false, "Update the examples in the repository")

// TestConvertExamplesToJSON finds all of the `.json` and `.yaml` files in the
// package and attempts to convert the to JSON Envelopes.
func TestConvertExamplesToJSON(t *testing.T) {
	// Find all .yaml files in subdirectories
	var files []string
	err := filepath.Walk("./", func(path string, _ os.FileInfo, _ error) error {
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
	data, err := os.ReadFile(path)
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
		doc := new(schema.Object)
		if err := yaml.Unmarshal(data, doc); err != nil {
			return fmt.Errorf("invalid contents: %w", err)
		}
		env, err = gobl.Envelop(doc)
		if err != nil {
			return fmt.Errorf("failed to envelop: %w", err)
		}
	}

	// override the UUID
	env.Head.UUID = uuid.MustParse("8a51fd30-2a27-11ee-be56-0242ac120002")

	if err := env.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Output to the filesystem in the /out/ directory
	out, err := json.MarshalIndent(env, "", "	")
	if err != nil {
		return fmt.Errorf("marshalling output: %w", err)
	}

	dir := filepath.Join(filepath.Dir(path), "out")
	of := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) + ".json"
	np := filepath.Join(dir, of)
	if _, err := os.Stat(np); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("checking file: %s: %w", np, err)
		}
		if !*updateExamples {
			return fmt.Errorf("output file missing, run tests with `--update` flag to create")
		}
	}

	if *updateExamples {
		if err := os.WriteFile(np, out, 0644); err != nil {
			return fmt.Errorf("saving file data: %w", err)
		}
		t.Logf("wrote file: %v", np)
	} else {
		// Compare to existing file
		existing, err := os.ReadFile(np)
		if err != nil {
			return fmt.Errorf("reading existing file: %w", err)
		}
		t.Run(np, func(t *testing.T) {
			assert.JSONEq(t, string(existing), string(out), "output file does not match, run tests with `--update` flag to update")
		})
	}

	return nil
}
