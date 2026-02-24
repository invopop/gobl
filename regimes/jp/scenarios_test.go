package jp_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoiceWithTags(t *testing.T, tags ...cbc.Key) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Regime: tax.WithRegime("JP"),
		Series: "TEST",
		Code:   "0001",
		Tags:   tax.WithTags(tags...),
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "JP",
				Code:    "5010401067252",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "JP",
				Code:    "1130001011420",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 0),
					Unit:  org.UnitPackage,
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "general",
					},
				},
			},
		},
	}
}

func TestScenarios(t *testing.T) {
	t.Run("reverse charge", func(t *testing.T) {
		inv := testInvoiceWithTags(t, tax.TagReverseCharge)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		assert.Len(t, inv.Notes, 1)
		assert.Equal(t, tax.TagReverseCharge, inv.Notes[0].Src)
		assert.Equal(t, "Reverse Charge: The recipient is liable for the consumption tax on this transaction", inv.Notes[0].Text)
	})

	t.Run("simplified invoice", func(t *testing.T) {
		inv := testInvoiceWithTags(t, tax.TagSimplified)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		assert.Len(t, inv.Notes, 1)
		assert.Equal(t, tax.TagSimplified, inv.Notes[0].Src)
		assert.Equal(t, "Simplified Invoice", inv.Notes[0].Text)
	})

	t.Run("no tags - no notes", func(t *testing.T) {
		inv := testInvoiceWithTags(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		assert.Empty(t, inv.Notes)
	})
}
