package org

import (
	"context"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Address defines a globally acceptable set of attributes that describes
// a postal or fiscal address.
// Attribute names loosely based on the xCard file format.
type Address struct {
	uuid.Identify
	// Useful identifier, such as home, work, etc.
	Label string `json:"label,omitempty" jsonschema:"title=Label,example=Office"`
	// Box number or code for the post office box located at the address.
	PostOfficeBox string `json:"po_box,omitempty" jsonschema:"title=Post Office Box"`
	// House or building number in the street.
	Number string `json:"num,omitempty" jsonschema:"title=Number"`
	// Floor number within the building.
	Floor string `json:"floor,omitempty" jsonschema:"title=Floor"`
	// Block number within the building.
	Block string `json:"block,omitempty" jsonschema:"title=Block"`
	// Door number within the building.
	Door string `json:"door,omitempty" jsonschema:"title=Door"`
	// First line of street.
	Street string `json:"street,omitempty" jsonschema:"title=Street"`
	// Additional street address details.
	StreetExtra string `json:"street_extra,omitempty" jsonschema:"title=Extended Street"`
	// Village, town, district, or city, typically inside a region.
	Locality string `json:"locality,omitempty" jsonschema:"title=Locality"`
	// Province, county, or state, inside a country.
	Region string `json:"region,omitempty" jsonschema:"title=Region"`
	// Post or ZIP code.
	Code string `json:"code,omitempty" jsonschema:"title=Code"`
	// ISO country code.
	Country l10n.CountryCode `json:"country,omitempty" jsonschema:"title=Country"`
	// When the postal address is not sufficient, coordinates help locate the address more precisely.
	Coordinates *Coordinates `json:"coords,omitempty" jsonschema:"title=Coordinates"`
	// Any additional semi-structure details about the address.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate checks that an address looks okay.
func (a *Address) Validate() error {
	return a.ValidateWithContext(context.Background())
}

// ValidateWithContext checks that an address looks okay in the given context.
func (a *Address) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithRegime(ctx, a,
		validation.Field(&a.UUID),
		validation.Field(&a.Country),
		validation.Field(&a.Coordinates),
		validation.Field(&a.Meta),
	)
}

// JSONSchemaExtend adds extra details to the Address schema.
func (Address) JSONSchemaExtend(js *jsonschema.Schema) {
	js.Extras = map[string]any{
		schema.Recommended: []string{
			"number", "street", "locality", "region", "code", "country",
		},
	}
}
