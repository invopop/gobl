package in_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
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
		Tax: &bill.Tax{
			Ext: tax.Extensions{
				"in-supply-place": "Ciudad prueba",
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
					Price: num.MakeAmount(10000, 2),
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
}

func TestInvoiceValidation(t *testing.T) {
	inv := validInvoice()
	require.NoError(t, inv.Calculate())
	assert.NoError(t, inv.Validate())

	inv = validInvoice()
	inv.Supplier = nil
	require.NoError(t, inv.Calculate())
	assert.ErrorContains(t, inv.Validate(), "supplier: cannot be blank.")

}
