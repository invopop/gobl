package nz_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	_ "github.com/invopop/gobl/regimes/nz"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityRules(t *testing.T) {
	validate := func(code cbc.Code) error {
		return rules.Validate(&tax.Identity{Country: "NZ", Code: code})
	}

	t.Run("valid 8-digit IRD", func(t *testing.T) {
		assert.NoError(t, validate("49091850"))
	})

	t.Run("valid 9-digit IRD using secondary weights", func(t *testing.T) {
		assert.NoError(t, validate("136410132"))
	})

	t.Run("empty code is allowed", func(t *testing.T) {
		assert.NoError(t, validate(""))
	})

	t.Run("invalid check digit", func(t *testing.T) {
		err := validate("49091851")
		assert.ErrorContains(t, err, "IDENTITY-02")
	})

	t.Run("below valid range", func(t *testing.T) {
		err := validate("09999999")
		assert.ErrorContains(t, err, "IDENTITY-02")
	})

	t.Run("too short", func(t *testing.T) {
		err := validate("4909185")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})

	t.Run("too long", func(t *testing.T) {
		err := validate("1364101320")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})

	t.Run("contains letters", func(t *testing.T) {
		err := validate("4909185A")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})

	t.Run("special characters", func(t *testing.T) {
		err := validate("49-091-850")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})
}

func TestTaxIdentityNormalization(t *testing.T) {
	t.Run("nil identity is safe", func(t *testing.T) {
		var tID *tax.Identity
		assert.NotPanics(t, func() { norm.Normalize(tID) })
	})

	t.Run("normalizes identity", func(t *testing.T) {
		tID := &tax.Identity{Country: "NZ", Code: "49091850"}
		norm.Normalize(tID)
		assert.Equal(t, cbc.Code("49091850"), tID.Code)
	})

	t.Run("empty code is left untouched", func(t *testing.T) {
		tID := &tax.Identity{Country: "NZ", Code: ""}
		norm.Normalize(tID)
		assert.Equal(t, cbc.Code(""), tID.Code)
	})

	t.Run("strips whitespace and country prefix", func(t *testing.T) {
		tID := &tax.Identity{Country: "NZ", Code: " NZ 49-091-850 "}
		norm.Normalize(tID)
		assert.Equal(t, cbc.Code("49091850"), tID.Code)
	})
}
