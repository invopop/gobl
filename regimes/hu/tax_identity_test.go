package hu_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/hu"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{Code: "13895459", Expected: "13895459"},
		{Code: "HU13895459", Expected: "13895459"},
		{Code: "1389 5459", Expected: "13895459"},
		{Code: "13895459-2-41", Expected: "13895459"},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: "HU", Code: ts.Code}
		hu.Normalize(tID)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestTaxIdentityRules(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		// Valid tax codes — verified against python-stdnum and farkasdezso.hu
		{name: "valid - farkasdezso.hu example", code: "13895459"},
		{name: "valid - OTP Bank", code: "10537914"},
		{name: "valid - MOL Group", code: "10625790"},
		{name: "valid - Magyar Telekom", code: "10773381"},
		{name: "valid - Richter Gedeon", code: "10484878"},
		{name: "valid - check digit zero", code: "10000070"},
		{
			name: "too short",
			code: "1234567",
			err:  "IDENTITY-01",
		},
		{
			name: "too long",
			code: "123456789",
			err:  "IDENTITY-01",
		},
		{
			name: "contains letters",
			code: "1234567A",
			err:  "IDENTITY-01",
		},
		{
			name: "starts with zero",
			code: "01234567",
			err:  "IDENTITY-01",
		},
		{
			name: "bad checksum",
			code: "13895450",
			err:  "IDENTITY-01",
		},
		{
			name: "bad checksum - off by one",
			code: "10537915",
			err:  "IDENTITY-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "HU", Code: tt.code}
			err := rules.Validate(tID)
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
