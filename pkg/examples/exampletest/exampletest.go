// Package exampletest provides a single test entry point for verifying a
// directory of GOBL example documents: each source is converted to a calculated,
// validated JSON envelope and compared against its golden output (or the golden
// is rewritten when update is true).
//
// It depends only on the standard library and pkg/examples, so importing it
// from a module's _test.go adds no third-party test dependencies. The addon(s)
// under test must already be registered (e.g. via a blank import) in the test
// binary.
package exampletest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/invopop/gobl/pkg/examples"
)

// Run discovers the example sources under root, converts each, and either
// updates the golden output (when update is true) or compares against it. Extra
// skip fragments are passed through to examples.Sources. Each example runs as a
// named subtest.
func Run(t *testing.T, root string, update bool, skip ...string) {
	t.Helper()

	files, err := examples.Sources(root, skip...)
	if err != nil {
		t.Fatalf("discovering examples under %s: %v", root, err)
	}
	if len(files) == 0 {
		t.Fatalf("no example sources found under %s", root)
	}

	for _, path := range files {
		t.Run(path, func(t *testing.T) {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("reading %s: %v", path, err)
			}
			out, err := examples.Convert(data, examples.IsEnvelope(path))
			if err != nil {
				t.Fatalf("converting %s: %v", path, err)
			}

			golden := examples.GoldenPath(path)
			if update {
				if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
					t.Fatalf("creating golden dir for %s: %v", golden, err)
				}
				if err := os.WriteFile(golden, out, 0o644); err != nil {
					t.Fatalf("writing golden %s: %v", golden, err)
				}
				return
			}

			existing, err := os.ReadFile(golden)
			if err != nil {
				t.Fatalf("reading golden %s (run tests with -update to create): %v", golden, err)
			}
			if !jsonEqual(existing, out) {
				t.Errorf("output for %s does not match %s; run tests with -update and review the diff", path, golden)
			}
		})
	}
}

// jsonEqual reports whether two JSON payloads are semantically equal, ignoring
// formatting and object key order.
func jsonEqual(a, b []byte) bool {
	var av, bv any
	if err := json.Unmarshal(a, &av); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &bv); err != nil {
		return false
	}
	return reflect.DeepEqual(av, bv)
}
