package gobl_test

import (
	"flag"
	"testing"

	"github.com/invopop/gobl/pkg/examples/exampletest"
)

// skipExamplePaths excludes non-example files from the repo-wide walk. (The
// out/ directories and .git are skipped by the helper itself.)
var skipExamplePaths = []string{
	"build/",
	"data/",
	".github",
	".golangci.yaml",
	"wasm/",
	".claude/",
	"internal/",
	"pkg/",
}

var updateExamples = flag.Bool("update", false, "Update the examples in the repository")

// TestConvertExamplesToJSON finds all of the `.json` and `.yaml` files in the
// package and converts them to JSON Envelopes, comparing against (or updating)
// their golden output. The conversion and comparison logic is shared with
// external addon modules via pkg/examples.
func TestConvertExamplesToJSON(t *testing.T) {
	exampletest.Run(t, ".", *updateExamples, skipExamplePaths...)
}
