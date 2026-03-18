package in_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	i := &bill.Invoice{
		Series:   "TEST",
		Code:     "0002",
		Currency: "INR",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "IN",
				Code:    "27AAPFU0939F1ZV",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "IN",
				Code:    "27AAPFU0939F1ZV",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Development services",
					Price: num.NewAmount(10000, 2),
					Unit:  org.UnitPackage,
					Identities: []*org.Identity{
						{
							Type: "HSN",
							Code: "12345678",
						},
					},
				},
				Taxes: tax.Set{
					{
						Category: "CGST",
						Percent:  num.NewPercentage(9, 0),
					},
					{
						Category: "SGST",
						Percent:  num.NewPercentage(9, 0),
					},
				},
			},
		},
	}
	return i
}

func TestInvoiceDocumentScenarios(t *testing.T) {
	inv := testInvoiceStandard(t)
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))

	inv = testInvoiceStandard(t)
	inv.SetTags(tax.TagSimplified)
	inv.Customer = nil
	require.NoError(t, inv.Calculate())
	assert.Len(t, inv.Notes, 1)
	assert.Equal(t, inv.Notes[0].Src, tax.TagSimplified)
	assert.Equal(t, inv.Notes[0].Text, "Simplified Tax Invoice")
}
