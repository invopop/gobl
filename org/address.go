package org

import (
	"context"
	"strings"

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
	// Name of a village, town, district, or city, typically inside a region.
	Locality string `json:"locality,omitempty" jsonschema:"title=Locality"`
	// Name of a city, province, county, or state, inside a country.
	Region string `json:"region,omitempty" jsonschema:"title=Region"`
	// State or province code for countries that require it.
	State cbc.Code `json:"state,omitempty" jsonschema:"title=State"`
	// Post or ZIP code.
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`
	// ISO country code.
	Country l10n.ISOCountryCode `json:"country,omitempty" jsonschema:"title=Country"`
	// When the postal address is not sufficient, coordinates help locate the address more precisely.
	Coordinates *Coordinates `json:"coords,omitempty" jsonschema:"title=Coordinates"`
	// Any additional semi-structure details about the address.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Normalize will perform basic normalization of the address's data.
func (a *Address) Normalize(normalizers tax.Normalizers) {
	if a == nil {
		return
	}
	uuid.Normalize(&a.UUID)
	a.PostOfficeBox = strings.TrimSpace(a.PostOfficeBox)
	a.Number = strings.TrimSpace(a.Number)
	a.Floor = strings.TrimSpace(a.Floor)
	a.Block = strings.TrimSpace(a.Block)
	a.Door = strings.TrimSpace(a.Door)
	a.Street = strings.TrimSpace(a.Street)
	a.StreetExtra = strings.TrimSpace(a.StreetExtra)
	a.Locality = strings.TrimSpace(a.Locality)
	a.Region = strings.TrimSpace(a.Region)
	a.State = cbc.NormalizeAlphanumericalCode(a.State)
	a.Code = cbc.NormalizeCode(a.Code)
	normalizers.Each(a)
}

// Validate checks that an address looks okay.
func (a *Address) Validate() error {
	return a.ValidateWithContext(context.Background())
}

// ValidateWithContext checks that an address looks okay in the given context.
func (a *Address) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithContext(ctx, a,
		validation.Field(&a.UUID),
		validation.Field(&a.State),
		validation.Field(&a.Code),
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
