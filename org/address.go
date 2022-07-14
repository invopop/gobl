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
	UUID          *uuid.UUID   `json:"uuid,omitempty" jsonschema:"title=UUID"`
	Label         string       `json:"label,omitempty" jsonschema:"title=Label,description=Useful identifier, such as home, work, etc."`
	PostOfficeBox string       `json:"po_box,omitempty" jsonschema:"title=Post Office Box,description=Box number or code for the post office box located at the address."`
	Number        string       `json:"num,omitempty" jsonschema:"title=Number,description=House or building number in the street."`
	Floor         string       `json:"floor,omitempty" jsonschema:"title=Floor,description=Floor number within the building."`
	Block         string       `json:"block,omitempty" jsonschema:"title=Block,description=Block number within the building."`
	Door          string       `json:"door,omitempty" jsonschema:"title=Door,description=Door number within the building."`
	Street        string       `json:"street,omitempty" jsonschema:"title=Street,description=Fist line of street."`
	StreetExtra   string       `json:"street_extra,omitempty" jsonschema:"title=Extended Street,description=Additional street address details."`
	Locality      string       `json:"locality" jsonschema:"title=Locality,description=The village, town, district, or city."`
	Region        string       `json:"region" jsonschema:"title=Region,description=Province, County, or State."`
	Code          string       `json:"code,omitempty" jsonschema:"title=Code,description=Post or ZIP code."`
	Country       l10n.Code    `json:"country,omitempty" jsonschema:"title=Country,description=ISO country code."`
	Coordinates   *Coordinates `json:"coords,omitempty" jsonschema:"title=Coordinates,description=For when the postal address is not sufficient, coordinates help locate the address more precisely."`
	Meta          Meta         `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Coordinates describes an exact geographical location in the world. We provide support
// for a set of different options beyond regular latitude and longitude.
type Coordinates struct {
	Latitude  float64 `json:"lat,omitempty" jsonschema:"title=Latitude,description=Decimal latitude coordinate."`
	Longitude float64 `json:"lon,omitempty" jsonschema:"title=Longitude,description=Decimal longitude coordinate."`
	W3W       string  `json:"w3w,omitempty" jsonschema:"title=What 3 Words,description=Text coordinates compose of three words."`
	Geohash   string  `json:"geohash,omitempty" jsonschema:"title=Geohash,description=Single string coordinate based on geohash standard."`
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
