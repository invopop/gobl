package es_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validTicketBAIInvoice() *bill.Invoice {
	return &bill.Invoice{
		Code: "123",
		Tax: &bill.Tax{
			Tags: []cbc.Key{es.TagTicketBAI},
		},
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "NL",
				Code:    "000099995B57",
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
						Rate:     "exempt",
						Ext: tax.Extensions{
							es.ExtKeyTBAIExemption: "E1",
						},
					},
				},
			},
		},
	}
}

func TestBasqueLineValidation(t *testing.T) {
	inv := validTicketBAIInvoice()
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	inv.Lines[0].Taxes[0].Ext[es.ExtKeyTBAIProduct] = "services"
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	inv.Lines[0].Taxes[0].Ext = nil
	assertValidationError(t, inv, "es-tbai-exemption: require")
}

func assertValidationError(t *testing.T, inv *bill.Invoice, expected string) {
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), expected)
}
