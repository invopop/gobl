package tbai

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeInvoiceRegimeDefensive(t *testing.T) {
	t.Run("nil line is skipped", func(t *testing.T) {
		inv := &bill.Invoice{Lines: []*bill.Line{nil}}
		assert.NotPanics(t, func() { normalizeInvoiceRegime(inv) })
	})

	t.Run("nil tax combo is skipped", func(t *testing.T) {
		inv := &bill.Invoice{
			Lines: []*bill.Line{
				{Taxes: tax.Set{nil}},
			},
		}
		assert.NotPanics(t, func() { normalizeInvoiceRegime(inv) })
	})

	t.Run("non VAT/IGIC category is skipped", func(t *testing.T) {
		tc := &tax.Combo{Category: tax.CategoryGST}
		inv := &bill.Invoice{
			Lines: []*bill.Line{{Taxes: tax.Set{tc}}},
		}
		normalizeInvoiceRegime(inv)
		assert.True(t, tc.Ext.IsZero())
	})
}

func TestNormalizeBillLineNoVAT(t *testing.T) {
	line := &bill.Line{
		Quantity: num.MakeAmount(1, 0),
		Item:     &org.Item{Name: "x", Price: num.NewAmount(100, 2)},
		Taxes: tax.Set{
			{Category: tax.CategoryGST},
		},
	}
	assert.NotPanics(t, func() { normalizeBillLine(line) })
	assert.True(t, line.Taxes[0].Ext.IsZero())
}

func TestNotesHasGeneralKeyWrongType(t *testing.T) {
	assert.False(t, notesHasGeneralKey("not a slice"))
	assert.False(t, notesHasGeneralKey(nil))
}

func TestNotesHasGeneralKeyNoGeneralNote(t *testing.T) {
	notes := []*org.Note{
		{Key: org.NoteKeyLegal, Text: "legal"},
	}
	assert.False(t, notesHasGeneralKey(notes))
}

func TestNormalizeInvoicePartyIdentityNilCustomer(t *testing.T) {
	assert.NotPanics(t, func() { normalizeInvoicePartyIdentity(nil) })
}

func TestNormalizeInvoicePartyIdentityUnkeyedNoExt(t *testing.T) {
	cus := &org.Party{
		Identities: []*org.Identity{
			{Code: "X"},
		},
	}
	normalizeInvoicePartyIdentity(cus)
	assert.True(t, cus.Identities[0].Ext.IsZero())
}

func TestNormalizeInvoicePartyIdentitySpanishNIFShortCircuits(t *testing.T) {
	cus := &org.Party{
		TaxID: &tax.Identity{Country: "ES", Code: "B12345678"},
		Identities: []*org.Identity{
			{Key: org.IdentityKeyPassport, Code: "AA"},
		},
	}
	normalizeInvoicePartyIdentity(cus)
	assert.True(t, cus.Identities[0].Ext.IsZero())
}

func TestNormalizeInvoicePartyIdentityEmptyIdentities(t *testing.T) {
	cus := &org.Party{}
	assert.NotPanics(t, func() { normalizeInvoicePartyIdentity(cus) })
}

func TestNormalizeInvoiceNil(t *testing.T) {
	assert.NotPanics(t, func() { normalizeInvoice(nil) })
}

func TestIsBizkaiaIndividualWrongType(t *testing.T) {
	assert.False(t, isBizkaiaIndividual("not an invoice"))
	assert.False(t, isBizkaiaIndividual(nil))
}
