package fr_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestIdentityNormalization(t *testing.T) {
	r := tax.RegimeDefFor("FR")

	t.Run("normalize SIREN", func(t *testing.T) {
		p1 := &org.Identity{
			Key:  fr.IdentityKeySiren,
			Code: "FR3412000. 68",
		}
		r.NormalizeObject(p1)
		assert.Equal(t, "341200068", p1.Code.String())
	})

	t.Run("normalize SIRET", func(t *testing.T) {
		p1 := &org.Identity{
			Key:  fr.IdentityKeySiret,
			Code: "FR 341200068-00001",
		}
		r.NormalizeObject(p1)
		assert.Equal(t, "34120006800001", p1.Code.String())
	})
}

func TestIdentityValidation(t *testing.T) {

	t.Run("validate SIREN", func(t *testing.T) {
		tests := []struct {
			name string
			code cbc.Code
			err  string
		}{
			{name: "good 1", code: "356000000"},
			{name: "good 2", code: "732829320"},
			{name: "good 3", code: "391838042"},
			{
				name: "empty",
				code: "",
				err:  "cannot be blank",
			},
			{
				name: "too long",
				code: "356000000000",
				err:  "invalid format",
			},
			{
				name: "too short",
				code: "123456",
				err:  "invalid format",
			},
			{
				name: "not normalized",
				code: "12.449.965-4",
				err:  "invalid format",
			},
			{
				name: "bad checksum",
				code: "999999991",
				err:  "checksum mismatch",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				id := &org.Identity{Key: fr.IdentityKeySiren, Code: tt.code}
				err := fr.Validate(id)
				if tt.err == "" {
					assert.NoError(t, err)
				} else {
					if assert.Error(t, err) {
						assert.Contains(t, err.Error(), tt.err)
					}
				}
			})
		}
	})

	t.Run("validate SIRET", func(t *testing.T) {
		tests := []struct {
			name string
			code cbc.Code
			err  string
		}{
			{name: "good 1", code: "35600000012345"},
			{name: "good 2", code: "73282932000000"},
			{name: "good 3", code: "39183804200000"},
			{
				name: "empty",
				code: "",
				err:  "cannot be blank",
			},
			{
				name: "too long",
				code: "57201358347049000000",
				err:  "invalid format",
			},
			{
				name: "too short",
				code: "123456",
				err:  "invalid format",
			},
			{
				name: "not normalized",
				code: "12.449.965-409872",
				err:  "invalid format",
			},
			{
				name: "bad checksum",
				code: "58252940504934",
				err:  "checksum mismatch",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				id := &org.Identity{Key: fr.IdentityKeySiret, Code: tt.code}
				err := fr.Validate(id)
				if tt.err == "" {
					assert.NoError(t, err)
				} else {
					if assert.Error(t, err) {
						assert.Contains(t, err.Error(), tt.err)
					}
				}
			})
		}
	})

}
