// Package template provides a common set of tools around Go templates
// that help with converting data in other formats to GOBL
// documents.
package template

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/invopop/gobl"
	"github.com/invopop/yaml"
)

// Template contains a GOBL template document prepared for interpolating
// with incoming rows or objects of data.
type Template struct {
	tmpl *template.Template
}

// New defines a new template with the given name and data.
func New(name, data string) (*Template, error) {
	t := new(Template)
	t.tmpl = template.New(name).
		Option("missingkey=zero").
		Funcs(template.FuncMap{
			"indent":   Indent,
			"optional": Optional,
		})

	var err error
	t.tmpl, err = t.tmpl.Parse(data)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Must is a helper function that wraps a call to a function returning
// (*Template, error) and panics if the error is non-nil. It is intended
// for use in variable initializations such as
//
//	var t = template.Must(template.New("name", "..data.."))
func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}

// Execute takes the given data and interpolates it into the
// template to generate a GOBL Envelope or Schema Object according
// to the schema defined in the template.
func (t *Template) Execute(data any) (any, error) {
	buf := new(strings.Builder)
	if err := t.tmpl.Execute(buf, data); err != nil {
		return nil, err
	}

	out, err := yaml.YAMLToJSON([]byte(buf.String()))
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	res, err := gobl.Parse(out)
	if err != nil {
		return nil, fmt.Errorf("parsing GOBL: %w", err)
	}

	return res, nil
}
