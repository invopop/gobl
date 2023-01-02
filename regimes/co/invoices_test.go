package co_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func baseInvoice() *bill.Invoice {
	inv := &bill.Invoice{
		Currency:  currency.COP,
		Code:      "TEST",
		IssueDate: cal.MakeDate(2022, 12, 27),
		Supplier: &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: l10n.CO,
				Code:    "412615332",
				Zone:    "11001",
			},
			Addresses: []*org.Address{
				{
					Locality: "Foo",
				},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: l10n.CO,
				Code:    "124499654",
				Zone:    "08638",
			},
			Addresses: []*org.Address{
				{
					Locality: "Foo",
				},
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 3),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.MakeAmount(1000, 3),
				},
			},
		},
	}
	return inv
}

func TestBasicInvoiceValidation(t *testing.T) {
	inv := baseInvoice()
	err := inv.Calculate()
	require.NoError(t, err)
	err = inv.Validate()
	assert.NoError(t, err)
	assert.Equal(t, inv.Supplier.Addresses[0].Locality, "BOGOTÁ, D.C.")
	assert.Equal(t, inv.Supplier.Addresses[0].Region, "Bogotá")
	assert.Equal(t, inv.Customer.Addresses[0].Locality, "SABANALARGA")
	assert.Equal(t, inv.Customer.Addresses[0].Region, "Atlántico")

	inv.Supplier.TaxID.Zone = ""
	err = inv.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "zone: cannot be blank")

	inv.Supplier.TaxID.Zone = "1100X"
	err = inv.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "zone: must be a valid value")
}
