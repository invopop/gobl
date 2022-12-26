//go:build mage
// +build mage

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/internal/currency"
	"github.com/invopop/gobl/internal/schemas"
	"github.com/invopop/gobl/tax"
)

// Schema generates the JSON Schema from the base models
func Schema() error {
	return schemas.Generate()
}

// Regimes generates JSON version of each regimes's data.
func Regimes() error {
	for _, r := range tax.AllRegimes() {
		doc, err := gobl.NewDocument(r)
		if err != nil {
			return err
		}
		data, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			return err
		}
		n := string(r.Country)
		if r.Zone != "" {
			n = n + "_" + string(r.Zone)
		}
		n = strings.ToLower(n)
		f := filepath.Join("build", "regimes", n+".json")
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

// Samples runs through all the `.yaml` samples and generates complete GOBL
// Envelopes of each in `.json` files.
func Samples() error {

	return nil
}
