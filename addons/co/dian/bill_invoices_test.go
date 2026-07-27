package dian_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/co/dian"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
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
		Code:      "12345",
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
			Ext: tax.ExtensionsOf(cbc.CodeMap{
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
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				dian.ExtKeyMunicipality:         "08638",
				dian.ExtKeyFiscalResponsibility: "O-47",
			}),
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1000, 3),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.NewAmount(1000, 3),
				},
			},
		},
		Payment: &bill.PaymentDetails{
			Terms: &pay.Terms{
				DueDates: []*pay.DueDate{
					{Date: cal.NewDate(2023, 1, 27), Percent: num.NewPercentage(1, 0)},
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
		Code:      "12346",
		Type:      bill.InvoiceTypeCreditNote,
		IssueDate: cal.MakeDate(2022, 12, 29),
		Preceding: []*org.DocumentRef{
			{
				Code:      "12345",
				IssueDate: cal.NewDate(2022, 12, 27),
				Stamps: []*head.Stamp{
					{Provider: dian.StampCUDE, Value: "a1b2c3d4e5f6"},
				},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
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
			Ext: tax.ExtensionsOf(cbc.CodeMap{
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
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				dian.ExtKeyMunicipality:         "08638",
				dian.ExtKeyFiscalResponsibility: "O-47",
			}),
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1000, 3),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.NewAmount(1000, 3),
				},
			},
		},
		Payment: &bill.PaymentDetails{
			Terms: &pay.Terms{
				DueDates: []*pay.DueDate{
					{Date: cal.NewDate(2023, 1, 29), Percent: num.NewPercentage(1, 0)},
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

	t.Run("handles nil invoice", func(t *testing.T) {
		var inv *bill.Invoice
		assert.NotPanics(t, func() {
			norm.Normalize(inv, tax.AddonContext(dian.V2))
		})
		assert.Nil(t, inv)
	})

	t.Run("sets default tax responsibility for Colombian supplier", func(t *testing.T) {
		inv := baseInvoice()
		// Remove existing tax responsibility
		inv.Supplier.Ext = inv.Supplier.Ext.Delete(dian.ExtKeyFiscalResponsibility)

		norm.Normalize(inv, tax.AddonContext(dian.V2))

		assert.Equal(t, cbc.Code("R-99-PN"), inv.Supplier.Ext.Get(dian.ExtKeyFiscalResponsibility))
	})

	t.Run("sets default tax responsibility for Colombian customer", func(t *testing.T) {
		inv := baseInvoice()
		// Remove existing tax responsibility
		inv.Customer.Ext = inv.Customer.Ext.Delete(dian.ExtKeyFiscalResponsibility)

		norm.Normalize(inv, tax.AddonContext(dian.V2))

		assert.Equal(t, cbc.Code("R-99-PN"), inv.Customer.Ext.Get(dian.ExtKeyFiscalResponsibility))
	})

	t.Run("keeps existing tax responsibility for supplier", func(t *testing.T) {
		inv := baseInvoice()
		// Set a specific tax responsibility
		inv.Supplier.Ext = inv.Supplier.Ext.Set(dian.ExtKeyFiscalResponsibility, "O-13")

		norm.Normalize(inv, tax.AddonContext(dian.V2))

		assert.Equal(t, cbc.Code("O-13"), inv.Supplier.Ext.Get(dian.ExtKeyFiscalResponsibility))
	})

	t.Run("keeps existing tax responsibility for customer", func(t *testing.T) {
		inv := baseInvoice()
		// Set a specific tax responsibility
		inv.Customer.Ext = inv.Customer.Ext.Set(dian.ExtKeyFiscalResponsibility, "O-47")

		norm.Normalize(inv, tax.AddonContext(dian.V2))

		assert.Equal(t, cbc.Code("O-47"), inv.Customer.Ext.Get(dian.ExtKeyFiscalResponsibility))
	})

	t.Run("does not set tax responsibility for non-Colombian supplier", func(t *testing.T) {
		inv := baseInvoice()
		inv.Supplier.TaxID.Country = "ES"
		inv.Supplier.Ext = inv.Supplier.Ext.Delete(dian.ExtKeyFiscalResponsibility)

		norm.Normalize(inv, tax.AddonContext(dian.V2))

		assert.Empty(t, inv.Supplier.Ext.Get(dian.ExtKeyFiscalResponsibility))
	})

	t.Run("does not set tax responsibility for non-Colombian customer", func(t *testing.T) {
		inv := baseInvoice()
		inv.Customer.TaxID.Country = "ES"
		inv.Customer.Ext = inv.Customer.Ext.Delete(dian.ExtKeyFiscalResponsibility)

		norm.Normalize(inv, tax.AddonContext(dian.V2))

		assert.Empty(t, inv.Customer.Ext.Get(dian.ExtKeyFiscalResponsibility))
	})

	t.Run("handles nil supplier", func(t *testing.T) {
		inv := baseInvoice()
		inv.Supplier = nil

		assert.NotPanics(t, func() {
			norm.Normalize(inv, tax.AddonContext(dian.V2))
		})
	})

	t.Run("handles nil customer", func(t *testing.T) {
		inv := baseInvoice()
		inv.Customer = nil

		assert.NotPanics(t, func() {
			norm.Normalize(inv, tax.AddonContext(dian.V2))
		})
	})

	t.Run("handles nil extensions", func(t *testing.T) {
		inv := baseInvoice()
		inv.Supplier.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.Extensions{}

		norm.Normalize(inv, tax.AddonContext(dian.V2))

		assert.Equal(t, cbc.Code("R-99-PN"), inv.Supplier.Ext.Get(dian.ExtKeyFiscalResponsibility))
		assert.Equal(t, cbc.Code("R-99-PN"), inv.Customer.Ext.Get(dian.ExtKeyFiscalResponsibility))
	})
}

func TestInvoiceCodeValidation(t *testing.T) {
	t.Run("rejects non-numeric codes", func(t *testing.T) {
		inv := baseInvoice()
		inv.Code = "TEST"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "GOBL-CO-DIAN-BILL-INVOICE-15")
	})

	t.Run("rejects leading zeroes", func(t *testing.T) {
		inv := baseInvoice()
		inv.Code = "0123"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "GOBL-CO-DIAN-BILL-INVOICE-15")
	})

	t.Run("allows empty draft codes", func(t *testing.T) {
		inv := baseInvoice()
		inv.Code = ""
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})
}

func TestCreditNoteWithoutReferenceValidation(t *testing.T) {
	// unreferenced returns a credit note whose preceding ref has no dian-cude stamp (91-22).
	unreferenced := func() *bill.Invoice {
		inv := creditNote()
		inv.Preceding[0].Stamps = nil
		inv.Preceding[0].Reason = "Correcting an error"
		inv.Preceding[0].Ext = inv.Preceding[0].Ext.Set(dian.ExtKeyCreditCode, "3")
		return inv
	}

	t.Run("rejects revoking without a referenced document", func(t *testing.T) {
		inv := unreferenced()
		inv.Preceding[0].Ext = inv.Preceding[0].Ext.Set(dian.ExtKeyCreditCode, "2")
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "GOBL-CO-DIAN-BILL-INVOICE-16")
	})

	t.Run("treats empty stamp values as unreferenced", func(t *testing.T) {
		inv := creditNote()
		inv.Preceding[0].Reason = "Correcting an error"
		inv.Preceding[0].Stamps[0].Value = ""
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "GOBL-CO-DIAN-BILL-INVOICE-16")
	})

	t.Run("allows revoking with a referenced document", func(t *testing.T) {
		inv := creditNote()
		inv.Preceding[0].Reason = "Correcting an error"
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("rejects periods crossing calendar months", func(t *testing.T) {
		inv := unreferenced()
		inv.Preceding[0].Period = &cal.Period{Start: cal.MakeDate(2022, 11, 15), End: cal.MakeDate(2022, 12, 15)}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "GOBL-CO-DIAN-BILL-INVOICE-17")
	})

	t.Run("allows periods within one calendar month", func(t *testing.T) {
		inv := unreferenced()
		inv.Preceding[0].Period = &cal.Period{Start: cal.MakeDate(2022, 12, 1), End: cal.MakeDate(2022, 12, 28)}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("drops nil preceding entries on calculation", func(t *testing.T) {
		inv := creditNote()
		inv.Preceding = append(inv.Preceding, nil)
		inv.Preceding[0].Reason = "Correcting an error"
		require.NoError(t, inv.Calculate())
		assert.Len(t, inv.Preceding, 1)
		assert.NoError(t, rules.Validate(inv))
	})
}

func TestPaymentDueDateValidation(t *testing.T) {
	t.Run("rejects unpaid invoices without due dates", func(t *testing.T) {
		inv := baseInvoice()
		inv.Payment = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "GOBL-CO-DIAN-BILL-INVOICE-18")
	})

	t.Run("rejects unpaid invoices with only nil due dates", func(t *testing.T) {
		inv := baseInvoice()
		inv.Payment.Terms.DueDates = []*pay.DueDate{nil}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "GOBL-CO-DIAN-BILL-INVOICE-18")
	})

	t.Run("allows fully paid invoices without due dates", func(t *testing.T) {
		inv := baseInvoice()
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Record{
				{Description: "prepaid", Percent: num.NewPercentage(1, 0)},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("skips proforma invoices", func(t *testing.T) {
		inv := baseInvoice()
		inv.Type = bill.InvoiceTypeProforma
		inv.Payment = nil
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("skips uncalculated drafts", func(t *testing.T) {
		inv := baseInvoice()
		inv.Payment = nil
		err := rules.Validate(inv)
		require.Error(t, err) // core totals faults, but no due date fault
		assert.NotContains(t, err.Error(), "GOBL-CO-DIAN-BILL-INVOICE-18")
	})
}
