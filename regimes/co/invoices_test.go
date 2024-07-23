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
					Locality: "Bogotá, D.C.",
					Region:   "Bogotá",
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
					Locality: "Sabanalarga",
					Region:   "Atlántico",
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
				Ext: tax.Extensions{
					co.ExtKeyDIANCreditCode: "2", // revoked
				},
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
					Locality: "Bogotá, D.C.",
					Region:   "Bogotá",
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
					Locality: "Sabanalarga",
					Region:   "Atlántico",
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

	delete(inv.Supplier.Ext, co.ExtKeyDIANMunicipality)
	err := inv.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "supplier: (ext: cannot be blank.).")

	inv.Supplier.Ext[co.ExtKeyDIANMunicipality] = "110011"
	err = inv.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "supplier: (ext: (co-dian-municipality: does not match pattern.).")

	inv = baseInvoice()
	inv.Supplier.TaxID.Code = ""
	require.NoError(t, inv.Calculate())
	err = inv.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "supplier: (tax_id: (code: cannot be blank.).).")

	inv = baseInvoice()
	inv.Tax = &bill.Tax{Tags: []cbc.Key{tax.TagSimplified}}
	inv.Customer.TaxID.Code = ""
	inv.Customer.Identities = org.AddIdentity(inv.Customer.Identities,
		&org.Identity{
			Key:  co.IdentityKeyCitizenID,
			Code: "124499654",
		},
	)
	require.NoError(t, inv.Calculate())
	err = inv.Validate()
	assert.NoError(t, err)

	inv = baseInvoice()
	inv.Customer.TaxID.Country = l10n.ES
	inv.Customer.TaxID.Code = "A13180492"
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
	assert.Contains(t, inv.Preceding[0].Ext, co.ExtKeyDIANCreditCode)
	assert.Equal(t, inv.Preceding[0].Ext[co.ExtKeyDIANCreditCode], tax.ExtValue("2"))

	inv.Preceding[0].Ext["foo"] = "bar"
	err = inv.Validate()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "preceding: (0: (ext: (foo: undefined.).).)")
	}

}

func TestNormalizeParty(t *testing.T) {
	p := &org.Party{
		Name: "Test Party",
		TaxID: &tax.Identity{
			Country: l10n.CO,
			Code:    "412615332",
			Zone:    "11001",
		},
	}
	err := co.Calculate(p)
	assert.NoError(t, err)
	assert.Empty(t, p.TaxID.Zone) //nolint:staticcheck
	assert.Equal(t, p.Ext[co.ExtKeyDIANMunicipality].String(), "11001")
}
