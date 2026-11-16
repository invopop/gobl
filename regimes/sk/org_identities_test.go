package sk_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/sk"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdentityTypeDefinitions(t *testing.T) {
	assert.Equal(t, "ICO", string(sk.IdentityTypeICO))

	r := tax.RegimeDefFor("SK")
	require.NotNil(t, r)
	require.Len(t, r.Identities, 1)
	assert.Equal(t, sk.IdentityTypeICO, r.Identities[0].Code)
}

func TestValidateICO(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid company", code: "36396567"},
		{name: "valid Statistical Office", code: "00166197"},
		{name: "valid remainder 0 gives check digit 1", code: "00000001"},
		{name: "valid remainder 1 gives check digit 0", code: "00000060"},
		{name: "wrong check digit", code: "36396568", err: "[GOBL-SK-ORG-IDENTITY-02]"},
		{name: "too short", code: "3639656", err: "[GOBL-SK-ORG-IDENTITY-01]"},
		{name: "too long", code: "363965670", err: "[GOBL-SK-ORG-IDENTITY-01]"},
		{name: "letter in the middle", code: "36A96567", err: "[GOBL-SK-ORG-IDENTITY-01]"},
		{name: "empty code", code: "", err: "[GOBL-ORG-IDENTITY-01]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := validInvoice()
			inv.Supplier.Identities = []*org.Identity{{Type: sk.IdentityTypeICO, Code: tt.code}}
			require.NoError(t, inv.Calculate())
			err := rules.Validate(inv)
			if tt.err == "" {
				assert.NoError(t, err)
				return
			}
			assert.ErrorContains(t, err, tt.err)
		})
	}
}

// A malformed code fails the format rule only, never the check digit as well.
func TestValidateICOReportsOneRule(t *testing.T) {
	inv := validInvoice()
	inv.Supplier.Identities = []*org.Identity{{Type: sk.IdentityTypeICO, Code: "1234"}}
	require.NoError(t, inv.Calculate())
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "[GOBL-SK-ORG-IDENTITY-01]")
	assert.NotContains(t, err.Error(), "[GOBL-SK-ORG-IDENTITY-02]")
}

// Other regimes keep their own meaning of an identity typed ICO.
func TestValidateICOOutsideSlovakia(t *testing.T) {
	id := &org.Identity{Type: sk.IdentityTypeICO, Code: "123"}
	assert.NoError(t, rules.Validate(id))
}

func TestNormalizeICO(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		want cbc.Code
	}{
		{name: "groups of three", code: "36 396 567", want: "36396567"},
		{name: "no-break spaces", code: "36 396 567", want: "36396567"},
		{name: "dots", code: "36.396.567", want: "36396567"},
		{name: "dash and slash", code: "36-396/567", want: "36396567"},
		{name: "letters kept", code: "36A396567", want: "36A396567"},
		{name: "already normalized", code: "36396567", want: "36396567"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Type: sk.IdentityTypeICO, Code: tt.code}
			norm.Normalize(id, tax.RegimeContext("SK"))
			assert.Equal(t, tt.want, id.Code)
		})
	}

	t.Run("letters still fail after normalizing", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Identities = []*org.Identity{{Type: sk.IdentityTypeICO, Code: "36A396567"}}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-SK-ORG-IDENTITY-01]")
	})

	t.Run("other types left alone", func(t *testing.T) {
		id := &org.Identity{Type: "OTHER", Code: "36 396 567"}
		norm.Normalize(id, tax.RegimeContext("SK"))
		assert.Equal(t, cbc.Code("36 396 567"), id.Code)
	})

	t.Run("scope left alone", func(t *testing.T) {
		id := &org.Identity{Type: sk.IdentityTypeICO, Code: "36396567"}
		norm.Normalize(id, tax.RegimeContext("SK"))
		assert.Equal(t, cbc.KeyEmpty, id.Scope)
	})

	t.Run("nil identity", func(t *testing.T) {
		assert.NotPanics(t, func() {
			norm.Normalize((*org.Identity)(nil), tax.RegimeContext("SK"))
		})
	})
}
