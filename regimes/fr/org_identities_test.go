package fr_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateIdentity(t *testing.T) {
	opts := []rules.WithContext{
		tax.RegimeContext(fr.CountryCode),
	}
	t.Run("valid SIREN", func(t *testing.T) {
		id := &org.Identity{
			Type: fr.IdentityTypeSIREN,
			Code: "123456789",
		}
		err := rules.Validate(id, opts...)
		assert.NoError(t, err)
	})

	t.Run("invalid SIREN - too short", func(t *testing.T) {
		id := &org.Identity{
			Type: fr.IdentityTypeSIREN,
			Code: "12345",
		}
		err := rules.Validate(id, opts...)
		assert.ErrorContains(t, err, "[GOBL-FR-ORG-IDENTITY-01] ($.code) identity code for type SIREN must be valid")
	})

	t.Run("invalid SIREN - too long", func(t *testing.T) {
		id := &org.Identity{
			Type: fr.IdentityTypeSIREN,
			Code: "1234567890",
		}
		err := rules.Validate(id, opts...)
		assert.ErrorContains(t, err, "[GOBL-FR-ORG-IDENTITY-01]")
	})

	t.Run("invalid SIREN - non-numeric", func(t *testing.T) {
		id := &org.Identity{
			Type: fr.IdentityTypeSIREN,
			Code: "12345678A",
		}
		err := rules.Validate(id, opts...)
		assert.ErrorContains(t, err, "[GOBL-FR-ORG-IDENTITY-01]")
	})

	t.Run("valid SIRET", func(t *testing.T) {
		id := &org.Identity{
			Type: fr.IdentityTypeSIRET,
			Code: "12345678901234",
		}
		err := rules.Validate(id, opts...)
		assert.NoError(t, err)
	})

	t.Run("invalid SIRET - too short", func(t *testing.T) {
		id := &org.Identity{
			Type: fr.IdentityTypeSIRET,
			Code: "123456789",
		}
		err := rules.Validate(id, opts...)
		assert.ErrorContains(t, err, "[GOBL-FR-ORG-IDENTITY-02] ($.code) identity code for type SIRET must be valid")
	})

	t.Run("invalid SIRET - too long", func(t *testing.T) {
		id := &org.Identity{
			Type: fr.IdentityTypeSIRET,
			Code: "123456789012345",
		}
		err := rules.Validate(id, opts...)
		assert.ErrorContains(t, err, "[GOBL-FR-ORG-IDENTITY-02]")
	})

	t.Run("invalid SIRET - non-numeric", func(t *testing.T) {
		id := &org.Identity{
			Type: fr.IdentityTypeSIRET,
			Code: "1234567890123A",
		}
		err := rules.Validate(id, opts...)
		assert.ErrorContains(t, err, "[GOBL-FR-ORG-IDENTITY-02]")
	})

	t.Run("other identity types are not validated", func(t *testing.T) {
		id := &org.Identity{
			Type: "OTHER",
			Code: "anything",
		}
		err := rules.Validate(id, opts...)
		assert.NoError(t, err)
	})

	t.Run("nil identity", func(t *testing.T) {
		var id *org.Identity
		err := rules.Validate(id, opts...)
		assert.NoError(t, err)
	})
}
