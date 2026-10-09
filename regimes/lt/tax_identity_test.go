package lt_test

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
			name:      "valid 9 digits, first pass",
			inputCode: "119511515",
		},
		{
			name:      "valid 9 digits, second pass",
			inputCode: "100538411",
		},
		{
			name:      "valid 12 digits, first pass",
			inputCode: "100004278519",
		},
		{
			name:      "valid 12 digits, first pass, another code",
			inputCode: "100001919017",
		},
		{
			name:      "valid 12 digits, second pass",
			inputCode: "100008194913",
		},
		{
			name:      "valid 12 digits, both passes 10, check digit 0",
			inputCode: "100004801610",
		},
		{
			name:      "empty code",
			inputCode: "",
		},
		{
			name:         "9 digits, wrong check digit",
			inputCode:    "100538412",
			expectedErrs: []string{"IDENTITY-02"},
		},
		{
			name:         "12 digits, wrong check digit",
			inputCode:    "100001919018",
			expectedErrs: []string{"IDENTITY-02"},
		},
		{
			name:         "8th digit is not 1, checksum fails",
			inputCode:    "119511525",
			expectedErrs: []string{"IDENTITY-01", "IDENTITY-02"},
		},
		{
			name:         "8th digit is not 1, checksum still valid",
			inputCode:    "100538427",
			expectedErrs: []string{"IDENTITY-01"},
		},
		{
			name:         "11th digit is not 1, checksum still valid",
			inputCode:    "100004278520",
			expectedErrs: []string{"IDENTITY-01"},
		},
		{
			name:         "wrong length",
			inputCode:    "12345678",
			expectedErrs: []string{"IDENTITY-01", "IDENTITY-02"},
		},
		{
			name:         "contains letters",
			inputCode:    "1005384A1",
			expectedErrs: []string{"IDENTITY-01", "IDENTITY-02"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "LT", Code: tt.inputCode}
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
			name:         "strips LT prefix and space",
			inputCode:    "LT 100538411",
			expectedCode: "100538411",
		},
		{
			name:         "strips LT prefix",
			inputCode:    "LT100538411",
			expectedCode: "100538411",
		},
		{
			name:         "strips lowercase prefix",
			inputCode:    "lt100538411",
			expectedCode: "100538411",
		},
		{
			name:         "strips hyphens",
			inputCode:    "100-538-411",
			expectedCode: "100538411",
		},
		{
			name:         "already normalized",
			inputCode:    "100538411",
			expectedCode: "100538411",
		},
		{
			name:         "empty",
			inputCode:    "",
			expectedCode: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "LT", Code: tt.inputCode}
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
