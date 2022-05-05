package samples_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var signingKey = dsig.NewES256Key()

func TestConvertSamplesToJSON(t *testing.T) {
	// Find all .yaml files in subdirectories
	var files []string
	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".yaml" {
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
	if strings.Contains(path, ".env.yaml") {
		// Handle Envelopes
		env = new(gobl.Envelope)
		if err := yaml.Unmarshal(data, env); err != nil {
			return fmt.Errorf("invalid contents: %w", err)
		}
		if err := env.Complete(); err != nil {
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

	if err := env.Sign(signingKey); err != nil {
		return fmt.Errorf("failed to sign the doc: %w", err)
	}

	// Output to the filesystem
	np := strings.TrimSuffix(path, filepath.Ext(path)) + ".json"
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
