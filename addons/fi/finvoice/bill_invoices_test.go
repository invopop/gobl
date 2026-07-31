package finvoice_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/fi/finvoice"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Regime:    tax.WithRegime("FI"),
		Addons:    tax.WithAddons(finvoice.V3),
		IssueDate: cal.MakeDate(2026, 7, 1),
		Type:      "standard",
		Currency:  "EUR",
		Series:    "2026",
		Code:      "1000",
		Supplier: &org.Party{
			Name: "Myyjä Oy",
			TaxID: &tax.Identity{
				Country: "FI",
				Code:    "23456780",
			},
			Addresses: []*org.Address{
				{
					Street:   "Esimerkkikatu 1",
					Locality: "Helsinki",
					Code:     "00100",
					Country:  "FI",
				},
			},
		},
		Customer: &org.Party{
			Name: "Ostaja Oy",
			TaxID: &tax.Identity{
				Country: "FI",
				Code:    "01120389",
			},
			Addresses: []*org.Address{
				{
					Street:   "Testitie 2",
					Locality: "Espoo",
					Code:     "02100",
					Country:  "FI",
				},
			},
		},
		Payment: &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key: "credit-transfer+sepa",
				Ref: "RF18539007547034",
				CreditTransfer: []*pay.CreditTransfer{
					{
						IBAN: "FI2112345600000785",
					},
				},
			},
			Terms: &pay.Terms{
				DueDates: []*pay.DueDate{
					{
						Date:   cal.NewDate(2026, 7, 31),
						Amount: num.MakeAmount(12750, 2),
					},
				},
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(1000, 2),
					Unit:  "item",
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "general",
					},
				},
			},
		},
	}
}

func TestBillInvoiceRules(t *testing.T) {
	t.Run("valid invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("valid invoice without BIC", func(t *testing.T) {
		// Finvoice recommends the BIC alongside the IBAN but the schema does
		// not require it, so its absence must not fail validation.
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.Empty(t, inv.Payment.Instructions.CreditTransfer[0].BIC)
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("missing customer", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(tax.TagSimplified)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "customer is required")
	})

	t.Run("missing customer name", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Name = ""
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "customer name is required")
	})

	t.Run("missing payment", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = nil
		require.NoError(t, inv.Calculate())
		faults := rules.Validate(inv)
		require.NotNil(t, faults)
		assert.True(t, faults.HasCode("GOBL-FI-FINVOICE-BILL-INVOICE-03"))
	})

	t.Run("credit note requires payment too", func(t *testing.T) {
		// en16931 only requires payment details when an amount is due, which
		// exempts credit notes. Finvoice's EpiDetails is mandatory on every
		// invoice, so the addon must fault here on its own.
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Preceding = []*org.DocumentRef{
			{
				Series:    "2026",
				Code:      "0999",
				IssueDate: cal.NewDate(2026, 6, 1),
			},
		}
		inv.Payment = nil
		require.NoError(t, inv.Calculate())
		faults := rules.Validate(inv)
		require.NotNil(t, faults)
		assert.True(t, faults.HasCode("GOBL-FI-FINVOICE-BILL-INVOICE-03"))
	})

	t.Run("missing payment instructions", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment.Instructions = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "payment instructions are required")
	})

	t.Run("instructions key must be credit transfer", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment.Instructions.Key = pay.MeansKeyCash
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "credit-transfer")
	})

	t.Run("missing payment reference", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment.Instructions.Ref = ""
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "payment reference is required")
	})

	t.Run("missing credit transfer", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment.Instructions.CreditTransfer = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "credit transfer details are required")
	})

	t.Run("missing IBAN in first credit transfer", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment.Instructions.CreditTransfer[0].IBAN = ""
		inv.Payment.Instructions.CreditTransfer[0].Number = "12345-678"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "IBAN is required")
	})

	t.Run("invalid IBAN checksum", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment.Instructions.CreditTransfer[0].IBAN = "FI2112345600000786"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "IBAN is not valid")
	})

	t.Run("malformed IBAN", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment.Instructions.CreditTransfer[0].IBAN = "NOT-AN-IBAN"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "IBAN is not valid")
	})

	t.Run("missing payment terms", func(t *testing.T) {
		// A credit note has no amount due, so en16931's BR-CO-25 does not
		// require terms here — only the addon does.
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Preceding = []*org.DocumentRef{
			{
				Series:    "2026",
				Code:      "0999",
				IssueDate: cal.NewDate(2026, 6, 1),
			},
		}
		inv.Payment.Terms = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "payment terms are required")
	})

	t.Run("terms with notes but no due dates", func(t *testing.T) {
		// en16931 accepts terms with notes only (BR-CO-25); Finvoice needs a
		// dated entry to fill EpiDateOptionDate.
		inv := testInvoiceStandard(t)
		inv.Payment.Terms = &pay.Terms{
			Notes: "Payable on receipt",
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "due date is required")
	})
}

func TestNormalization(t *testing.T) {
	t.Run("strips spaces and uppercases IBAN", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment.Instructions.CreditTransfer[0].IBAN = "fi21 1234 5600 0007 85"
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "FI2112345600000785", inv.Payment.Instructions.CreditTransfer[0].IBAN)
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("strips spaces from payment reference", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment.Instructions.Ref = "RF18 5390 0754 7034"
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "RF18539007547034", inv.Payment.Instructions.Ref.String())
	})
}
