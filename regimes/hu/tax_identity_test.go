package hu_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/hu"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	r := hu.New()
	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{
			Code:     "12345678-1-12",
			Expected: "12345678112",
		},
		{
			Code:     "HU12345678-1-12",
			Expected: "12345678112",
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: "HU", Code: ts.Code}
		r.NormalizeObject(tID)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{"Empty code", "", ""},
		{"Invalid length (5)", "12345", "invalid length"},
		{"Invalid length (10)", "1234567890", "invalid length"},
		{"Invalid check digit", "12345678", "checksum mismatch"},
		{"Invalid VAT code", "21114445623", "invalid VAT code"},
		{"Invalid area code", "82713452101", "invalid area code"},
		{"Valid code (8 chars)", "98109858", ""},
		{"Valid code (11 chars)", "88212131103", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "HU", Code: tt.code}
			err := hu.Validate(tID)
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
