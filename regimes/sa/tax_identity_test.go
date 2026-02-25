package sa_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/sa"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "good 1", code: "300075588700003"},
		{name: "good 2", code: "310122393500003"},
		{name: "good 3", code: "399999999999993"},
		{name: "empty", code: ""},

		// Invalid formats
		{name: "too short", code: "30007558870003", err: "must be a 15-digit number starting and ending with 3"},
		{name: "too long", code: "3000755887000030", err: "must be a 15-digit number starting and ending with 3"},
		{name: "non-numeric", code: "3000755ABCD0003", err: "must be a 15-digit number starting and ending with 3"},
		{name: "not normalized", code: "3000-7558-8700-003", err: "must be a 15-digit number starting and ending with 3"},
		{name: "wrong first digit", code: "100075588700003", err: "must be a 15-digit number starting and ending with 3"},
		{name: "wrong last digit", code: "300075588700001", err: "must be a 15-digit number starting and ending with 3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "SA", Code: tt.code}
			err := sa.Validate(tID)
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
	tests := []struct {
		name     string
		code     cbc.Code
		expected cbc.Code
	}{
		{
			name:     "with SA prefix",
			code:     "SA300075588700003",
			expected: "300075588700003",
		},
		{
			name:     "with spaces",
			code:     "300 075 588 700 003",
			expected: "300075588700003",
		},
		{
			name:     "already normalized",
			code:     "300075588700003",
			expected: "300075588700003",
		},
		{
			name:     "empty",
			code:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "SA", Code: tt.code}
			sa.Normalize(tID)
			assert.Equal(t, tt.expected, tID.Code)
		})
	}
}
