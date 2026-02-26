package ro_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceScenarios(t *testing.T) {
	tests := []struct {
		name string
		inv  *bill.Invoice
		src  cbc.Key
		text string
	}{
		{
			name: "reverse charge",
			inv:  reverseChargeInvoice(),
			src:  tax.TagReverseCharge,
			text: "Reverse Charge / Taxare inversă.",
		},
		{
			name: "simplified",
			inv:  simplifiedInvoice(),
			src:  tax.TagSimplified,
			text: "Factură simplificată / Simplified invoice.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, tt.inv.Calculate())
			require.NoError(t, tt.inv.Validate())
			assert.Len(t, tt.inv.Notes, 1)
			assert.Equal(t, tt.src, tt.inv.Notes[0].Src)
			assert.Equal(t, tt.text, tt.inv.Notes[0].Text)
		})
	}
}

func TestReverseChargeZeroVAT(t *testing.T) {
	inv := reverseChargeInvoice()
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	// Reverse charge must produce zero VAT
	require.NotNil(t, inv.Totals)
	assert.True(t, inv.Totals.Tax.IsZero(), "tax total should be zero for reverse charge")
	assert.True(t, inv.Totals.Payable.Equals(inv.Totals.Sum), "payable should equal sum (no VAT added)")

	// Verify the tax category shows zero amount
	require.NotNil(t, inv.Totals.Taxes)
	require.Len(t, inv.Totals.Taxes.Categories, 1)
	cat := inv.Totals.Taxes.Categories[0]
	assert.Equal(t, tax.CategoryVAT, cat.Code)
	require.Len(t, cat.Rates, 1)
	assert.True(t, cat.Rates[0].Amount.IsZero(), "VAT rate amount should be zero")
	assert.Equal(t, tax.KeyReverseCharge, cat.Rates[0].Key)
}

func TestCreditNoteScenario(t *testing.T) {
	t.Run("valid credit note with preceding", func(t *testing.T) {
		inv := creditNoteInvoice()
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("valid debit note with preceding", func(t *testing.T) {
		inv := creditNoteInvoice()
		inv.Type = bill.InvoiceTypeDebitNote
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})
}

func reverseChargeInvoice() *bill.Invoice {
	return &bill.Invoice{
		Series: "TEST",
		Code:   "0004",
		Tags:   tax.WithTags(tax.TagReverseCharge),
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "RO",
				Code:    "14399840",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "RO",
				Code:    "18547290",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Construction services",
					Price: num.NewAmount(10000, 2),
					Unit:  org.UnitPackage,
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Key:      tax.KeyReverseCharge,
					},
				},
			},
		},
	}
}

func creditNoteInvoice() *bill.Invoice {
	return &bill.Invoice{
		Series: "CN",
		Code:   "0005",
		Type:   bill.InvoiceTypeCreditNote,
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "RO",
				Code:    "14399840",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "RO",
				Code:    "18547290",
			},
		},
		Preceding: []*org.DocumentRef{
			{
				Type:      "standard",
				Series:    "SAMPLE",
				Code:      "001",
				IssueDate: cal.NewDate(2024, 2, 13),
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Development services",
					Price: num.NewAmount(50000, 2),
					Unit:  org.UnitHour,
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
