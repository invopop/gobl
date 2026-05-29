package flow2

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/catalogues/dgfip"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeBillingMode(t *testing.T) {
	t.Run("keeps caller-supplied mode", func(t *testing.T) {
		inv := &bill.Invoice{Tax: &bill.Tax{
			Ext: tax.ExtensionsOf(cbc.CodeMap{dgfip.ExtKeyBillingMode: dgfip.BillingModeB2}),
		}}
		normalizeBillingMode(inv)
		assert.Equal(t, dgfip.BillingModeB2, inv.Tax.Ext.Get(dgfip.ExtKeyBillingMode))
	})

	t.Run("defaults to M1 when unpaid", func(t *testing.T) {
		inv := &bill.Invoice{}
		normalizeBillingMode(inv)
		require.NotNil(t, inv.Tax)
		assert.Equal(t, dgfip.BillingModeM1, inv.Tax.Ext.Get(dgfip.ExtKeyBillingMode))
	})

	t.Run("defaults to M2 when fully paid", func(t *testing.T) {
		due := num.MakeAmount(0, 2)
		inv := &bill.Invoice{Totals: &bill.Totals{
			TotalWithTax: num.MakeAmount(12100, 2),
			Due:          &due,
		}}
		normalizeBillingMode(inv)
		assert.Equal(t, dgfip.BillingModeM2, inv.Tax.Ext.Get(dgfip.ExtKeyBillingMode))
	})
}

func TestNormalizeSTCNote(t *testing.T) {
	stcSupplier := func() *org.Party {
		return &org.Party{Identities: []*org.Identity{
			{Code: "1", Ext: tax.ExtensionsOf(cbc.CodeMap{iso.ExtKeySchemeID: "0231"})},
		}}
	}

	t.Run("non-STC supplier is a no-op", func(t *testing.T) {
		inv := &bill.Invoice{Supplier: &org.Party{}}
		normalizeSTCNote(inv)
		assert.Empty(t, inv.Notes)
	})

	t.Run("adds TXD note for STC supplier", func(t *testing.T) {
		inv := &bill.Invoice{Supplier: stcSupplier()}
		normalizeSTCNote(inv)
		require.Len(t, inv.Notes, 1)
		assert.Equal(t, stcMembreAssujettiUnique, inv.Notes[0].Text)
		assert.Equal(t, noteSubjectTXD, inv.Notes[0].Ext.Get("untdid-text-subject"))
	})

	t.Run("does not duplicate an existing TXD note", func(t *testing.T) {
		inv := &bill.Invoice{
			Supplier: stcSupplier(),
			Notes: []*org.Note{{
				Text: stcMembreAssujettiUnique,
				Ext:  tax.ExtensionsOf(cbc.CodeMap{"untdid-text-subject": noteSubjectTXD}),
			}},
		}
		normalizeSTCNote(inv)
		assert.Len(t, inv.Notes, 1)
	})
}

func TestAttachmentsAtMostOneLISIBLE(t *testing.T) {
	assert.True(t, attachmentsAtMostOneLISIBLE("wrong-type"))
	assert.True(t, attachmentsAtMostOneLISIBLE([]*org.Attachment{}))
	assert.True(t, attachmentsAtMostOneLISIBLE([]*org.Attachment{
		{Description: attachmentFormatLisible}, nil, {Description: "RIB"},
	}))
	assert.False(t, attachmentsAtMostOneLISIBLE([]*org.Attachment{
		{Description: attachmentFormatLisible}, {Description: attachmentFormatLisible},
	}))
}

func TestInvoiceCodeValid(t *testing.T) {
	assert.True(t, invoiceCodeValid("wrong-type"))
	assert.True(t, invoiceCodeValid((*bill.Invoice)(nil)))
	assert.True(t, invoiceCodeValid(&bill.Invoice{}))
	assert.True(t, invoiceCodeValid(&bill.Invoice{Code: "2024-00001"}))
	assert.True(t, invoiceCodeValid(&bill.Invoice{Series: "FAC", Code: "001"}))
	assert.False(t, invoiceCodeValid(&bill.Invoice{Code: "bad code!"}))
}

func TestInvoiceDueDatesValid(t *testing.T) {
	assert.True(t, invoiceDueDatesValid("wrong-type"))
	assert.True(t, invoiceDueDatesValid((*bill.Invoice)(nil)))
	// no payment terms → passes
	assert.True(t, invoiceDueDatesValid(&bill.Invoice{}))

	issue := cal.MakeDate(2024, 6, 13)
	mk := func(d *cal.Date) *bill.Invoice {
		return &bill.Invoice{
			IssueDate: issue,
			Payment:   &bill.PaymentDetails{Terms: &pay.Terms{DueDates: []*pay.DueDate{{Date: d}}}},
		}
	}
	t.Run("due date after issue passes", func(t *testing.T) {
		assert.True(t, invoiceDueDatesValid(mk(cal.NewDate(2024, 7, 13))))
	})
	t.Run("due date before issue fails", func(t *testing.T) {
		assert.False(t, invoiceDueDatesValid(mk(cal.NewDate(2024, 5, 13))))
	})
	t.Run("nil due date entry skipped", func(t *testing.T) {
		inv := &bill.Invoice{
			IssueDate: issue,
			Payment:   &bill.PaymentDetails{Terms: &pay.Terms{DueDates: []*pay.DueDate{nil, {}}}},
		}
		assert.True(t, invoiceDueDatesValid(inv))
	})
}

func TestNotesHaveRequired(t *testing.T) {
	mk := func(subjects ...cbc.Code) []*org.Note {
		notes := make([]*org.Note, 0, len(subjects))
		for _, s := range subjects {
			notes = append(notes, &org.Note{Ext: tax.ExtensionsOf(cbc.CodeMap{"untdid-text-subject": s})})
		}
		return notes
	}
	assert.False(t, notesHaveRequired("wrong-type"))
	assert.True(t, notesHaveRequired(mk("PMT", "PMD", "AAB")))
	assert.False(t, notesHaveRequired(mk("PMT", "PMD")))
	// nil + zero-ext entries are skipped without panic
	assert.False(t, notesHaveRequired([]*org.Note{nil, {}}))
}

func TestNotesNoDuplicates(t *testing.T) {
	mk := func(subjects ...cbc.Code) []*org.Note {
		notes := make([]*org.Note, 0, len(subjects))
		for _, s := range subjects {
			notes = append(notes, &org.Note{Ext: tax.ExtensionsOf(cbc.CodeMap{"untdid-text-subject": s})})
		}
		return notes
	}
	assert.True(t, notesNoDuplicates("wrong-type"))
	assert.True(t, notesNoDuplicates([]*org.Note{}))
	assert.True(t, notesNoDuplicates(mk("PMT", "PMD", "AAB", "TXD")))
	assert.False(t, notesNoDuplicates(mk("PMT", "PMT")))
}

func TestInvoiceTaxExtInGuards(t *testing.T) {
	inv := &bill.Invoice{Tax: &bill.Tax{
		Ext: tax.ExtensionsOf(cbc.CodeMap{dgfip.ExtKeyBillingMode: dgfip.BillingModeB4}),
	}}
	in := invoiceTaxExtIn(dgfip.ExtKeyBillingMode, dgfip.BillingModeB4)
	notIn := invoiceTaxExtNotIn(dgfip.ExtKeyBillingMode, dgfip.BillingModeB4)

	assert.True(t, in.Check(inv))
	assert.False(t, notIn.Check(inv))

	// wrong type / missing tax
	assert.False(t, in.Check("wrong-type"))
	assert.False(t, in.Check(&bill.Invoice{}))
	assert.True(t, notIn.Check("wrong-type"))
	assert.True(t, notIn.Check(&bill.Invoice{}))
}
