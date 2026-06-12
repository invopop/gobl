package org

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// Attribute describes a named feature or property of the parent object,
// such as the colour or size of an item.
type Attribute struct {
	// Name of the attribute or property.
	Name string `json:"name" jsonschema:"title=Name"`
	// Value of the attribute or property.
	Value string `json:"value" jsonschema:"title=Value"`
}

func attributeRules() *rules.Set {
	return rules.For(new(Attribute),
		rules.Field("name",
			rules.Assert("01", "attribute name is required", is.Present),
		),
		rules.Field("value",
			rules.Assert("02", "attribute value is required", is.Present),
		),
	)
}

func normalizeAttribute(a *Attribute) {
	a.Name = cbc.NormalizeString(a.Name)
	a.Value = cbc.NormalizeString(a.Value)
}

// CleanAttributes removes any nil or empty attributes from the list,
// returning nil if none remain.
func CleanAttributes(attrs []*Attribute) []*Attribute {
	cleaned := make([]*Attribute, 0, len(attrs))
	for _, a := range attrs {
		if a == nil || (a.Name == "" && a.Value == "") {
			continue
		}
		cleaned = append(cleaned, a)
	}
	if len(cleaned) == 0 {
		return nil
	}
	return cleaned
}
