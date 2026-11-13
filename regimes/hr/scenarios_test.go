package hr_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/hr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoice(t *testing.T, tags ...cbc.Key) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Code:     "1-1-1",
		Currency: "EUR",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "HR",
				Code:    "12345678903",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "HR",
				Code:    "11111111119",
			},
		},
		IssueDate: cal.MakeDate(2026, 1, 15),
		Tags:      tax.WithTags(tags...),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Service",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     tax.RateGeneral,
					},
				},
			},
		},
	}
}

func TestInvoiceTags(t *testing.T) {
	regime := hr.New()
	require.Len(t, regime.Tags, 1)
	assert.Equal(t, bill.ShortSchemaInvoice, regime.Tags[0].Schema)
	require.Len(t, regime.Tags[0].List, 1)
	assert.Equal(t, hr.TagCashBasis, regime.Tags[0].List[0].Key)
}

func TestInvoiceScenarios(t *testing.T) {
	t.Run("cash basis", func(t *testing.T) {
		inv := testInvoice(t, hr.TagCashBasis)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
		require.NotNil(t, inv.Tax)
		require.Len(t, inv.Tax.Notes, 1)
		assert.Equal(t, tax.CategoryVAT, inv.Tax.Notes[0].Category)
		assert.Equal(t, hr.TagCashBasis, inv.Tax.Notes[0].Key)
		assert.Equal(t, "Obračun prema naplaćenim naknadama", inv.Tax.Notes[0].Text)
	})

	t.Run("reverse charge", func(t *testing.T) {
		inv := testInvoice(t, tax.TagReverseCharge)
		require.NoError(t, inv.Calculate())
		require.NotNil(t, inv.Tax)
		require.Len(t, inv.Tax.Notes, 1)
		assert.Equal(t, tax.KeyReverseCharge, inv.Tax.Notes[0].Key)
		assert.Equal(t, "Reverse charge: Customer to account for VAT to the relevant tax authority.", inv.Tax.Notes[0].Text)
	})

	t.Run("no tags", func(t *testing.T) {
		inv := testInvoice(t)
		require.NoError(t, inv.Calculate())
		if inv.Tax != nil {
			assert.Empty(t, inv.Tax.Notes)
		}
	})
}
