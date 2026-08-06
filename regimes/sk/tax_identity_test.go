package sk_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityRules(t *testing.T) {
	tests := []struct {
		name        string
		inputCode   cbc.Code
		expectedErr string
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
			name:        "leading zero",
			inputCode:   "0000000011",
			expectedErr: "IDENTITY-01",
		},
		{
			name:        "all zeros",
			inputCode:   "0000000000",
			expectedErr: "IDENTITY-01",
		},
		{
			name:        "not divisible by 11",
			inputCode:   "2020273894",
			expectedErr: "IDENTITY-01",
		},
		{
			name:        "too short",
			inputCode:   "202027389",
			expectedErr: "IDENTITY-01",
		},
		{
			name:        "too long",
			inputCode:   "20202738931",
			expectedErr: "IDENTITY-01",
		},
		{
			name:        "contains letters",
			inputCode:   "202027389A",
			expectedErr: "IDENTITY-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "SK", Code: tt.inputCode}
			err := rules.Validate(tID)
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.expectedErr)
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
