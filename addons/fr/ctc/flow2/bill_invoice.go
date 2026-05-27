package flow2

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

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
)

// invoiceCodeRegexp enforces BR-FR-01/02 invoice-code format: max 35
// characters, alphanumeric plus -+_/.
var invoiceCodeRegexp = regexp.MustCompile(`^[A-Za-z0-9\-\+_/]{1,35}$`)

// allowedAttachmentDescriptions enumerates the BR-FR-17 attachment
// descriptions accepted on a French CTC invoice.
var allowedAttachmentDescriptions = []any{
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

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Assert("01", "invoice must be in EUR or provide an exchange rate to EUR",
			currency.CanConvertTo(currency.EUR),
		),
		rules.Assert("02", "invoice code (with series, when set) must be 1-35 characters of alphanumerics plus -+_/ (BR-FR-01/02)",
			is.Func("valid invoice code", invoiceCodeValid),
		),
		rules.Field("preceding",
			rules.Each(
				rules.Assert("03", "invoice preceding code (with series, when set) must be 1-35 characters of alphanumerics plus -+_/ (BR-FR-01/02)",
					is.Func("valid preceding code", precedingDocCodeValid),
				),
			),
		),
		rules.When(
			invoiceTaxExtIn(untdid.ExtKeyDocumentType, correctiveInvoiceTypes...),
			rules.Field("preceding",
				rules.Assert("04", "invoice preceding is required for corrective invoices (BR-FR-CO-04)",
					is.Present,
				),
				rules.Assert("05", "invoice preceding must contain exactly one reference for corrective invoices (BR-FR-CO-04)",
					is.Length(1, 1),
				),
			),
		),
		rules.When(
			invoiceTaxExtIn(untdid.ExtKeyDocumentType, creditNoteTypes...),
			rules.Field("preceding",
				rules.Assert("06", "invoice preceding must contain at least one reference for credit notes (BR-FR-CO-05)",
					is.Present,
				),
			),
		),
		rules.Field("tax",
			rules.Assert("07", "invoice tax is required", is.Present),
			rules.Field("ext",
				rules.Assert("08", "invoice tax ext untdid-document-type must be a valid Flow 2 document type (BR-FR-04)",
					tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, allowedInvoiceDocumentTypes...),
				),
				rules.Assert("09", "invoice tax ext dgfip-billing-mode is required",
					tax.ExtensionsRequire(dgfip.ExtKeyBillingMode),
				),
			),
		),
		rules.When(
			invoiceTaxExtIn(dgfip.ExtKeyBillingMode, dgfip.BillingModeB4, dgfip.BillingModeS4, dgfip.BillingModeM4),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("10", "invoice tax ext untdid-document-type must not be an advance-payment code (386, 500, 503) when billing mode is factoring (B4, S4, M4) (BR-FR-CO-08)",
						tax.ExtensionsExcludeCodes(untdid.ExtKeyDocumentType, advancePaymentDocumentTypes...),
					),
				),
			),
		),
		rules.Field("supplier",
			rules.Field("inboxes",
				rules.Assert("11", "invoice supplier inboxes are required for French B2B invoices (BR-FR-13)",
					is.Present,
				),
			),
			rules.Field("identities",
				rules.Assert("12", "invoice supplier identities must include a SIREN identity with iso-scheme-id 0002 and scope legal (BR-FR-10/11)",
					is.Func("has SIREN (legal scope)", identitiesHasLegalSIREN),
				),
			),
		),
		rules.When(
			invoiceTaxExtNotIn(untdid.ExtKeyDocumentType, selfBilledDocumentTypes...),
			rules.Field("supplier",
				rules.Assert("13", "invoice supplier must have an inbox with scheme 0225 matching the SIREN code (BR-FR-21/22)",
					is.Func("has SIREN inbox", partyHasSIRENInbox),
				),
			),
		),
		rules.Field("customer",
			rules.Field("inboxes",
				rules.Assert("14", "invoice customer inboxes are required for French B2B invoices (BR-FR-13)",
					is.Present,
				),
			),
			rules.Field("identities",
				rules.Assert("15", "invoice customer identities must include a SIREN identity with iso-scheme-id 0002 and scope legal (BR-FR-10/11)",
					is.Func("has SIREN (legal scope)", identitiesHasLegalSIREN),
				),
			),
		),
		rules.When(
			invoiceTaxExtIn(untdid.ExtKeyDocumentType, selfBilledDocumentTypes...),
			rules.Field("customer",
				rules.Assert("16", "invoice customer must have an inbox with scheme 0225 matching the SIREN code (BR-FR-21/22)",
					is.Func("has SIREN inbox", partyHasSIRENInbox),
				),
			),
		),
		rules.Field("ordering",
			rules.Field("identities",
				rules.Assert("17", "invoice ordering identities must not contain more than one entry with UNTDID reference 'AFL' (BR-FR-30)",
					identitiesNoDupExt(untdid.ExtKeyReference, "AFL"),
				),
				rules.Assert("18", "invoice ordering identities must not contain more than one entry with UNTDID reference 'AWW' (BR-FR-30)",
					identitiesNoDupExt(untdid.ExtKeyReference, "AWW"),
				),
			),
		),
		rules.When(
			invoiceSupplierIsSTC(),
			rules.Field("ordering",
				rules.Assert("19", "invoice ordering is required when supplier is under STC scheme (BR-FR-CO-15)",
					is.Present,
				),
				rules.Field("seller",
					rules.Assert("20", "invoice ordering seller is required when supplier is under STC scheme (BR-FR-CO-15)",
						is.Present,
					),
					rules.Field("tax_id",
						rules.Assert("21", "invoice ordering seller tax_id is required when supplier is under STC scheme (BR-FR-CO-15)",
							is.Present,
						),
						rules.Field("code",
							rules.Assert("22", "invoice ordering seller tax_id code is required when supplier is under STC scheme (BR-FR-CO-15)",
								is.Present,
							),
						),
					),
				),
			),
			rules.Field("notes",
				rules.Assert("23", "invoice notes must include an entry with subject TXD and text MEMBRE_ASSUJETTI_UNIQUE when supplier is under STC scheme (BR-FR-CO-14)",
					is.Func("has TXD note", notesHaveTXD),
				),
			),
		),
		rules.When(
			invoiceTaxExtIn(untdid.ExtKeyDocumentType, "262"),
			rules.Field("ordering",
				rules.Assert("24", "invoice ordering is required for consolidated credit notes (BR-FR-CO-03)",
					is.Present,
				),
				rules.Field("contracts",
					rules.Assert("25", "invoice ordering contracts is required for consolidated credit notes (BR-FR-CO-03)",
						is.Present,
					),
					rules.Assert("26", "invoice ordering contracts must contain at least one entry for consolidated credit notes (BR-FR-CO-03)",
						is.Length(1, 0),
					),
				),
			),
			rules.Field("delivery",
				rules.Assert("27", "invoice delivery is required for consolidated credit notes (BR-FR-CO-03)",
					is.Present,
				),
				rules.Field("period",
					rules.Assert("28", "invoice delivery period is required for consolidated credit notes (BR-FR-CO-03)",
						is.Present,
					),
				),
			),
		),
		rules.When(
			invoiceTaxExtNotIn(untdid.ExtKeyDocumentType, advancePaymentDocumentTypes...),
			rules.When(
				invoiceTaxExtNotIn(dgfip.ExtKeyBillingMode, dgfip.BillingModeB2, dgfip.BillingModeS2, dgfip.BillingModeM2),
				rules.Assert("29", "invoice payment terms due_dates must not be before the invoice issue date (BR-FR-CO-07)",
					is.Func("due dates valid", invoiceDueDatesValid),
				),
			),
		),
		rules.When(
			invoiceTaxExtIn(dgfip.ExtKeyBillingMode, dgfip.BillingModeB2, dgfip.BillingModeS2, dgfip.BillingModeM2),
			rules.Field("payment",
				rules.Assert("30", "invoice payment is required for final invoices (BR-FR-CO-09)",
					is.Present,
				),
				rules.Field("terms",
					rules.Assert("31", "invoice payment terms is required for final invoices (BR-FR-CO-09)",
						is.Present,
					),
					rules.Field("due_dates",
						rules.Assert("32", "invoice payment terms due_dates must contain at least one entry for final invoices (BR-FR-CO-09)",
							is.Present,
						),
					),
				),
			),
			rules.Field("totals",
				rules.Field("advance",
					rules.Assert("33", "invoice totals advance is required for final invoices (BR-FR-CO-09)",
						is.Present,
					),
				),
				rules.Assert("34", "invoice totals advance must equal total_with_tax for final invoices (BR-FR-CO-09)",
					is.Func("advances match total_with_tax", finalInvoiceAdvancesMatch),
				),
				rules.Assert("35", "invoice totals payable must be zero for final invoices (BR-FR-CO-09)",
					is.Func("payable is zero", finalInvoicePayableZero),
				),
			),
		),
		rules.Field("notes",
			rules.Assert("36", "invoice notes are required for French CTC invoices (BR-FR-05)",
				is.Present,
			),
			rules.Assert("37", "invoice notes must include entries with subjects PMT, PMD and AAB (BR-FR-05)",
				is.Func("has required note subjects", notesHaveRequired),
			),
			rules.Assert("38", "invoice notes must not contain duplicate subjects PMT, PMD, AAB or TXD (BR-FR-06/BR-FR-30)",
				is.Func("no duplicate note subjects", notesNoDuplicates),
			),
		),
		rules.Field("attachments",
			rules.Each(
				rules.Field("description",
					rules.Assert("39", "invoice attachment description is required (BR-FR-17)",
						is.Present,
					),
					rules.Assert("40", "invoice attachment description must be one of the allowed values (BR-FR-17)",
						is.In(allowedAttachmentDescriptions...),
					),
				),
			),
			rules.Assert("41", "invoice attachments must not contain more than one entry with description 'LISIBLE' (BR-FR-18)",
				is.Func("at most one LISIBLE", attachmentsAtMostOneLISIBLE),
			),
		),
	)
}

// -- Rule-level guards ----------------------------------------------------

// invoiceTaxExtIn returns a Test that passes when bill.Invoice.Tax.Ext[key]
// matches one of the provided codes. Used to gate per-document-type or
// per-billing-mode branches without writing a dedicated predicate.
func invoiceTaxExtIn(key cbc.Key, codes ...cbc.Code) rules.Test {
	return is.Func(
		fmt.Sprintf("invoice tax ext %s in [%s]", key, joinCodes(codes)),
		func(v any) bool {
			inv, ok := v.(*bill.Invoice)
			if !ok || inv == nil || inv.Tax == nil {
				return false
			}
			return slices.Contains(codes, inv.Tax.Ext.Get(key))
		},
	)
}

// invoiceTaxExtNotIn is the negation of invoiceTaxExtIn. A missing
// Tax or empty value counts as "not in" — the guarded rule fires.
func invoiceTaxExtNotIn(key cbc.Key, codes ...cbc.Code) rules.Test {
	return is.Func(
		fmt.Sprintf("invoice tax ext %s not in [%s]", key, joinCodes(codes)),
		func(v any) bool {
			inv, ok := v.(*bill.Invoice)
			if !ok || inv == nil || inv.Tax == nil {
				return true
			}
			return !slices.Contains(codes, inv.Tax.Ext.Get(key))
		},
	)
}

// invoiceSupplierIsSTC returns a Test that passes when the invoice's
// Supplier carries an STC (iso-scheme-id 0231) identity.
func invoiceSupplierIsSTC() rules.Test {
	return is.Func("invoice supplier is under STC scheme", func(v any) bool {
		inv, ok := v.(*bill.Invoice)
		return ok && inv != nil && isPartyIdentitySTC(inv.Supplier)
	})
}

// identitiesNoDupExt returns a Test for an []*org.Identity slice that
// passes unless more than one identity carries the (key, code) pair.
func identitiesNoDupExt(key cbc.Key, code cbc.Code) rules.Test {
	return is.Func(
		fmt.Sprintf("identities have at most one ext %s=%s", key, code),
		func(v any) bool {
			identities, ok := v.([]*org.Identity)
			if !ok {
				return true
			}
			count := 0
			for _, id := range identities {
				if id == nil || id.Ext.IsZero() {
					continue
				}
				if id.Ext.Get(key) == code {
					count++
					if count > 1 {
						return false
					}
				}
			}
			return true
		},
	)
}

// -- Imperative predicates ------------------------------------------------

// invoiceCodeValid checks the {series, code} pair against BR-FR-01/02.
// Series.Join produces a derived string that no field-level matcher can
// reach, so this stays as an imperative predicate.
func invoiceCodeValid(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Code == cbc.CodeEmpty {
		return true
	}
	id := string(inv.Code)
	if inv.Series != cbc.CodeEmpty {
		id = string(inv.Series.Join(inv.Code))
	}
	return invoiceCodeRegexp.MatchString(id)
}

func precedingDocCodeValid(val any) bool {
	ref, ok := val.(*org.DocumentRef)
	if !ok || ref == nil || ref.Code == cbc.CodeEmpty {
		return true
	}
	id := string(ref.Code)
	if ref.Series != cbc.CodeEmpty {
		id = string(ref.Series.Join(ref.Code))
	}
	return invoiceCodeRegexp.MatchString(id)
}

func notesHaveTXD(val any) bool {
	notes, ok := val.([]*org.Note)
	if !ok || len(notes) == 0 {
		return false
	}
	for _, n := range notes {
		if n == nil || n.Ext.IsZero() {
			continue
		}
		if n.Ext.Get(untdid.ExtKeyTextSubject) == noteSubjectTXD && n.Text == stcMembreAssujettiUnique {
			return true
		}
	}
	return false
}

func notesHaveRequired(val any) bool {
	notes, ok := val.([]*org.Note)
	if !ok {
		return false
	}
	seen := make(map[cbc.Code]bool, 3)
	for _, n := range notes {
		if n == nil || n.Ext.IsZero() {
			continue
		}
		seen[n.Ext.Get(untdid.ExtKeyTextSubject)] = true
	}
	for _, code := range [...]cbc.Code{"PMT", "PMD", "AAB"} {
		if !seen[code] {
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
	counts := make(map[cbc.Code]int, len(notes))
	for _, n := range notes {
		if n == nil || n.Ext.IsZero() {
			continue
		}
		counts[n.Ext.Get(untdid.ExtKeyTextSubject)]++
	}
	for _, code := range [...]cbc.Code{"PMT", "PMD", "AAB", "TXD"} {
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
	t, ok := val.(*bill.Totals)
	if !ok || t == nil || t.Advances == nil {
		return true
	}
	return t.Advances.Equals(t.TotalWithTax)
}

func finalInvoicePayableZero(val any) bool {
	t, ok := val.(*bill.Totals)
	if !ok || t == nil {
		return true
	}
	if t.Due != nil {
		return t.Due.Equals(num.AmountZero)
	}
	return t.Payable.Equals(num.AmountZero)
}

func attachmentsAtMostOneLISIBLE(val any) bool {
	attachments, ok := val.([]*org.Attachment)
	if !ok || len(attachments) == 0 {
		return true
	}
	count := 0
	for _, a := range attachments {
		if a != nil && a.Description == attachmentFormatLisible {
			count++
			if count > 1 {
				return false
			}
		}
	}
	return true
}

// joinCodes is a small helper used only by the guard-name strings.
func joinCodes(codes []cbc.Code) string {
	ss := make([]string, len(codes))
	for i, c := range codes {
		ss[i] = string(c)
	}
	return strings.Join(ss, ", ")
}
