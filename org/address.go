package org

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/uuid"
)

// Address defines a globally acceptable set of attributes that describes
// a postal or fiscal address.
// Attribute names loosly based on the xCard file format.
type Address struct {
	// Internal ID used to identify the party inside a document.
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Useful identifier, such as home, work, etc.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
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
	// The village, town, district, or city.
	Locality string `json:"locality" jsonschema:"title=Locality"`
	// Province, County, or State.
	Region string `json:"region" jsonschema:"title=Region"`
	// Post or ZIP code.
	Code string `json:"code,omitempty" jsonschema:"title=Code"`
	// ISO country code.
	Country l10n.Code `json:"country,omitempty" jsonschema:"title=Country"`
	// When the postal address is not sufficient, coordinates help locate the address more precisely.
	Coordinates *Coordinates `json:"coords,omitempty" jsonschema:"title=Coordinates"`
	// Any additional semi-structure details about the address.
	Meta Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Coordinates describes an exact geographical location in the world. We provide support
// for a set of different options beyond regular latitude and longitude.
type Coordinates struct {
	// Decimal latitude coordinate.
	Latitude float64 `json:"lat,omitempty" jsonschema:"title=Latitude"`
	// Decimal longitude coordinate.
	Longitude float64 `json:"lon,omitempty" jsonschema:"title=Longitude"`
	// Text coordinates compose of three words.
	W3W string `json:"w3w,omitempty" jsonschema:"title=What 3 Words"`
	// Single string coordinate based on geohash standard.
	Geohash string `json:"geohash,omitempty" jsonschema:"title=Geohash"`
}

// Validate checks that an address looks okay.
func (a *Address) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.UUID),
		validation.Field(&a.Country, l10n.IsCountry),
		validation.Field(&a.Coordinates),
		validation.Field(&a.Meta),
	)
}
