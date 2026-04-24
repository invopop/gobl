package dian_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/co/dian"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func baseInvoice() *bill.Invoice {
	inv := &bill.Invoice{
		Regime:    tax.WithRegime("CO"),
		Addons:    tax.WithAddons(dian.V2),
		Currency:  currency.COP,
		Code:      "TEST",
		IssueDate: cal.MakeDate(2022, 12, 27),
		Type:      bill.InvoiceTypeStandard,
		Supplier: &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "CO",
				Code:    "412615332",
			},
			Addresses: []*org.Address{
				{
					Locality: "Bogotá, D.C.",
					Region:   "Bogotá",
				},
			},
			Ext: tax.ExtensionsOf(tax.ExtMap{
				dian.ExtKeyMunicipality:         "11001",
				dian.ExtKeyFiscalResponsibility: "O-13",
			}),
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "CO",
				Code:    "124499654",
			},
			Addresses: []*org.Address{
				{
					Locality: "Sabanalarga",
					Region:   "Atlántico",
				},
			},
			Ext: tax.ExtensionsOf(tax.ExtMap{
				dian.ExtKeyMunicipality:         "08638",
				dian.ExtKeyFiscalResponsibility: "O-47",
			}),
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 3),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.NewAmount(1000, 3),
				},
			},
		},
	}
	return inv
}

func creditNote() *bill.Invoice {
	inv := &bill.Invoice{
		Regime:    tax.WithRegime("CO"),
		Addons:    tax.WithAddons(dian.V2),
		Currency:  currency.COP,
		Code:      "TEST",
		Type:      bill.InvoiceTypeCreditNote,
		IssueDate: cal.MakeDate(2022, 12, 29),
		Preceding: []*org.DocumentRef{
			{
				Code:      "TEST",
				IssueDate: cal.NewDate(2022, 12, 27),
				Ext: tax.ExtensionsOf(tax.ExtMap{
					dian.ExtKeyCreditCode: "2", // revoked
				}),
			},
		},
		Supplier: &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "CO",
				Code:    "412615332",
			},
			Addresses: []*org.Address{
				{
					Locality: "Bogotá, D.C.",
					Region:   "Bogotá",
				},
			},
			Ext: tax.ExtensionsOf(tax.ExtMap{
				dian.ExtKeyMunicipality:         "11001",
				dian.ExtKeyFiscalResponsibility: "O-47",
			}),
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "CO",
				Code:    "124499654",
			},
			Addresses: []*org.Address{
				{
					Locality: "Sabanalarga",
					Region:   "Atlántico",
				},
			},
			Ext: tax.ExtensionsOf(tax.ExtMap{
				dian.ExtKeyMunicipality:         "08638",
				dian.ExtKeyFiscalResponsibility: "O-47",
			}),
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 3),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.NewAmount(1000, 3),
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
	require.NoError(t, rules.Validate(inv))
	assert.Equal(t, inv.Supplier.Addresses[0].Locality, "Bogotá, D.C.")
	assert.Equal(t, inv.Supplier.Addresses[0].Region, "Bogotá")
	assert.Equal(t, inv.Customer.Addresses[0].Locality, "Sabanalarga")
	assert.Equal(t, inv.Customer.Addresses[0].Region, "Atlántico")

	inv.Supplier.Ext = inv.Supplier.Ext.Delete(dian.ExtKeyMunicipality)
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "extension 'co-dian-municipality' is required")

	inv = baseInvoice()
	inv.SetTags(tax.TagSimplified)
	inv.Customer.TaxID.Code = ""
	inv.Customer.Identities = org.AddIdentity(inv.Customer.Identities,
		&org.Identity{
			Key:  dian.IdentityKeyCitizenID,
			Code: "124499654",
		},
	)
	require.NoError(t, inv.Calculate())
	err = rules.Validate(inv)
	assert.NoError(t, err)

	inv = baseInvoice()
	inv.Customer.TaxID.Country = "ES"
	inv.Customer.TaxID.Code = "A13180492"
	require.NoError(t, inv.Calculate())
	err = rules.Validate(inv)
	assert.NoError(t, err)
}

func TestInvoiceCurrencyValidation(t *testing.T) {
	t.Run("non-COP currency without exchange rates", func(t *testing.T) {
		inv := baseInvoice()
		inv.Currency = "USD"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-CO-DIAN-BILL-INVOICE-14] invoice must be in COP or provide exchange rate for conversion")
	})

	t.Run("non-COP currency with exchange rates", func(t *testing.T) {
		inv := baseInvoice()
		inv.Currency = "USD"
		inv.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   "USD",
				To:     "COP",
				Amount: num.MakeAmount(4100, 0),
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.NoError(t, err)
	})
}

func TestFiscalResponsibilityExtensionValidation(t *testing.T) {
	// Colombian parties
	inv := baseInvoice()
	require.NoError(t, inv.Calculate()) // calculate before delete to avoid normalization
	inv.Supplier.Ext = inv.Supplier.Ext.Delete(dian.ExtKeyFiscalResponsibility)
	inv.Customer.Ext = inv.Customer.Ext.Delete(dian.ExtKeyFiscalResponsibility)
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "extension 'co-dian-fiscal-responsibility' is required")

	// Non-Colombian parties
	inv = baseInvoice()
	inv.Supplier.TaxID.Code = "E47180476"
	inv.Supplier.TaxID.Country = "ES"
	inv.Customer.TaxID.Code = "C87547287"
	inv.Customer.TaxID.Country = "ES"
	inv.Supplier.Ext = inv.Supplier.Ext.Delete(dian.ExtKeyFiscalResponsibility)
	inv.Customer.Ext = inv.Customer.Ext.Delete(dian.ExtKeyFiscalResponsibility)
	require.NoError(t, inv.Calculate())
	err = rules.Validate(inv)
	assert.NoError(t, err)
}

func TestBasicCreditNoteValidation(t *testing.T) {
	inv := creditNote()
	inv.Preceding[0].Reason = "Correcting an error"
	err := inv.Calculate()
	require.NoError(t, err)
	err = rules.Validate(inv)
	assert.NoError(t, err)
	assert.True(t, inv.Preceding[0].Ext.Has(dian.ExtKeyCreditCode))
	assert.Equal(t, inv.Preceding[0].Ext.Get(dian.ExtKeyCreditCode), cbc.Code("2"))
}

func TestNormalizeInvoice(t *testing.T) {
	addon := tax.AddonForKey(dian.V2)

	t.Run("handles nil invoice", func(t *testing.T) {
		var inv *bill.Invoice
		assert.NotPanics(t, func() {
			addon.Normalizer(inv)
		})
		assert.Nil(t, inv)
	})

	t.Run("sets default tax responsibility for Colombian supplier", func(t *testing.T) {
		inv := baseInvoice()
		// Remove existing tax responsibility
		inv.Supplier.Ext = inv.Supplier.Ext.Delete(dian.ExtKeyFiscalResponsibility)

		addon.Normalizer(inv)

		assert.Equal(t, cbc.Code("R-99-PN"), inv.Supplier.Ext.Get(dian.ExtKeyFiscalResponsibility))
	})

	t.Run("sets default tax responsibility for Colombian customer", func(t *testing.T) {
		inv := baseInvoice()
		// Remove existing tax responsibility
		inv.Customer.Ext = inv.Customer.Ext.Delete(dian.ExtKeyFiscalResponsibility)

		addon.Normalizer(inv)

		assert.Equal(t, cbc.Code("R-99-PN"), inv.Customer.Ext.Get(dian.ExtKeyFiscalResponsibility))
	})

	t.Run("keeps existing tax responsibility for supplier", func(t *testing.T) {
		inv := baseInvoice()
		// Set a specific tax responsibility
		inv.Supplier.Ext = inv.Supplier.Ext.Set(dian.ExtKeyFiscalResponsibility, "O-13")

		addon.Normalizer(inv)

		assert.Equal(t, cbc.Code("O-13"), inv.Supplier.Ext.Get(dian.ExtKeyFiscalResponsibility))
	})

	t.Run("keeps existing tax responsibility for customer", func(t *testing.T) {
		inv := baseInvoice()
		// Set a specific tax responsibility
		inv.Customer.Ext = inv.Customer.Ext.Set(dian.ExtKeyFiscalResponsibility, "O-47")

		addon.Normalizer(inv)

		assert.Equal(t, cbc.Code("O-47"), inv.Customer.Ext.Get(dian.ExtKeyFiscalResponsibility))
	})

	t.Run("does not set tax responsibility for non-Colombian supplier", func(t *testing.T) {
		inv := baseInvoice()
		inv.Supplier.TaxID.Country = "ES"
		inv.Supplier.Ext = inv.Supplier.Ext.Delete(dian.ExtKeyFiscalResponsibility)

		addon.Normalizer(inv)

		assert.Empty(t, inv.Supplier.Ext.Get(dian.ExtKeyFiscalResponsibility))
	})

	t.Run("does not set tax responsibility for non-Colombian customer", func(t *testing.T) {
		inv := baseInvoice()
		inv.Customer.TaxID.Country = "ES"
		inv.Customer.Ext = inv.Customer.Ext.Delete(dian.ExtKeyFiscalResponsibility)

		addon.Normalizer(inv)

		assert.Empty(t, inv.Customer.Ext.Get(dian.ExtKeyFiscalResponsibility))
	})

	t.Run("handles nil supplier", func(t *testing.T) {
		inv := baseInvoice()
		inv.Supplier = nil

		assert.NotPanics(t, func() {
			addon.Normalizer(inv)
		})
	})

	t.Run("handles nil customer", func(t *testing.T) {
		inv := baseInvoice()
		inv.Customer = nil

		assert.NotPanics(t, func() {
			addon.Normalizer(inv)
		})
	})

	t.Run("handles nil extensions", func(t *testing.T) {
		inv := baseInvoice()
		inv.Supplier.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.Extensions{}

		addon.Normalizer(inv)

		assert.Equal(t, cbc.Code("R-99-PN"), inv.Supplier.Ext.Get(dian.ExtKeyFiscalResponsibility))
		assert.Equal(t, cbc.Code("R-99-PN"), inv.Customer.Ext.Get(dian.ExtKeyFiscalResponsibility))
	})
}
