//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
)

func main() {
	if error := generate(); error != nil {
		panic(error)
	}
}

// generate will output the JSON definitions of the catalogues to the data directory.
// Please not that in the case of Catalogues specifically, the source data is the JSON
// output. This implies that any changes to structures or refactoring will be reflected
// in the output, despite having the same source.
func generate() error {
	for _, cd := range tax.AllCatalogueDefs() {
		doc, err := schema.NewObject(cd)
		if err != nil {
			return err
		}
		data, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			return err
		}
		n := string(cd.Key)
		f := filepath.Join("data", "catalogues", n+".json")
		if err := os.WriteFile(f, data, 0644); err != nil {
			return err
		}
		fmt.Printf("Processed %v\n", f)
	}
	return nil
}
