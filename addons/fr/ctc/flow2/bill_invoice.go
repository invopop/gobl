package flow2

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/dgfip"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

const (
	// noteSubjectTXD is the UNTDID 4451 text-subject code carried on
	// the STC mention added by the normalizer.
	noteSubjectTXD cbc.Code = "TXD"

	// stcMembreAssujettiUnique is the fixed text that pairs with TXD.
	stcMembreAssujettiUnique = "MEMBRE_ASSUJETTI_UNIQUE"
)

// -- Normalisation --------------------------------------------------------

func normalizeInvoice(inv *bill.Invoice) {
	if inv == nil {
		return
	}
	normalizeParty(inv.Supplier)
	normalizeParty(inv.Customer)
	if inv.Tax == nil {
		inv.Tax = &bill.Tax{}
	}
	inv.Tax.Rounding = tax.RoundingRuleCurrency
	normalizeBillingMode(inv)
	normalizeRequiredNotes(inv)
	normalizeSTCNote(inv)
}

func normalizeBillingMode(inv *bill.Invoice) {
	if inv.Tax != nil && !inv.Tax.Ext.IsZero() && inv.Tax.Ext.Get(dgfip.ExtKeyBillingMode) != "" {
		return
	}
	mode := dgfip.BillingModeM1
	if inv.Totals != nil && inv.Totals.Paid() {
		mode = dgfip.BillingModeM2
	}
	if inv.Tax == nil {
		inv.Tax = &bill.Tax{}
	}
	inv.Tax.Ext = inv.Tax.Ext.Set(dgfip.ExtKeyBillingMode, mode)
}

func normalizeSTCNote(inv *bill.Invoice) {
	if !isPartyIdentitySTC(inv.Supplier) {
		return
	}
	for _, n := range inv.Notes {
		if n != nil && n.Ext.Get(untdid.ExtKeyTextSubject) == noteSubjectTXD && n.Text == stcMembreAssujettiUnique {
			return
		}
	}
	inv.Notes = append(inv.Notes, &org.Note{
		Key:  org.NoteKeyLegal,
		Text: stcMembreAssujettiUnique,
		Ext:  tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: noteSubjectTXD}),
	})
}

var defaultRequiredNotes = []*org.Note{
	{
		Key:  org.NoteKeyPayment,
		Text: "Conditions de paiement selon les conditions générales de vente.",
		Ext:  tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: "PMT"}),
	},
	{
		Key:  org.NoteKeyPaymentMethod,
		Text: "Pénalités et indemnités de retard applicables conformément aux conditions générales de vente.",
		Ext:  tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: "PMD"}),
	},
	{
		Key:  org.NoteKeyPaymentTerm,
		Text: "Aucun escompte n'est accordé pour paiement anticipé.",
		Ext:  tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: "AAB"}),
	},
}

func normalizeRequiredNotes(inv *bill.Invoice) {
	for _, def := range defaultRequiredNotes {
		want := def.Ext.Get(untdid.ExtKeyTextSubject)
		if invoiceHasNoteWithSubject(inv, want) {
			continue
		}
		clone := *def
		inv.Notes = append(inv.Notes, &clone)
	}
}

func invoiceHasNoteWithSubject(inv *bill.Invoice, subject cbc.Code) bool {
	for _, n := range inv.Notes {
		if n != nil && n.Ext.Get(untdid.ExtKeyTextSubject) == subject {
			return true
		}
	}
	return false
}

// -- Rule set -------------------------------------------------------------

// billInvoiceRules validates only the integrity of the addon's own
// extensions: the UNTDID document type must be one Flow 2 recognises.
// The French CTC format/business rules (BR-FR-*) are the converter's
// responsibility — see the package doc.
func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Field("tax",
			rules.Field("ext",
				rules.Assert("01", "invoice tax ext untdid-document-type must be a recognised Flow 2 document type",
					tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, allowedInvoiceDocumentTypes...),
				),
			),
		),
	)
}
