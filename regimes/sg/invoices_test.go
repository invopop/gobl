package sg_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/sg"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
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
					Country: l10n.SG.ISO(),
				},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			Addresses: []*org.Address{
				{
					Street:  "Test Street",
					Code:    "123456",
					Country: l10n.SG.ISO(),
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
	require.NoError(t, inv.Validate())
}

func TestValidInvoiceWithUEN(t *testing.T) {
	inv := validInvoice()
	inv.Supplier.Identities = []*org.Identity{
		{
			Type: sg.IdentityTypeUEN,
			Code: "199912345A",
		},
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestNilSupplier(t *testing.T) {
	inv := validInvoice()
	inv.Supplier = nil
	require.Error(t, inv.Validate())
}

func TestMissingSupplierTaxID(t *testing.T) {
	inv := validInvoice()
	inv.Supplier.TaxID = nil
	require.Error(t, inv.Validate())
}

func TestMissingSupplierName(t *testing.T) {
	inv := validInvoice()
	inv.Supplier.Name = ""
	require.Error(t, inv.Validate())
}
