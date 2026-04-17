package zatca_test

import (
	"testing"

	"github.com/invopop/gobl/addons/sa/zatca"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPayInstructionsNormalize(t *testing.T) {
	ad := tax.AddonForKey(zatca.V1)

	t.Run("nil does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			ad.Normalizer((*pay.Instructions)(nil))
		})
	})

	t.Run("credit-transfer maps to 30", func(t *testing.T) {
		m := &pay.Instructions{
			Key: pay.MeansKeyCreditTransfer,
		}
		ad.Normalizer(m)
		assert.Equal(t, "30", m.Ext[untdid.ExtKeyPaymentMeans].String())
	})

	t.Run("card maps to 48", func(t *testing.T) {
		m := &pay.Instructions{
			Key: pay.MeansKeyCard,
		}
		ad.Normalizer(m)
		assert.Equal(t, "48", m.Ext[untdid.ExtKeyPaymentMeans].String())
	})

	t.Run("cash maps to 10", func(t *testing.T) {
		m := &pay.Instructions{
			Key: pay.MeansKeyCash,
		}
		ad.Normalizer(m)
		assert.Equal(t, "10", m.Ext[untdid.ExtKeyPaymentMeans].String())
	})

	t.Run("unknown key leaves ext unchanged", func(t *testing.T) {
		m := &pay.Instructions{
			Key: "unknown-method",
		}
		ad.Normalizer(m)
		assert.Empty(t, m.Ext)
	})
}

func TestPayInstructionsRules(t *testing.T) {
	t.Run("valid instructions with ext passes", func(t *testing.T) {
		inv := validStandardInvoice()
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "30", inv.Payment.Instructions.Ext[untdid.ExtKeyPaymentMeans].String())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("missing payment means ext fails", func(t *testing.T) {
		inv := validStandardInvoice()
		require.NoError(t, inv.Calculate())
		delete(inv.Payment.Instructions.Ext, untdid.ExtKeyPaymentMeans)
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "payment means extension is required (BR-49)")
	})
}