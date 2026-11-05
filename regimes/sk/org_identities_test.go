package sk_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/org"
	_ "github.com/invopop/gobl/regimes"
	"github.com/invopop/gobl/regimes/sk"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestICOIdentity(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
		err  string
	}{
		{name: "valid", code: "36396567", want: "36396567"},
		{name: "written in groups", code: "36 396 567", want: "36396567"},
		{name: "with dashes", code: "36-396-567", want: "36396567"},
		{name: "too short", code: "3639656", want: "3639656", err: "[GOBL-SK-ORG-IDENTITY-01]"},
		{name: "too long", code: "363965670", want: "363965670", err: "[GOBL-SK-ORG-IDENTITY-01]"},
		{name: "with dots and slashes", code: "36.396/567", want: "36396567"},
		{name: "letter at the end", code: "3639656A", want: "3639656A", err: "[GOBL-SK-ORG-IDENTITY-01]"},
		{name: "letter in the middle", code: "36A396567", want: "36A396567", err: "[GOBL-SK-ORG-IDENTITY-01]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := validInvoice()
			inv.Supplier.Identities = []*org.Identity{
				{Type: sk.IdentityTypeICO, Code: cbc.Code(tt.code)},
			}
			require.NoError(t, inv.Calculate())
			assert.Equal(t, tt.want, inv.Supplier.Identities[0].Code.String())
			err := rules.Validate(inv)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}
		})
	}

	t.Run("gets the legal scope", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{Type: sk.IdentityTypeICO, Code: "36396567"},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, org.IdentityScopeLegal, inv.Supplier.Identities[0].Scope)
	})

	t.Run("other identity keys are not checked", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{Key: org.IdentityKeyOther, Code: "1"},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("not applied outside the SK regime", func(t *testing.T) {
		id := &org.Identity{Type: sk.IdentityTypeICO, Code: "36 396 567"}
		norm.Normalize(id, tax.RegimeContext("ES"))
		assert.Equal(t, "36 396 567", id.Code.String())
	})
}
