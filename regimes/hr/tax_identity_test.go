package hr_test

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
	allCodes := []string{"GOBL-HR-TAX-IDENTITY-01", "GOBL-HR-TAX-IDENTITY-02"}

	tests := []struct {
		name         string
		code         cbc.Code
		expectedErrs []string
	}{
		{
			name: "empty code",
			code: "",
		},
		{
			name: "valid OIB 1",
			code: "12345678903",
		},
		{
			name: "valid OIB 2",
			code: "11111111119",
		},
		{
			name: "valid OIB, check digit 0 from a remainder of 10",
			code: "00000000010",
		},
		{
			name:         "too short",
			code:         "1234567890",
			expectedErrs: []string{"GOBL-HR-TAX-IDENTITY-01", "GOBL-HR-TAX-IDENTITY-02"},
		},
		{
			name:         "too long",
			code:         "123456789012",
			expectedErrs: []string{"GOBL-HR-TAX-IDENTITY-01", "GOBL-HR-TAX-IDENTITY-02"},
		},
		{
			name:         "contains letters",
			code:         "1234567890A",
			expectedErrs: []string{"GOBL-HR-TAX-IDENTITY-01", "GOBL-HR-TAX-IDENTITY-02"},
		},
		{
			name:         "invalid checksum",
			code:         "12345678900",
			expectedErrs: []string{"GOBL-HR-TAX-IDENTITY-02"},
		},
		{
			name:         "invalid checksum, placeholder used in the official HR CIUS examples",
			code:         "12345678901",
			expectedErrs: []string{"GOBL-HR-TAX-IDENTITY-02"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "HR", Code: tt.code}
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
		name     string
		code     cbc.Code
		expected cbc.Code
	}{
		{
			name:     "already normalized",
			code:     "12345678903",
			expected: "12345678903",
		},
		{
			name:     "with spaces",
			code:     "1234 5678 903",
			expected: "12345678903",
		},
		{
			name:     "strips HR prefix",
			code:     "HR12345678903",
			expected: "12345678903",
		},
		{
			name:     "strips lowercase prefix",
			code:     "hr12345678903",
			expected: "12345678903",
		},
		{
			name:     "empty",
			code:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "HR", Code: tt.code}
			norm.Normalize(tID)
			assert.Equal(t, tt.expected, tID.Code)
		})
	}
}
