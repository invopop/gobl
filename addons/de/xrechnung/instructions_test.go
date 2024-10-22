package xrechnung_test

import (
	"testing"

	"github.com/invopop/gobl/addons/de/xrechnung"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/stretchr/testify/assert"
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
				Key: cbc.Key("sepa-credit-transfer"),
				CreditTransfer: []*pay.CreditTransfer{
					{
						IBAN: "DE89370400440532013000",
						BIC:  "DEUTDEFF",
					},
				},
			},
		}
		assert.NoError(t, xrechnung.ValidateInvoice(inv))
	})

	t.Run("invalid invoice with missing IBAN for SEPA credit transfer", func(t *testing.T) {
		inv := invoiceTemplate(t)
		inv.Payment = &bill.Payment{
			Instructions: &pay.Instructions{
				Key: xrechnung.KeyPaymentMeansSEPACreditTransfer,
				CreditTransfer: []*pay.CreditTransfer{
					{
						BIC: "DEUTDEFF",
					},
				},
			},
		}
		assert.Error(t, xrechnung.ValidateInvoice(inv))
	})

	t.Run("valid invoice with card payment", func(t *testing.T) {
		inv := invoiceTemplate(t)
		inv.Payment = &bill.Payment{
			Instructions: &pay.Instructions{
				Key:  pay.MeansKeyCard,
				Card: &pay.Card{},
			},
		}
		assert.NoError(t, xrechnung.ValidateInvoice(inv))
	})

	t.Run("valid invoice with SEPA direct debit", func(t *testing.T) {
		inv := invoiceTemplate(t)
		inv.Payment = &bill.Payment{
			Instructions: &pay.Instructions{
				Key: xrechnung.KeyPaymentMeansSEPADirectDebit,
				DirectDebit: &pay.DirectDebit{
					Ref:      "MANDATE123",
					Creditor: "DE98ZZZ09999999999",
					Account:  "DE89370400440532013000",
				},
			},
		}
		assert.NoError(t, xrechnung.ValidateInvoice(inv))
	})

	t.Run("invalid invoice with missing mandate reference for direct debit", func(t *testing.T) {
		inv := invoiceTemplate(t)
		inv.Payment = &bill.Payment{
			Instructions: &pay.Instructions{
				Key: xrechnung.KeyPaymentMeansSEPADirectDebit,
				DirectDebit: &pay.DirectDebit{
					Creditor: "DE98ZZZ09999999999",
					Account:  "DE89370400440532013000",
				},
			},
		}
		assert.Error(t, xrechnung.ValidateInvoice(inv))
	})

	t.Run("invalid invoice with invalid payment key", func(t *testing.T) {
		inv := invoiceTemplate(t)
		inv.Payment = &bill.Payment{
			Instructions: &pay.Instructions{
				Key: cbc.Key("invalid-key"),
			},
		}
		assert.Error(t, xrechnung.ValidateInvoice(inv))
	})
}
