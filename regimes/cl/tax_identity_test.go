package cl_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	_ "github.com/invopop/gobl/regimes/cl"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityRules(t *testing.T) {
	validate := func(code cbc.Code) error {
		return rules.Validate(&tax.Identity{Country: "CL", Code: code})
	}

	t.Run("valid RUT, valid DV", func(t *testing.T) {
		assert.NoError(t, validate("123456785"))
	})

	t.Run("valid RUT, DV = K", func(t *testing.T) {
		assert.NoError(t, validate("1000005K"))
	})

	t.Run("invalid RUT", func(t *testing.T) {
		assert.Error(t, validate("123456781"))
	})

	t.Run("too short", func(t *testing.T) {
		err := validate("5")
		assert.ErrorContains(t, err, "CL-TAX-IDENTITY-01")
	})

	t.Run("too long", func(t *testing.T) {
		err := validate("1234567890")
		assert.ErrorContains(t, err, "CL-TAX-IDENTITY-01")
	})

	t.Run("contains letters that is not K", func(t *testing.T) {
		err := validate("1234567A")
		assert.ErrorContains(t, err, "CL-TAX-IDENTITY-01")
	})
}

func TestTaxIdentityNormalization(t *testing.T) {
	t.Run("nil identity is safe", func(t *testing.T) {
		var tID *tax.Identity
		assert.NotPanics(t, func() { norm.Normalize(tID) })
	})

	t.Run("normalizes identity", func(t *testing.T) {
		tID := &tax.Identity{Country: "CL", Code: "49091850"}
		norm.Normalize(tID)
		assert.Equal(t, cbc.Code("49091850"), tID.Code)
	})

	t.Run("empty code is left untouched", func(t *testing.T) {
		tID := &tax.Identity{Country: "CL", Code: ""}
		norm.Normalize(tID)
		assert.Equal(t, cbc.Code(""), tID.Code)
	})

	t.Run("strips whitespace and country prefix", func(t *testing.T) {
		tID := &tax.Identity{Country: "CL", Code: " CL 12.345.678-5 "}
		norm.Normalize(tID)
		assert.Equal(t, cbc.Code("123456785"), tID.Code)
	})

	t.Run("uppercases the K check digit", func(t *testing.T) {
		tID := &tax.Identity{Country: "CL", Code: "12345678k"}
		norm.Normalize(tID)
		assert.Equal(t, cbc.Code("12345678K"), tID.Code)
	})
}
