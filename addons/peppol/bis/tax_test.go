package bis

import (
	"testing"

	"github.com/invopop/gobl/catalogues/cef"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func combo(vatex, cat cbc.Code) *tax.Combo {
	c := &tax.Combo{Ext: tax.Extensions{}}
	if vatex != "" {
		c.Ext[cef.ExtKeyVATEX] = vatex
	}
	if cat != "" {
		c.Ext[untdid.ExtKeyTaxCategory] = cat
	}
	return c
}

func TestVatexCategoryCoherent(t *testing.T) {
	t.Run("nil passes", func(t *testing.T) {
		assert.True(t, vatexCategoryCoherent(nil))
	})
	t.Run("no vatex passes", func(t *testing.T) {
		assert.True(t, vatexCategoryCoherent(combo("", "S")))
	})
	t.Run("unknown vatex passes", func(t *testing.T) {
		assert.True(t, vatexCategoryCoherent(combo("VATEX-OTHER", "S")))
	})
	t.Run("vatex matches required category", func(t *testing.T) {
		cases := map[cbc.Code]cbc.Code{
			"VATEX-EU-G":  "G",
			"VATEX-EU-O":  "O",
			"VATEX-EU-IC": "K",
			"VATEX-EU-AE": "AE",
			"VATEX-EU-D":  "E",
			"VATEX-EU-F":  "E",
			"VATEX-EU-I":  "E",
			"VATEX-EU-J":  "E",
		}
		for vatex, want := range cases {
			assert.True(t, vatexCategoryCoherent(combo(vatex, want)), string(vatex))
		}
	})
	t.Run("vatex with wrong category fails", func(t *testing.T) {
		assert.False(t, vatexCategoryCoherent(combo("VATEX-EU-G", "S")))
		assert.False(t, vatexCategoryCoherent(combo("VATEX-EU-AE", "K")))
	})
}
