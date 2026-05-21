package flow2

import (
	"regexp"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/dgfip"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"

	"slices"
)

// invoiceCodeRegexp enforces BR-FR-01/02 invoice-code format: max 35
// characters, alphanumeric plus -+_/.
var invoiceCodeRegexp = regexp.MustCompile(`^[A-Za-z0-9\-\+_/]{1,35}$`)

// allowedAttachmentDescriptions enumerates the BR-FR-17 attachment
// descriptions accepted on a French CTC invoice.
var allowedAttachmentDescriptions = []string{
	"RIB",
	"LISIBLE",
	"FEUILLE_DE_STYLE",
	"PJA",
	"BON_LIVRAISON",
	"BON_COMMANDE",
	"DOCUMENT_ANNEXE",
	"BORDEREAU_SUIVI",
	"BORDEREAU_SUIVI_VALIDATION",
	"ETAT_ACOMPTE",
	"FACTURE_PAIEMENT_DIRECT",
	"RECAPITULATIF_COTRAITANCE",
}

// vatKeyToUNTDIDCategory maps each supported GOBL VAT rate key to its
// UNTDID 5305 category code.
var vatKeyToUNTDIDCategory = map[cbc.Key]cbc.Code{
	tax.KeyStandard:       "S",
	tax.KeyZero:           "Z",
	tax.KeyExempt:         "E",
	tax.KeyReverseCharge:  "AE",
	tax.KeyIntraCommunity: "K",
	tax.KeyExport:         "G",
	tax.KeyOutsideScope:   "O",
}

const (
	// attachmentFormatLisible is the attachment format category for
	// BR-FR-18.
	attachmentFormatLisible = "LISIBLE"

	// noteSubjectTXD is the UNTDID 4451 text-subject code carried on
	// the BR-FR-CO-14 STC mention.
	noteSubjectTXD cbc.Code = "TXD"

	// stcMembreAssujettiUnique is the fixed text that pairs with TXD.
	stcMembreAssujettiUnique = "MEMBRE_ASSUJETTI_UNIQUE"
)

// -- Normalisation --------------------------------------------------------

func normalizeInvoice(inv *bill.Invoice) {
	if inv == nil {
		return
	}
	normalizeInvoiceTaxCategories(inv)
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

func normalizeInvoiceTaxCategories(inv *bill.Invoice) {
	for _, line := range inv.Lines {
		if line == nil {
			continue
		}
		for _, combo := range line.Taxes {
			if combo == nil || combo.Category != tax.CategoryVAT {
				continue
			}
			if combo.Ext.Get(untdid.ExtKeyTaxCategory) != "" {
				continue
			}
			if code, ok := vatKeyToUNTDIDCategory[combo.Key]; ok {
				combo.Ext = combo.Ext.Set(untdid.ExtKeyTaxCategory, code)
			}
		}
	}
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

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Assert("01", "invoice must be in EUR or provide an exchange rate to EUR",
			currency.CanConvertTo(currency.EUR),
		),
		rules.Assert("02", "must be 1-35 characters, alphanumeric plus -+_/ (BR-FR-01/02), including the series",
			is.Func("valid invoice code", invoiceCodeValid),
		),
		rules.Field("preceding",
			rules.Each(
				rules.Assert("03", "preceding code must be 1-35 characters, alphanumeric plus -+_/ (BR-FR-01/02), including the series",
					is.Func("valid preceding code", precedingDocCodeValid),
				),
			),
		),
		rules.When(
			is.Func("corrective invoice", invoiceIsCorrectiveAny),
			rules.Field("preceding",
				rules.Assert("04", "corrective invoices must reference the original invoice in preceding (BR-FR-CO-04)",
					is.Present,
				),
				rules.Assert("05", "corrective invoices must reference exactly one preceding invoice — multiple references are not allowed (BR-FR-CO-04)",
					is.Length(1, 1),
				),
			),
		),
		rules.When(
			is.Func("credit note", invoiceIsCreditNoteAny),
			rules.Field("preceding",
				rules.Assert("06", "credit notes must have at least one preceding invoice reference (BR-FR-CO-05)",
					is.Present,
				),
			),
		),
		rules.Field("tax",
			rules.Assert("07", "tax is required", is.Present),
			rules.Field("ext",
				rules.Assert("08", "UNTDID document type must be valid (BR-FR-04)",
					tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, allowedInvoiceDocumentTypes...),
				),
				rules.Assert("09", "billing mode extension is required",
					tax.ExtensionsRequire(dgfip.ExtKeyBillingMode),
				),
			),
		),
		rules.When(
			is.Func("factoring mode", invoiceIsFactoringAny),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("10", "advance payment document types (386, 500, 503) are not allowed for factoring billing modes (B4, S4, M4) (BR-FR-CO-08)",
						tax.ExtensionsExcludeCodes(untdid.ExtKeyDocumentType, advancePaymentDocumentTypes...),
					),
				),
			),
		),
		rules.Field("supplier",
			rules.Field("inboxes",
				rules.Assert("11", "seller electronic address required for French B2B invoices (BR-FR-13)",
					is.Present,
				),
			),
			rules.Field("identities",
				rules.Assert("12", "SIREN identity required for French parties with scheme 0002 and scope legal (BR-FR-10/11)",
					is.Func("has SIREN (legal scope)", identitiesHasLegalSIREN),
				),
			),
		),
		rules.When(
			is.Func("not self-billed", invoiceIsNotSelfBilledAny),
			rules.Field("supplier",
				rules.Assert("13", "party must have endpoint ID with scheme 0225 (SIREN) (BR-FR-21/22)",
					is.Func("has SIREN inbox", partyHasSIRENInbox),
				),
			),
		),
		rules.Field("customer",
			rules.Field("inboxes",
				rules.Assert("14", "buyer electronic address required for French B2B invoices (BR-FR-13)",
					is.Present,
				),
			),
			rules.Field("identities",
				rules.Assert("15", "SIREN identity required for French parties with scheme 0002 and scope legal (BR-FR-10/11)",
					is.Func("has SIREN (legal scope)", identitiesHasLegalSIREN),
				),
			),
		),
		rules.When(
			is.Func("self-billed", invoiceIsSelfBilledAny),
			rules.Field("customer",
				rules.Assert("16", "party must have endpoint ID with scheme 0225 (SIREN) (BR-FR-21/22)",
					is.Func("has SIREN inbox", partyHasSIRENInbox),
				),
			),
		),
		rules.Field("ordering",
			rules.Field("identities",
				rules.Assert("17", "only one ordering identity with UNTDID reference 'AFL' is allowed (BR-FR-30)",
					is.Func("no dup AFL", orderingIdentitiesNoDupAFL),
				),
				rules.Assert("18", "only one ordering identity with UNTDID reference 'AWW' is allowed (BR-FR-30)",
					is.Func("no dup AWW", orderingIdentitiesNoDupAWW),
				),
			),
		),
		rules.When(
			is.Func("supplier STC", invoiceSupplierIsSTC),
			rules.Field("ordering",
				rules.Assert("19", "ordering with seller is required when supplier is under STC scheme (BR-FR-CO-15)",
					is.Present,
				),
				rules.Field("seller",
					rules.Assert("20", "seller is required when supplier is under STC scheme (BR-FR-CO-15)",
						is.Present,
					),
					rules.Field("tax_id",
						rules.Assert("21", "tax ID is required when supplier is under STC scheme (BR-FR-CO-15)",
							is.Present,
						),
						rules.Field("code",
							rules.Assert("22", "code is required when supplier is under STC scheme (BR-FR-CO-15)",
								is.Present,
							),
						),
					),
				),
			),
			rules.Field("notes",
				rules.Assert("23", "for sellers with STC scheme (0231), a note with code 'TXD' and text 'MEMBRE_ASSUJETTI_UNIQUE' is required (BR-FR-CO-14)",
					is.Func("has TXD note", notesHaveTXD),
				),
			),
		),
		rules.When(
			is.Func("consolidated credit note", invoiceIsConsolidatedCreditNoteAny),
			rules.Field("ordering",
				rules.Assert("24", "ordering with contracts is required for consolidated credit notes (BR-FR-CO-03)",
					is.Present,
				),
				rules.Field("contracts",
					rules.Assert("25", "ordering.contracts is required for consolidated credit notes (BR-FR-CO-03)",
						is.Present,
					),
					rules.Assert("26", "ordering.contracts must contain at least one entry for consolidated credit notes (BR-FR-CO-03)",
						is.Length(1, 0),
					),
				),
			),
			rules.Field("delivery",
				rules.Assert("27", "delivery details are required for consolidated credit notes (BR-FR-CO-03)",
					is.Present,
				),
				rules.Field("period",
					rules.Assert("28", "delivery period is required for consolidated credit notes (BR-FR-CO-03)",
						is.Present,
					),
				),
			),
		),
		rules.When(
			is.Func("not advance or final", invoiceIsNotAdvanceOrFinalAny),
			rules.Assert("29", "due dates must not be before invoice issue date (BR-FR-CO-07)",
				is.Func("due dates valid", invoiceDueDatesValid),
			),
		),
		rules.When(
			is.Func("final invoice", invoiceIsFinalAny),
			rules.Field("payment",
				rules.Assert("30", "payment details are required for final invoices (BR-FR-CO-09)",
					is.Present,
				),
				rules.Field("terms",
					rules.Assert("31", "payment terms required for final invoices (BR-FR-CO-09)",
						is.Present,
					),
					rules.Field("due_dates",
						rules.Assert("32", "at least one due date required for final invoices (BR-FR-CO-09)",
							is.Present,
						),
					),
				),
			),
			rules.Field("totals",
				rules.Field("advance",
					rules.Assert("33", "advance amount is required for already-paid invoices (BR-FR-CO-09)",
						is.Present,
					),
				),
				rules.Assert("34", "advance amount must equal total with tax for final invoices (BR-FR-CO-09)",
					is.Func("advances match", finalInvoiceAdvancesMatch),
				),
				rules.Assert("35", "payable amount must be zero for final invoices (BR-FR-CO-09)",
					is.Func("payable zero", finalInvoicePayableZero),
				),
			),
		),
		rules.Field("notes",
			rules.Assert("36", "notes are required for French CTC invoices (BR-FR-05)", is.Present),
			rules.Assert("37", "missing required note codes (BR-FR-05)",
				is.Func("has required notes", notesHaveRequired),
			),
			rules.Assert("38", "duplicate note codes found (BR-FR-06/BR-FR-30)",
				is.Func("no duplicate notes", notesNoDuplicates),
			),
		),
		rules.Field("attachments",
			rules.Each(
				rules.Field("description",
					rules.Assert("39", "must be one of the allowed attachment descriptions (BR-FR-17)",
						is.Present,
					),
					rules.Assert("40", "must be one of the allowed attachment descriptions (BR-FR-17)",
						is.In(toAnySlice(allowedAttachmentDescriptions)...),
					),
				),
			),
			rules.Assert("41", "only one attachment with description 'LISIBLE' is allowed per invoice (BR-FR-18)",
				is.Func("unique LISIBLE", attachmentsUniqueLISIBLE),
			),
		),
	)
}

// -- Predicates -----------------------------------------------------------

func invoiceCodeValid(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Code == cbc.CodeEmpty {
		return true
	}
	invoiceID := string(inv.Code)
	if inv.Series != cbc.CodeEmpty {
		invoiceID = string(inv.Series.Join(inv.Code))
	}
	return invoiceCodeRegexp.MatchString(invoiceID)
}

func precedingDocCodeValid(val any) bool {
	docRef, ok := val.(*org.DocumentRef)
	if !ok || docRef == nil || docRef.Code == cbc.CodeEmpty {
		return true
	}
	invoiceID := string(docRef.Code)
	if docRef.Series != cbc.CodeEmpty {
		invoiceID = string(docRef.Series.Join(docRef.Code))
	}
	return invoiceCodeRegexp.MatchString(invoiceID)
}

func invoiceIsCorrectiveAny(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && isCorrectiveInvoice(inv)
}

func invoiceIsCreditNoteAny(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && isCreditNote(inv)
}

func invoiceIsFactoringAny(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
		return false
	}
	return isFactoringExtension(inv.Tax.Ext.Get(dgfip.ExtKeyBillingMode))
}

func invoiceIsNotSelfBilledAny(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && !isSelfBilledInvoice(inv)
}

func invoiceIsSelfBilledAny(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && isSelfBilledInvoice(inv)
}

func invoiceSupplierIsSTC(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && isPartyIdentitySTC(inv.Supplier)
}

func invoiceIsConsolidatedCreditNoteAny(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && isConsolidatedCreditNote(inv)
}

func invoiceIsNotAdvanceOrFinalAny(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && !isAdvancedInvoice(inv) && !isFinalInvoice(inv)
}

func invoiceIsFinalAny(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && isFinalInvoice(inv)
}

func isSelfBilledInvoice(inv *bill.Invoice) bool {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
		return false
	}
	docType := inv.Tax.Ext.Get(untdid.ExtKeyDocumentType)
	if docType == "" {
		return false
	}
	return slices.Contains(selfBilledDocumentTypes, docType)
}

func isCorrectiveInvoice(inv *bill.Invoice) bool {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
		return false
	}
	docType := inv.Tax.Ext.Get(untdid.ExtKeyDocumentType)
	if docType == "" {
		return false
	}
	return slices.Contains(correctiveInvoiceTypes, docType)
}

func isCreditNote(inv *bill.Invoice) bool {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
		return false
	}
	docType := inv.Tax.Ext.Get(untdid.ExtKeyDocumentType)
	return slices.Contains(creditNoteTypes, docType)
}

func isConsolidatedCreditNote(inv *bill.Invoice) bool {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
		return false
	}
	docType := inv.Tax.Ext.Get(untdid.ExtKeyDocumentType)
	return docType == "262"
}

func isAdvancedInvoice(inv *bill.Invoice) bool {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
		return false
	}
	docType := inv.Tax.Ext.Get(untdid.ExtKeyDocumentType)
	return slices.Contains(advancePaymentDocumentTypes, docType)
}

func isFinalInvoice(inv *bill.Invoice) bool {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
		return false
	}
	bm := inv.Tax.Ext.Get(dgfip.ExtKeyBillingMode)
	return bm == dgfip.BillingModeB2 || bm == dgfip.BillingModeS2 || bm == dgfip.BillingModeM2
}

func isFactoringExtension(bm cbc.Code) bool {
	return bm == dgfip.BillingModeB4 || bm == dgfip.BillingModeS4 || bm == dgfip.BillingModeM4
}

func orderingIdentitiesNoDupAFL(val any) bool {
	return orderingIdentitiesNoDup(val, "AFL")
}

func orderingIdentitiesNoDupAWW(val any) bool {
	return orderingIdentitiesNoDup(val, "AWW")
}

func orderingIdentitiesNoDup(val any, ref string) bool {
	identities, ok := val.([]*org.Identity)
	if !ok {
		return true
	}
	count := 0
	for _, id := range identities {
		if id == nil || id.Ext.IsZero() {
			continue
		}
		if id.Ext.Get(untdid.ExtKeyReference).String() == ref {
			count++
			if count > 1 {
				return false
			}
		}
	}
	return true
}

func notesHaveTXD(val any) bool {
	notes, ok := val.([]*org.Note)
	if !ok || len(notes) == 0 {
		return false
	}
	for _, note := range notes {
		if note != nil && !note.Ext.IsZero() {
			if code := note.Ext.Get(untdid.ExtKeyTextSubject); code == noteSubjectTXD && note.Text == stcMembreAssujettiUnique {
				return true
			}
		}
	}
	return false
}

func notesHaveRequired(val any) bool {
	notes, ok := val.([]*org.Note)
	if !ok || len(notes) == 0 {
		return false
	}
	required := []cbc.Code{"PMT", "PMD", "AAB"}
	counts := make(map[cbc.Code]int)
	for _, note := range notes {
		if note != nil && !note.Ext.IsZero() {
			if code := note.Ext.Get(untdid.ExtKeyTextSubject); code != cbc.CodeEmpty {
				counts[code]++
			}
		}
	}
	for _, code := range required {
		if counts[code] == 0 {
			return false
		}
	}
	return true
}

func notesNoDuplicates(val any) bool {
	notes, ok := val.([]*org.Note)
	if !ok || len(notes) == 0 {
		return true
	}
	counts := make(map[cbc.Code]int)
	for _, note := range notes {
		if note != nil && !note.Ext.IsZero() {
			if code := note.Ext.Get(untdid.ExtKeyTextSubject); code != cbc.CodeEmpty {
				counts[code]++
			}
		}
	}
	checkUnique := []cbc.Code{"PMT", "PMD", "AAB", "TXD"}
	for _, code := range checkUnique {
		if counts[code] > 1 {
			return false
		}
	}
	return true
}

func invoiceDueDatesValid(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	if inv.Payment == nil || inv.Payment.Terms == nil || len(inv.Payment.Terms.DueDates) == 0 {
		return true
	}
	for _, dd := range inv.Payment.Terms.DueDates {
		if dd == nil || dd.Date == nil {
			continue
		}
		if inv.IssueDate.DaysSince(dd.Date.Date) > 0 {
			return false
		}
	}
	return true
}

func finalInvoiceAdvancesMatch(val any) bool {
	totals, ok := val.(*bill.Totals)
	if !ok || totals == nil || totals.Advances == nil {
		return true
	}
	return totals.Advances.Equals(totals.TotalWithTax)
}

func finalInvoicePayableZero(val any) bool {
	totals, ok := val.(*bill.Totals)
	if !ok || totals == nil {
		return true
	}
	if totals.Due != nil {
		return totals.Due.Equals(num.AmountZero)
	}
	return totals.Payable.Equals(num.AmountZero)
}

func attachmentsUniqueLISIBLE(val any) bool {
	attachments, ok := val.([]*org.Attachment)
	if !ok || len(attachments) == 0 {
		return true
	}
	count := 0
	for _, att := range attachments {
		if att != nil && att.Description == attachmentFormatLisible {
			count++
		}
	}
	return count <= 1
}

func toAnySlice(ss []string) []any {
	out := make([]any, len(ss))
	for i, s := range ss {
		out[i] = s
	}
	return out
}
