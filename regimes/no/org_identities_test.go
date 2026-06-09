package no_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/no"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeOrgIdentity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    cbc.Code
		expected cbc.Code
	}{
		{name: "already clean", input: "923456783", expected: "923456783"},
		{name: "with spaces", input: "923 456 783", expected: "923456783"},
		{name: "with dashes", input: "923-456-783", expected: "923456783"},
		{name: "empty code", input: "", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Type: no.IdentityTypeOrgNr, Code: tt.input}
			norm.Normalize(id, tax.RegimeContext(no.CountryCode))
			assert.Equal(t, tt.expected, id.Code)
		})
	}

	t.Run("unknown type ignored", func(t *testing.T) {
		id := &org.Identity{Type: "OTHER", Code: "123 456"}
		norm.Normalize(id, tax.RegimeContext(no.CountryCode))
		assert.Equal(t, cbc.Code("123 456"), id.Code)
	})
}

func TestValidateOrgIdentity(t *testing.T) {
	t.Parallel()

	opts := []rules.WithContext{
		tax.RegimeContext(no.CountryCode),
	}

	tests := []struct {
		name  string
		code  cbc.Code
		valid bool
	}{
		{name: "valid code", code: "923456783", valid: true},
		{name: "empty code", code: ""},
		{name: "bad check digit", code: "923456780"},
		{name: "too short", code: "92345678"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Type: no.IdentityTypeOrgNr, Code: tt.code}
			err := rules.Validate(id, opts...)
			if tt.valid {
				assert.NoError(t, err)
			} else if assert.Error(t, err) {
				assert.Contains(t, err.Error(), "invalid organisasjonsnummer")
			}
		})
	}

	t.Run("nil identity", func(t *testing.T) {
		assert.NoError(t, rules.Validate((*org.Identity)(nil), opts...))
	})

	t.Run("unknown type skipped", func(t *testing.T) {
		id := &org.Identity{Type: "OTHER", Code: "923456783"}
		assert.NoError(t, rules.Validate(id, opts...))
	})
}
