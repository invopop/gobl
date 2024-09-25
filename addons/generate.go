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
	for _, ao := range tax.AllAddonDefs() {
		doc, err := schema.NewObject(ao)
		if err != nil {
			return err
		}
		data, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			return err
		}
		n := string(ao.Key)
		f := filepath.Join("data", "addons", n+".json")
		if err := os.WriteFile(f, data, 0644); err != nil {
			return err
		}
		fmt.Printf("Processed %v\n", f)
	}
	return nil
}
