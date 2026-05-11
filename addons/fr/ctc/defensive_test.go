package ctc

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

// Defensive-branch coverage: nil / zero / wrong-type / empty-slice
// inputs to the predicate helpers, so the "ok / nil" guards report as
// exercised rather than dead.

// --- bill_invoice predicates --------------------------------------------

func TestInvoiceCodeValidNonInvoice(t *testing.T) {
	assert.True(t, invoiceCodeValid(42))
}

func TestInvoiceCodeValidEmptyCode(t *testing.T) {
	assert.True(t, invoiceCodeValid(&bill.Invoice{}))
}

func TestPrecedingDocCodeValidNonDocumentRef(t *testing.T) {
	assert.True(t, precedingDocCodeValid(42))
}

func TestPrecedingDocCodeValidNil(t *testing.T) {
	assert.True(t, precedingDocCodeValid((*org.DocumentRef)(nil)))
}

func TestInvoiceIsFactoringAnyNonInvoice(t *testing.T) {
	assert.False(t, invoiceIsFactoringAny(42))
}

func TestInvoiceIsFactoringAnyEmptyTax(t *testing.T) {
	assert.False(t, invoiceIsFactoringAny(&bill.Invoice{}))
}

func TestIsCorrectiveInvoiceEmptyExt(t *testing.T) {
	assert.False(t, isCorrectiveInvoice(&bill.Invoice{Tax: &bill.Tax{}}))
}

func TestIsCreditNoteEmptyExt(t *testing.T) {
	assert.False(t, isCreditNote(&bill.Invoice{Tax: &bill.Tax{}}))
}

func TestIsConsolidatedCreditNoteEmptyExt(t *testing.T) {
	assert.False(t, isConsolidatedCreditNote(&bill.Invoice{Tax: &bill.Tax{}}))
}

func TestIsAdvancedInvoiceEmptyExt(t *testing.T) {
	assert.False(t, isAdvancedInvoice(&bill.Invoice{Tax: &bill.Tax{}}))
}

func TestIsFinalInvoiceEmptyExt(t *testing.T) {
	assert.False(t, isFinalInvoice(&bill.Invoice{Tax: &bill.Tax{}}))
}

func TestIsPartyIdentitySTCNilIdentity(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{nil}}
	assert.False(t, isPartyIdentitySTC(p))
}

func TestIsPartyIdentitySTCEmptyExt(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{{Code: "X"}}}
	assert.False(t, isPartyIdentitySTC(p))
}

func TestGetPartySIRENEmpty(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{{Code: "X"}}}
	assert.Equal(t, "", getPartySIREN(p))
}

func TestIdentitiesHasLegalSIRENNilEntry(t *testing.T) {
	assert.False(t, identitiesHasLegalSIREN([]*org.Identity{nil}))
}

func TestPartyHasSIRENInboxNoSIREN(t *testing.T) {
	p := &org.Party{Inboxes: []*org.Inbox{{Scheme: inboxSchemeSIREN, Code: "X"}}}
	assert.True(t, partyHasSIRENInbox(p))
}

func TestOrderingIdentitiesNoDupWrongType(t *testing.T) {
	assert.True(t, orderingIdentitiesNoDup("x", "AFL"))
}

func TestOrderingIdentitiesNoDupNilEntry(t *testing.T) {
	assert.True(t, orderingIdentitiesNoDup([]*org.Identity{nil}, "AFL"))
}

func TestNotesHaveRequiredEmpty(t *testing.T) {
	assert.False(t, notesHaveRequired([]*org.Note{}))
}

func TestNotesHaveRequiredNilEntry(t *testing.T) {
	assert.False(t, notesHaveRequired([]*org.Note{nil}))
}

func TestInvoiceHasNoteWithSubjectNilNote(t *testing.T) {
	inv := &bill.Invoice{Notes: []*org.Note{nil}}
	assert.False(t, invoiceHasNoteWithSubject(inv, "PMT"))
}

func TestNormalizeRequiredNotesNoOpWhenPresent(t *testing.T) {
	inv := &bill.Invoice{
		Notes: []*org.Note{
			{Key: org.NoteKeyPayment, Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: "PMT"})},
			{Key: org.NoteKeyPaymentMethod, Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: "PMD"})},
			{Key: org.NoteKeyPaymentTerm, Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: "AAB"})},
		},
	}
	before := len(inv.Notes)
	normalizeRequiredNotes(inv)
	assert.Equal(t, before, len(inv.Notes))
}

func TestNormalizeB2CCategoryOnInvoicePreservesExisting(t *testing.T) {
	inv := &bill.Invoice{Tax: &bill.Tax{
		Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyB2CCategory: B2CCategoryGoods}),
	}}
	normalizeB2CCategoryOnInvoice(inv)
	assert.Equal(t, B2CCategoryGoods, inv.Tax.Ext.Get(ExtKeyB2CCategory))
}

func TestNormalizeInvoiceTaxCategoriesNilLine(t *testing.T) {
	inv := &bill.Invoice{Lines: []*bill.Line{nil}}
	assert.NotPanics(t, func() { normalizeInvoiceTaxCategories(inv) })
}

func TestNormalizeInvoiceTaxCategoriesNilCombo(t *testing.T) {
	inv := &bill.Invoice{Lines: []*bill.Line{{Taxes: tax.Set{nil}}}}
	assert.NotPanics(t, func() { normalizeInvoiceTaxCategories(inv) })
}

// --- bill_status predicates ---------------------------------------------

func TestSetPartyRoleDefaultNilParty(t *testing.T) {
	assert.NotPanics(t, func() { setPartyRoleDefault(nil, RoleSE) })
}

func TestSetPartyRoleDefaultExistingNotOverridden(t *testing.T) {
	p := &org.Party{Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyRole: RoleBY})}
	setPartyRoleDefault(p, RoleSE)
	assert.Equal(t, RoleBY, p.Ext.Get(ExtKeyRole))
}

func TestPartyHasRoleWrongType(t *testing.T) {
	assert.False(t, partyHasRole("x"))
}

func TestPartyHasRoleEmptyExt(t *testing.T) {
	assert.False(t, partyHasRole(&org.Party{}))
}

func TestPartyHasInboxWhenRequiredWrongType(t *testing.T) {
	assert.True(t, partyHasInboxWhenRequired("x"))
}

func TestPartyHasInboxWhenRequiredWKRole(t *testing.T) {
	p := &org.Party{Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyRole: RoleWK})}
	assert.True(t, partyHasInboxWhenRequired(p))
}

func TestStatusPartiesIdentitySchemesAllowedWrongType(t *testing.T) {
	assert.True(t, statusPartiesIdentitySchemesAllowed("x"))
}

func TestStatusReasonCodesAllowedWrongType(t *testing.T) {
	assert.True(t, statusReasonCodesAllowed("x"))
}

func TestStatusReasonCodesAllowedNilReason(t *testing.T) {
	st := &bill.Status{
		Type: bill.StatusTypeResponse,
		Lines: []*bill.StatusLine{{
			Key:     bill.StatusEventRejected,
			Reasons: []*bill.Reason{nil},
		}},
	}
	assert.True(t, statusReasonCodesAllowed(st))
}

// --- ensureSIRENOnSupplier covers the "supplier already carries the
// SIREN" early-return path that the happy-path tests don't reach
// (since the test fixture already aligns SIRENs).

func TestEnsureSIRENOnSupplierAlreadyCarries(t *testing.T) {
	siren := &org.Identity{
		Code: "356000000",
		Ext:  tax.ExtensionsOf(cbc.CodeMap{"iso-scheme-id": "0002"}),
	}
	p := &org.Party{Identities: []*org.Identity{
		{Code: "356000000", Ext: tax.ExtensionsOf(cbc.CodeMap{"iso-scheme-id": "0002"})},
	}}
	got := ensureSIRENOnSupplier(p, siren)
	assert.Same(t, p, got)
	assert.Len(t, got.Identities, 1)
}

// --- ReasonCodeAllowedForProcessCode -----------------------------------

func TestReasonCodeAllowedForProcessCodeUnknownProcess(t *testing.T) {
	// An unknown process code means we have no allow-list — should pass
	// (rule defers to the bucket consistency check).
	assert.True(t, ReasonCodeAllowedForProcessCode("DEST_INC", "999"))
}

// --- org.go sirenFromFrenchTaxID + partyCarriesSIREN ------------------

func TestSirenFromFrenchTaxIDNilSIRETEntry(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{nil}}
	// Falls back to TaxID digits.
	got := sirenFromFrenchTaxID("FR39356000000", p)
	assert.Equal(t, "356000000", got)
}

func TestPartyCarriesSIRENNilParty(t *testing.T) {
	assert.False(t, partyCarriesSIREN(nil))
}

func TestPartyCarriesSIRENNilIdentity(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{nil}}
	assert.False(t, partyCarriesSIREN(p))
}

// TestNormalizeIdentityMapsSIRENTypeToScheme confirms the addon
// normaliser sets iso-scheme-id=0002 on a SIREN-typed identity even
// without eu-en16931 declared. Downstream validators (e.g.
// statusPartyHasSIRENIdentity) can therefore rely on the scheme-id
// extension being present after normalisation.
func TestNormalizeIdentityMapsSIRENTypeToScheme(t *testing.T) {
	id := &org.Identity{Type: "SIREN", Code: "356000000"}
	normalizeIdentity(id)
	assert.Equal(t, identitySchemeIDSIREN, id.Ext.Get("iso-scheme-id").String())
}

func TestNormalizeIdentityMapsSIRETTypeToScheme(t *testing.T) {
	id := &org.Identity{Type: "SIRET", Code: "35600000000011"}
	normalizeIdentity(id)
	assert.Equal(t, identitySchemeIDSIRET, id.Ext.Get("iso-scheme-id").String())
}

func TestNormalizeIdentityPreservesExistingScheme(t *testing.T) {
	id := &org.Identity{
		Type: "SIREN",
		Code: "356000000",
		Ext:  tax.ExtensionsOf(cbc.CodeMap{"iso-scheme-id": "9999"}),
	}
	normalizeIdentity(id)
	assert.Equal(t, "9999", id.Ext.Get("iso-scheme-id").String())
}
