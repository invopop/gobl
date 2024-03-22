//go:build ignore

package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/pkg/here"
)

func main() {
	if err := GenerateCodes(); err != nil {
		fmt.Println(err)
	}
}

var codeTemplate = here.Doc(`
	// Code generated by "go run ./generate.go"; DO NOT EDIT.

	package currency

	// List of ISO4217 and common currency codes.
	const (
		{{- range .Defs }}
		// {{ .Name }} ({{ .Symbol }})
		{{ .ISOCode }} Code = "{{ .ISOCode }}"
		{{- end }}
	)
`)

// GenerateCodes is a special tool function used to convert the source XML
// data into an array of currency definitions.
func GenerateCodes() error {
	tmpl, err := template.New("codes").Parse(codeTemplate)
	if err != nil {
		return err
	}

	fields := map[string]any{
		"Defs": currency.Definitions(),
	}

	f, err := os.Create("./codes.go")
	if err != nil {
		return err
	}
	defer f.Close() // nolint:errcheck
	if err := tmpl.Execute(f, fields); err != nil {
		return err
	}

	return nil
}
