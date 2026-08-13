package sg_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/sg"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime: tax.WithRegime(sg.CountryCode),
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Code:    "M91234567X",
				Country: "SG",
			},
			Name: "Test Supplier",
			Addresses: []*org.Address{
				{
					Street:  "Test Street",
					Code:    "123456",
					Country: sg.CountryCode,
				},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			Addresses: []*org.Address{
				{
					Street:  "Test Street",
					Code:    "123456",
					Country: sg.CountryCode,
				},
			},
		},
		Code:     "0001",
		Currency: "SGD",
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryGST,
						Rate:     tax.RateGeneral,
					},
				},
			},
		},
	}
}

func TestValidInvoice(t *testing.T) {
	inv := validInvoice()
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
}

func TestValidInvoiceWithUEN(t *testing.T) {
	inv := validInvoice()
	inv.Supplier.TaxID = nil
	inv.Supplier.Identities = []*org.Identity{
		{
			Type: sg.IdentityTypeUEN,
			Code: "199912345A",
		},
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
}

func TestMissingSupplierTaxID(t *testing.T) {
	inv := validInvoice()
	inv.Supplier.TaxID = nil
	require.NoError(t, inv.Calculate())
	require.ErrorContains(t, rules.Validate(inv),
		"[GOBL-SG-BILL-INVOICE-01] ($.supplier) invoice supplier in Singapore must have a GST tax ID code or a UEN identity")
}

func TestMissingSupplierName(t *testing.T) {
	inv := validInvoice()
	inv.Supplier.Name = ""
	require.NoError(t, inv.Calculate())
	require.ErrorContains(t, rules.Validate(inv),
		"[GOBL-BILL-INVOICE-06] ($.supplier.name)")
}
