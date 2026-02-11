package fr_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/stretchr/testify/assert"
)

func TestValidateIdentity(t *testing.T) {
	t.Run("valid SIREN", func(t *testing.T) {
		id := &org.Identity{
			Type: fr.IdentityTypeSIREN,
			Code: "123456789",
		}
		err := fr.Validate(id)
		assert.NoError(t, err)
	})

	t.Run("invalid SIREN - too short", func(t *testing.T) {
		id := &org.Identity{
			Type: fr.IdentityTypeSIREN,
			Code: "12345",
		}
		err := fr.Validate(id)
		assert.ErrorContains(t, err, "must be exactly 9 digits")
	})

	t.Run("invalid SIREN - too long", func(t *testing.T) {
		id := &org.Identity{
			Type: fr.IdentityTypeSIREN,
			Code: "1234567890",
		}
		err := fr.Validate(id)
		assert.ErrorContains(t, err, "must be exactly 9 digits")
	})

	t.Run("invalid SIREN - non-numeric", func(t *testing.T) {
		id := &org.Identity{
			Type: fr.IdentityTypeSIREN,
			Code: "12345678A",
		}
		err := fr.Validate(id)
		assert.ErrorContains(t, err, "must be exactly 9 digits")
	})

	t.Run("valid SIRET", func(t *testing.T) {
		id := &org.Identity{
			Type: fr.IdentityTypeSIRET,
			Code: "12345678901234",
		}
		err := fr.Validate(id)
		assert.NoError(t, err)
	})

	t.Run("invalid SIRET - too short", func(t *testing.T) {
		id := &org.Identity{
			Type: fr.IdentityTypeSIRET,
			Code: "123456789",
		}
		err := fr.Validate(id)
		assert.ErrorContains(t, err, "must be exactly 14 digits")
	})

	t.Run("invalid SIRET - too long", func(t *testing.T) {
		id := &org.Identity{
			Type: fr.IdentityTypeSIRET,
			Code: "123456789012345",
		}
		err := fr.Validate(id)
		assert.ErrorContains(t, err, "must be exactly 14 digits")
	})

	t.Run("invalid SIRET - non-numeric", func(t *testing.T) {
		id := &org.Identity{
			Type: fr.IdentityTypeSIRET,
			Code: "1234567890123A",
		}
		err := fr.Validate(id)
		assert.ErrorContains(t, err, "must be exactly 14 digits")
	})

	t.Run("other identity types are not validated", func(t *testing.T) {
		id := &org.Identity{
			Type: "OTHER",
			Code: "anything",
		}
		err := fr.Validate(id)
		assert.NoError(t, err)
	})

	t.Run("nil identity", func(t *testing.T) {
		err := fr.Validate(nil)
		assert.NoError(t, err)
	})
}
