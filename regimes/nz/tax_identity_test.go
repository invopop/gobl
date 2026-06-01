package nz_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	_ "github.com/invopop/gobl/regimes/nz"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaxIdentityValidation(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		// Valid IRD numbers
		{name: "valid 9-digit", code: "123456785"},
		{name: "valid 8-digit", code: "12345674"},
		{name: "valid 9-digit requiring pass-2 weights", code: "100000305"},

		// Format failures
		{name: "too short", code: "1234567", err: "IDENTITY-01"},
		{name: "too long", code: "1234567890", err: "IDENTITY-01"},
		{name: "non-numeric characters", code: "1234abc89", err: "IDENTITY-01"},
		{name: "letters in code", code: "1234abc89", err: "IDENTITY-01"},

		// Range failures
		{name: "below minimum", code: "09999999", err: "IDENTITY-02"},
		{name: "above maximum", code: "150000001", err: "IDENTITY-02"},

		// Checksum failures
		{name: "wrong check digit", code: "123456789", err: "IDENTITY-03"},
		{name: "wrong check digit 8-digit", code: "12345670", err: "IDENTITY-03"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &tax.Identity{
				Country: "NZ",
				Code:    tt.code,
			}
			err := rules.Validate(id)
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

func TestTaxIdentityNormalization(t *testing.T) {
	tests := []struct {
		name string
		in   cbc.Code
		out  cbc.Code
	}{
		{name: "hyphens stripped", in: "123-456-785", out: "123456785"},
		{name: "spaces stripped", in: "123 456 785", out: "123456785"},
		{name: "NZ prefix stripped", in: "NZ123456785", out: "123456785"},
		{name: "mixed formatting", in: "NZ 123-456-785", out: "123456785"},
		{name: "already clean", in: "123456785", out: "123456785"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &tax.Identity{
				Country: "NZ",
				Code:    tt.in,
			}
			regime := tax.RegimeDefFor("NZ")
			require.NotNil(t, regime)
			regime.Normalizer(id)
			assert.Equal(t, tt.out, id.Code)
		})
	}
}
