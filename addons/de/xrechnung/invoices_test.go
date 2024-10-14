package xrechnung_test

import (
	"testing"

	"github.com/invopop/gobl/addons/de/xrechnung"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"

	// "github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := &bill.Invoice{
		Addons:    tax.WithAddons(xrechnung.V3),
		IssueDate: cal.MakeDate(2024, 1, 1),
		Currency:  "EUR",
		Series:    "A",
		Code:      "1000",
		Supplier: &org.Party{
			Name: "Cursor AG",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "82741168",
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
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "111111125",
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
				Quantity: num.MakeAmount(1, 3),
				Item: &org.Item{
					Name:  "Cursor Subscription",
					Price: num.MakeAmount(1000, 3),
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

	t.Run("missing customer tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (tax_id: cannot be blank.)")
	})

}

func assertValidationError(t *testing.T, inv *bill.Invoice, expectedError string) {
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	assert.ErrorContains(t, err, expectedError)
}
