package no_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/no"
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
			no.Normalize(id)
			assert.Equal(t, tt.expected, id.Code)
		})
	}

	t.Run("nil identity", func(t *testing.T) {
		assert.NotPanics(t, func() {
			no.Normalize((*org.Identity)(nil))
		})
	})

	t.Run("unknown type ignored", func(t *testing.T) {
		id := &org.Identity{Type: "OTHER", Code: "123 456"}
		no.Normalize(id)
		assert.Equal(t, cbc.Code("123 456"), id.Code)
	})
}

func TestValidateOrgIdentity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid code", code: "923456783"},
		{name: "empty code", code: ""},
		{
			name: "bad check digit",
			code: "923456780",
			err:  "checksum mismatch",
		},
		{
			name: "too short",
			code: "92345678",
			err:  "must have 9 digits",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Type: no.IdentityTypeOrgNr, Code: tt.code}
			err := no.Validate(id)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.err)
				}
			}
		})
	}

	t.Run("nil identity", func(t *testing.T) {
		assert.NoError(t, no.Validate((*org.Identity)(nil)))
	})

	t.Run("unknown type skipped", func(t *testing.T) {
		id := &org.Identity{Type: "OTHER", Code: "invalid"}
		assert.NoError(t, no.Validate(id))
	})
}
