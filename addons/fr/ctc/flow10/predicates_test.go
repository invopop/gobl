package flow10

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestPartyHasVATCode(t *testing.T) {
	assert.False(t, partyHasVATCode(nil))
	assert.False(t, partyHasVATCode(&org.Party{}))
	assert.False(t, partyHasVATCode(&org.Party{TaxID: &tax.Identity{Country: "FR"}}))
	assert.True(t, partyHasVATCode(&org.Party{TaxID: &tax.Identity{Country: "FR", Code: "44732829320"}}))
}

func TestInvoiceHasSellerVATIDForExempt(t *testing.T) {
	t.Run("wrong type / nil", func(t *testing.T) {
		assert.False(t, invoiceHasSellerVATIDForExempt("nope"))
		assert.False(t, invoiceHasSellerVATIDForExempt((*bill.Invoice)(nil)))
	})
	t.Run("supplier carries VAT code", func(t *testing.T) {
		inv := &bill.Invoice{Supplier: &org.Party{TaxID: &tax.Identity{Country: "FR", Code: "44732829320"}}}
		assert.True(t, invoiceHasSellerVATIDForExempt(inv))
	})
	t.Run("tax representative (ordering.seller) carries VAT code", func(t *testing.T) {
		inv := &bill.Invoice{
			Supplier: &org.Party{},
			Ordering: &bill.Ordering{Seller: &org.Party{TaxID: &tax.Identity{Country: "FR", Code: "44732829320"}}},
		}
		assert.True(t, invoiceHasSellerVATIDForExempt(inv))
	})
	t.Run("neither has a VAT code", func(t *testing.T) {
		inv := &bill.Invoice{Supplier: &org.Party{}}
		assert.False(t, invoiceHasSellerVATIDForExempt(inv))
	})
}

func TestInvoiceHasExemptTaxNote(t *testing.T) {
	t.Run("wrong type / nil / no tax", func(t *testing.T) {
		assert.False(t, invoiceHasExemptTaxNote("nope"))
		assert.False(t, invoiceHasExemptTaxNote((*bill.Invoice)(nil)))
		assert.False(t, invoiceHasExemptTaxNote(&bill.Invoice{}))
	})
	t.Run("has exempt note with text", func(t *testing.T) {
		inv := &bill.Invoice{Tax: &bill.Tax{Notes: []*tax.Note{
			{Key: tax.KeyExempt, Text: "Exonération de TVA"},
		}}}
		assert.True(t, invoiceHasExemptTaxNote(inv))
	})
	t.Run("exempt note without text", func(t *testing.T) {
		inv := &bill.Invoice{Tax: &bill.Tax{Notes: []*tax.Note{
			{Key: tax.KeyExempt},
		}}}
		assert.False(t, invoiceHasExemptTaxNote(inv))
	})
	t.Run("non-exempt note ignored", func(t *testing.T) {
		inv := &bill.Invoice{Tax: &bill.Tax{Notes: []*tax.Note{
			{Key: "other", Text: "x"},
		}}}
		assert.False(t, invoiceHasExemptTaxNote(inv))
	})
}
