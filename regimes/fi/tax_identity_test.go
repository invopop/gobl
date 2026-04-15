package fi_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/fi"
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
			name:      "valid check digit zero",
			inputCode: "23456780",
		},
		{
			name:      "valid check digit non-zero",
			inputCode: "01120389",
		},
		{
			name:      "valid small business ID",
			inputCode: "07375462",
		},
		{
			name:      "valid large business ID",
			inputCode: "12345671",
		},
		{
			name:      "empty code",
			inputCode: "",
		},
		{
			name:        "too short",
			inputCode:   "1234567",
			expectedErr: "IDENTITY-01",
		},
		{
			name:        "too long",
			inputCode:   "123456789",
			expectedErr: "IDENTITY-01",
		},
		{
			name:        "contains letters",
			inputCode:   "1234567A",
			expectedErr: "IDENTITY-01",
		},
		{
			name:        "bad checksum",
			inputCode:   "01120388",
			expectedErr: "IDENTITY-01",
		},
		{
			name:        "remainder 1 - no valid check digit",
			inputCode:   "80000000",
			expectedErr: "IDENTITY-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "FI", Code: tt.inputCode}
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
			name:         "strips FI prefix",
			inputCode:    "FI01120389",
			expectedCode: "01120389",
		},
		{
			name:         "strips spaces",
			inputCode:    "0112 0389",
			expectedCode: "01120389",
		},
		{
			name:         "strips hyphen",
			inputCode:    "0112038-9",
			expectedCode: "01120389",
		},
		{
			name:         "already normalized",
			inputCode:    "01120389",
			expectedCode: "01120389",
		},
		{
			name:         "empty",
			inputCode:    "",
			expectedCode: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "FI", Code: tt.inputCode}
			fi.Normalize(tID)
			assert.Equal(t, tt.expectedCode, tID.Code)
		})
	}

	t.Run("nil identity", func(t *testing.T) {
		assert.NotPanics(t, func() {
			fi.Normalize((*tax.Identity)(nil))
		})
	})
}
