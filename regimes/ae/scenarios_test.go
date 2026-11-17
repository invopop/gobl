package ae_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/ae"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoiceReverseCharge(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Code:     "123TEST",
		Currency: "AED",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "AE",
				Code:    "123456789700003",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "AE",
				Code:    "187654321400003",
			},
		},
		IssueDate: cal.MakeDate(2023, 1, 15),
		Tags:      tax.WithTags(tax.TagReverseCharge),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(5, 0),
				Item: &org.Item{
					Name:  "Service Item",
					Price: num.NewAmount(5000, 2),
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

func testInvoiceSimplified(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Code:     "123TEST",
		Currency: "AED",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "AE",
				Code:    "123456789700003",
			},
		},
		IssueDate: cal.MakeDate(2023, 1, 15),
		Tags:      tax.WithTags(tax.TagSimplified),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(3, 0),
				Item: &org.Item{
					Name:  "Product Item",
					Price: num.NewAmount(2000, 2),
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

func TestInvoiceScenarios(t *testing.T) {
	i := testInvoiceReverseCharge(t)
	require.NoError(t, i.Calculate())
	require.NoError(t, rules.Validate(i))
	require.NotNil(t, i.Tax)
	assert.Len(t, i.Tax.Notes, 1)
	assert.Equal(t, tax.KeyReverseCharge, i.Tax.Notes[0].Key)
	assert.Equal(t, "Reverse Charge", i.Tax.Notes[0].Text)

	i = testInvoiceSimplified(t)
	require.NoError(t, i.Calculate())
	require.NoError(t, rules.Validate(i))
	require.NotNil(t, i.Tax)
	assert.Len(t, i.Tax.Notes, 1)
	assert.Equal(t, "Simplified Tax Invoice", i.Tax.Notes[0].Text)
}

func TestMarginSchemeScenario(t *testing.T) {
	marginNotes := func(inv *bill.Invoice) []*tax.Note {
		var notes []*tax.Note
		if inv.Tax == nil {
			return nil
		}
		for _, n := range inv.Tax.Notes {
			if n.Key == ae.TagMarginScheme {
				notes = append(notes, n)
			}
		}
		return notes
	}
	testInvoice := func(t *testing.T) *bill.Invoice {
		t.Helper()
		inv := testInvoiceSimplified(t)
		inv.Tags = tax.WithTags(ae.TagMarginScheme)
		inv.Lines[0].Taxes = tax.Set{{Category: tax.CategoryVAT, Key: tax.KeyExempt}}
		return inv
	}

	t.Run("adds the statement with the tag", func(t *testing.T) {
		inv := testInvoice(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
		notes := marginNotes(inv)
		require.Len(t, notes, 1)
		assert.Equal(t, tax.CategoryVAT, notes[0].Category)
		assert.Equal(t, "VAT charged under the Profit Margin Scheme", notes[0].Text)
		assert.Equal(t, "0.00", inv.Totals.Tax.String())
		assert.Equal(t, "60.00", inv.Totals.Payable.String())
	})

	t.Run("no statement without the tag", func(t *testing.T) {
		inv := testInvoice(t)
		inv.Tags = tax.Tags{}
		require.NoError(t, inv.Calculate())
		assert.Empty(t, marginNotes(inv))
	})

	t.Run("calculating twice keeps one statement", func(t *testing.T) {
		inv := testInvoice(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Calculate())
		assert.Len(t, marginNotes(inv), 1)
	})
}
