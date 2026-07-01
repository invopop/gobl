package is_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/regimes/is"
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
		{name: "already normalized", input: "5012031220", expected: "5012031220"},
		{name: "lowercase prefix", input: "is5012031220", expected: "5012031220"},
		{name: "uppercase prefix", input: "IS5012031220", expected: "5012031220"},
		{name: "hyphenated with prefix", input: "IS-1234567890", expected: "1234567890"},
		{name: "spaces and prefix", input: "  IS 1234567890  ", expected: "1234567890"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "IS", Code: tt.input}
			norm.Normalize(tID)
			assert.Equal(t, tt.expected, tID.Code)
		})
	}
}

func TestTaxIdentityRules(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid company kennitala", code: "5012031220"},
		{name: "empty code", code: ""},
		{name: "valid person kennitala rejected for VAT", code: "0902862349", err: "person kennitala is not valid for VAT"},
		{name: "temporary kennitala starting with 8", code: "8101850150", err: "temporary kennitala is not valid for VAT"},
		{name: "wrong checksum", code: "5012031230", err: "invalid checksum for kennitala"},
		{name: "wrong length", code: "501203122", err: "invalid kennitala length"},
		{name: "non-numeric characters", code: "501203122A", err: "invalid kennitala format"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "IS", Code: tt.code}
			err := rules.Validate(tID)
			if tt.err == "" {
				assert.NoError(t, err)
			} else if assert.Error(t, err) {
				assert.Contains(t, err.Error(), tt.err)
			}
		})
	}
}

// TestKennitalaChecksum covers the verification anchor from both angles: its
// checksum is valid, yet it is a person kennitala and so must fail VAT-ID checks.
func TestKennitalaChecksum(t *testing.T) {
	t.Parallel()

	// Anchor: 0902862349 — valid checksum, but a PERSON (first digit 0).
	assert.True(t, is.ValidKennitala("0902862349"), "anchor checksum should be valid")
	assert.True(t, is.Person("0902862349"))
	assert.False(t, is.Company("0902862349"))

	// Valid company kennitala.
	assert.True(t, is.ValidKennitala("5012031220"))
	assert.True(t, is.Company("5012031220"))

	// Same number, tampered check digit -> invalid checksum.
	assert.False(t, is.ValidKennitala("5012031230"))

	// Temporary kennitala (first digit 8) is neither person nor company form.
	assert.False(t, is.Company("8101850150"))
	assert.False(t, is.Person("8101850150"))
}
