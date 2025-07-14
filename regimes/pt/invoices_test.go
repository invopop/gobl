package pt_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Code:    "123456789",
				Country: "PT",
			},
			Name: "Test Supplier",
		},
		Customer: &org.Party{
			Name: "Test Customer",
		},
		Code:      "INV/1",
		Currency:  "EUR",
		IssueDate: cal.MakeDate(2023, 1, 1),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
			},
		},
	}
}

func TestValidInvoice(t *testing.T) {
	reg := tax.RegimeDefFor(l10n.PT)

	inv := validInvoice()
	require.NoError(t, reg.Validator(inv))
}

func TestValidSimplifiedInvoice(t *testing.T) {
	reg := tax.RegimeDefFor(l10n.PT)

	inv := validInvoice()
	inv.SetTags(tax.TagSimplified, pt.TagInvoiceReceipt)
	inv.Customer = nil
	require.NoError(t, reg.Validator(inv))
}

func TestLineValidation(t *testing.T) {
	reg := tax.RegimeDefFor(l10n.PT)

	t.Run("negative quantity", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Quantity = num.MakeAmount(-1, 0)
		assert.ErrorContains(t, reg.Validator(inv), "lines: (0: (quantity: must be no less than 0.).)")
	})

	t.Run("negative price", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Item.Price = num.NewAmount(-1, 0)
		assert.ErrorContains(t, reg.Validator(inv), "lines: (0: (item: (price: must be no less than 0.).).)")
	})
}

func TestInvoicePaymentValidation(t *testing.T) {
	reg := tax.RegimeDefFor(l10n.PT)

	t.Run("empty advances", func(t *testing.T) {
		inv := validInvoice()
		inv.Payment = &bill.PaymentDetails{}
		require.NoError(t, reg.Validator(inv))
	})

	t.Run("advance with past date", func(t *testing.T) {
		inv := validInvoice()
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Advance{
				{
					Date: cal.NewDate(2022, 12, 31),
				},
			},
		}
		require.NoError(t, reg.Validator(inv))
	})

	t.Run("advance with current date", func(t *testing.T) {
		inv := validInvoice()
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Advance{
				{
					Date: cal.NewDate(2023, 1, 1),
				},
			},
		}
		require.NoError(t, reg.Validator(inv))
	})

	t.Run("advance with future date", func(t *testing.T) {
		inv := validInvoice()
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Advance{
				{
					Date: cal.NewDate(2023, 1, 2),
				},
			},
		}
		assert.ErrorContains(t, reg.Validator(inv), "advances: (0: (date: too late")
	})

	t.Run("nil advance", func(t *testing.T) {
		inv := validInvoice()
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Advance{nil},
		}
		require.NoError(t, reg.Validator(inv))
	})
}
