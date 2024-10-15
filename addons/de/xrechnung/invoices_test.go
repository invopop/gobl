package xrechnung_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/invopop/gobl/addons/de/xrechnung"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"

	// "github.com/invopop/gobl/cbc"
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
				Code:    "449674701",
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
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
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
		invJSON, err := json.MarshalIndent(inv, "", "  ")
		require.NoError(t, err)
		fmt.Println(string(invJSON))
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
