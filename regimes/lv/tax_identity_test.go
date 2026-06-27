package lv_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	_ "github.com/invopop/gobl/regimes/lv"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityNormalization(t *testing.T) {
	tests := []struct {
		name     string
		input    cbc.Code
		expected cbc.Code
	}{
		{name: "legal entity with uppercase LV prefix", input: "LV40003032065", expected: "40003032065"},
		{name: "legacy personal with uppercase-lowercase LV prefix", input: "Lv01019012345", expected: "01019012345"},
		{name: "modern personal with lowecase LV prefix", input: "lv32123456789", expected: "32123456789"},
		{name: "legal entity without prefix", input: "40003032065", expected: "40003032065"},
		{name: "legacy personal without prefix", input: "01019012345", expected: "01019012345"},
		{name: "modern personal without prefix", input: "32123456789", expected: "32123456789"},
		{name: "legacy personal with hyphen and LV prefix", input: "LV010190-12345", expected: "01019012345"},
		{name: "legacy personal with hyphen no prefix", input: "010190-12345", expected: "01019012345"},
		{name: "legacy personal with hyphen 15 June", input: "LV150690-12345", expected: "15069012345"},
		{name: "legacy personal with hyphen no prefix 15 June", input: "150690-12345", expected: "15069012345"},
		{name: "modern personal with hyphen and LV prefix", input: "LV3212345-6789", expected: "32123456789"},
		{name: "modern personal with hyphen no prefix", input: "3212345-6789", expected: "32123456789"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "LV", Code: tt.input}
			norm.Normalize(tID)
			assert.Equal(t, tt.expected, tID.Code)
		})
	}
}

func TestTaxIdentityRules(t *testing.T) {
	tests := []struct {
		name     string
		input    cbc.Code
		expected string // error code or empty for valid
	}{
		// Legal entity codes (starting with 4-9)
		{name: "valid legal entity with prefix 4", input: "LV40003032065", expected: ""},
		{name: "valid legal entity with prefix 9", input: "LV99000000001", expected: ""},
		{name: "valid legal entity with prefix 5", input: "LV52000000000", expected: ""},
		{name: "valid legal entity with prefix 6", input: "LV64000000000", expected: ""},
		{name: "valid legal entity with prefix 4 - alt", input: "LV40003521600", expected: ""},

		// Legacy personal codes (DDMMYYCXXXX format, pre-2017)
		// Format: DDMMYYCXXXX where C is century (0=18xx, 1=19xx, 2=20xx)
		{name: "valid legacy personal - 1 Jan 1990", input: "LV01019012345", expected: ""},
		{name: "valid legacy personal - 15 June 1990", input: "LV15069012345", expected: ""},
		{name: "valid legacy personal - 29 Feb 2000 leap year", input: "LV29020021234", expected: ""},

		// Invalid legacy personal codes - month validation
		{name: "invalid legacy - month 0", input: "LV01009012345", expected: "IDENTITY-01"},
		{name: "invalid legacy - month 13", input: "LV01139012345", expected: "IDENTITY-01"},
		// Invalid legacy personal codes - day validation
		{name: "invalid legacy - day 0", input: "LV00019012345", expected: "IDENTITY-01"},
		{name: "invalid legacy - April 31", input: "LV31040012345", expected: "IDENTITY-01"},
		{name: "invalid legacy - June 31", input: "LV31060012345", expected: "IDENTITY-01"},
		{name: "invalid legacy - Sep 31", input: "LV31090012345", expected: "IDENTITY-01"},
		{name: "invalid legacy - Nov 31", input: "LV31110012345", expected: "IDENTITY-01"},
		// Feb 29 in non-leap years
		{name: "invalid legacy - Feb 29 non-leap 1900", input: "LV29020012345", expected: "IDENTITY-01"},
		{name: "invalid legacy - Feb 29 non-leap 2014", input: "LV29021421234", expected: "IDENTITY-01"},

		// Modern personal codes (starting with 3X, post-1 July 2017)
		{name: "valid modern personal code 32X", input: "LV32123456789", expected: ""},
		{name: "valid modern personal code 33X", input: "LV33000000000", expected: ""},
		{name: "valid modern personal code 34X", input: "LV34000000000", expected: ""},
		{name: "valid modern personal code 35X", input: "LV35000000000", expected: ""},
		{name: "valid modern personal code 36X", input: "LV36000000000", expected: ""},
		{name: "valid modern personal code 37X", input: "LV37000000000", expected: ""},
		{name: "valid modern personal code 38X", input: "LV38000000000", expected: ""},
		{name: "valid modern personal code 39X", input: "LV39000000000", expected: ""},

		// Invalid modern personal codes - second digit must be 2-9
		{name: "invalid modern personal - second digit 0", input: "LV30000000000", expected: "IDENTITY-01"},
		{name: "invalid modern personal - second digit 1", input: "LV31000000000", expected: "IDENTITY-01"},

		// Invalid checksum for legal entities
		{name: "bad checksum for legal entity", input: "LV40003032066", expected: "IDENTITY-01"},

		// Invalid formats
		{name: "too long", input: "LV123456789012", expected: "IDENTITY-01"},
		{name: "too short", input: "LV1234567890", expected: "IDENTITY-01"},
		{name: "not enough digits", input: "LV123456789", expected: "IDENTITY-01"},
		{name: "letters in code", input: "LVAAAAAAAAAA", expected: "IDENTITY-01"},

		// Empty code is valid (not required)
		{name: "empty code", input: "", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "LV", Code: tt.input}
			err := rules.Validate(tID)
			if tt.expected == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.expected)
				}
			}
		})
	}
}