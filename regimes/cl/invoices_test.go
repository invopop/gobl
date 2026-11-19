package cl_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	_ "github.com/invopop/gobl/regimes/cl"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Series: "TEST",
		Code:   "0001",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "CL",
				Code:    "760864285",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "CL",
				Code:    "187654327",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.NewAmount(10000, 0),
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

func TestInvoiceValidation(t *testing.T) {
	inv := validInvoice()
	require.NoError(t, inv.Calculate())
	assert.NoError(t, rules.Validate(inv))

	inv = validInvoice()
	inv.Customer.TaxID.Code = "invalid"
	require.NoError(t, inv.Calculate())
	assert.ErrorContains(t, rules.Validate(inv), "invalid Chilean RUT")
}
