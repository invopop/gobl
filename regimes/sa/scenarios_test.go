package sa_test

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

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Code:     "INV-001",
		Currency: "SAR",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "SA",
				Code:    "300075588700003",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
		},
		IssueDate: cal.MakeDate(2024, 1, 15),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(5, 0),
				Item: &org.Item{
					Name:  "Service Item",
					Price: num.NewAmount(200, 0),
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

func testInvoiceSimplified(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Code:     "INV-002",
		Currency: "SAR",
		Tags:     tax.WithTags(tax.TagSimplified),
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "SA",
				Code:    "300075588700003",
			},
		},
		IssueDate: cal.MakeDate(2024, 1, 15),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(3, 0),
				Item: &org.Item{
					Name:  "Product Item",
					Price: num.NewAmount(100, 0),
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

func TestInvoiceScenarios(t *testing.T) {
	// Standard invoice with reverse charge tag should produce a note
	i := testInvoiceStandard(t)
	i.Tags = tax.WithTags(tax.TagReverseCharge)
	require.NoError(t, i.Calculate())
	require.NoError(t, i.Validate())
	require.Len(t, i.Notes, 1)
	assert.Equal(t, tax.TagReverseCharge, i.Notes[0].Src)
	assert.Contains(t, i.Notes[0].Text, "Reverse Charge")

	// Simplified invoice should produce a note and validate without customer
	i = testInvoiceSimplified(t)
	require.NoError(t, i.Calculate())
	require.NoError(t, i.Validate())
	require.Len(t, i.Notes, 1)
	assert.Equal(t, tax.TagSimplified, i.Notes[0].Src)
	assert.Contains(t, i.Notes[0].Text, "Simplified Tax Invoice")
}
