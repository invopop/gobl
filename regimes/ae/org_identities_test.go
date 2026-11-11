package ae_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/ae"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdentityDefinitions(t *testing.T) {
	r := tax.RegimeDefFor(l10n.AE)
	require.NotNil(t, r)

	keys := []cbc.Key{
		ae.IdentityKeyTIN,
		ae.IdentityKeyTradeLicense,
		ae.IdentityKeyEmiratesID,
		ae.IdentityKeyCabinetDecision,
	}
	assert.Len(t, r.Identities, len(keys))
	for _, key := range keys {
		def := cbc.GetKeyDefinition(key, r.Identities)
		if assert.NotNil(t, def, key) {
			assert.NotEmpty(t, def.Name.In(i18n.EN), key)
			assert.NotEmpty(t, def.Name.In(i18n.AR), key)
		}
	}
}

func TestTINValidation(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid", code: "1234567890"},
		{name: "first 10 digits of a TRN", code: "1345678901"},

		{name: "too short", code: "123456789", err: "[GOBL-AE-ORG-IDENTITY-01]"},
		{name: "too long", code: "12345678901", err: "[GOBL-AE-ORG-IDENTITY-01]"},
		{name: "a full TRN", code: "134567890123003", err: "[GOBL-AE-ORG-IDENTITY-01]"},
		{name: "non-numeric", code: "12345ABCDE", err: "[GOBL-AE-ORG-IDENTITY-01]"},
	}

	ctx := tax.RegimeContext(ae.CountryCode)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Key: ae.IdentityKeyTIN, Code: tt.code}
			err := rules.Validate(id, ctx)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}
		})
	}

	t.Run("the rule only applies under the AE regime", func(t *testing.T) {
		id := &org.Identity{Key: ae.IdentityKeyTIN, Code: "123"}
		err := rules.Validate(id, tax.RegimeContext("DE"))
		if err != nil {
			assert.NotContains(t, err.Error(), "GOBL-AE-ORG-IDENTITY-01")
		}
	})
}

func TestEmiratesIDValidation(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid", code: "784199012345671"},

		{name: "too short", code: "78419901234567", err: "[GOBL-AE-ORG-IDENTITY-02]"},
		{name: "too long", code: "7841990123456712", err: "[GOBL-AE-ORG-IDENTITY-02]"},
		{name: "non-numeric", code: "784ABCD12345671", err: "[GOBL-AE-ORG-IDENTITY-02]"},
	}

	ctx := tax.RegimeContext(ae.CountryCode)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Key: ae.IdentityKeyEmiratesID, Code: tt.code}
			err := rules.Validate(id, ctx)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}
		})
	}
}

func TestIdentityNormalization(t *testing.T) {
	ctx := tax.RegimeContext(ae.CountryCode)

	t.Run("Emirates ID separators are removed", func(t *testing.T) {
		id := &org.Identity{Key: ae.IdentityKeyEmiratesID, Code: "784-1990-1234567-1"}
		norm.Normalize(id, ctx)
		assert.Equal(t, cbc.Code("784199012345671"), id.Code)
		assert.NoError(t, rules.Validate(id, ctx))
	})

	t.Run("TIN separators are removed", func(t *testing.T) {
		id := &org.Identity{Key: ae.IdentityKeyTIN, Code: "1234 567 890"}
		norm.Normalize(id, ctx)
		assert.Equal(t, cbc.Code("1234567890"), id.Code)
	})

	t.Run("other identities are left as they are", func(t *testing.T) {
		id := &org.Identity{Key: ae.IdentityKeyTradeLicense, Code: "CN-1234567"}
		norm.Normalize(id, ctx)
		assert.Equal(t, cbc.Code("CN-1234567"), id.Code)
	})

	t.Run("TIN gets the tax scope", func(t *testing.T) {
		id := &org.Identity{Key: ae.IdentityKeyTIN, Code: "1234567890"}
		norm.Normalize(id, ctx)
		assert.Equal(t, org.IdentityScopeTax, id.Scope)
	})

	t.Run("commercial registration numbers get the legal scope", func(t *testing.T) {
		for _, key := range []cbc.Key{
			ae.IdentityKeyTradeLicense,
			ae.IdentityKeyEmiratesID,
			org.IdentityKeyPassport,
			ae.IdentityKeyCabinetDecision,
		} {
			id := &org.Identity{Key: key, Code: "123"}
			norm.Normalize(id, ctx)
			assert.Equal(t, org.IdentityScopeLegal, id.Scope, key)
		}
	})

	t.Run("an explicit scope is kept", func(t *testing.T) {
		id := &org.Identity{Key: ae.IdentityKeyTradeLicense, Scope: org.IdentityScopeTax, Code: "123"}
		norm.Normalize(id, ctx)
		assert.Equal(t, org.IdentityScopeTax, id.Scope)
	})

	t.Run("identities without a UAE key get no scope", func(t *testing.T) {
		id := &org.Identity{Type: "XYZ", Code: "123"}
		norm.Normalize(id, ctx)
		assert.Empty(t, id.Scope)
	})

	t.Run("scope is only set under the AE regime", func(t *testing.T) {
		id := &org.Identity{Key: org.IdentityKeyPassport, Code: "123"}
		norm.Normalize(id, tax.RegimeContext("DE"))
		assert.Empty(t, id.Scope)
	})
}

func TestIdentitiesWithoutFormatRules(t *testing.T) {
	ctx := tax.RegimeContext(ae.CountryCode)
	for _, key := range []cbc.Key{ae.IdentityKeyTradeLicense, ae.IdentityKeyCabinetDecision} {
		id := &org.Identity{Key: key, Code: "CN-1234567"}
		assert.NoError(t, rules.Validate(id, ctx), key)
	}
}
