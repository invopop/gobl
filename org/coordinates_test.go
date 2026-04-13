package org_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
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
			err: "[GOBL-ORG-COORDINATES-03] ($.w3w) what3words coordinate must be valid",
		},
		{
			name: "invalid w3w",
			c: &org.Coordinates{
				W3W: "speak.duck.green.green",
			},
			err: "[GOBL-ORG-COORDINATES-03]",
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
			err: "[GOBL-ORG-COORDINATES-01]",
		},
		{
			name: "invalid longitude",
			c: &org.Coordinates{
				Latitude:  newFloat64(40.416775),
				Longitude: newFloat64(-300.703790),
			},
			err: "[GOBL-ORG-COORDINATES-02]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rules.Validate(tt.c)
			if tt.err != "" {
				assert.ErrorContains(t, err, tt.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func newFloat64(f float64) *float64 {
	return &f
}
