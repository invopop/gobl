package ctc

import (
	"regexp"
	"slices"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
)

// BR-FR-01/02: Invoice code validation
// Max 35 characters, alphanumeric plus -+_/
var invoiceCodeRegexp = regexp.MustCompile(`^[A-Za-z0-9\-\+_/]{1,35}$`)

// BR-FR-04: Allowed UNTDID document types for French CTC
var allowedDocumentTypes = []cbc.Code{
	"380", // Commercial invoice
	"389", // Self-billed invoice
	"393", // Factoring invoice
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
	"396", // Factoring credit note
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
	"396", // Factoring credit note
	"502", // Self-billed corrective
	"503", // Self-billed credit for claim
}

var advancePaymentDocumentTypes = []cbc.Code{
	"386", // Advance payment invoice
	"500", // Self-billed advance payment
	"503", // Self-billed credit for claim
}

// Allowed attachment description values for French CTC (BR-FR-17)
var allowedAttachmentDescriptions = []string{
	"RIB",                        // Bank account details (Relevé d'Identité Bancaire)
	"LISIBLE",                    // Human-readable representation of the invoice
	"FEUILLE_DE_STYLE",           // Style sheet for document presentation
	"PJA",                        // Additional supporting document (Pièce Jointe Additionnelle)
	"BON_LIVRAISON",              // Delivery note
	"BON_COMMANDE",               // Purchase order
	"DOCUMENT_ANNEXE",            // Annex document
	"BORDEREAU_SUIVI",            // Follow-up form
	"BORDEREAU_SUIVI_VALIDATION", // Validated follow-up form
	"ETAT_ACOMPTE",               // Payment status statement
	"FACTURE_PAIEMENT_DIRECT",    // Direct payment invoice
	"RECAPITULATIF_COTRAITANCE",  // Co-contracting summary
}

const (
	// attachmentFormatLisible is the attachment format category for BR-FR-18
	attachmentFormatLisible = "LISIBLE"
)

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

// isB2BTransaction determines if the transaction is B2B (business to business)
// by checking for a note with code "BAR" and text containing "B2B"
func isB2BTransaction(inv *bill.Invoice) bool {
	if inv == nil || len(inv.Notes) == 0 {
		return false
	}

	for _, note := range inv.Notes {
		if note != nil && !note.Ext.IsZero() {
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
	if inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
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
	if inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
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
		if id != nil && !id.Ext.IsZero() {
			if code := id.Ext.Get(iso.ExtKeySchemeID); code == "0231" {
				return true
			}
		}
	}
	return false
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
	return docType == "262" // Consolidated credit note
}

func isAdvancedInvoice(inv *bill.Invoice) bool {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
		return false
	}

	docType := inv.Tax.Ext.Get(untdid.ExtKeyDocumentType)
	return slices.Contains(advancePaymentDocumentTypes, docType)
}

// isFinalInvoice checks if the invoice is a final invoice based on billing mode (B2, S2, M2)
func isFinalInvoice(inv *bill.Invoice) bool {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
		return false
	}

	bm := inv.Tax.Ext.Get(ExtKeyBillingMode)
	return bm == BillingModeB2 || bm == BillingModeS2 || bm == BillingModeM2
}

func isFactoringExtension(bm cbc.Code) bool {
	return bm == BillingModeB4 || bm == BillingModeS4 || bm == BillingModeM4
}

// getPartySIREN extracts the SIREN from the party's SIREN identity
func getPartySIREN(party *org.Party) string {
	if party == nil {
		return ""
	}

	// SIREN identity - check by type or ISO scheme ID 0002
	for _, id := range party.Identities {
		if id != nil && (id.Type == fr.IdentityTypeSIREN || (!id.Ext.IsZero() && id.Ext.Get(iso.ExtKeySchemeID) == identitySchemeIDSIREN)) {
			return string(id.Code)
		}
	}

	return ""
}
