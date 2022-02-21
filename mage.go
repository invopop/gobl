//go:build mage
// +build mage

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/internal/currency"
	"github.com/invopop/gobl/internal/schemas"
	"github.com/invopop/gobl/region"
)

// Schema generates the JSON Schema from the base models
func Schema() error {
	return schemas.Generate()
}

func RegionData() error {
	for c, r := range region.All() {
		doc := new(gobl.Document)
		if err := doc.Insert(r.Taxes()); err != nil {
			return err
		}
		data, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			return err
		}
		f := filepath.Join("build", "data", "tax", string(c)+".json")
		if err := ioutil.WriteFile(f, data, 0644); err != nil {
			return err
		}
		fmt.Printf("Processed %v\n", f)
	}
	return nil
}

// Currencies generates the Go definition files from the raw list of
// XML ISO data.
func Currencies() error {
	return currency.GenerateCodes()
}
