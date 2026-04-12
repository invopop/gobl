package uy_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/uy"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		// Valid RUT numbers (from python-stdnum test suite)
		{name: "valid 1", code: "211003420017"},
		{name: "valid 2", code: "020164180014"},
		{name: "valid 3", code: "020334150013"},
		{name: "valid 4", code: "040003080012"},
		{name: "valid 5", code: "040005970015"},
		{name: "valid 6", code: "211004160019"},
		{name: "valid 7", code: "211049510019"},
		{name: "valid 8", code: "211073320011"},
		{name: "valid 9", code: "216893210012"},
		{name: "valid 10", code: "217055850011"},
		{name: "valid 11", code: "220018800014"},

		// Empty and nil codes should pass (not required)
		{name: "empty code", code: ""},

		// Invalid - wrong length
		{
			name: "too short",
			code: "21100342001",
			err:  "must have 12 digits",
		},
		{
			name: "too long",
			code: "2142184200106",
			err:  "must have 12 digits",
		},

		// Invalid - non-numeric characters
		{
			name: "contains letters",
			code: "FF1599340019",
			err:  "must contain only digits",
		},

		// Invalid - registration type out of range
		{
			name: "prefix 00",
			code: "001599340019",
			err:  "invalid registration type",
		},
		{
			name: "prefix too high",
			code: "991599340011",
			err:  "invalid registration type",
		},

		// Invalid - sequence number all zeros
		{
			name: "zero sequence number",
			code: "210000000019",
			err:  "invalid sequence number",
		},

		// Invalid - fixed field not 001
		{
			name: "invalid fixed field",
			code: "211599345519",
			err:  "invalid fixed field",
		},

		// Invalid - wrong check digit
		{
			name: "checksum mismatch",
			code: "211599340010",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "UY", Code: tt.code}
			err := uy.Validate(tID)
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

func TestNormalizeTaxIdentity(t *testing.T) {
	var tID *tax.Identity
	assert.NotPanics(t, func() {
		uy.Normalize(tID)
	})

	tests := []struct {
		name string
		code cbc.Code
		want cbc.Code
	}{
		{
			name: "already normalized",
			code: "211003420017",
			want: "211003420017",
		},
		{
			name: "with hyphens",
			code: "21-100342-001-7",
			want: "211003420017",
		},
		{
			name: "with spaces",
			code: "21 100342 001 7",
			want: "211003420017",
		},
		{
			name: "with UY prefix",
			code: "UY211003420017",
			want: "211003420017",
		},
		{
			name: "with UY prefix and spaces",
			code: "UY 21 140634 001 1",
			want: "211406340011",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "UY", Code: tt.code}
			uy.Normalize(tID)
			assert.Equal(t, tt.want, tID.Code)
		})
	}
}
