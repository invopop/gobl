package co_test

import (
	"testing"

	"github.com/invopop/gobl/addons/co/dian"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/co"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeParty(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var p *org.Party
		assert.NotPanics(t, func() {
			co.Normalize(p)
		})
	})

	t.Run("basic", func(t *testing.T) {
		p := &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "CO",
				Code:    "412615332",
				Zone:    "11001",
			},
		}
		co.Normalize(p)
		assert.Empty(t, p.TaxID.Zone) //nolint:staticcheck
		assert.Equal(t, p.Ext[dian.ExtKeyMunicipality].String(), "11001")
	})
}
