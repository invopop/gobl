package fr_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	r := fr.New()
	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{
			Code:     "356000000", // SIREN to VAT
			Expected: "39356000000",
		},
		{
			Code:     "44 73282 9320 ",
			Expected: "44732829320",
		},
		{
			Code:     "391-838-042",
			Expected: "44391838042",
		},
		{
			Code:     "FR391838042",
			Expected: "44391838042",
		},
		{
			Code:     "FR44391838042",
			Expected: "44391838042",
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: l10n.FR, Code: ts.Code}
		err := r.CalculateObject(tID)
		assert.NoError(t, err)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "good 1", code: "39356000000"},
		{name: "good 2", code: "44732829320"},
		{name: "good 3", code: "44391838042"},
		{
			name: "empty",
			code: "",
		},
		{
			name: "too long",
			code: "44123456789100",
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
			code: "44999999991",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: l10n.FR, Code: tt.code}
			err := fr.Validate(tID)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.err)
				}
			}
		})
	}
}
