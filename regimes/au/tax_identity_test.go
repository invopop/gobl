package au_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	_ "github.com/invopop/gobl/regimes/au"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityRules(t *testing.T) {
	validate := func(code cbc.Code) error {
		return rules.Validate(&tax.Identity{Country: "AU", Code: code})
	}

	t.Run("valid ABN (ATO)", func(t *testing.T) {
		assert.NoError(t, validate("51824753556"))
	})

	t.Run("valid ABN", func(t *testing.T) {
		assert.NoError(t, validate("83914571673"))
	})

	t.Run("empty code is allowed", func(t *testing.T) {
		assert.NoError(t, validate(""))
	})

	t.Run("invalid check digits", func(t *testing.T) {
		err := validate("51824753557")
		assert.ErrorContains(t, err, "IDENTITY-02")
	})

	t.Run("too short", func(t *testing.T) {
		err := validate("5182475355")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})

	t.Run("too long", func(t *testing.T) {
		err := validate("518247535560")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})

	t.Run("contains letters", func(t *testing.T) {
		err := validate("5182475355A")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})

	t.Run("special characters", func(t *testing.T) {
		err := validate("51-824753-556")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})
}

func TestTaxIdentityNormalization(t *testing.T) {
	t.Run("nil identity is safe", func(t *testing.T) {
		var tID *tax.Identity
		assert.NotPanics(t, func() { norm.Normalize(tID) })
	})

	t.Run("normalizes identity", func(t *testing.T) {
		tID := &tax.Identity{Country: "AU", Code: "51824753556"}
		norm.Normalize(tID)
		assert.Equal(t, cbc.Code("51824753556"), tID.Code)
	})

	t.Run("empty code is left untouched", func(t *testing.T) {
		tID := &tax.Identity{Country: "AU", Code: ""}
		norm.Normalize(tID)
		assert.Equal(t, cbc.Code(""), tID.Code)
	})

	t.Run("strips whitespace and country prefix", func(t *testing.T) {
		tID := &tax.Identity{Country: "AU", Code: " AU 51 824 753 556 "}
		norm.Normalize(tID)
		assert.Equal(t, cbc.Code("51824753556"), tID.Code)
	})
}
