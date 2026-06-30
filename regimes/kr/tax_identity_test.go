package kr_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/regimes/kr"
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
		{name: "already normalized", input: "1208147521", expected: "1208147521"},
		{name: "with hyphens", input: "120-81-47521", expected: "1208147521"},
		{name: "with spaces", input: "120 81 47521", expected: "1208147521"},
		{name: "with country prefix", input: "KR1208147521", expected: "1208147521"},
		{name: "lowercase country prefix", input: "kr1208147521", expected: "1208147521"},
		{name: "empty code", input: "", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: kr.CountryCode, Code: tt.input}
			norm.Normalize(tID)
			assert.Equal(t, tt.expected, tID.Code)
		})
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		code  cbc.Code
		valid bool
	}{
		// Real, verifiable business registration numbers.
		{name: "valid BRN (Kakao)", code: "1208147521", valid: true},
		{name: "valid BRN (Naver)", code: "2208162517", valid: true},
		// AssertIfPresent skips empty codes; presence is enforced elsewhere.
		{name: "empty code", code: "", valid: true},
		{name: "bad check digit", code: "1208147520"},
		{name: "too short", code: "120814752"},
		{name: "too long", code: "12081475210"},
		{name: "non-numeric", code: "120814752A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: kr.CountryCode, Code: tt.code}
			err := rules.Validate(tID)
			if tt.valid {
				assert.NoError(t, err)
			} else if assert.Error(t, err) {
				assert.Contains(t, err.Error(), "invalid Korean business registration number")
			}
		})
	}
}
