package cl

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestInvoiceNotSimplified(t *testing.T) {
	t.Run("nil value", func(t *testing.T) {
		got := invoiceNotSimplified(nil)
		assert.True(t, got)
	})

	t.Run("wrong type", func(t *testing.T) {
		got := invoiceNotSimplified("not an invoice")
		assert.True(t, got)
	})

	t.Run("nil invoice pointer", func(t *testing.T) {
		var inv *bill.Invoice
		got := invoiceNotSimplified(inv)
		assert.True(t, got)
	})

	t.Run("standard invoice is not simplified", func(t *testing.T) {
		inv := &bill.Invoice{}
		got := invoiceNotSimplified(inv)
		assert.True(t, got)
	})

	t.Run("simplified invoice", func(t *testing.T) {
		inv := &bill.Invoice{}
		inv.SetTags(tax.TagSimplified)
		got := invoiceNotSimplified(inv)
		assert.False(t, got)
	})
}
