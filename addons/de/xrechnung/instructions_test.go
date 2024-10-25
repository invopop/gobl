package xrechnung_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func invoiceTemplate(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := testInvoiceStandard(t)
	inv.Payment = nil
	return inv
}

func TestValidateInvoice(t *testing.T) {
	t.Run("valid invoice with SEPA credit transfer", func(t *testing.T) {
		inv := invoiceTemplate(t)
		inv.Payment = &bill.Payment{
			Instructions: &pay.Instructions{
				Key: "credit-transfer+sepa",
				CreditTransfer: []*pay.CreditTransfer{
					{
						IBAN: "DE89370400440532013000",
						BIC:  "DEUTDEFF",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("invalid invoice with missing IBAN for SEPA credit transfer", func(t *testing.T) {
		inv := invoiceTemplate(t)
		inv.Payment = &bill.Payment{
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCreditTransfer.With(pay.MeansKeySEPA),
				CreditTransfer: []*pay.CreditTransfer{
					{
						BIC: "DEUTDEFF",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "payment: (instructions: (credit_transfer: (0: (number: cannot be blank.).).).)")
	})

	t.Run("valid invoice with card payment", func(t *testing.T) {
		inv := invoiceTemplate(t)
		inv.Payment = &bill.Payment{
			Instructions: &pay.Instructions{
				Key:  pay.MeansKeyCard,
				Card: &pay.Card{},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("valid invoice with SEPA direct debit", func(t *testing.T) {
		inv := invoiceTemplate(t)
		inv.Payment = &bill.Payment{
			Instructions: &pay.Instructions{
				Key: "direct-debit+sepa",
				DirectDebit: &pay.DirectDebit{
					Ref:      "MANDATE123",
					Creditor: "DE98ZZZ09999999999",
					Account:  "DE89370400440532013000",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("invalid invoice with missing mandate reference for direct debit", func(t *testing.T) {
		inv := invoiceTemplate(t)
		inv.Payment = &bill.Payment{
			Instructions: &pay.Instructions{
				Key: "direct-debit+sepa",
				DirectDebit: &pay.DirectDebit{
					Creditor: "DE98ZZZ09999999999",
					Account:  "DE89370400440532013000",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "payment: (instructions: (direct_debit: (ref: cannot be blank.).).)")
	})

	t.Run("invalid invoice with invalid payment key", func(t *testing.T) {
		inv := invoiceTemplate(t)
		inv.Payment = &bill.Payment{
			Instructions: &pay.Instructions{
				Key: cbc.Key("invalid-key"),
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "payment: (instructions: (key: must be or start with a valid key.).)")
	})
}
