package org

import (
	"context"
	"regexp"

	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// RegexpPatternW3W is a regular expression that matches a what3words address.
var RegexpPatternW3W = "^/*(?:(?:\\p{L}\\p{M}*)+[.｡。･・︒។։။۔።।](?:\\p{L}\\p{M}*)+[.｡。･・︒។։။۔።।](?:\\p{L}\\p{M}*)+|(?:\\p{L}\\p{M}*)+([\u0020\u00A0](?:\\p{L}\\p{M}*)+){1,3}[.｡。･・︒។։။۔።।](?:\\p{L}\\p{M}*)+([\u0020\u00A0](?:\\p{L}\\p{M}*)+){1,3}[.｡。･・︒។։။۔።।](?:\\p{L}\\p{M}*)+([\u0020\u00A0](?:\\p{L}\\p{M}*)+){1,3})$"

var regexpW3W = regexp.MustCompile(RegexpPatternW3W)

// Coordinates describes an exact geographical location in the world. We provide support
// for a set of different options beyond regular latitude and longitude.
type Coordinates struct {
	// Decimal latitude coordinate.
	Latitude *float64 `json:"lat,omitempty" jsonschema:"title=Latitude"`
	// Decimal longitude coordinate.
	Longitude *float64 `json:"lon,omitempty" jsonschema:"title=Longitude"`
	// What 3 Words text coordinates.
	W3W string `json:"w3w,omitempty" jsonschema:"title=What 3 Words"`
	// Single string coordinate based on geohash standard.
	Geohash string `json:"geohash,omitempty" jsonschema:"title=Geohash"`
}

// Validate checks that coordinates look okay.
func (c *Coordinates) Validate() error {
	return c.ValidateWithContext(context.Background())
}

// ValidateWithContext checks that coordinates look okay in the given context.
func (c *Coordinates) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithContext(ctx, c,
		validation.Field(&c.Latitude, validation.Min(-90.0), validation.Max(90.0)),
		validation.Field(&c.Longitude, validation.Min(-180.0), validation.Max(180.0)),
		validation.Field(&c.W3W, validation.Match(regexpW3W)),
	)
}

// LatLon provides the Latitude and Longitude values as a pair,
// or 0, 0 if the coordinates are not set.
func (c *Coordinates) LatLon() (float64, float64) {
	if c == nil || c.Latitude == nil || c.Longitude == nil {
		return 0, 0
	}
	return *c.Latitude, *c.Longitude
}
