package bill_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/stretchr/testify/assert"
)

func TestPaymentTypeIn(t *testing.T) {
	test := bill.PaymentTypeIn(bill.PaymentTypeReceipt, bill.PaymentTypeAdvice)

	t.Run("matches", func(t *testing.T) {
		assert.True(t, test.Check(&bill.Payment{Type: bill.PaymentTypeReceipt}))
		assert.True(t, test.Check(&bill.Payment{Type: bill.PaymentTypeAdvice}))
	})
	t.Run("no match", func(t *testing.T) {
		assert.False(t, test.Check(&bill.Payment{Type: bill.PaymentTypeRequest}))
	})
	t.Run("wrong type", func(t *testing.T) {
		assert.False(t, test.Check("not-a-payment"))
	})
	t.Run("nil payment", func(t *testing.T) {
		assert.False(t, test.Check((*bill.Payment)(nil)))
	})
	t.Run("string", func(t *testing.T) {
		assert.Equal(t, "payment type in [receipt, advice]", test.String())
	})
}

func TestStatusTypeIn(t *testing.T) {
	test := bill.StatusTypeIn(bill.StatusTypeResponse, bill.StatusTypeUpdate)

	t.Run("matches", func(t *testing.T) {
		assert.True(t, test.Check(&bill.Status{Type: bill.StatusTypeResponse}))
		assert.True(t, test.Check(&bill.Status{Type: bill.StatusTypeUpdate}))
	})
	t.Run("no match", func(t *testing.T) {
		assert.False(t, test.Check(&bill.Status{Type: bill.StatusTypeSystem}))
	})
	t.Run("wrong type", func(t *testing.T) {
		assert.False(t, test.Check(42))
	})
	t.Run("nil status", func(t *testing.T) {
		assert.False(t, test.Check((*bill.Status)(nil)))
	})
	t.Run("string", func(t *testing.T) {
		assert.Equal(t, "status type in [response, update]", test.String())
	})
}

func TestStatusLineKeyIn(t *testing.T) {
	test := bill.StatusLineKeyIn(bill.StatusLineRejected, bill.StatusLineQuerying, bill.StatusLineError)

	t.Run("matches", func(t *testing.T) {
		assert.True(t, test.Check(&bill.StatusLine{Key: bill.StatusLineRejected}))
		assert.True(t, test.Check(&bill.StatusLine{Key: bill.StatusLineError}))
	})
	t.Run("no match", func(t *testing.T) {
		assert.False(t, test.Check(&bill.StatusLine{Key: bill.StatusLineAccepted}))
	})
	t.Run("wrong type", func(t *testing.T) {
		assert.False(t, test.Check("nope"))
	})
	t.Run("nil line", func(t *testing.T) {
		assert.False(t, test.Check((*bill.StatusLine)(nil)))
	})
	t.Run("string", func(t *testing.T) {
		assert.Equal(t, "status line key in [rejected, querying, error]", test.String())
	})
}
