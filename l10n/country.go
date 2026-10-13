package l10n

import (
	"fmt"
	"slices"
	"strings"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/jsonschema"
)

// ISOCountryCode defines an ISO 3166-1 alpha-2 country code.
type ISOCountryCode Code

// TaxCountryCode defines a code that may coincide with a country code,
// but may not always coincide.
type TaxCountryCode Code

func validISOCountryValues() []any {
	defs := countries.ISO()
	list := make([]any, len(defs))
	for i, d := range defs {
		list[i] = string(d.Code)
	}
	return list
}

func validTaxCountryValues() []any {
	defs := countries.Tax()
	list := make([]any, len(defs))
	for i, d := range defs {
		list[i] = string(d.Code)
	}
	return list
}

// Validate ensures the ISO country code is inside the known and valid
// list of countries.
func (c ISOCountryCode) Validate() error {
	if c.Empty() {
		return nil
	}
	if slices.Contains(validISOCountryValues(), any(string(c))) {
		return nil
	}
	return fmt.Errorf("must be a valid ISO country code")
}

// Validate ensures the tax country code is inside the known and valid
// list of country codes for taxes.
func (c TaxCountryCode) Validate() error {
	if c.Empty() {
		return nil
	}
	if slices.Contains(validTaxCountryValues(), any(string(c))) {
		return nil
	}
	return fmt.Errorf("must be a valid tax country code")
}

func normalizeISOCountryCode(c *ISOCountryCode) {
	*c = ISOCountryCode(strings.ToUpper(strings.TrimSpace(string(*c))))
}

func normalizeTaxCountryCode(c *TaxCountryCode) {
	*c = TaxCountryCode(strings.ToUpper(strings.TrimSpace(string(*c))))
}

func isoCountryCodeRules() *rules.Set {
	return rules.For(ISOCountryCode(""),
		rules.AssertIfPresent("01", "must be a valid ISO country code",
			is.In(validISOCountryValues()...),
		),
	)
}

func taxCountryCodeRules() *rules.Set {
	return rules.For(TaxCountryCode(""),
		rules.AssertIfPresent("01", "must be a valid tax country code",
			is.In(validTaxCountryValues()...),
		),
	)
}

// Empty returns true if the ISO country code is empty.
func (c ISOCountryCode) Empty() bool {
	return c == ""
}

// Empty returns true if the tax country code is empty.
func (c TaxCountryCode) Empty() bool {
	return c == ""
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (ISOCountryCode) JSONSchema() *jsonschema.Schema {
	defs := countries.ISO()
	s := &jsonschema.Schema{
		Title:       "ISO Country Code",
		Type:        "string",
		OneOf:       make([]*jsonschema.Schema, len(defs)),
		Description: `Defines an ISO 3166-1 alpha-2 country code.`,
	}
	for i, d := range defs {
		s.OneOf[i] = &jsonschema.Schema{
			Const: d.Code,
			Title: d.Name,
		}
	}
	return s
}

// In returns true if the country code is contained inside the provided set
func (c ISOCountryCode) In(set ...ISOCountryCode) bool {
	for _, x := range set {
		if c == x {
			return true
		}
	}
	return false
}

// String provides string representation of the ISO country code
func (c ISOCountryCode) String() string {
	return string(c)
}

// Name provides the ISO Country Name for the code
func (c ISOCountryCode) Name() string {
	if d := countries.Code(c.Code()); d != nil {
		return d.Name
	}
	return ""
}

// Alpha3 provides the ISO 3166-1 alpha-3 country code
func (c ISOCountryCode) Alpha3() string {
	if d := countries.Code(c.Code()); d != nil {
		return d.Alpha3
	}
	return ""
}

// Code provides the Code type for the ISO country code.
func (c ISOCountryCode) Code() Code {
	return Code(c)
}

// In returns true if the tax country code is contained inside the provided set
func (c TaxCountryCode) In(set ...TaxCountryCode) bool {
	for _, x := range set {
		if c == x {
			return true
		}
	}
	return false
}

// String provides string representation of the tax country code
func (c TaxCountryCode) String() string {
	return string(c)
}

// Name provides the ISO Country Name for the code
func (c TaxCountryCode) Name() string {
	if d := countries.Code(c.Code()); d != nil {
		return d.Name
	}
	return ""
}

// Code provides the Code type for the tax country code.
func (c TaxCountryCode) Code() Code {
	return Code(c)
}

// JSONSchema provides a representation of the type for usage in Schema.
func (TaxCountryCode) JSONSchema() *jsonschema.Schema {
	defs := countries.Tax()
	s := &jsonschema.Schema{
		Title:       "Tax Country Code",
		Type:        "string",
		OneOf:       make([]*jsonschema.Schema, len(defs)),
		Description: `Defines an ISO base country code used for tax purposes.`,
	}
	for i, d := range defs {
		s.OneOf[i] = &jsonschema.Schema{
			Const: d.Code,
			Title: d.Name,
		}
	}
	return s
}
