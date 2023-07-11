package org_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestCoordinates(t *testing.T) {
	// Use a table to test coordinate use-cases
	tests := []struct {
		name string
		c    *org.Coordinates
		err  string
	}{
		{
			name: "empty",
			c:    &org.Coordinates{},
		},
		{
			name: "valid w3w",
			c: &org.Coordinates{
				W3W: "speak.duck.green",
			},
		},
		{
			name: "invalid w3w",
			c: &org.Coordinates{
				W3W: "speak.duck",
			},
			err: "w3w: must be in a valid format.",
		},
		{
			name: "invalid w3w",
			c: &org.Coordinates{
				W3W: "speak.duck.green.green",
			},
			err: "w3w: must be in a valid format.",
		},
		{
			name: "valid coords",
			c: &org.Coordinates{
				Latitude:  newFloat64(40.416775),
				Longitude: newFloat64(-3.703790),
			},
		},
		{
			name: "invalid latitude",
			c: &org.Coordinates{
				Latitude:  newFloat64(140.416775),
				Longitude: newFloat64(-3.703790),
			},
			err: "lat: must be no greater than 90",
		},
		{
			name: "invalid longitude",
			c: &org.Coordinates{
				Latitude:  newFloat64(40.416775),
				Longitude: newFloat64(-300.703790),
			},
			err: "lon: must be no less than -180",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.Validate()
			if tt.err != "" {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func newFloat64(f float64) *float64 {
	return &f
}
