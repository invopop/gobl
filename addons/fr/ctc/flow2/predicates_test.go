package flow2

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestPrecedingDocCodeValid(t *testing.T) {
	t.Run("wrong type passes", func(t *testing.T) {
		assert.True(t, precedingDocCodeValid("nope"))
	})
	t.Run("nil and empty code pass", func(t *testing.T) {
		assert.True(t, precedingDocCodeValid((*org.DocumentRef)(nil)))
		assert.True(t, precedingDocCodeValid(&org.DocumentRef{}))
	})
	t.Run("valid code", func(t *testing.T) {
		assert.True(t, precedingDocCodeValid(&org.DocumentRef{Code: "INV-2026-001"}))
	})
	t.Run("valid code with series", func(t *testing.T) {
		assert.True(t, precedingDocCodeValid(&org.DocumentRef{Series: "A", Code: "001"}))
	})
	t.Run("invalid characters", func(t *testing.T) {
		assert.False(t, precedingDocCodeValid(&org.DocumentRef{Code: "INV 2026 001"}))
	})
}

func TestNotesHaveTXD(t *testing.T) {
	txd := &org.Note{
		Text: stcMembreAssujettiUnique,
		Ext:  tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: noteSubjectTXD}),
	}
	t.Run("wrong type / empty", func(t *testing.T) {
		assert.False(t, notesHaveTXD("nope"))
		assert.False(t, notesHaveTXD([]*org.Note{}))
	})
	t.Run("present", func(t *testing.T) {
		assert.True(t, notesHaveTXD([]*org.Note{txd}))
	})
	t.Run("nil entry and wrong subject skipped", func(t *testing.T) {
		other := &org.Note{Text: "x", Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: "PMT"})}
		assert.False(t, notesHaveTXD([]*org.Note{nil, other}))
	})
	t.Run("right subject wrong text", func(t *testing.T) {
		bad := &org.Note{Text: "other", Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: noteSubjectTXD})}
		assert.False(t, notesHaveTXD([]*org.Note{bad}))
	})
}

func TestFinalInvoiceAdvancesMatch(t *testing.T) {
	t.Run("wrong type / nil / no advances pass", func(t *testing.T) {
		assert.True(t, finalInvoiceAdvancesMatch("nope"))
		assert.True(t, finalInvoiceAdvancesMatch((*bill.Totals)(nil)))
		assert.True(t, finalInvoiceAdvancesMatch(&bill.Totals{}))
	})
	t.Run("matches total with tax", func(t *testing.T) {
		amt := num.MakeAmount(12000, 2)
		tot := &bill.Totals{Advances: &amt, TotalWithTax: num.MakeAmount(12000, 2)}
		assert.True(t, finalInvoiceAdvancesMatch(tot))
	})
	t.Run("does not match", func(t *testing.T) {
		amt := num.MakeAmount(5000, 2)
		tot := &bill.Totals{Advances: &amt, TotalWithTax: num.MakeAmount(12000, 2)}
		assert.False(t, finalInvoiceAdvancesMatch(tot))
	})
}

func TestIdentitiesNoDupExt(t *testing.T) {
	test := identitiesNoDupExt(untdid.ExtKeyReference, "AFL")
	mk := func(ref cbc.Code) *org.Identity {
		return &org.Identity{Code: "X", Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyReference: ref})}
	}
	t.Run("wrong type passes", func(t *testing.T) {
		assert.True(t, test.Check("nope"))
	})
	t.Run("none with the ref", func(t *testing.T) {
		assert.True(t, test.Check([]*org.Identity{mk("AWW"), nil, {Code: "Y"}}))
	})
	t.Run("exactly one", func(t *testing.T) {
		assert.True(t, test.Check([]*org.Identity{mk("AFL"), mk("AWW")}))
	})
	t.Run("duplicate rejected", func(t *testing.T) {
		assert.False(t, test.Check([]*org.Identity{mk("AFL"), mk("AFL")}))
	})
	t.Run("string", func(t *testing.T) {
		assert.Equal(t, "identities have at most one ext untdid-reference=AFL", test.String())
	})
}

func TestFinalInvoicePayableZero(t *testing.T) {
	t.Run("wrong type / nil pass", func(t *testing.T) {
		assert.True(t, finalInvoicePayableZero("nope"))
		assert.True(t, finalInvoicePayableZero((*bill.Totals)(nil)))
	})
	t.Run("payable zero", func(t *testing.T) {
		assert.True(t, finalInvoicePayableZero(&bill.Totals{Payable: num.AmountZero}))
	})
	t.Run("payable non-zero", func(t *testing.T) {
		assert.False(t, finalInvoicePayableZero(&bill.Totals{Payable: num.MakeAmount(100, 2)}))
	})
	t.Run("due takes precedence when present", func(t *testing.T) {
		due := num.MakeAmount(0, 2)
		assert.True(t, finalInvoicePayableZero(&bill.Totals{Due: &due, Payable: num.MakeAmount(100, 2)}))
		due2 := num.MakeAmount(100, 2)
		assert.False(t, finalInvoicePayableZero(&bill.Totals{Due: &due2}))
	})
}
