package lu_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/regimes/lu"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    cbc.Code
		expected cbc.Code
	}{
		{name: "already normalized", input: "26375245", expected: "26375245"},
		{name: "with LU prefix", input: "LU26375245", expected: "26375245"},
		{name: "lowercase prefix", input: "lu26375245", expected: "26375245"},
		{name: "with spaces", input: "LU 2637 5245", expected: "26375245"},
		{name: "dashes", input: "LU-2637-5245", expected: "26375245"},
		{name: "no prefix with spaces", input: "2637 5245", expected: "26375245"},
		{name: "empty", input: "", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: lu.CountryCode, Code: tt.input}
			norm.Normalize(tID)
			assert.Equal(t, tt.expected, tID.Code)
		})
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		code      cbc.Code
		valid     bool
		errorCode string
	}{
		// 263752 mod 89 = 45  →  valid
		{name: "valid code", code: "26375245", valid: true},
		// 123456 mod 89 = 13  →  valid
		{name: "valid code 2", code: "12345613", valid: true},
		// empty code is allowed (TVA registration is not mandatory for all businesses)
		{name: "empty code", code: "", valid: true},
		// wrong check digits
		{name: "bad checksum", code: "26375200", errorCode: "IDENTITY-03"},
		{name: "bad checksum 2", code: "12345699", errorCode: "IDENTITY-03"},
		// wrong length
		{name: "too short", code: "2637524", errorCode: "IDENTITY-01"},
		{name: "too long", code: "263752450", errorCode: "IDENTITY-01"},
		// non-numeric
		{name: "contains letter", code: "2637524A", errorCode: "IDENTITY-02"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: lu.CountryCode, Code: tt.code}
			err := rules.Validate(tID)
			if tt.valid {
				assert.NoError(t, err)
			} else if assert.Error(t, err) {
				assert.Contains(t, err.Error(), tt.errorCode)
			}
		})
	}
}
