package pay_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestMeansKey(t *testing.T) {
	i := new(pay.Instructions)
	i.Key = cbc.Key("foo")
	err := rules.Validate(i)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key must be valid")

	i.Key = pay.MeansKeyCard
	err = rules.Validate(i)
	assert.NoError(t, err)

	// Key with an extension
	i.Key = pay.MeansKeyCard.With("visa")
	err = rules.Validate(i)
	assert.NoError(t, err)
}

func TestLookupMeansCode(t *testing.T) {
	m := map[cbc.Key]cbc.Code{
		pay.MeansKeyCard:                              "04",
		pay.MeansKeyCard.With(pay.MeansKeyDebit):      "28",
		pay.MeansKeyCash:                              "01",
	}

	t.Run("exact match wins", func(t *testing.T) {
		assert.Equal(t, cbc.Code("28"), pay.LookupMeansCode(m, pay.MeansKeyCard.With(pay.MeansKeyDebit)))
	})

	t.Run("falls back through Pop", func(t *testing.T) {
		// `card+credit` is not registered; should fall back to `card`.
		assert.Equal(t, cbc.Code("04"), pay.LookupMeansCode(m, pay.MeansKeyCard.With(pay.MeansKeyCredit)))
	})

	t.Run("bare key match", func(t *testing.T) {
		assert.Equal(t, cbc.Code("01"), pay.LookupMeansCode(m, pay.MeansKeyCash))
	})

	t.Run("missing key returns empty", func(t *testing.T) {
		assert.Equal(t, cbc.Code(""), pay.LookupMeansCode(m, pay.MeansKeyCheque))
	})

	t.Run("empty key returns empty", func(t *testing.T) {
		assert.Equal(t, cbc.Code(""), pay.LookupMeansCode(m, cbc.KeyEmpty))
	})
}
