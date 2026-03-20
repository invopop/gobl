package pt_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime: tax.WithRegime("PT"),
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
						Rate:     "general",
					},
				},
			},
		},
	}
}

func TestValidInvoice(t *testing.T) {
	inv := validInvoice()
	assert.NoError(t, rules.Validate(inv))
}

func TestValidSimplifiedInvoice(t *testing.T) {
	inv := validInvoice()
	inv.SetTags(tax.TagSimplified, pt.TagInvoiceReceipt)
	inv.Customer = nil
	assert.NoError(t, rules.Validate(inv))
}

func TestInvoiceTypeValidation(t *testing.T) {
	t.Run("invalid type", func(t *testing.T) {
		inv := validInvoice()
		inv.Type = "unknown"
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-BILL-INVOICE-01]")
	})

	t.Run("credit note", func(t *testing.T) {
		inv := validInvoice()
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Preceding = []*org.DocumentRef{{Code: "INV/1"}}
		assert.NoError(t, rules.Validate(inv))
	})
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("value date after issue date", func(t *testing.T) {
		inv := validInvoice()
		inv.ValueDate = cal.NewDate(2023, 1, 2)
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-BILL-INVOICE-11]")
	})

	t.Run("value date on issue date", func(t *testing.T) {
		inv := validInvoice()
		inv.ValueDate = cal.NewDate(2023, 1, 1)
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("value date before issue date", func(t *testing.T) {
		inv := validInvoice()
		inv.ValueDate = cal.NewDate(2022, 12, 31)
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("operation date after issue date", func(t *testing.T) {
		inv := validInvoice()
		inv.OperationDate = cal.NewDate(2023, 1, 2)
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-BILL-INVOICE-12]")
	})

	t.Run("operation date on issue date", func(t *testing.T) {
		inv := validInvoice()
		inv.OperationDate = cal.NewDate(2023, 1, 1)
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("operation date before issue date", func(t *testing.T) {
		inv := validInvoice()
		inv.OperationDate = cal.NewDate(2022, 12, 31)
		require.NoError(t, rules.Validate(inv))
	})
}

func TestSupplierValidation(t *testing.T) {
	t.Run("nil supplier", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier = nil
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("missing tax ID", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = nil
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-BILL-INVOICE-03]")
	})
}

func TestLineValidation(t *testing.T) {
	t.Run("nil line", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines = append(inv.Lines, nil)
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("negative quantity", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Quantity = num.MakeAmount(-1, 0)
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-BILL-INVOICE-05]")
	})

	t.Run("nil item", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Item = nil
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("negative price", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Item.Price = num.NewAmount(-1, 0)
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-BILL-INVOICE-06]")
	})
}

func TestInvoicePaymentValidation(t *testing.T) {
	t.Run("empty advances", func(t *testing.T) {
		inv := validInvoice()
		inv.Payment = &bill.PaymentDetails{}
		require.NoError(t, rules.Validate(inv))
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
		require.NoError(t, rules.Validate(inv))
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
		require.NoError(t, rules.Validate(inv))
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
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-BILL-INVOICE-14]")
	})

	t.Run("nil advance", func(t *testing.T) {
		inv := validInvoice()
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Advance{nil},
		}
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("empty terms", func(t *testing.T) {
		inv := validInvoice()
		inv.Payment = &bill.PaymentDetails{}
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("due date with past date", func(t *testing.T) {
		inv := validInvoice()
		inv.Payment = &bill.PaymentDetails{
			Terms: &pay.Terms{
				DueDates: []*pay.DueDate{
					{
						Date: cal.NewDate(2022, 12, 31),
					},
				},
			},
		}
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-BILL-INVOICE-15]")
	})

	t.Run("due date with current date", func(t *testing.T) {
		inv := validInvoice()
		inv.Payment = &bill.PaymentDetails{
			Terms: &pay.Terms{
				DueDates: []*pay.DueDate{
					{
						Date: cal.NewDate(2023, 1, 1),
					},
				},
			},
		}
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("due date with future date", func(t *testing.T) {
		inv := validInvoice()
		inv.Payment = &bill.PaymentDetails{
			Terms: &pay.Terms{
				DueDates: []*pay.DueDate{
					{
						Date: cal.NewDate(2023, 1, 2),
					},
				},
			},
		}
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("nil due date", func(t *testing.T) {
		inv := validInvoice()
		inv.Payment = &bill.PaymentDetails{
			Terms: &pay.Terms{
				DueDates: []*pay.DueDate{nil},
			},
		}
		require.NoError(t, rules.Validate(inv))
	})
}

func TestInvoicePrecedingValidation(t *testing.T) {
	t.Run("empty preceding", func(t *testing.T) {
		inv := validInvoice()
		inv.Preceding = nil
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("empty preceding with credit note", func(t *testing.T) {
		inv := validInvoice()
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Preceding = nil
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-BILL-INVOICE-02]")
	})

	t.Run("preceding document with no date", func(t *testing.T) {
		inv := validInvoice()
		inv.Preceding = []*org.DocumentRef{
			{
				Code:      "INV/1",
				IssueDate: nil,
			},
		}
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("preceding document with past date", func(t *testing.T) {
		inv := validInvoice()
		inv.Preceding = []*org.DocumentRef{
			{
				Code:      "INV/1",
				IssueDate: cal.NewDate(2022, 12, 31),
			},
		}
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("preceding document with same date", func(t *testing.T) {
		inv := validInvoice()
		inv.Preceding = []*org.DocumentRef{
			{
				Code:      "INV/1",
				IssueDate: cal.NewDate(2023, 1, 1),
			},
		}
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("preceding document with future date", func(t *testing.T) {
		inv := validInvoice()
		inv.Preceding = []*org.DocumentRef{
			{
				Code:      "INV/1",
				IssueDate: cal.NewDate(2023, 1, 2),
			},
		}
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-BILL-INVOICE-13]")
	})

	t.Run("nil preceding", func(t *testing.T) {
		inv := validInvoice()
		inv.Preceding = []*org.DocumentRef{nil}
		require.NoError(t, rules.Validate(inv))
	})
}

func TestInvoiceTotalsValidation(t *testing.T) {
	t.Run("negative due amount", func(t *testing.T) {
		inv := validInvoice()
		inv.Totals = &bill.Totals{
			Due: num.NewAmount(-1, 2),
		}
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-BILL-INVOICE-07]")
	})

	t.Run("zero due amount", func(t *testing.T) {
		inv := validInvoice()
		inv.Totals = &bill.Totals{
			Due: num.NewAmount(0, 2),
		}
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("positive due amount", func(t *testing.T) {
		inv := validInvoice()
		inv.Totals = &bill.Totals{
			Due: num.NewAmount(1, 2),
		}
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("nil totals", func(t *testing.T) {
		inv := validInvoice()
		inv.Totals = nil
		require.NoError(t, rules.Validate(inv))
	})
}
