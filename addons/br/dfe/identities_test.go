package dfe_test

import (
	"testing"

	"github.com/invopop/gobl/addons/br/dfe"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeIdentities(t *testing.T) {
	addon := tax.AddonForKey(dfe.V1)

	t.Run("nil identity", func(t *testing.T) {
		var ident *org.Identity
		addon.Normalizer(ident)
		assert.Nil(t, ident)
	})

	t.Run("migrates old key br-nfse-municipal-reg", func(t *testing.T) {
		ident := &org.Identity{
			Key:  "br-nfse-municipal-reg",
			Code: "1234567890",
		}
		addon.Normalizer(ident)
		assert.Equal(t, cbc.Key("br-dfe-municipal-reg"), ident.Key)
		assert.Equal(t, cbc.Code("1234567890"), ident.Code)
	})

	t.Run("migrates old key br-nfse-national-reg", func(t *testing.T) {
		ident := &org.Identity{
			Key:  "br-nfse-national-reg",
			Code: "1234567890",
		}
		addon.Normalizer(ident)
		assert.Equal(t, cbc.Key("br-dfe-state-reg"), ident.Key)
		assert.Equal(t, cbc.Code("1234567890"), ident.Code)
	})
}
