// Package schemas helps generate JSON Schema files from the main GOBL packages.
package schemas

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/jsonschema"
)

const outPath = "./build/schemas"

// Generate is used to generate a set of schema files from the GOBL bases.
func Generate() error {
	r := new(jsonschema.Reflector)
	r.AllowAdditionalProperties = true

	if err := r.AddGoComments("github.com/invopop/gobl", "./"); err != nil {
		return fmt.Errorf("reading comments: %w", err)
	}

	typs := schema.Types()
	r.Lookup = func(t reflect.Type) jsonschema.ID {
		id, ok := typs[t]
		if ok {
			return jsonschema.ID(id.String())
		}
		return jsonschema.EmptyID
	}

	comment := fmt.Sprintf("Generated with GOBL %v", gobl.VERSION)

	// Cleanup the old
	if err := os.RemoveAll(outPath); err != nil {
		return fmt.Errorf("unable to remove old data: %w", err)
	}

	for t, id := range schema.Types() {
		fmt.Printf("processing %v... ", id)
		s := r.ReflectFromType(t)
		s.Comments = comment

		f := strings.TrimPrefix(id.String(), schema.GOBL.String()) + ".json"
		f = filepath.Join(outPath, f)

		d, err := json.MarshalIndent(s, "", "  ")
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(f), os.ModePerm); err != nil {
			return fmt.Errorf("unable to create directories: %w", err)
		}

		if err := ioutil.WriteFile(f, d, 0644); err != nil {
			return fmt.Errorf("unable to write file '%v': %w", f, err)
		}

		fmt.Printf("wrote: %v\n", f)
	}

	return nil
}
