package en16931_test

import (
	"testing"
	"time"

	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPayInstructions(t *testing.T) {
	ad := tax.AddonForKey(en16931.V2017)

	t.Run("valid", func(t *testing.T) {
		m := &pay.Instructions{
			Key: pay.MeansKeyCreditTransfer,
		}
		ad.Normalizer(m)
		assert.Equal(t, "30", m.Ext.Get(untdid.ExtKeyPaymentMeans).String())
	})

	t.Run("bare card maps to 48", func(t *testing.T) {
		m := &pay.Instructions{Key: pay.MeansKeyCard}
		ad.Normalizer(m)
		assert.Equal(t, "48", m.Ext.Get(untdid.ExtKeyPaymentMeans).String())
	})

	t.Run("card+credit maps to 54", func(t *testing.T) {
		m := &pay.Instructions{Key: pay.MeansKeyCard.With(pay.MeansKeyCredit)}
		ad.Normalizer(m)
		assert.Equal(t, "54", m.Ext.Get(untdid.ExtKeyPaymentMeans).String())
	})

	t.Run("card+debit maps to 55", func(t *testing.T) {
		m := &pay.Instructions{Key: pay.MeansKeyCard.With(pay.MeansKeyDebit)}
		ad.Normalizer(m)
		assert.Equal(t, "55", m.Ext.Get(untdid.ExtKeyPaymentMeans).String())
	})

	t.Run("nil", func(t *testing.T) {
		var m *pay.Instructions
		assert.NotPanics(t, func() {
			ad.Normalizer(m)
		})
	})

	t.Run("validation", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCreditTransfer,
			},
			Terms: &pay.Terms{
				DueDates: []*pay.DueDate{
					{
						Date:   cal.NewDate(2025, time.January, 1),
						Amount: num.MakeAmount(1000, 2),
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "30", inv.Payment.Instructions.Ext.Get(untdid.ExtKeyPaymentMeans).String())
		err := rules.Validate(inv)
		assert.NoError(t, err)
	})
}

func TestPayTerms(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		p := &pay.Terms{
			DueDates: []*pay.DueDate{
				{
					Date:   cal.NewDate(2025, time.January, 1),
					Amount: num.MakeAmount(1000, 2),
				},
			},
		}
		err := rules.Validate(p, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})

	t.Run("with empty due date", func(t *testing.T) {
		p := &pay.Terms{
			DueDates: []*pay.DueDate{},
		}
		err := rules.Validate(p, tax.AddonContext(en16931.V2017))
		assert.ErrorContains(t, err, "either due_dates or notes must be provided")
	})

	t.Run("with empty notes", func(t *testing.T) {
		p := &pay.Terms{
			Notes: "",
		}
		err := rules.Validate(p, tax.AddonContext(en16931.V2017))
		assert.ErrorContains(t, err, "either due_dates or notes must be provided")
	})
}
