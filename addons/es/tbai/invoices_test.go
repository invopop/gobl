package tbai_test

import (
	"testing"

	"github.com/invopop/gobl/addons/es/tbai"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/require"
)

func validTicketBAIInvoice() *bill.Invoice {
	return &bill.Invoice{
		Addons: tax.WithAddons(tbai.V1),
		Series: "ABC",
		Code:   "123",
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
							tbai.ExtKeyExemption: "E1",
						},
					},
				},
			},
		},
		Notes: []*cbc.Note{
			{
				Key:  cbc.NoteKeyGeneral,
				Text: "This is a test invoice",
			},
		},
	}
}

func TestInvoiceValidation(t *testing.T) {
	inv := validTicketBAIInvoice()
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	inv.Lines[0].Taxes[0].Ext[tbai.ExtKeyProduct] = "services"
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	inv.Lines[0].Taxes[0].Ext = nil
	assertValidationError(t, inv, "es-tbai-exemption: required")

	t.Run("without series", func(t *testing.T) {
		inv := validTicketBAIInvoice()
		inv.Series = ""
		assertValidationError(t, inv, "series: cannot be blank")
	})

	t.Run("without notes", func(t *testing.T) {
		inv := validTicketBAIInvoice()
		inv.Notes = nil
		assertValidationError(t, inv, "notes: with key 'general' missing")
	})
}

func assertValidationError(t *testing.T, inv *bill.Invoice, expected string) {
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.ErrorContains(t, err, expected)
}
