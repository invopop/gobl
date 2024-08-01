package pl_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/pl"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{
			Code:     "9551893317",
			Expected: "9551893317",
		},
		{
			Code:     "PL9551893317",
			Expected: "9551893317",
		},
		{
			Code:     "955-189-33.17",
			Expected: "9551893317",
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: "PL", Code: ts.Code}
		err := pl.Calculate(tID)
		assert.NoError(t, err)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestTaxIdentityValidation(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "good 1", code: "9551893317"},
		{name: "good 2", code: "1132191233"},
		{name: "good 3", code: "5841896486"},
		{name: "good 4", code: "7010009325"},
		{
			name: "bad mid length",
			code: "12345678910",
			err:  "invalid format",
		},
		{
			name: "too long",
			code: "1234567890123",
			err:  "invalid format",
		},
		{
			name: "too short",
			code: "123456",
			err:  "invalid format",
		},
		{
			name: "not normalized",
			code: "12.449.965-4",
			err:  "invalid format",
		},
		{
			name: "bad format",
			code: "1002191233",
			err:  "invalid format",
		},
		{
			name: "bad checksum",
			code: "9551893318",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "PL", Code: tt.code}
			err := pl.Validate(tID)
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
