package xrechnung_test

import (
	"testing"

	"github.com/invopop/gobl/addons/de/xrechnung"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/stretchr/testify/assert"
)

func TestPaymentInstructions(t *testing.T) {
	t.Run("valid SEPA credit transfer", func(t *testing.T) {
		inv := &bill.Invoice{
			Payment: &bill.Payment{
				Instructions: &pay.Instructions{
					Key: xrechnung.KeyPaymentMeansSEPACreditTransfer,
					CreditTransfer: []*pay.CreditTransfer{
						{
							IBAN: "DE89370400440532013000",
							BIC:  "DEUTDEFF",
						},
					},
				},
			},
		}
		assert.NoError(t, xrechnung.ValidatePaymentInstructions(inv))
	})

	t.Run("missing IBAN for SEPA credit transfer", func(t *testing.T) {
		inv := &bill.Invoice{
			Payment: &bill.Payment{
				Instructions: &pay.Instructions{
					Key: xrechnung.KeyPaymentMeansSEPACreditTransfer,
					CreditTransfer: []*pay.CreditTransfer{
						{
							BIC: "DEUTDEFF",
						},
					},
				},
			},
		}
		assert.Error(t, xrechnung.ValidatePaymentInstructions(inv))
	})

	t.Run("valid card payment", func(t *testing.T) {
		inv := &bill.Invoice{
			Payment: &bill.Payment{
				Instructions: &pay.Instructions{
					Key:  pay.MeansKeyCard,
					Card: &pay.Card{},
				},
			},
		}
		assert.NoError(t, xrechnung.ValidatePaymentInstructions(inv))
	})

	t.Run("valid SEPA direct debit", func(t *testing.T) {
		inv := &bill.Invoice{
			Payment: &bill.Payment{
				Instructions: &pay.Instructions{
					Key: xrechnung.KeyPaymentMeansSEPADirectDebit,
					DirectDebit: &pay.DirectDebit{
						Ref:      "MANDATE123",
						Creditor: "DE98ZZZ09999999999",
						Account:  "DE89370400440532013000",
					},
				},
			},
		}
		assert.NoError(t, xrechnung.ValidatePaymentInstructions(inv))
	})

	t.Run("missing mandate reference for direct debit", func(t *testing.T) {
		inv := &bill.Invoice{
			Payment: &bill.Payment{
				Instructions: &pay.Instructions{
					Key: xrechnung.KeyPaymentMeansSEPADirectDebit,
					DirectDebit: &pay.DirectDebit{
						Creditor: "DE98ZZZ09999999999",
						Account:  "DE89370400440532013000",
					},
				},
			},
		}
		assert.Error(t, xrechnung.ValidatePaymentInstructions(inv))
	})

	t.Run("invalid payment key", func(t *testing.T) {
		inv := &bill.Invoice{
			Payment: &bill.Payment{
				Instructions: &pay.Instructions{
					Key: cbc.Key("invalid-key"),
				},
			},
		}
		assert.Error(t, xrechnung.ValidatePaymentInstructions(inv))
	})
}
