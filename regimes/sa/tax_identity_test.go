package sa_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	_ "github.com/invopop/gobl/regimes/sa"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityRules(t *testing.T) {
	validate := func(code cbc.Code) error {
		return rules.Validate(&tax.Identity{Country: "SA", Code: code})
	}

	t.Run("valid 15-digit starting and ending with 3", func(t *testing.T) {
		assert.NoError(t, validate("312345678912343"))
	})

	t.Run("valid another pattern", func(t *testing.T) {
		assert.NoError(t, validate("399999999900003"))
	})

	t.Run("valid all zeros in middle", func(t *testing.T) {
		assert.NoError(t, validate("300000000000003"))
	})

	t.Run("empty code is allowed", func(t *testing.T) {
		assert.NoError(t, validate(""))
	})

	t.Run("does not start with 3", func(t *testing.T) {
		err := validate("212345678912343")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})

	t.Run("does not end with 3", func(t *testing.T) {
		err := validate("312345678912341")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})

	t.Run("too short", func(t *testing.T) {
		err := validate("31234567893")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})

	t.Run("too long", func(t *testing.T) {
		err := validate("3123456789123430")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})

	t.Run("contains letters", func(t *testing.T) {
		err := validate("31234567891234A")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})

	t.Run("only digits but wrong prefix and suffix", func(t *testing.T) {
		err := validate("112345678912341")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})

	t.Run("14 digits", func(t *testing.T) {
		err := validate("31234567891233")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})

	t.Run("16 digits", func(t *testing.T) {
		err := validate("3123456789123433")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})

	t.Run("special characters", func(t *testing.T) {
		err := validate("3-1234567891234-3")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})
}

func TestTaxIdentityNormalization(t *testing.T) {
	t.Run("normalizes identity", func(t *testing.T) {
		tID := &tax.Identity{Country: "SA", Code: "312345678912343"}
		rd := tax.RegimeDefFor("SA")
		assert.NotNil(t, rd)
		rd.NormalizeObject(tID)
		assert.Equal(t, cbc.Code("312345678912343"), tID.Code)
	})
}
