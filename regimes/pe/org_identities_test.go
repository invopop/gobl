package pe_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/pe"
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
		{name: "already clean", input: "40123456", expected: "40123456"},
		{name: "with dots", input: "40.123.456", expected: "40123456"},
		{name: "with spaces", input: "40 123 456", expected: "40123456"},
		{name: "empty code", input: "", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Type: pe.IdentityTypeDNI, Code: tt.input}
			norm.Normalize(id, tax.RegimeContext(pe.CountryCode))
			assert.Equal(t, tt.expected, id.Code)
		})
	}

	t.Run("other types left as typed", func(t *testing.T) {
		id := &org.Identity{Type: pe.IdentityTypeCE, Code: "X-123456"}
		norm.Normalize(id, tax.RegimeContext(pe.CountryCode))
		assert.Equal(t, cbc.Code("X-123456"), id.Code)
	})

	t.Run("unknown type ignored", func(t *testing.T) {
		id := &org.Identity{Type: "OTHER", Code: "40 123 456"}
		norm.Normalize(id, tax.RegimeContext(pe.CountryCode))
		assert.Equal(t, cbc.Code("40 123 456"), id.Code)
	})
}

func TestValidateOrgIdentity(t *testing.T) {
	t.Parallel()

	opts := []rules.WithContext{
		tax.RegimeContext(pe.CountryCode),
	}

	tests := []struct {
		name  string
		code  cbc.Code
		valid bool
	}{
		{name: "valid DNI", code: "40123456", valid: true},
		{name: "valid DNI with leading zero", code: "04123456", valid: true},
		{name: "empty code"},
		{name: "too short", code: "4012345"},
		{name: "too long", code: "401234567"},
		{name: "non-numeric", code: "4012345A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Type: pe.IdentityTypeDNI, Code: tt.code}
			err := rules.Validate(id, opts...)
			if tt.valid {
				assert.NoError(t, err)
			} else if assert.Error(t, err) {
				assert.Contains(t, err.Error(), "invalid DNI")
			}
		})
	}

	t.Run("CE not strictly validated", func(t *testing.T) {
		id := &org.Identity{Type: pe.IdentityTypeCE, Code: "X123456"}
		assert.NoError(t, rules.Validate(id, opts...))
	})

	t.Run("passport not strictly validated", func(t *testing.T) {
		id := &org.Identity{Type: pe.IdentityTypePassport, Code: "PA1234567"}
		assert.NoError(t, rules.Validate(id, opts...))
	})
}
