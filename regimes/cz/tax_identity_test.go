package cz_test

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
		// Legal persons: an 8-digit IČO with a check digit.
		{
			name:      "valid IČO, EET sample identity",
			inputCode: "00000019",
		},
		{
			name:      "valid IČO, Česká spořitelna",
			inputCode: "45244782",
		},
		{
			name:      "valid IČO, ČEZ",
			inputCode: "45274649",
		},
		{
			name:      "valid IČO, remainder 0 gives check digit 1",
			inputCode: "00000001",
		},
		{
			name:      "valid IČO, remainder 1 gives check digit 0",
			inputCode: "00000060",
		},
		{
			name:        "bad IČO check digit, rejected by the EET Playground",
			inputCode:   "00000018",
			expectedErr: "IDENTITY-02",
		},
		// Natural persons: a 10-digit birth number with a check digit.
		{
			name:      "valid birth number, EET sample identity",
			inputCode: "8551015704",
		},
		{
			name:      "birth number with the 1954-1985 check digit 0",
			inputCode: "6001010060",
		},
		{
			name:        "birth number with check digit 0 outside the 1954-1985 window",
			inputCode:   "8601010100",
			expectedErr: "IDENTITY-02",
		},
		{
			name:        "birth number not divisible by 11, rejected by the EET Playground",
			inputCode:   "1234567890",
			expectedErr: "IDENTITY-02",
		},
		// Natural persons: 9 digits carry no check digit.
		{
			name:      "tax office identifier, EET sample identity",
			inputCode: "683555118",
		},
		{
			name:      "birth number issued before 1954",
			inputCode: "530101123",
		},
		// Shape.
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
			inputCode:   "12345678901",
			expectedErr: "IDENTITY-01",
		},
		{
			name:        "contains letters",
			inputCode:   "4524478A",
			expectedErr: "IDENTITY-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "CZ", Code: tt.inputCode}
			err := rules.Validate(tID)
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.expectedErr)
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
			name:         "strips CZ prefix",
			inputCode:    "CZ45244782",
			expectedCode: "45244782",
		},
		{
			name:         "strips lowercase prefix",
			inputCode:    "cz45244782",
			expectedCode: "45244782",
		},
		{
			name:         "strips spaces",
			inputCode:    "452 447 82",
			expectedCode: "45244782",
		},
		{
			name:         "already normalized",
			inputCode:    "45244782",
			expectedCode: "45244782",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "CZ", Code: tt.inputCode}
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
