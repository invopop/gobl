package mt_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	_ "github.com/invopop/gobl/regimes/mt"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityValidation(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid 1", code: "12701906"},
		{name: "valid 2", code: "12357210"},
		{name: "valid 3", code: "13043536"},
		{name: "valid mod-37 edge (check digits 37)", code: "00000037"},
		{name: "checksum near-miss", code: "12345678", err: "IDENTITY-01"},
		{name: "too short", code: "1270190", err: "IDENTITY-01"},
		{name: "too long", code: "127019066", err: "IDENTITY-01"},
		{name: "non-numeric", code: "1270190A", err: "IDENTITY-01"},
		{name: "empty is skipped", code: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &tax.Identity{
				Country: "MT",
				Code:    tt.code,
			}
			err := rules.Validate(id)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.err)
				}
			}
		})
	}

	t.Run("nil", func(t *testing.T) {
		var id *tax.Identity
		err := rules.Validate(id)
		assert.NoError(t, err)
	})
}

func TestTaxIdentityNormalization(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		want cbc.Code
	}{
		{name: "strips prefix and spaces", code: "MT 1270 1906", want: "12701906"},
		{name: "lowercase prefix", code: "mt12701906", want: "12701906"},
		{name: "already clean", code: "12701906", want: "12701906"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &tax.Identity{Country: "MT", Code: tt.code}
			norm.Normalize(id)
			assert.Equal(t, tt.want, id.Code)
		})
	}
}
