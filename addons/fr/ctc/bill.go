package ctc

import (
	"errors"
	"regexp"
	"slices"
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// BR-FR-01/02: Invoice code validation
// Max 35 characters, alphanumeric plus -+_/
var invoiceCodeRegexp = regexp.MustCompile(`^[A-Za-z0-9\-\+_/]{1,35}$`)

// BR-FR-04: Allowed UNTDID document types for French CTC
var allowedDocumentTypes = []cbc.Code{
	"380", // Commercial invoice
	"389", // Self-billed invoice
	"393", // Factored invoice
	"501", // Final invoice
	"386", // Advance payment invoice
	"500", // Self-billed advance payment
	"384", // Corrective invoice
	"471", // Prepaid amount invoice
	"472", // Self-billed prepaid amount
	"473", // Stand-alone credit note
	"261", // Self-billed credit note
	"262", // Consolidated credit note
	"381", // Credit note
	"396", // Factored credit note
	"502", // Self-billed corrective
	"503", // Self-billed credit for claim
}

// Allowed BAR treatment values for French CTC
var allowedBARTreatments = []string{
	"B2B",
	"B2BINT",
	"B2C",
	"OUTOFSCOPE",
	"ARCHIVEONLY",
}

// Self-billed document types (used for BR-FR-21/BR-FR-22)
var selfBilledDocumentTypes = []cbc.Code{
	"389", // Self-billed invoice
	"501", // Final invoice (self-billed context)
	"500", // Self-billed advance payment
	"471", // Prepaid amount invoice (self-billed context)
	"473", // Stand-alone credit note (self-billed context)
	"261", // Self-billed credit note
	"502", // Self-billed corrective
}

// Corrective invoice document types (BR-FR-CO-04)
var correctiveInvoiceTypes = []cbc.Code{
	"384", // Corrective invoice
	"471", // Prepaid amount invoice
	"472", // Self-billed prepaid amount
	"473", // Stand-alone credit note
}

// Credit note document types (BR-FR-CO-05)
var creditNoteTypes = []cbc.Code{
	"261", // Self-billed credit note
	"381", // Credit note
	"396", // Factored credit note
	"502", // Self-billed corrective
	"503", // Self-billed credit for claim
}

var advancePaymentDocumentTypes = []cbc.Code{
	"386", // Advance payment invoice
	"500", // Self-billed advance payment
	"503", // Self-billed credit for claim
}

// normalizeInvoice ensures invoice settings comply with French CTC requirements
func normalizeInvoice(inv *bill.Invoice) {
	if inv == nil {
		return
	}

	// Ensure Tax object exists
	if inv.Tax == nil {
		inv.Tax = &bill.Tax{}
	}

	// Always set rounding to currency for French CTC
	inv.Tax.Rounding = tax.RoundingRuleCurrency
}

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Code,
			validation.By(
				validateCode(inv.Series),
			),
			validation.Skip,
		),
		validation.Field(&inv.Preceding,
			validation.Each(
				validation.By(
					validatePrecedingDocument,
				),
			),
			validation.When(
				isCorrectiveInvoice(inv),
				validation.Required.Error("corrective invoices must have exactly one preceding invoice reference (BR-FR-CO-04)"),
				validation.Length(1, 1).Error("corrective invoices must have exactly one preceding invoice reference (BR-FR-CO-04)"),
			),
			validation.When(
				isCreditNote(inv),
				validation.Required.Error("credit notes must have at least one preceding invoice reference (BR-FR-CO-05)"),
			),
			validation.Skip,
		),
		validation.Field(&inv.Tax,
			validation.By(
				validateBillInvoiceTax,
			),
			validation.Required,
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.By(
				validateSupplier,
			),
			validation.When(
				isB2BTransaction(inv) && !isSelfBilledInvoice(inv),
				validation.By(
					validateSirenInbox,
				),
			),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.By(
				validateCustomer,
			),
			validation.When(
				isB2BTransaction(inv),
				validation.By(
					validateB2BCustomer,
				),
			),
			validation.When(
				isB2BTransaction(inv) && isSelfBilledInvoice(inv),
				validation.By(
					validateSirenInbox,
				),
			),
			validation.Skip,
		),
		validation.Field(&inv.Ordering,
			validation.By(
				validateOrdering,
			),
			validation.When(
				isPartyIdentitySTC(inv.Supplier),
				validation.By(
					validateOrderingSeller,
				),
				validation.Required.Error("ordering with seller is required when supplier is under STC scheme (BR-FR-CO-15)"),
			),
			validation.When(
				isConsolidatedCreditNote(inv),
				validation.By(
					validateOrderingContracts,
				),
				validation.Required.Error("ordering with contracts is required for consolidated credit notes (BR-FR-CO-03)"),
			),
			validation.Skip,
		),
		validation.Field(&inv.Payment,
			validation.When(
				!isAdvancedInvoice(inv) && !isFinalInvoice(inv),
				validation.By(
					validatePayment(inv.IssueDate),
				),
			),
			validation.When(
				isFinalInvoice(inv),
				validation.By(
					validatePaymentDueDatePresent,
				),
				validation.Required.Error("payment details are required for final invoices (BR-FR-CO-09)"),
			),
			validation.Skip,
		),
		validation.Field(&inv.Delivery,
			validation.When(
				isConsolidatedCreditNote(inv),
				validation.By(
					validateDelivery,
				),
				validation.Required.Error("delivery details are required for consolidated credit notes (BR-FR-CO-03)"),
			),
			validation.Skip,
		),
		validation.Field(&inv.Totals,
			validation.When(
				isFinalInvoice(inv),
				validation.By(
					validateTotals,
				),
			),
			validation.Skip,
		),
		validation.Field(&inv.Notes,
			validation.By(
				validateMandatoryNotes,
			),
			validation.When(
				isPartyIdentitySTC(inv.Supplier),
				validation.By(
					validateNoteTXD,
				),
			),
			validation.Required.Error("notes are required for French CTC invoices (BR-FR-05)"),
			validation.Skip,
		),
	)
}

// validateInvoiceCode validates the code according to BR-FR-01/02
// The invoice ID is either just the code, or series-code if series is present
func validateCode(series cbc.Code) validation.RuleFunc {
	return func(value any) error {
		code, ok := value.(cbc.Code)
		if !ok || code == cbc.CodeEmpty {
			return nil // Let required validation handle empty codes
		}

		// Construct the full invoice ID
		invoiceID := string(code)
		if series != cbc.CodeEmpty {
			invoiceID = string(series.Join(code))
		}

		if !invoiceCodeRegexp.MatchString(invoiceID) {
			return errors.New("must be 1-35 characters, alphanumeric plus -+_/ (BR-FR-01/02), including the series")
		}
		return nil
	}
}

// validatePrecedingDocument validates preceding document codes
func validatePrecedingDocument(value any) error {
	docRef, ok := value.(*org.DocumentRef)
	if !ok || docRef == nil {
		return nil
	}

	return validation.ValidateStruct(docRef,
		validation.Field(&docRef.Code,
			validation.By(
				validateCode(docRef.Series),
			),
			validation.Skip,
		),
	)
}

func validateBillInvoiceTax(value any) error {
	tx, ok := value.(*bill.Tax)
	if !ok || tx == nil {
		return nil
	}
	return validation.ValidateStruct(tx,
		validation.Field(&tx.Ext,
			tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, allowedDocumentTypes...),
			tax.ExtensionsRequire(ExtKeyBillingMode),
			validation.When(
				// BR-FR-CO-08
				isFactoredExtension(tx.Ext.Get(ExtKeyBillingMode)),
				tax.ExtensionsExcludeCodes(untdid.ExtKeyDocumentType, advancePaymentDocumentTypes...),
			),
			validation.Skip,
		),
	)
}

// validateSupplier validates supplier requirements for French invoices
func validateSupplier(value any) error {
	supplier, ok := value.(*org.Party)
	if !ok || supplier == nil {
		return nil
	}

	return validation.ValidateStruct(supplier,
		validation.Field(&supplier.Inboxes,
			validation.Required.Error("seller electronic address required for French B2B invoices (BR-FR-13)"),
			validation.Skip,
		),
		validation.Field(&supplier.Identities,
			validation.By(
				validateSirenPresent,
			),
			validation.Skip,
		),
	)

}

// validateSupplierInbox validates supplier inbox for B2B non-self-billed invoices (BR-FR-21)
func validateSirenInbox(value any) error {
	supplier, ok := value.(*org.Party)
	if !ok || supplier == nil {
		return nil
	}

	// Get supplier's SIREN
	siren := getPartySIREN(supplier)
	if siren == "" {
		return nil // SIREN validation handled elsewhere
	}

	// Check for SIREN inbox (scheme 0225)
	hasSIRENInbox := false
	for _, inbox := range supplier.Inboxes {
		if inbox != nil && inbox.Scheme == inboxSchemeSIREN {
			hasSIRENInbox = true
			// Validate that inbox code starts with SIREN
			if !strings.HasPrefix(string(inbox.Code), siren) {
				return errors.New("party endpoint ID scheme inbox (0225) must start with SIREN (BR-FR-21/22)")
			}
		}
	}

	if !hasSIRENInbox {
		return errors.New("party must have endpoint ID with scheme 0225 (SIREN) (BR-FR-21/22)")
	}

	return nil

}

// validateFrenchCustomer validates customer requirements for French B2B invoices
func validateCustomer(value any) error {
	customer, ok := value.(*org.Party)
	if !ok || customer == nil {
		return nil
	}

	// BR-FR-13: Buyer electronic address is required for B2B
	return validation.ValidateStruct(customer,
		validation.Field(&customer.Inboxes,
			validation.Required.Error("buyer electronic address required for French B2B invoices (BR-FR-13)"),
			validation.Skip,
		),
	)
}

func validateB2BCustomer(value any) error {
	customer, ok := value.(*org.Party)
	if !ok || customer == nil {
		return nil
	}

	// BR-FR-14: For B2B transactions, customer must have a SIREN identity (iso-scheme-id: 0002)
	return validation.ValidateStruct(customer,
		validation.Field(&customer.Identities,
			validation.By(
				validateSirenPresent,
			),
			validation.Skip,
		),
	)
}

func validateSirenPresent(value any) error {
	identities, ok := value.([]*org.Identity)
	if !ok {
		return nil
	}

	for _, id := range identities {
		if id != nil && id.Ext != nil {
			if code := id.Ext.Get(iso.ExtKeySchemeID); code == "0002" && id.Scope.Has(org.IdentityScopeLegal) {
				return nil
			}
		}
	}

	return errors.New("SIREN identity required for French parties with scheme 0002 and scope legal (BR-FR-10/11)")
}

// validateOrdering validates the ordering details (BR-FR-29)
func validateOrdering(value any) error {
	ordering, ok := value.(*bill.Ordering)
	if !ok || ordering == nil {
		return nil
	}

	return validation.ValidateStruct(ordering,
		validation.Field(&ordering.Identities,
			validation.By(
				validateOrderingIdentities,
			),
			validation.Skip,
		),
	)
}

func validateOrderingIdentities(value any) error {
	identities, ok := value.([]*org.Identity)
	if !ok {
		return nil
	}

	// BR-FR-30: Check for duplicate UNTDID reference schemes AFL and AWW
	// The value requirement for code is forced at an identity level
	var afl, aww bool

	for _, id := range identities {
		if id == nil || id.Ext == nil {
			continue
		}

		// Check UNTDID reference extension
		ref := id.Ext.Get(untdid.ExtKeyReference)
		switch ref.String() {
		case "AFL":
			// BR-FR-30: Only one identity with UNTDID reference AFL allowed
			if afl {
				return errors.New("only one ordering identity with UNTDID reference 'AFL' is allowed (BR-FR-30)")
			}
			afl = true
		case "AWW":
			// BR-FR-30: Only one identity with UNTDID reference AWW allowed
			if aww {
				return errors.New("only one ordering identity with UNTDID reference 'AWW' is allowed (BR-FR-30)")
			}
			aww = true
		}
	}

	return nil
}

func validateOrderingSeller(value any) error {
	ordering, ok := value.(*bill.Ordering)
	if !ok || ordering == nil {
		return nil
	}

	// BR-FR-29: If the supplier is STC, then the seller tax ID must be present
	return validation.ValidateStruct(ordering,
		validation.Field(&ordering.Seller,
			validation.By(
				validateSeller,
			),
			validation.Required.Error("seller is required when supplier is under STC scheme (BR-FR-CO-15)"),
			validation.Skip,
		),
	)
}

func validateSeller(value any) error {
	seller, ok := value.(*org.Party)
	if !ok || seller == nil {
		return nil
	}

	return validation.ValidateStruct(seller,
		validation.Field(&seller.TaxID,
			validation.By(
				validateSellerTaxID,
			),
			validation.Required.Error("tax ID is required when supplier is under STC scheme (BR-FR-CO-15)"),
			validation.Skip,
		),
	)
}

func validateSellerTaxID(value any) error {
	taxID, ok := value.(*tax.Identity)
	if !ok || taxID == nil {
		return nil
	}

	return validation.ValidateStruct(taxID,
		validation.Field(&taxID.Code,
			validation.Required.Error("code is required when supplier is under STC scheme (BR-FR-CO-15)"),
		),
	)
}

func validateOrderingContracts(value any) error {
	ordering, ok := value.(*bill.Ordering)
	if !ok || ordering == nil {
		return nil
	}

	// BR-FR-CO-03: For consolidated credit notes, at least one contract reference is required in the ordering details
	return validation.ValidateStruct(ordering,
		validation.Field(&ordering.Contracts,
			validation.Required.Error("at least one contract reference is required in ordering details for consolidated credit notes (BR-FR-CO-03)"),
			validation.Length(1, 0).Error("at least one contract reference is required in ordering details for consolidated credit notes (BR-FR-CO-03)"),
			validation.Skip,
		),
	)
}

// validatePayment validates that due date is on or after issue date (BR-FR-CO-07)
func validatePayment(issueDate cal.Date) validation.RuleFunc {
	return func(value any) error {
		payment, ok := value.(*bill.PaymentDetails)
		if !ok || payment == nil {
			return nil // No terms or due dates, rule doesn't apply
		}

		return validation.ValidateStruct(payment,
			validation.Field(&payment.Terms,
				validation.By(
					validateTerms(issueDate),
				),
				validation.Skip,
			),
		)
	}
}

func validateTerms(issueDate cal.Date) validation.RuleFunc {
	return func(value any) error {
		terms, ok := value.(*pay.Terms)
		if !ok || terms == nil {
			return nil // No due dates, rule doesn't apply
		}

		return validation.ValidateStruct(terms,
			validation.Field(&terms.DueDates,
				validation.Each(
					validation.By(
						validateDueDate(issueDate),
					),
				),
				validation.Skip,
			),
		)
	}
}

func validateDueDate(issueDate cal.Date) validation.RuleFunc {
	return func(value any) error {
		dueDate, ok := value.(*pay.DueDate)
		if !ok || dueDate == nil {
			return nil // No due date, rule doesn't apply
		}

		return validation.ValidateStruct(dueDate,
			validation.Field(&dueDate.Date,
				cal.DateAfter(issueDate),
				validation.Skip,
			),
		)
	}
}

func validatePaymentDueDatePresent(value any) error {
	payment, ok := value.(*bill.PaymentDetails)
	if !ok || payment == nil {
		return nil // Let required validation handle missing payment details
	}

	return validation.ValidateStruct(payment,
		validation.Field(&payment.Terms,
			validation.By(
				validateTermsDueDatePresent,
			),
			validation.Required.Error("payment terms required for final invoices (BR-FR-CO-09)"),
			validation.Skip,
		),
	)
}

func validateTermsDueDatePresent(value any) error {
	terms, ok := value.(*pay.Terms)
	if !ok || terms == nil {
		return nil // Let required validation handle missing terms
	}

	return validation.ValidateStruct(terms,
		validation.Field(&terms.DueDates,
			validation.Required.Error("at least one due date required for final invoices (BR-FR-CO-09)"),
			validation.Skip,
		),
	)
}

func validateDelivery(value any) error {
	delivery, ok := value.(*bill.DeliveryDetails)
	if !ok || delivery == nil {
		return nil
	}

	return validation.ValidateStruct(delivery,
		validation.Field(&delivery.Period,
			validation.Required.Error("delivery period is required for consolidated credit notes (BR-FR-CO-03)"),
			validation.Skip,
		),
	)
}

// validateTotals validates that for final invoices (B2, S2, M2):
// - BR-FR-CO-09 BT-23-1: the advance amount (BT-113) must equal the tax-inclusive total (BT-112)
// - BR-FR-CO-09 BT-23-2: the payable amount (BT-115) must be 0
// BT-115 maps to Due if present, otherwise Payable.
func validateTotals(value any) error {
	totals, ok := value.(*bill.Totals)
	if !ok || totals == nil {
		return nil
	}

	return validation.ValidateStruct(totals,
		// BT-23-1: PrepaidAmount (Advances) must equal TaxInclusiveAmount (TotalWithTax)
		validation.Field(&totals.Advances,
			validation.Required.Error("advance amount is required for already-paid invoices (BR-FR-CO-09)"),
			num.Equals(totals.TotalWithTax),
			validation.Skip,
		),
		// BT-23-2: PayableAmount must be 0
		// PayableAmount maps to Due if present, otherwise Payable
		validation.Field(&totals.Due,
			num.Equals(num.AmountZero),
			validation.Skip,
		),
		validation.Field(&totals.Payable,
			validation.When(totals.Due == nil,
				num.Equals(num.AmountZero),
			),
			validation.Skip,
		),
	)
}

// validateMandatoryNotes validates that required notes are present and unique (BR-FR-05, BR-FR-06, BR-FR-30)
// BR-FR-05: French CTC requires three mandatory note types:
// - PMT: for the mention of a flat-rate penalty of 40 EUROS for collection costs (org.NoteKeyPayment)
// - PMD: penalty corresponding to the payment terms specific to each company (org.NoteKeyPaymentMethod)
// - AAB: mention of discount or no discount (in BT-22) (org.NoteKeyPaymentTerm)
// BR-FR-06: Each code (PMT, PMD, AAB, TXD) should appear at most once.
// BR-FR-30: BAR code should appear at most once.
func validateMandatoryNotes(value any) error {
	notes, ok := value.([]*org.Note)
	if !ok {
		return nil
	}

	required := []cbc.Code{"PMT", "PMD", "AAB"}
	counts := make(map[cbc.Code]int)

	for _, note := range notes {
		if note != nil && note.Ext != nil {
			if code := note.Ext.Get(untdid.ExtKeyTextSubject); code != cbc.CodeEmpty {
				counts[code]++
			}
		}
	}

	// BR-FR-05: Check for missing required codes
	var missing []string
	for _, code := range required {
		if counts[code] == 0 {
			missing = append(missing, string(code))
		}
	}

	if len(missing) > 0 {
		return errors.New("missing required note codes: " + strings.Join(missing, ", ") + " (BR-FR-05)")
	}

	// BR-FR-06/BR-FR-30: Check for duplicate codes (required codes + TXD + BAR)
	checkUnique := append(required, "TXD", "BAR")
	var duplicates []string
	for _, code := range checkUnique {
		if counts[code] > 1 {
			duplicates = append(duplicates, string(code))
		}
	}

	if len(duplicates) > 0 {
		return errors.New("duplicate note codes found: " + strings.Join(duplicates, ", ") + " (BR-FR-06/BR-FR-30)")
	}

	// Validate BAR note text if present
	for _, note := range notes {
		if note != nil && note.Ext != nil {
			if note.Ext.Get(untdid.ExtKeyTextSubject) == "BAR" {
				if note.Text != "" && !slices.Contains(allowedBARTreatments, note.Text) {
					return errors.New("BAR note text must be one of: B2B, B2BINT, B2C, OUTOFSCOPE, ARCHIVEONLY")
				}
			}
		}
	}

	return nil
}

func validateNoteTXD(value any) error {
	notes, ok := value.([]*org.Note)
	if !ok || len(notes) == 0 {
		return nil
	}

	for _, note := range notes {
		if note != nil && note.Ext != nil {
			if code := note.Ext.Get(untdid.ExtKeyTextSubject); code == "TXD" && note.Text == "MEMBRE_ASSUJETTI_UNIQUE" {
				return nil
			}
		}
	}

	return errors.New("for sellers with STC scheme (0231), a note with code 'TXD' and text 'MEMBRE_ASSUJETTI_UNIQUE' is required (BR-FR-CO-14)")
}

// isB2BTransaction determines if the transaction is B2B (business to business)
// by checking for a note with code "BAR" and text containing "B2B"
func isB2BTransaction(inv *bill.Invoice) bool {
	if inv == nil || len(inv.Notes) == 0 {
		return false
	}

	for _, note := range inv.Notes {
		if note != nil && note.Ext != nil {
			if note.Ext.Get(untdid.ExtKeyTextSubject) == "BAR" && note.Text == "B2B" {
				// Check if note text indicates B2B transaction (B2B or B2BINT)
				return true
			}
		}
	}

	return false
}

// isSelfBilledInvoice checks if the invoice is self-billed based on document type
func isSelfBilledInvoice(inv *bill.Invoice) bool {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext == nil {
		return false
	}

	docType := inv.Tax.Ext.Get(untdid.ExtKeyDocumentType)
	if docType == "" {
		return false
	}

	return slices.Contains(selfBilledDocumentTypes, docType)
}

// isCorrectiveInvoice checks if the invoice is corrective based on document type
func isCorrectiveInvoice(inv *bill.Invoice) bool {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext == nil {
		return false
	}

	docType := inv.Tax.Ext.Get(untdid.ExtKeyDocumentType)
	if docType == "" {
		return false
	}

	return slices.Contains(correctiveInvoiceTypes, docType)
}

func isPartyIdentitySTC(party *org.Party) bool {
	if party == nil || len(party.Identities) == 0 {
		return false
	}

	for _, id := range party.Identities {
		if id != nil && id.Ext != nil {
			if code := id.Ext.Get(iso.ExtKeySchemeID); code == "0231" {
				return true
			}
		}
	}
	return false
}

func isCreditNote(inv *bill.Invoice) bool {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext == nil {
		return false
	}
	docType := inv.Tax.Ext.Get(untdid.ExtKeyDocumentType)
	return slices.Contains(creditNoteTypes, docType)
}

func isConsolidatedCreditNote(inv *bill.Invoice) bool {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext == nil {
		return false
	}
	docType := inv.Tax.Ext.Get(untdid.ExtKeyDocumentType)
	return docType == "262" // Consolidated credit note
}

func isAdvancedInvoice(inv *bill.Invoice) bool {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext == nil {
		return false
	}

	docType := inv.Tax.Ext.Get(untdid.ExtKeyDocumentType)
	return slices.Contains(advancePaymentDocumentTypes, docType)
}

// isFinalInvoice checks if the invoice is a final invoice based on billing mode (B2, S2, M2)
func isFinalInvoice(inv *bill.Invoice) bool {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext == nil {
		return false
	}

	bm := inv.Tax.Ext.Get(ExtKeyBillingMode)
	return bm == BillingModeB2 || bm == BillingModeS2 || bm == BillingModeM2
}

func isFactoredExtension(bm cbc.Code) bool {
	return bm == BillingModeB4 || bm == BillingModeS4 || bm == BillingModeM4
}

// getPartySIREN extracts the SIREN from the party's SIREN identity
func getPartySIREN(party *org.Party) string {
	if party == nil {
		return ""
	}

	// SIREN identity - check by type or ISO scheme ID 0002
	for _, id := range party.Identities {
		if id != nil && (id.Type == fr.IdentityTypeSIREN || (id.Ext != nil && id.Ext[iso.ExtKeySchemeID] == identitySchemeIDSIREN)) {
			return string(id.Code)
		}
	}

	return ""
}
