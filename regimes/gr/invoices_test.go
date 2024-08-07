package gr_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Series: "TEST",
		Code:   "0002",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "EL",
				Code:    "728089281",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "EL",
				Code:    "841442160",
			},
			Addresses: []*org.Address{
				{
					Locality: "Athens",
					Code:     "11528",
				},
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.MakeAmount(10000, 2),
					Unit:  org.UnitPackage,
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
			},
		},
		Payment: &bill.Payment{
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCreditTransfer,
			},
		},
	}
}

func TestInvoiceValidation(t *testing.T) {
	inv := validInvoice()
	require.NoError(t, inv.Calculate())
	assert.NoError(t, inv.Validate())

	// Make it invalid
	inv.Series = ""
	inv.Supplier.TaxID.Code = ""
	inv.Customer.Addresses = nil
	inv.Lines[0].Quantity = num.MakeAmount(0, 0)
	inv.Payment.Instructions.Key = "debit-transfer"

	require.NoError(t, inv.Calculate())

	err := inv.Validate()
	assert.ErrorContains(t, err, "series: cannot be blank")
	assert.ErrorContains(t, err, "supplier: (tax_id: (code: cannot be blank")
	assert.ErrorContains(t, err, "customer: (addresses: cannot be blank")
	assert.ErrorContains(t, err, "lines: (0: (total: must be greater than 0")
	assert.ErrorContains(t, err, "payment: (instructions: (key: must be a valid value")
}

func TestSimplifiedInvoiceValidation(t *testing.T) {
	inv := validInvoice()
	inv.Tax = &bill.Tax{
		Tags: []cbc.Key{tax.TagSimplified},
	}
	inv.Customer.TaxID = nil
	inv.Customer.Addresses = nil

	require.NoError(t, inv.Calculate())
	assert.NoError(t, inv.Validate())
}
