package org

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/uuid"
)

// Registration is used in countries that require additional information to be associated
// with a company usually related to a specific registration office.
// The definition found here is based on the details required for spain.
// If your country requires additional fields, please let us know.
type Registration struct {
	uuid.Identify
	Label    string        `json:"label,omitempty" jsonschema:"title=Label,example=Registration"`
	Capital  *num.Amount   `json:"capital,omitempty" jsonschema:"title=Capital"`
	Currency currency.Code `json:"currency,omitempty" jsonschema:"title=Currency"`
	Office   string        `json:"office,omitempty" jsonschema:"title=Office"`
	Book     string        `json:"book,omitempty" jsonschema:"title=Book"`
	Volume   string        `json:"volume,omitempty" jsonschema:"title=Volume"`
	Sheet    string        `json:"sheet,omitempty" jsonschema:"title=Sheet"`
	Section  string        `json:"section,omitempty" jsonschema:"title=Section"`
	Page     string        `json:"page,omitempty" jsonschema:"title=Page"`
	Entry    string        `json:"entry,omitempty" jsonschema:"title=Entry"`
	Other    string        `json:"other,omitempty" jsonschema:"title=Other"`
}

// Normalize ensures the registration is in a canonical format.
func (r *Registration) Normalize() {
	if r == nil {
		return
	}
	uuid.Normalize(&r.UUID)
	r.Label = cbc.NormalizeString(r.Label)
	r.Office = cbc.NormalizeString(r.Office)
	r.Book = cbc.NormalizeString(r.Book)
	r.Volume = cbc.NormalizeString(r.Volume)
	r.Sheet = cbc.NormalizeString(r.Sheet)
	r.Section = cbc.NormalizeString(r.Section)
	r.Page = cbc.NormalizeString(r.Page)
	r.Entry = cbc.NormalizeString(r.Entry)
	r.Other = cbc.NormalizeString(r.Other)
}

func registrationRules() *rules.Set {
	return rules.For(new(Registration),
		rules.Field("currency",
			rules.AssertIfPresent("01", "registration currency must be a valid ISO 4217 code",
				is.FuncError("is valid currency", func(val any) error {
					c, _ := val.(currency.Code)
					return c.Validate()
				}),
			),
		),
	)
}
