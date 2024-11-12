package ae_test

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
		Currency: "AED",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "AE",
				Code:    "123456789012345",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "AE",
				Code:    "123456789012346",
			},
		},
		IssueDate: cal.MakeDate(2023, 1, 15),
		Tags:      tax.WithTags(tax.TagReverseCharge),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(5, 0),
				Item: &org.Item{
					Name:  "Service Item",
					Price: num.MakeAmount(5000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
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
				Code:    "123456789012345",
			},
		},
		IssueDate: cal.MakeDate(2023, 1, 15),
		Tags:      tax.WithTags(tax.TagSimplified),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(3, 0),
				Item: &org.Item{
					Name:  "Product Item",
					Price: num.MakeAmount(2000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
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
	assert.Equal(t, i.Notes[0].Text, "Reverse Charge / التحويل العكسي")

	i = testInvoiceSimplified(t)
	require.NoError(t, i.Calculate())
	require.NoError(t, i.Validate())
	assert.Len(t, i.Notes, 1)
	assert.Equal(t, i.Notes[0].Src, tax.TagSimplified)
	assert.Equal(t, i.Notes[0].Text, "Simplified Tax Invoice / فاتورة ضريبية مبسطة")
}
