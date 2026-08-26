package jp_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	_ "github.com/invopop/gobl/regimes/jp"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		want cbc.Code
	}{
		{name: "already normalized", code: "T1234567890123", want: "T1234567890123"},
		{name: "lowercase", code: "t1234567890123", want: "T1234567890123"},
		{name: "separators", code: "T-1234-567890123", want: "T1234567890123"},
		{name: "spaces", code: " T 1234567890123 ", want: "T1234567890123"},
		{name: "country prefix", code: "JPT1234567890123", want: "T1234567890123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &tax.Identity{Country: "JP", Code: tt.code}
			norm.Normalize(id)
			assert.Equal(t, tt.want, id.Code)
		})
	}
}

func TestTaxIdentityRules(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid corporate-derived", code: "T8700110005901"},
		{name: "valid", code: "T1234567890123"},
		{name: "missing T prefix", code: "1234567890123", err: "IDENTITY-01"},
		{name: "too short", code: "T123456789012", err: "IDENTITY-01"},
		{name: "too long", code: "T12345678901234", err: "IDENTITY-01"},
		{name: "letters in numeric part", code: "T123456789012A", err: "IDENTITY-01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &tax.Identity{Country: "JP", Code: tt.code}
			err := rules.Validate(id)
			if tt.err == "" {
				assert.NoError(t, err)
			} else if assert.Error(t, err) {
				assert.Contains(t, err.Error(), tt.err)
			}
		})
	}

	t.Run("empty code", func(t *testing.T) {
		id := &tax.Identity{Country: "JP"}
		assert.NoError(t, rules.Validate(id))
	})

	t.Run("nil", func(t *testing.T) {
		var id *tax.Identity
		assert.NoError(t, rules.Validate(id))
	})
}
