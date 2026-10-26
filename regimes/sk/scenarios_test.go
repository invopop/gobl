package sk_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/sk"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func noteTexts(t *testing.T, tags ...cbc.Key) []string {
	t.Helper()
	inv := validInvoice()
	inv.SetTags(tags...)
	require.NoError(t, inv.Calculate())
	var out []string
	if inv.Tax != nil {
		for _, n := range inv.Tax.Notes {
			out = append(out, n.Text)
		}
	}
	return out
}

func TestInvoiceScenarios(t *testing.T) {
	tests := []struct {
		name string
		tag  cbc.Key
		text string
	}{
		{"reverse charge", tax.TagReverseCharge, "Prenesenie daňovej povinnosti"},
		{"self-billed", tax.TagSelfBilled, "Vyhotovenie faktúry odberateľom"},
		{"travel agency", sk.TagTravelAgency, "Úprava zdaňovania prirážky - cestovné kancelárie"},
		{"second-hand goods", sk.TagSecondHandGoods, "Úprava zdaňovania prirážky - použitý tovar"},
		{"works of art", sk.TagArt, "Úprava zdaňovania prirážky - umelecké diela"},
		{"antiques", sk.TagAntiques, "Úprava zdaňovania prirážky - zberateľské predmety a starožitnosti"},
		{"cash basis", sk.TagCashBasis, "Daň sa uplatňuje na základe prijatia platby"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, []string{tt.text}, noteTexts(t, tt.tag))
		})
	}

	t.Run("no tag adds no note", func(t *testing.T) {
		assert.Empty(t, noteTexts(t))
	})
}
