package in_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/in"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		name     string
		input    cbc.Code
		expected cbc.Code
	}{
		{name: "already normalized", input: "27AAPFU0939F1ZV", expected: "27AAPFU0939F1ZV"},
		{name: "lowercase input", input: "27aapfu0939f1zv", expected: "27AAPFU0939F1ZV"},
		{name: "mixed case input", input: "27AaPfU0939F1zV", expected: "27AAPFU0939F1ZV"},
		{name: "extra spaces", input: "  27AAPFU0939F1ZV  ", expected: "27AAPFU0939F1ZV"},
		{name: "special characters", input: "27-AAPFU0939F1-ZV", expected: "27AAPFU0939F1ZV"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "IN", Code: tt.input}

			in.Normalize(tID)

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
		{name: "valid GSTIN 1", code: "27AAPFU0939F1ZV"},
		{name: "valid GSTIN 2", code: "29AAGCB7383J1Z4"},
		{name: "valid GSTIN 3", code: "10AABCU9355J1Z9"},
		{name: "valid GSTIN 4", code: "09AABCU9355J1ZS"},

		{
			name: "too short",
			code: "27AAPFU0939F",
			err:  "invalid GSTIN format",
		},
		{
			name: "state code not numeric",
			code: "AAAPFU0939F1ZV",
			err:  "invalid GSTIN format",
		},
		{
			name: "invalid checksum",
			code: "27AAPFU0939F1Z0",
			err:  "checksum mismatch",
		},
		{
			name: "too long",
			code: "27AAPFU0939F1ZV12",
			err:  "invalid GSTIN format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "IN", Code: tt.code}

			err := in.Validate(tID)

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
