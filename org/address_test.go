package org_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestAddressNormalize(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *org.Address
		assert.NotPanics(t, func() {
			a.Normalize(nil)
		})
	})

	t.Run("normalize fields", func(t *testing.T) {
		a := &org.Address{
			PostOfficeBox: " 20 ",
			Number:        " 12 ",
			Floor:         "  3  ",
			Block:         " A  ",
			Door:          " 1 ",
			Street:        " Main St. ",
			StreetExtra:   "  Apt. 3  ",
			Locality:      "  Town ",
			Region:        "  City ",
			State:         "  MAD  ",
			Code:          " HG12 2AB ",
		}
		a.Normalize(nil)

		assert.Equal(t, "20", a.PostOfficeBox)
		assert.Equal(t, "12", a.Number)
		assert.Equal(t, "3", a.Floor)
		assert.Equal(t, "A", a.Block)
		assert.Equal(t, "1", a.Door)
		assert.Equal(t, "Main St.", a.Street)
		assert.Equal(t, "Apt. 3", a.StreetExtra)
		assert.Equal(t, "Town", a.Locality)
		assert.Equal(t, "City", a.Region)
		assert.Equal(t, "MAD", a.State.String())
		assert.Equal(t, "HG12 2AB", a.Code.String())
	})
}

func TestAddressValidation(t *testing.T) {
	t.Run("valid address", func(t *testing.T) {
		a := &org.Address{
			Number:   "12",
			Street:   "Main St.",
			Locality: "Town",
			Region:   "City",
			State:    "MAD",
			Code:     "HG12 2AB",
			Country:  "GB",
		}
		assert.NoError(t, a.Validate())
	})

	t.Run("invalid UUID", func(t *testing.T) {
		a := &org.Address{
			Number:   "12",
			Street:   "Main St.",
			Locality: "Town",
			Region:   "City",
			State:    "MAD",
			Code:     "HG12 2AB",
			Country:  "GB",
		}
		a.UUID = "invalid"
		assert.ErrorContains(t, a.Validate(), "uuid: invalid UUID length: 7")
	})
}

func TestAddressLineOne(t *testing.T) {
	t.Run("number first", func(t *testing.T) {
		a := &org.Address{
			Number:  "12",
			Street:  "Main St.",
			Country: "GB",
		}
		assert.Equal(t, "12 Main St.", a.LineOne())
	})

	t.Run("street first", func(t *testing.T) {
		a := &org.Address{
			Number:  "12",
			Street:  "Gran Vía",
			Country: "ES",
		}
		assert.Equal(t, "Gran Vía 12", a.LineOne())
	})

	t.Run("combine number details", func(t *testing.T) {
		a := &org.Address{
			Number:  "12",
			Block:   "esc. 1",
			Floor:   "10",
			Door:    "C",
			Street:  "Gran Vía",
			Country: "ES",
		}
		assert.Equal(t, "Gran Vía 12 esc. 1 10 C", a.LineOne())
	})
}

func TestAddressLineTwo(t *testing.T) {
	t.Run("number first", func(t *testing.T) {
		a := &org.Address{
			Number:      "12",
			Street:      "Main St.",
			StreetExtra: "Apt. 3",
			Country:     "GB",
		}
		assert.Equal(t, "12 Main St.", a.LineOne())
		assert.Equal(t, "Apt. 3", a.LineTwo())
	})
}
