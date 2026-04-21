package bis

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func taxExt(code cbc.Code) *bill.Tax {
	if code == "" {
		return &bill.Tax{}
	}
	return &bill.Tax{Ext: tax.Extensions{untdid.ExtKeyDocumentType: code}}
}

func partyInCountry(c l10n.Code) *org.Party {
	return &org.Party{TaxID: &tax.Identity{Country: l10n.TaxCountryCode(c)}}
}

func TestNotesCardinalityValid(t *testing.T) {
	t.Run("nil/wrong type passes", func(t *testing.T) {
		assert.True(t, notesCardinalityValid(nil))
		assert.True(t, notesCardinalityValid("x"))
	})
	t.Run("zero or one note passes", func(t *testing.T) {
		assert.True(t, notesCardinalityValid(&bill.Invoice{}))
		assert.True(t, notesCardinalityValid(&bill.Invoice{Notes: []*org.Note{{Text: "a"}}}))
	})
	t.Run("multiple notes fails for non-DK pair", func(t *testing.T) {
		inv := &bill.Invoice{
			Notes:    []*org.Note{{Text: "a"}, {Text: "b"}},
			Supplier: partyInCountry(l10n.SE),
			Customer: partyInCountry(l10n.SE),
		}
		assert.False(t, notesCardinalityValid(inv))
	})
	t.Run("multiple notes pass when both DK", func(t *testing.T) {
		inv := &bill.Invoice{
			Notes:    []*org.Note{{Text: "a"}, {Text: "b"}},
			Supplier: partyInCountry(l10n.DK),
			Customer: partyInCountry(l10n.DK),
		}
		assert.True(t, notesCardinalityValid(inv))
	})
}

func TestHasBuyerReferenceOrPO(t *testing.T) {
	t.Run("nil/wrong type passes", func(t *testing.T) {
		assert.True(t, hasBuyerReferenceOrPO(nil))
	})
	t.Run("no ordering fails", func(t *testing.T) {
		assert.False(t, hasBuyerReferenceOrPO(&bill.Invoice{}))
	})
	t.Run("ordering code present", func(t *testing.T) {
		assert.True(t, hasBuyerReferenceOrPO(&bill.Invoice{
			Ordering: &bill.Ordering{Code: "BR-1"},
		}))
	})
	t.Run("purchase code present", func(t *testing.T) {
		assert.True(t, hasBuyerReferenceOrPO(&bill.Invoice{
			Ordering: &bill.Ordering{Purchases: []*org.DocumentRef{{Code: "PO-1"}}},
		}))
	})
	t.Run("nil purchase entry skipped", func(t *testing.T) {
		assert.False(t, hasBuyerReferenceOrPO(&bill.Invoice{
			Ordering: &bill.Ordering{Purchases: []*org.DocumentRef{nil, {}}},
		}))
	})
}

func TestInvoiceTypeCodeValid(t *testing.T) {
	t.Run("nil passes", func(t *testing.T) {
		assert.True(t, invoiceTypeCodeValid(nil))
	})
	t.Run("credit note skipped here", func(t *testing.T) {
		inv := &bill.Invoice{Type: bill.InvoiceTypeCreditNote, Tax: taxExt("999")}
		assert.True(t, invoiceTypeCodeValid(inv))
	})
	t.Run("missing code passes", func(t *testing.T) {
		assert.True(t, invoiceTypeCodeValid(&bill.Invoice{Type: bill.InvoiceTypeStandard, Tax: &bill.Tax{}}))
	})
	t.Run("allowed code", func(t *testing.T) {
		assert.True(t, invoiceTypeCodeValid(&bill.Invoice{Type: bill.InvoiceTypeStandard, Tax: taxExt("380")}))
	})
	t.Run("disallowed code", func(t *testing.T) {
		assert.False(t, invoiceTypeCodeValid(&bill.Invoice{Type: bill.InvoiceTypeStandard, Tax: taxExt("325")}))
	})
}

func TestCreditNoteTypeCodeValid(t *testing.T) {
	t.Run("non-credit-note skipped", func(t *testing.T) {
		assert.True(t, creditNoteTypeCodeValid(&bill.Invoice{Type: bill.InvoiceTypeStandard, Tax: taxExt("380")}))
	})
	t.Run("missing code passes", func(t *testing.T) {
		assert.True(t, creditNoteTypeCodeValid(&bill.Invoice{Type: bill.InvoiceTypeCreditNote, Tax: &bill.Tax{}}))
	})
	t.Run("allowed credit-note code", func(t *testing.T) {
		assert.True(t, creditNoteTypeCodeValid(&bill.Invoice{Type: bill.InvoiceTypeCreditNote, Tax: taxExt("381")}))
	})
	t.Run("disallowed credit-note code", func(t *testing.T) {
		assert.False(t, creditNoteTypeCodeValid(&bill.Invoice{Type: bill.InvoiceTypeCreditNote, Tax: taxExt("380")}))
	})
	t.Run("nil/wrong-type passes", func(t *testing.T) {
		assert.True(t, creditNoteTypeCodeValid(nil))
	})
}

func TestPartialCorrectiveITOnly(t *testing.T) {
	t.Run("non-326/384 passes", func(t *testing.T) {
		assert.True(t, partialCorrectiveITOnly(&bill.Invoice{Tax: taxExt("380")}))
	})
	t.Run("326 with both IT", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax:      taxExt("326"),
			Supplier: partyInCountry(l10n.IT),
			Customer: partyInCountry(l10n.IT),
		}
		assert.True(t, partialCorrectiveITOnly(inv))
	})
	t.Run("326 with non-IT supplier", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax:      taxExt("326"),
			Supplier: partyInCountry(l10n.DE),
			Customer: partyInCountry(l10n.IT),
		}
		assert.False(t, partialCorrectiveITOnly(inv))
	})
	t.Run("nil/wrong-type passes", func(t *testing.T) {
		assert.True(t, partialCorrectiveITOnly(nil))
	})
	t.Run("invoice without tax passes", func(t *testing.T) {
		assert.True(t, partialCorrectiveITOnly(&bill.Invoice{}))
	})
}

func TestHasPaymentInstructions(t *testing.T) {
	t.Run("nil passes", func(t *testing.T) {
		assert.True(t, hasPaymentInstructions(nil))
	})
	t.Run("no payment fails", func(t *testing.T) {
		assert.False(t, hasPaymentInstructions(&bill.Invoice{}))
	})
	t.Run("no instructions fails", func(t *testing.T) {
		assert.False(t, hasPaymentInstructions(&bill.Invoice{Payment: &bill.PaymentDetails{}}))
	})
	t.Run("instructions present", func(t *testing.T) {
		assert.True(t, hasPaymentInstructions(&bill.Invoice{
			Payment: &bill.PaymentDetails{Instructions: &pay.Instructions{Key: "credit-transfer"}},
		}))
	})
}

func TestHasOrderingCode(t *testing.T) {
	assert.True(t, hasOrderingCode(nil))
	assert.False(t, hasOrderingCode(&bill.Invoice{}))
	assert.False(t, hasOrderingCode(&bill.Invoice{Ordering: &bill.Ordering{}}))
	assert.True(t, hasOrderingCode(&bill.Invoice{Ordering: &bill.Ordering{Code: "X"}}))
}
