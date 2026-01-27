package au_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/au"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		name     string
		code     cbc.Code
		expected cbc.Code
	}{
		{
			name:     "already normalized",
			code:     "51824753556",
			expected: "51824753556",
		},
		{
			name:     "with AU prefix",
			code:     "AU51824753556",
			expected: "51824753556",
		},
		{
			name:     "with l10n.AU prefix",
			code:     "AU51824753556",
			expected: "51824753556",
		},
		{
			name:     "with spaces",
			code:     "51 824 753 556",
			expected: "51824753556",
		},
		{
			name:     "with hyphens",
			code:     "51-824-753-556",
			expected: "51824753556",
		},
		{
			name:     "with dots",
			code:     "51.824.753.556",
			expected: "51824753556",
		},
		{
			name:     "with mixed formatting",
			code:     "AU 51-824.753 556",
			expected: "51824753556",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "AU", Code: tt.code}
			au.Normalize(tID)
			assert.Equal(t, tt.expected, tID.Code)
		})
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		// Valid ABNs (verified using modulus 89 algorithm)
		{
			name: "valid ABN 1",
			code: "51824753556",
		},
		{
			name: "valid ABN 2",
			code: "53004085616",
		},
		{
			name: "valid ABN 3",
			code: "33102417032",
		},

		// Format errors
		{
			name: "too short",
			code: "5182475355",
			err:  "invalid format",
		},
		{
			name: "too long",
			code: "518247535567",
			err:  "invalid format",
		},
		{
			name: "contains letters",
			code: "5182475355A",
			err:  "invalid format",
		},
		{
			name: "contains spaces (not normalized)",
			code: "51 824 753 556",
			err:  "invalid format",
		},
		{
			name: "empty",
			code: "",
		},

		// Checksum errors
		{
			name: "bad checksum - last digit wrong",
			code: "51824753557",
			err:  "checksum mismatch",
		},
		{
			name: "bad checksum - first digit wrong",
			code: "61824753556",
			err:  "checksum mismatch",
		},
		{
			name: "bad checksum - middle digit wrong",
			code: "51824853556",
			err:  "checksum mismatch",
		},
		{
			name: "all zeros",
			code: "00000000000",
			err:  "checksum mismatch",
		},
		{
			name: "all ones",
			code: "11111111111",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "AU", Code: tt.code}
			err := au.Validate(tID)
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

func TestABNChecksumAlgorithm(t *testing.T) {
	// Test the modulus 89 algorithm with a known ABN: 51824753556
	//
	// Step 1: Subtract 1 from first digit: 5-1 = 4
	// Digits: [4, 1, 8, 2, 4, 7, 5, 3, 5, 5, 6]
	//
	// Step 2: Apply weights [10, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19]
	// Products: [40, 1, 24, 10, 28, 63, 55, 39, 75, 85, 114]
	//
	// Step 3: Sum = 40+1+24+10+28+63+55+39+75+85+114 = 534
	//
	// Step 4: 534 % 89 = 0 ✓
	//
	// This test documents the algorithm for future reference
	t.Run("algorithm verification", func(t *testing.T) {
		tID := &tax.Identity{Country: "AU", Code: "51824753556"}
		err := au.Validate(tID)
		assert.NoError(t, err, "51824753556 should be valid (534 % 89 = 0)")
	})
}
