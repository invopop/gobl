package ctc

import (
	"slices"
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Assert("42", "invoice must be in EUR or provide exchange rate for conversion", currency.CanConvertTo(currency.EUR)),
		// Invoice code validation (BR-FR-01/02) - cross-field: series + code
		rules.Assert("01", "must be 1-35 characters, alphanumeric plus -+_/ (BR-FR-01/02), including the series",
			is.Func("valid invoice code", invoiceCodeValid),
		),
		// Preceding document code validation
		rules.Field("preceding",
			rules.Each(
				rules.Assert("02", "preceding code must be 1-35 characters, alphanumeric plus -+_/ (BR-FR-01/02), including the series",
					is.Func("valid preceding code", precedingDocCodeValid),
				),
			),
		),
		// Corrective invoice preceding (BR-FR-CO-04)
		rules.When(
			is.Func("corrective invoice", invoiceIsCorrectiveAny),
			rules.Field("preceding",
				rules.Assert("03", "corrective invoices must have exactly one preceding invoice reference (BR-FR-CO-04)",
					is.Present,
				),
				rules.Assert("04", "corrective invoices must have exactly one preceding invoice reference (BR-FR-CO-04)",
					is.Length(1, 1),
				),
			),
		),
		// Credit note preceding (BR-FR-CO-05)
		rules.When(
			is.Func("credit note", invoiceIsCreditNoteAny),
			rules.Field("preceding",
				rules.Assert("05", "credit notes must have at least one preceding invoice reference (BR-FR-CO-05)",
					is.Present,
				),
			),
		),
		// Tax validation
		rules.Field("tax",
			rules.Assert("06", "tax is required", is.Present),
			rules.Field("ext",
				rules.Assert("07", "UNTDID document type must be valid (BR-FR-04)",
					tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, allowedDocumentTypes...),
				),
				rules.Assert("08", "billing mode extension is required",
					tax.ExtensionsRequire(ExtKeyBillingMode),
				),
			),
		),
		// Factoring restriction (BR-FR-CO-08)
		rules.When(
			is.Func("factoring mode", invoiceIsFactoringAny),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("09", "advance payment document types not allowed for factoring billing modes (BR-FR-CO-08)",
						tax.ExtensionsExcludeCodes(untdid.ExtKeyDocumentType, advancePaymentDocumentTypes...),
					),
				),
			),
		),
		// Supplier validation
		rules.Field("supplier",
			rules.Field("inboxes",
				rules.Assert("10", "seller electronic address required for French B2B invoices (BR-FR-13)",
					is.Present,
				),
			),
			rules.Field("identities",
				rules.Assert("11", "SIREN identity required for French parties with scheme 0002 and scope legal (BR-FR-10/11)",
					is.Func("has SIREN", identitiesHasSIREN),
				),
			),
		),
		// Supplier SIREN inbox for B2B non-self-billed (BR-FR-21)
		rules.When(
			is.Func("B2B non-self-billed", invoiceIsB2BNonSelfBilledAny),
			rules.Field("supplier",
				rules.Assert("12", "party must have endpoint ID with scheme 0225 (SIREN) (BR-FR-21/22)",
					is.Func("has SIREN inbox", partyHasSIRENInbox),
				),
			),
		),
		// Customer validation
		rules.Field("customer",
			rules.Field("inboxes",
				rules.Assert("13", "buyer electronic address required for French B2B invoices (BR-FR-13)",
					is.Present,
				),
			),
		),
		// B2B customer requires SIREN (BR-FR-14)
		rules.When(
			is.Func("B2B transaction", invoiceIsB2BAny),
			rules.Field("customer",
				rules.Field("identities",
					rules.Assert("14", "SIREN identity required for French parties with scheme 0002 and scope legal (BR-FR-10/11)",
						is.Func("has SIREN", identitiesHasSIREN),
					),
				),
			),
		),
		// B2B self-billed customer SIREN inbox (BR-FR-22)
		rules.When(
			is.Func("B2B self-billed", invoiceIsB2BSelfBilledAny),
			rules.Field("customer",
				rules.Assert("15", "party must have endpoint ID with scheme 0225 (SIREN) (BR-FR-21/22)",
					is.Func("has SIREN inbox", partyHasSIRENInbox),
				),
			),
		),
		// Ordering validation (BR-FR-30)
		rules.Field("ordering",
			rules.Field("identities",
				rules.Assert("16", "only one ordering identity with UNTDID reference 'AFL' is allowed (BR-FR-30)",
					is.Func("no dup AFL", orderingIdentitiesNoDupAFL),
				),
				rules.Assert("17", "only one ordering identity with UNTDID reference 'AWW' is allowed (BR-FR-30)",
					is.Func("no dup AWW", orderingIdentitiesNoDupAWW),
				),
			),
		),
		// STC supplier ordering (BR-FR-CO-15)
		rules.When(
			is.Func("supplier STC", invoiceSupplierIsSTC),
			rules.Field("ordering",
				rules.Assert("18", "ordering with seller is required when supplier is under STC scheme (BR-FR-CO-15)",
					is.Present,
				),
				rules.Field("seller",
					rules.Assert("19", "seller is required when supplier is under STC scheme (BR-FR-CO-15)",
						is.Present,
					),
					rules.Field("tax_id",
						rules.Assert("20", "tax ID is required when supplier is under STC scheme (BR-FR-CO-15)",
							is.Present,
						),
						rules.Field("code",
							rules.Assert("21", "code is required when supplier is under STC scheme (BR-FR-CO-15)",
								is.Present,
							),
						),
					),
				),
			),
			// TXD note requirement (BR-FR-CO-14)
			rules.Field("notes",
				rules.Assert("22", "for sellers with STC scheme (0231), a note with code 'TXD' and text 'MEMBRE_ASSUJETTI_UNIQUE' is required (BR-FR-CO-14)",
					is.Func("has TXD note", notesHaveTXD),
				),
			),
		),
		// Consolidated credit note ordering (BR-FR-CO-03)
		rules.When(
			is.Func("consolidated credit note", invoiceIsConsolidatedCreditNoteAny),
			rules.Field("ordering",
				rules.Assert("23", "ordering with contracts is required for consolidated credit notes (BR-FR-CO-03)",
					is.Present,
				),
				rules.Field("contracts",
					rules.Assert("24", "at least one contract reference is required in ordering details for consolidated credit notes (BR-FR-CO-03)",
						is.Present,
					),
					rules.Assert("25", "at least one contract reference is required in ordering details for consolidated credit notes (BR-FR-CO-03)",
						is.Length(1, 0),
					),
				),
			),
			rules.Field("delivery",
				rules.Assert("26", "delivery details are required for consolidated credit notes (BR-FR-CO-03)",
					is.Present,
				),
				rules.Field("period",
					rules.Assert("27", "delivery period is required for consolidated credit notes (BR-FR-CO-03)",
						is.Present,
					),
				),
			),
		),
		// Payment due date validation (BR-FR-CO-07)
		rules.When(
			is.Func("not advance or final", invoiceIsNotAdvanceOrFinalAny),
			rules.Assert("28", "due dates must not be before invoice issue date (BR-FR-CO-07)",
				is.Func("due dates valid", invoiceDueDatesValid),
			),
		),
		// Final invoice payment (BR-FR-CO-09)
		rules.When(
			is.Func("final invoice", invoiceIsFinalAny),
			rules.Field("payment",
				rules.Assert("29", "payment details are required for final invoices (BR-FR-CO-09)",
					is.Present,
				),
				rules.Field("terms",
					rules.Assert("30", "payment terms required for final invoices (BR-FR-CO-09)",
						is.Present,
					),
					rules.Field("due_dates",
						rules.Assert("31", "at least one due date required for final invoices (BR-FR-CO-09)",
							is.Present,
						),
					),
				),
			),
			// Totals for final invoices
			rules.Field("totals",
				rules.Field("advance",
					rules.Assert("32", "advance amount is required for already-paid invoices (BR-FR-CO-09)",
						is.Present,
					),
				),
				rules.Assert("33", "advance amount must equal total with tax for final invoices (BR-FR-CO-09)",
					is.Func("advances match", finalInvoiceAdvancesMatch),
				),
				rules.Assert("34", "payable amount must be zero for final invoices (BR-FR-CO-09)",
					is.Func("payable zero", finalInvoicePayableZero),
				),
			),
		),
		// Notes validation
		rules.Field("notes",
			rules.Assert("35", "notes are required for French CTC invoices (BR-FR-05)", is.Present),
			rules.Assert("36", "missing required note codes (BR-FR-05)",
				is.Func("has required notes", notesHaveRequired),
			),
			rules.Assert("37", "duplicate note codes found (BR-FR-06/BR-FR-30)",
				is.Func("no duplicate notes", notesNoDuplicates),
			),
			rules.Assert("38", "BAR note text must be one of: B2B, B2BINT, B2C, OUTOFSCOPE, ARCHIVEONLY",
				is.Func("valid BAR text", notesValidBARText),
			),
		),
		// Attachment validation
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

// toAnySlice converts a []string to []any for is.In
func toAnySlice(ss []string) []any {
	out := make([]any, len(ss))
	for i, s := range ss {
		out[i] = s
	}
	return out
}

// --- Invoice-level helper functions ---

func invoiceCodeValid(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Code == cbc.CodeEmpty {
		return true // let required validation handle empty
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
	return isFactoringExtension(inv.Tax.Ext.Get(ExtKeyBillingMode))
}

func invoiceIsB2BAny(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && isB2BTransaction(inv)
}

func invoiceIsB2BNonSelfBilledAny(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && isB2BTransaction(inv) && !isSelfBilledInvoice(inv)
}

func invoiceIsB2BSelfBilledAny(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && isB2BTransaction(inv) && isSelfBilledInvoice(inv)
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

// --- Field-level helper functions ---

func identitiesHasSIREN(val any) bool {
	identities, ok := val.([]*org.Identity)
	if !ok {
		return true // nil/empty passes
	}
	for _, id := range identities {
		if id != nil && !id.Ext.IsZero() {
			if code := id.Ext.Get(iso.ExtKeySchemeID); code == "0002" && id.Scope.Has(org.IdentityScopeLegal) {
				return true
			}
		}
	}
	return false
}

func partyHasSIRENInbox(val any) bool {
	party, ok := val.(*org.Party)
	if !ok || party == nil {
		return true
	}
	siren := getPartySIREN(party)
	if siren == "" {
		return true // SIREN validation handled elsewhere
	}
	for _, inbox := range party.Inboxes {
		if inbox != nil && inbox.Scheme == inboxSchemeSIREN {
			if !strings.HasPrefix(string(inbox.Code), siren) {
				return false // "party endpoint ID scheme inbox (0225) must start with SIREN (BR-FR-21/22)"
			}
			return true
		}
	}
	return false // no SIREN inbox found
}

func orderingIdentitiesNoDupAFL(val any) bool {
	identities, ok := val.([]*org.Identity)
	if !ok {
		return true
	}
	count := 0
	for _, id := range identities {
		if id == nil || id.Ext.IsZero() {
			continue
		}
		if id.Ext.Get(untdid.ExtKeyReference).String() == "AFL" {
			count++
			if count > 1 {
				return false
			}
		}
	}
	return true
}

func orderingIdentitiesNoDupAWW(val any) bool {
	identities, ok := val.([]*org.Identity)
	if !ok {
		return true
	}
	count := 0
	for _, id := range identities {
		if id == nil || id.Ext.IsZero() {
			continue
		}
		if id.Ext.Get(untdid.ExtKeyReference).String() == "AWW" {
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
			if code := note.Ext.Get(untdid.ExtKeyTextSubject); code == "TXD" && note.Text == "MEMBRE_ASSUJETTI_UNIQUE" {
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
	checkUnique := []cbc.Code{"PMT", "PMD", "AAB", "TXD", "BAR"}
	for _, code := range checkUnique {
		if counts[code] > 1 {
			return false
		}
	}
	return true
}

func notesValidBARText(val any) bool {
	notes, ok := val.([]*org.Note)
	if !ok || len(notes) == 0 {
		return true
	}
	for _, note := range notes {
		if note != nil && !note.Ext.IsZero() {
			if note.Ext.Get(untdid.ExtKeyTextSubject) == "BAR" {
				if note.Text != "" && !slices.Contains(allowedBARTreatments, note.Text) {
					return false
				}
			}
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
		return true // handled by required check
	}
	return totals.Advances.Equals(totals.TotalWithTax)
}

func finalInvoicePayableZero(val any) bool {
	totals, ok := val.(*bill.Totals)
	if !ok || totals == nil {
		return true
	}
	// PayableAmount maps to Due if present, otherwise Payable
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
