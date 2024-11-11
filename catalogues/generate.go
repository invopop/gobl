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
