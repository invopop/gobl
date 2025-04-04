package sdi_test

import (
	"testing"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTest(t *testing.T) {
	ad := tax.AddonForKey(sdi.V1)

	t.Run("nil address", func(t *testing.T) {
		var addr *org.Address
		assert.NotPanics(t, func() {
			ad.Normalizer(addr)
		})
	})
	t.Run("normalize short code", func(t *testing.T) {
		addr := &org.Address{
			Country: l10n.IT.ISO(),
			Code:    cbc.Code("123"),
		}
		ad.Normalizer(addr)
		assert.Equal(t, cbc.Code("00123"), addr.Code)
	})
	t.Run("missing code", func(t *testing.T) {
		addr := &org.Address{
			Country: l10n.IT.ISO(),
			Code:    cbc.Code(""),
		}
		ad.Normalizer(addr)
		assert.Equal(t, cbc.Code(""), addr.Code)
	})
	t.Run("ignore invalid code", func(t *testing.T) {
		addr := &org.Address{
			Country: l10n.IT.ISO(),
			Code:    cbc.Code("1A3"),
		}
		ad.Normalizer(addr)
		assert.Equal(t, cbc.Code("1A3"), addr.Code)
	})
	t.Run("ignore invalid code", func(t *testing.T) {
		addr := &org.Address{
			Country: l10n.IT.ISO(),
			Code:    cbc.Code("1A3"),
		}
		ad.Normalizer(addr)
		assert.Equal(t, cbc.Code("1A3"), addr.Code)
	})
	t.Run("ignore other countries", func(t *testing.T) {
		addr := &org.Address{
			Country: l10n.ES.ISO(),
			Code:    cbc.Code("1A3"),
		}
		ad.Normalizer(addr)
		assert.Equal(t, cbc.Code("1A3"), addr.Code)
	})

}
