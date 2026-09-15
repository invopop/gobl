package za_test

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
			name:      "valid, starts with 4",
			inputCode: "4123456789",
		},
		{
			name:      "valid, another example",
			inputCode: "4480152117",
		},
		{
			name:      "empty code",
			inputCode: "",
		},
		{
			name:        "too short",
			inputCode:   "412345678",
			expectedErr: "IDENTITY-01",
		},
		{
			name:        "too long",
			inputCode:   "41234567890",
			expectedErr: "IDENTITY-01",
		},
		{
			name:        "does not start with 4",
			inputCode:   "5123456789",
			expectedErr: "IDENTITY-01",
		},
		{
			name:        "contains letters",
			inputCode:   "412345678A",
			expectedErr: "IDENTITY-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "ZA", Code: tt.inputCode}
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
			name:         "strips ZA prefix",
			inputCode:    "ZA4123456789",
			expectedCode: "4123456789",
		},
		{
			name:         "strips spaces",
			inputCode:    "4123 456 789",
			expectedCode: "4123456789",
		},
		{
			name:         "already normalized",
			inputCode:    "4123456789",
			expectedCode: "4123456789",
		},
		{
			name:         "empty",
			inputCode:    "",
			expectedCode: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "ZA", Code: tt.inputCode}
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
