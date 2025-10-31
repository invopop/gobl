package ar_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/ar"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityValidation(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		// Valid CUIT/CUIL examples
		// These are valid CUIT/CUIL numbers following the modulo 11 algorithm
		{
			name: "valid CUIT - standard company",
			code: "30500010912", // Valid test CUIT
		},
		{
			name: "valid CUIL - male individual",
			code: "20172543597",
		},
		{
			name: "valid CUIL - female individual",
			code: "27123456780",
		},
		{
			name: "valid CUIT - prefix 23 (conflict resolution)",
			code: "23000000019",
		},
		{
			name: "valid CUIT - prefix 33 (company conflict resolution)",
			code: "33000000049",
		},

		// Invalid - wrong length
		{
			name: "too short",
			code: "2017254359",
			err:  "must have 11 digits",
		},
		{
			name: "too long",
			code: "201725435978",
			err:  "must have 11 digits",
		},

		// Invalid - non-numeric
		{
			name: "contains letters",
			code: "2017254A597",
			err:  "must contain only digits",
		},
		{
			name: "contains special characters",
			code: "20172543A97",
			err:  "must contain only digits",
		},

		// Invalid - wrong check digit
		{
			name: "wrong check digit",
			code: "20172543598",
			err:  "verification digit mismatch",
		},
		{
			name: "wrong check digit - company",
			code: "30500010911", // Changed last digit
			err:  "verification digit mismatch",
		},

		// Invalid - another wrong check digit test
		{
			name: "wrong check digit - female",
			code: "27123456781", // Wrong check digit
			err:  "verification digit mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "AR", Code: tt.code}
			err := ar.Validate(tID)
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
		code cbc.Code
		want cbc.Code
	}{
		{
			name: "already normalized",
			code: "30714589840",
			want: "30714589840",
		},
		{
			name: "with hyphens",
			code: "30-71458984-0",
			want: "30714589840",
		},
		{
			name: "with spaces",
			code: "30 71458984 0",
			want: "30714589840",
		},
		{
			name: "with hyphens and spaces",
			code: "30-71458984-0 ",
			want: "30714589840",
		},
		{
			name: "CUIL with hyphens",
			code: "20-17254359-7",
			want: "20172543597",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "AR", Code: tt.code}
			ar.Normalize(tID)
			assert.Equal(t, tt.want, tID.Code)
		})
	}
}

func TestEmptyTaxIdentity(t *testing.T) {
	tID := &tax.Identity{Country: "AR"}
	err := ar.Validate(tID)
	assert.NoError(t, err)
}
