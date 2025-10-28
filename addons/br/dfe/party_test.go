package dfe_test

import (
	"testing"

	"github.com/invopop/gobl/addons/br/dfe"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeParty(t *testing.T) {
	addon := tax.AddonForKey(dfe.V1)

	t.Run("nil party", func(t *testing.T) {
		var party *org.Party
		addon.Normalizer(party)
		assert.Nil(t, party)
	})

	t.Run("migrates old addon extension keys", func(t *testing.T) {
		party := &org.Party{
			Ext: tax.Extensions{
				"br-nfse-fiscal-incentive": "1",
				"br-nfse-municipality":     "1234567890",
				"br-nfse-simples":          "2",
				"br-nfse-special-regime":   "3",
			},
		}
		addon.Normalizer(party)
		assert.Len(t, party.Ext, 4)
		assert.Equal(t, cbc.Code("1"), party.Ext["br-dfe-fiscal-incentive"])
		assert.Equal(t, cbc.Code("1234567890"), party.Ext["br-dfe-municipality"])
		assert.Equal(t, cbc.Code("2"), party.Ext["br-dfe-simples"])
		assert.Equal(t, cbc.Code("3"), party.Ext["br-dfe-special-regime"])
	})
}
