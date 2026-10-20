package il_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoiceReverseCharge(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Code:     "123TEST",
		Currency: "ILS",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "IL",
				Code:    "123456789",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "IL",
				Code:    "987654321",
			},
		},
		IssueDate: cal.MakeDate(2025, 1, 15),
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
		Currency: "ILS",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "IL",
				Code:    "123456789",
			},
		},
		IssueDate: cal.MakeDate(2025, 1, 15),
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
	require.NoError(t, i.Validate())
	assert.Len(t, i.Notes, 1)
	assert.Equal(t, i.Notes[0].Src, tax.TagReverseCharge)
	assert.Equal(t, i.Notes[0].Text, "Reverse Charge")

	i = testInvoiceSimplified(t)
	require.NoError(t, i.Calculate())
	require.NoError(t, i.Validate())
	assert.Len(t, i.Notes, 1)
	assert.Equal(t, i.Notes[0].Src, tax.TagSimplified)
	assert.Equal(t, i.Notes[0].Text, "Simplified Tax Invoice")
}
