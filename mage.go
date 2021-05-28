// +build mage

package main

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/alecthomas/jsonschema"
	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/internal/currency"
)

// Schema generates the JSON Schema from the base models
func Schema() error {
	types := map[string]interface{}{
		"schema/envelope.json":     &gobl.Envelope{},
		"schema/bill/invoice.json": &bill.Invoice{},
	}
	ref := new(jsonschema.Reflector)
	// ref.FullyQualifyTypeNames = true
	ref.TypeNamer = typeNamer
	for f, t := range types {
		s := ref.Reflect(t)
		d, _ := json.MarshalIndent(s, "", "  ")
		ioutil.WriteFile(f, d, 0644)
	}
	return nil
}

func typeNamer(t reflect.Type) string {
	p := strings.Split(t.PkgPath(), "/")
	if len(p) > 2 {
		p = p[2:]
		p = append(p, t.Name())
		return strings.Join(p, "/")
	}
	return ""
}

// Currencies generates the Go definition files from the raw list of
// XML ISO data.
func Currencies() error {
	return currency.GenerateCodes()
}
