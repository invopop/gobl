package sk_test

import (
	"slices"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityRules(t *testing.T) {
	// The format and checksum rules are asserted separately, so a code may
	// trigger either or both. Codes not listed must NOT be reported.
	allCodes := []string{"IDENTITY-01", "IDENTITY-02"}

	tests := []struct {
		name         string
		inputCode    cbc.Code
		expectedErrs []string
	}{
		{
			name:      "valid divisible by 11",
			inputCode: "2020273893",
		},
		{
			name:      "valid divisible by 11, natural person",
			inputCode: "7120001713",
		},
		{
			name:      "empty code",
			inputCode: "",
		},
		{
			name:         "leading zero, checksum still valid",
			inputCode:    "0000000011",
			expectedErrs: []string{"IDENTITY-01"},
		},
		{
			name:         "all zeros, checksum still valid",
			inputCode:    "0000000000",
			expectedErrs: []string{"IDENTITY-01"},
		},
		{
			name:         "well formed but not divisible by 11",
			inputCode:    "2020273894",
			expectedErrs: []string{"IDENTITY-02"},
		},
		{
			name:         "too short",
			inputCode:    "202027389",
			expectedErrs: []string{"IDENTITY-01", "IDENTITY-02"},
		},
		{
			name:         "too long",
			inputCode:    "20202738931",
			expectedErrs: []string{"IDENTITY-01", "IDENTITY-02"},
		},
		{
			name:         "contains letters",
			inputCode:    "202027389A",
			expectedErrs: []string{"IDENTITY-01", "IDENTITY-02"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "SK", Code: tt.inputCode}
			err := rules.Validate(tID)
			if len(tt.expectedErrs) == 0 {
				assert.NoError(t, err)
				return
			}
			if !assert.Error(t, err) {
				return
			}
			for _, code := range allCodes {
				if slices.Contains(tt.expectedErrs, code) {
					assert.Contains(t, err.Error(), code)
				} else {
					assert.NotContains(t, err.Error(), code)
				}
			}
		})
	}
}

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		name         string
		inputCode    cbc.Code
		expectedCode cbc.Code
	}{
		{
			name:         "strips SK prefix",
			inputCode:    "SK2020273893",
			expectedCode: "2020273893",
		},
		{
			name:         "strips lowercase prefix",
			inputCode:    "sk2020273893",
			expectedCode: "2020273893",
		},
		{
			name:         "strips spaces",
			inputCode:    "2020 2738 93",
			expectedCode: "2020273893",
		},
		{
			name:         "strips hyphen",
			inputCode:    "20202-73893",
			expectedCode: "2020273893",
		},
		{
			name:         "already normalized",
			inputCode:    "2020273893",
			expectedCode: "2020273893",
		},
		{
			name:         "empty",
			inputCode:    "",
			expectedCode: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "SK", Code: tt.inputCode}
			norm.Normalize(tID)
			assert.Equal(t, tt.expectedCode, tID.Code)
		})
	}

	t.Run("nil identity", func(t *testing.T) {
		assert.NotPanics(t, func() {
			norm.Normalize((*tax.Identity)(nil))
		})
	})
}
