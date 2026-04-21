package bis

import (
	"testing"

	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func payExt(code cbc.Code) tax.Extensions {
	return tax.Extensions{untdid.ExtKeyPaymentMeans: code}
}

func TestDirectDebitMandatePresent(t *testing.T) {
	t.Run("nil/wrong type passes", func(t *testing.T) {
		assert.True(t, directDebitMandatePresent(nil))
		assert.True(t, directDebitMandatePresent("x"))
	})
	t.Run("non-direct-debit code passes", func(t *testing.T) {
		assert.True(t, directDebitMandatePresent(&pay.Instructions{Ext: payExt("30")}))
	})
	t.Run("code 49 without mandate fails", func(t *testing.T) {
		assert.False(t, directDebitMandatePresent(&pay.Instructions{Ext: payExt("49")}))
	})
	t.Run("code 49 with mandate passes", func(t *testing.T) {
		assert.True(t, directDebitMandatePresent(&pay.Instructions{
			Ext:         payExt("49"),
			DirectDebit: &pay.DirectDebit{Ref: "MANDATE-001"},
		}))
	})
	t.Run("code 59 without DirectDebit fails", func(t *testing.T) {
		assert.False(t, directDebitMandatePresent(&pay.Instructions{Ext: payExt("59")}))
	})
	t.Run("code 59 with empty ref fails", func(t *testing.T) {
		assert.False(t, directDebitMandatePresent(&pay.Instructions{
			Ext:         payExt("59"),
			DirectDebit: &pay.DirectDebit{},
		}))
	})
}
