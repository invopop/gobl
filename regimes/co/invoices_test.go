package co_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/co"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func baseInvoice() *bill.Invoice {
	inv := &bill.Invoice{
		Currency:  currency.COP,
		Code:      "TEST",
		IssueDate: cal.MakeDate(2022, 12, 27),
		Type:      bill.InvoiceTypeStandard,
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

func creditNote() *bill.Invoice {
	inv := &bill.Invoice{
		Currency:  currency.COP,
		Code:      "TEST",
		Type:      bill.InvoiceTypeCreditNote,
		IssueDate: cal.MakeDate(2022, 12, 29),
		Preceding: []*bill.Preceding{
			{
				Code:      "TEST",
				IssueDate: cal.NewDate(2022, 12, 27),
				Changes:   []cbc.Key{co.CorrectionKeyRevoked},
			},
		},
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
	require.NoError(t, inv.Calculate())
	assert.Equal(t, inv.Type, bill.InvoiceTypeStandard)
	require.NoError(t, inv.Validate())
	assert.Equal(t, inv.Supplier.Addresses[0].Locality, "Bogotá, D.C.")
	assert.Equal(t, inv.Supplier.Addresses[0].Region, "Bogotá")
	assert.Equal(t, inv.Customer.Addresses[0].Locality, "Sabanalarga")
	assert.Equal(t, inv.Customer.Addresses[0].Region, "Atlántico")

	inv.Supplier.TaxID.Zone = ""
	err := inv.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "zone: cannot be blank")

	inv.Supplier.TaxID.Zone = "1100X"
	err = inv.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "zone: must be a valid value")

	inv = baseInvoice()
	inv.Supplier.TaxID.Type = co.TaxIdentityTypeCitizen
	require.NoError(t, inv.Calculate())
	err = inv.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "supplier: (tax_id: (type: must be a valid value.).).")

	inv = baseInvoice()
	inv.Customer.TaxID.Type = co.TaxIdentityTypeCitizen
	inv.Customer.TaxID.Code = "100100100"
	inv.Customer.TaxID.Zone = ""
	require.NoError(t, inv.Calculate())
	err = inv.Validate()
	assert.NoError(t, err)
}

func TestBasicCreditNoteValidation(t *testing.T) {
	inv := creditNote()
	inv.Preceding[0].Reason = "Correcting an error"
	err := inv.Calculate()
	require.NoError(t, err)
	err = inv.Validate()
	assert.NoError(t, err)
	assert.Contains(t, inv.Preceding[0].Changes, co.CorrectionKeyRevoked)

	inv.Preceding[0].Changes = []cbc.Key{co.CorrectionKeyDiscount, co.CorrectionKeyOther}
	err = inv.Validate()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "changes: the length must be exactly 1.")
	}

	inv.Preceding[0].Changes = []cbc.Key{"fooo"}
	err = inv.Validate()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "changes: (0: must be a valid value.)")
	}

}
