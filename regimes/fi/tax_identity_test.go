package fi_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/fi"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
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
			expectedErr: "invalid format",
		},
		{
			name:        "too long",
			inputCode:   "123456789",
			expectedErr: "invalid format",
		},
		{
			name:        "contains letters",
			inputCode:   "1234567A",
			expectedErr: "invalid format",
		},
		{
			name:        "bad checksum",
			inputCode:   "01120388",
			expectedErr: "checksum mismatch",
		},
		{
			name:        "remainder 1 - no valid check digit",
			inputCode:   "80000000",
			expectedErr: "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "FI", Code: tt.inputCode}
			err := fi.Validate(tID)
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
		{
			name:         "nil",
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
}
