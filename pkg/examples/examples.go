// Package examples provides reusable helpers for converting GOBL example
// documents into calculated, validated JSON envelopes and locating their
// golden output files. It captures the conventions used by GOBL's own example
// suite (a fixed envelope UUID, envelope-vs-document detection, and tab-indented
// output) so that external addon and converter modules can ship example
// documents tested exactly the same way.
//
// The companion subpackage pkg/examples/exampletest wraps these helpers into a
// single test entry point.
package examples

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/yaml"
)

// TestUUID is assigned to every converted envelope so that golden output is
// deterministic regardless of the source document's own UUID.
var TestUUID = uuid.MustParse("8a51fd30-2a27-11ee-be56-0242ac120002")

// defaultSkip are path fragments always excluded from example discovery.
var defaultSkip = []string{"/out/", ".out.", "/.git"}

// Convert reads a GOBL example (YAML or JSON) and returns the calculated,
// validated envelope as tab-indented JSON. When asEnvelope is true the input is
// parsed as a full gobl.Envelope and calculated; otherwise it is treated as a
// bare document and wrapped with gobl.Envelop. The envelope UUID is forced to
// TestUUID so the output is stable for golden comparison.
func Convert(data []byte, asEnvelope bool) ([]byte, error) {
	var env *gobl.Envelope
	if asEnvelope {
		env = new(gobl.Envelope)
		if err := yaml.Unmarshal(data, env); err != nil {
			return nil, fmt.Errorf("invalid contents: %w", err)
		}
		if err := env.Calculate(); err != nil {
			return nil, fmt.Errorf("failed to complete: %w", err)
		}
	} else {
		doc := new(schema.Object)
		if err := yaml.Unmarshal(data, doc); err != nil {
			return nil, fmt.Errorf("invalid contents: %w", err)
		}
		var err error
		env, err = gobl.Envelop(doc)
		if err != nil {
			return nil, fmt.Errorf("failed to envelop: %w", err)
		}
	}

	env.Head.UUID = TestUUID

	if err := env.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	out, err := json.MarshalIndent(env, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("marshalling output: %w", err)
	}
	return out, nil
}

// IsEnvelope reports whether the source path holds a full envelope (by the
// `.env.` naming convention) rather than a bare document.
func IsEnvelope(path string) bool {
	return strings.Contains(path, ".env.")
}

// GoldenPath returns the expected output path for a source example: a JSON file
// of the same base name inside a sibling `out/` directory.
func GoldenPath(src string) string {
	base := strings.TrimSuffix(filepath.Base(src), filepath.Ext(src)) + ".json"
	return filepath.Join(filepath.Dir(src), "out", base)
}

// Sources walks root and returns the example source files (.yaml and .json),
// skipping generated output, hidden git paths, and any path containing one of
// the extra skip fragments. Results are returned in lexical order.
func Sources(root string, skip ...string) ([]string, error) {
	all := append(append([]string{}, defaultSkip...), skip...)
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		switch filepath.Ext(path) {
		case ".yaml", ".json":
			for _, s := range all {
				if strings.Contains(path, s) {
					return nil
				}
			}
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
