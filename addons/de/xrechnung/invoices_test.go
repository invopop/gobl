package xrechnung_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/de/xrechnung"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := &bill.Invoice{
		Regime:    tax.WithRegime("DE"),
		Addons:    tax.WithAddons(xrechnung.V3),
		IssueDate: cal.MakeDate(2024, 1, 1),
		Type:      "standard",
		Currency:  "EUR",
		Series:    "2024",
		Code:      "1000",
		Supplier: &org.Party{
			Name: "Cursor AG",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "505898911",
			},
			Inboxes: []*org.Inbox{
				{
					Scheme: "0204",
					Code:   "505898911",
				},
			},
			People: []*org.Person{
				{
					Name: &org.Name{
						Given:   "Peter",
						Surname: "Cursorstone",
					},
					Emails: []*org.Email{
						{
							Address: "peter@test.com",
						},
					},
					Telephones: []*org.Telephone{
						{
							Number: "+49100200300",
						},
					},
				},
			},
			Addresses: []*org.Address{
				{
					Street:   "Dietmar-Hopp-Allee",
					Locality: "Walldorf",
					Code:     "69190",
					Country:  "DE",
				},
			},
			Emails: []*org.Email{
				{
					Address: "billing@test.com",
				},
			},
			Telephones: []*org.Telephone{
				{
					Number: "+49100200300",
				},
			},
		},
		Customer: &org.Party{
			Name: "Sample Consumer",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "449674701",
			},
			Inboxes: []*org.Inbox{
				{
					Email: "billing@sample.com",
				},
			},
			People: []*org.Person{
				{
					Name: &org.Name{
						Given:   "Max",
						Surname: "Musterman",
					},
					Telephones: []*org.Telephone{
						{
							Number: "+49100200300",
						},
					},
				},
			},
			Addresses: []*org.Address{
				{
					Street:   "Werner-Heisenberg-Allee",
					Locality: "München",
					Code:     "80939",
					Country:  "DE",
				},
			},
		},
		Payment: &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key: "credit-transfer+sepa",
				CreditTransfer: []*pay.CreditTransfer{
					{
						IBAN: "DE89370400440532013000",
					},
				},
			},
			Terms: &pay.Terms{
				Notes: "Please pay within 10 days",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 2),
					Unit:  "item",
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "general",
					},
				},
				Discounts: []*bill.LineDiscount{
					{
						Reason:  "Testing",
						Percent: num.NewPercentage(10, 2),
					},
				},
			},
		},
		Ordering: &bill.Ordering{
			Code: "1234567890",
		},
	}
	return inv
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("standard invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})
	t.Run("missing supplier tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (identities: missing key 'de-tax-number'; tax_id: cannot be blank.).")
	})
	t.Run("missing supplier tax ID but has tax number", func(t *testing.T) {
		// this is validation is performed in the DE regime, but we're
		// leaving it here for completeness.
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID = nil
		inv.Supplier.Identities = []*org.Identity{
			{
				Key:  "de-tax-number",
				Code: "123/456/7890",
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("nil tax", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Tax = nil
		add := tax.AddonForKey(xrechnung.V3)
		err := add.Validator(inv)
		assert.NoError(t, err)
	})

	// Test supplier telephone scenarios
	t.Run("supplier with party telephones only", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.People[0].Telephones = nil
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("supplier with people telephones only", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Telephones = nil
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("supplier missing both party and people telephones", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Telephones = nil
		inv.Supplier.People[0].Telephones = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "either party.telephones or party.people[0].telephones is required")
	})

	// Test supplier email scenarios
	t.Run("supplier with party emails only", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.People[0].Emails = nil
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("supplier with people emails only", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Emails = nil
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("supplier missing both party and people emails", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Emails = nil
		inv.Supplier.People[0].Emails = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "either party.emails or party.people[0].emails is required")
	})

	t.Run("ordering missing both code and identities", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Ordering = &bill.Ordering{}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "code: cannot be blank.")
	})

	// Test delivery scenarios
	t.Run("delivery with valid receiver", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Delivery = &bill.DeliveryDetails{
			Receiver: &org.Party{
				Name: "Delivery Receiver",
				Addresses: []*org.Address{
					{
						Street:   "Delivery Street",
						Locality: "Berlin",
						Code:     "10115",
						Country:  "DE",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("delivery with missing receiver", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Delivery = &bill.DeliveryDetails{}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "receiver: cannot be blank.")
	})

	t.Run("delivery with receiver missing address", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Delivery = &bill.DeliveryDetails{
			Receiver: &org.Party{
				Name: "Delivery Receiver",
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "addresses: cannot be blank.")
	})

	t.Run("nil delivery", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Delivery = nil
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

}
