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

var i18nStringUsed = false

// Schema generates the JSON Schema from the base models
func Schema() error {
	return schemas.Generate()
}

/*
func Schema() error {
	types := map[string]interface{}{
		"build/schema/envelope.json":     &gobl.Envelope{},
		"build/schema/bill/invoice.json": &bill.Invoice{},
		"build/schema/tax/region.json":   &tax.Region{},
		"build/schema/note/message.json": &note.Message{},
	}
	ref := new(jsonschema.Reflector)
	// ref.FullyQualifyTypeNames = true
	ref.TypeMapper = typeMapper
	ref.TypeNamer = typeNamer
	if err := ref.AddGoComments("github.com/invopop/gobl", "./"); err != nil {
		return err
	}
	var ls i18n.String
	for f, t := range types {
		i18nStringUsed = false
		s := ref.Reflect(t)
		if i18nStringUsed {
			s.Definitions["i18n.String"] = ls.JSONSchemaType()
		}
		d, _ := json.MarshalIndent(s, "", "  ")
		if err := ioutil.WriteFile(f, d, 0644); err != nil {
			return err
		}
		fmt.Printf("Processed %v\n", f)
	}
	return nil
}


func typeNamer(t reflect.Type) string {
	p := strings.Split(t.PkgPath(), "/")
	if len(p) > 2 {
		p = p[3:]
		p = append(p, t.Name())
		return strings.Join(p, ".")
	}
	return ""
}

func typeMapper(t reflect.Type) *jsonschema.Type {
	var s i18n.String
	if t == reflect.TypeOf(s) {
		i18nStringUsed = true
		return &jsonschema.Type{
			Ref: "#/definitions/i18n.String",
		}
	}
	return nil
}
*/

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
