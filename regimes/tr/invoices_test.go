package tr_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime:    tax.WithRegime("TR"),
		Series:    "TEST",
		Code:      "001",
		IssueDate: cal.MakeDate(2024, 3, 15),
		Currency:  "TRY",
		Supplier: &org.Party{
			Name: "Test Supplier Co.",
			Addresses: []*org.Address{
				{
					Street:   "Test Street",
					Locality: "Levent",
					Region:   "Istanbul",
					Country:  "TR",
				},
			},
			TaxID: &tax.Identity{
				Country: "TR",
				Code:    "1234567890",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer Ltd.",
			Addresses: []*org.Address{
				{
					Street:   "Test Street",
					Locality: "Kizilay",
					Region:   "Ankara",
					Country:  "TR",
				},
			},
			TaxID: &tax.Identity{
				Country: "TR",
				Code:    "1234567890",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Software License",
					Price: num.NewAmount(50000, 2),
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

func TestInvoiceValidation(t *testing.T) {
	t.Run("valid invoice", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("missing supplier tax ID", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, inv.Validate(), "supplier: (tax_id: cannot be blank.)")
	})

	t.Run("supplier tax ID missing code", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID.Code = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, inv.Validate(), "supplier: (tax_id: (code: cannot be blank.)")
	})
}

