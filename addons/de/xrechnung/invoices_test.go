package xrechnung_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/invopop/gobl/addons/de/xrechnung"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"

	// "github.com/invopop/gobl/l10n"

	// "github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"

	// "github.com/invopop/gobl/regimes/de"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	p := num.MakePercentage(19, 2)
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
			// TaxID: &tax.Identity{
			// 	Country: "DE",
			// 	Code:    "505898911",
			// },
			People: []*org.Person{
				{
					Name: &org.Name{
						Given:   "Peter",
						Surname: "Cursorstone",
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
					Address: "billing@cursor.com",
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
			// TaxID: &tax.Identity{
			// 	Country: "DE",
			// 	Code:    "449674701",
			// },
			People: []*org.Person{
				{
					Name: &org.Name{
						Given:   "Max",
						Surname: "Musterman",
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
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						// Rate:     "standard",
						Percent: &p,
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
		Payment: &bill.Payment{
			Instructions: &pay.Instructions{
				Key: "credit-transfer",
				CreditTransfer: []*pay.CreditTransfer{
					{
						IBAN: "DE89370400440532013000",
					},
				},
			},
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
		invJSON, err := json.MarshalIndent(inv, "", "  ")
		require.NoError(t, err)
		fmt.Println(string(invJSON))

		require.NoError(t, inv.Calculate())
		errr := inv.Validate()
		assert.ErrorContains(t, errr, "supplier: (tax_id: cannot be blank.)")
	})

	t.Run("missing invoice type", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = ""
		err := inv.Validate()
		assert.ErrorContains(t, err, "type: cannot be blank.")
	})

	t.Run("missing payment instructions", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment.Instructions = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "payment: (instructions: cannot be blank.).")
	})

	t.Run("missing ordering code", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Ordering.Code = ""
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "ordering: (code: cannot be blank.).")
	})

	t.Run("missing supplier city", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Addresses[0].Locality = ""
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (addresses: (0: (locality: cannot be blank.).).).")
	})

	t.Run("missing supplier postcode", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Addresses[0].Code = ""
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (addresses: (0: (code: cannot be blank.).).).")
	})

	t.Run("missing customer city", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Addresses[0].Locality = ""
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (addresses: (0: (locality: cannot be blank.).).).")
	})

	t.Run("missing customer postcode", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Addresses[0].Code = ""
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (addresses: (0: (code: cannot be blank.).).).")
	})

	t.Run("missing supplier name", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Name = ""
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (name: cannot be blank.)")
	})

	t.Run("missing delivery address", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Delivery = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err, "Delivery address should be optional")
	})

	t.Run("incomplete delivery address", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Delivery = &bill.Delivery{
			Receiver: &org.Party{
				Addresses: []*org.Address{
					{
						Street:  "Delivery Street",
						Country: "DE",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "delivery: (party.addresses.0.locality: cannot be blank.)")
		assert.ErrorContains(t, err, "delivery: (party.addresses.0.code: cannot be blank.)")
	})

	t.Run("valid delivery address", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Delivery = &bill.Delivery{
			Receiver: &org.Party{
				Addresses: []*org.Address{
					{
						Street:   "Delivery Street",
						Locality: "Delivery City",
						Code:     "12345",
						Country:  "DE",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err, "Valid delivery address should not cause validation errors")
	})

}
