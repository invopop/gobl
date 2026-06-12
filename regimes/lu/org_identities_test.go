package lu_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/lu"
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
		{name: "already normalized", input: "B263475", expected: "B263475"},
		{name: "with spaces", input: "B 263 475", expected: "B263475"},
		{name: "lowercase", input: "b263475", expected: "B263475"},
		{name: "lowercase with spaces", input: "b 263 475", expected: "B263475"},
		{name: "F register", input: "F 12345", expected: "F12345"},
		{name: "G register", input: "G 45678", expected: "G45678"},
		{name: "H register", input: "H 1234", expected: "H1234"},
		{name: "empty", input: "", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Type: lu.IdentityTypeRCS, Code: tt.input}
			norm.Normalize(id, tax.RegimeContext(lu.CountryCode))
			assert.Equal(t, tt.expected, id.Code)
		})
	}

	t.Run("unknown type not normalised", func(t *testing.T) {
		id := &org.Identity{Type: "OTHER", Code: "b 123"}
		norm.Normalize(id, tax.RegimeContext(lu.CountryCode))
		assert.Equal(t, cbc.Code("b 123"), id.Code)
	})
}

func TestValidateOrgIdentity(t *testing.T) {
	t.Parallel()

	opts := []rules.WithContext{
		tax.RegimeContext(lu.CountryCode),
	}

	tests := []struct {
		name  string
		code  cbc.Code
		valid bool
	}{
		{name: "B register", code: "B263475", valid: true},
		{name: "F register", code: "F12345", valid: true},
		{name: "G register", code: "G45678", valid: true},
		{name: "H register", code: "H1234", valid: true},
		{name: "single digit", code: "B1", valid: true},
		// invalid cases
		{name: "wrong letter", code: "A263475"},
		{name: "lowercase", code: "b263475"},
		{name: "no digits", code: "B"},
		{name: "too many digits", code: "B1234567"},
		{name: "digits only", code: "263475"},
		{name: "empty code", code: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Type: lu.IdentityTypeRCS, Code: tt.code}
			err := rules.Validate(id, opts...)
			if tt.valid {
				assert.NoError(t, err)
			} else if assert.Error(t, err) {
				assert.Contains(t, err.Error(), "invalid RCS number")
			}
		})
	}

	t.Run("nil identity", func(t *testing.T) {
		assert.NoError(t, rules.Validate((*org.Identity)(nil), opts...))
	})

	t.Run("unknown type skipped", func(t *testing.T) {
		id := &org.Identity{Type: "OTHER", Code: "A9999999"}
		assert.NoError(t, rules.Validate(id, opts...))
	})
}
