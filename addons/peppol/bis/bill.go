package bis

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// validInvoiceUNTDIDCodes is the UNTDID 1001 subset allowed on invoice-type
// documents under Peppol BIS 3.0 (PEPPOL-EN16931-P0100).
var validInvoiceUNTDIDCodes = []cbc.Code{"380", "383", "386", "389", "751"}

// validCreditNoteUNTDIDCodes is the UNTDID 1001 subset allowed on credit-note
// documents under Peppol BIS 3.0 (PEPPOL-EN16931-P0101).
var validCreditNoteUNTDIDCodes = []cbc.Code{"381", "396", "261"}

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		// PEPPOL-EN16931-R002: at most one note unless both parties are DK.
		rules.Field("notes",
			rules.Assert("R002", "at most one note allowed unless both buyer and seller are Danish (PEPPOL-EN16931-R002)",
				is.Func("notes cardinality", notesCardinalityValid),
			),
		),
		// PEPPOL-EN16931-R003: buyer reference or purchase order reference required.
		rules.Assert("R003", "buyer reference or purchase order reference is required (PEPPOL-EN16931-R003)",
			is.Func("buyer reference", hasBuyerReferenceOrPO),
		),
		// PEPPOL-EN16931-R080: at most one project reference.
		rules.Field("ordering",
			rules.Field("projects",
				rules.Assert("R080", "only one project reference allowed (PEPPOL-EN16931-R080)",
					is.Length(0, 1),
				),
			),
		),
		// PEPPOL-EN16931-P0100/P0101: restrict UNTDID document type codes per invoice type.
		rules.Assert("P0100", "invoice type code must be one of 380, 383, 386, 389, 751 (PEPPOL-EN16931-P0100)",
			is.Func("invoice type code", invoiceTypeCodeValid),
		),
		rules.Assert("P0101", "credit note type code must be one of 381, 396, 261 (PEPPOL-EN16931-P0101)",
			is.Func("credit note type code", creditNoteTypeCodeValid),
		),
		// PEPPOL-EN16931-P0112: 326 and 384 only when both parties are IT.
		rules.Assert("P0112", "invoice type 326 or 384 only allowed when both buyer and seller are Italian (PEPPOL-EN16931-P0112)",
			is.Func("partial/corrective IT-only", partialCorrectiveITOnly),
		),
	)
}

// notesCardinalityValid returns true when the invoice has ≤1 note, OR when
// both supplier and customer postal addresses are in DK (schematron targets
// cac:PostalAddress/cbc:Country/cbc:IdentificationCode, not the tax country).
func notesCardinalityValid(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	if len(inv.Notes) <= 1 {
		return true
	}
	return partyAddressCountry(inv.Supplier) == l10n.DK &&
		partyAddressCountry(inv.Customer) == l10n.DK
}

// hasBuyerReferenceOrPO returns true when the invoice carries either an
// ordering code (buyer reference) or at least one purchase order reference.
func hasBuyerReferenceOrPO(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	if inv.Ordering == nil {
		return false
	}
	if inv.Ordering.Code != "" {
		return true
	}
	for _, p := range inv.Ordering.Purchases {
		if p != nil && p.Code != "" {
			return true
		}
	}
	return false
}

// invoiceTypeCodeValid checks the UNTDID 1001 document type code is allowed
// for the invoice type. Only applies when the invoice is a non-credit-note
// variant.
func invoiceTypeCodeValid(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Tax == nil {
		return true
	}
	if inv.Type == bill.InvoiceTypeCreditNote {
		return true // handled by creditNoteTypeCodeValid
	}
	code := inv.Tax.Ext.Get(untdid.ExtKeyDocumentType)
	if code == "" {
		return true // presence enforced elsewhere
	}
	return code.In(validInvoiceUNTDIDCodes...)
}

// creditNoteTypeCodeValid checks the UNTDID 1001 code is allowed for credit notes.
func creditNoteTypeCodeValid(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Tax == nil {
		return true
	}
	if inv.Type != bill.InvoiceTypeCreditNote {
		return true
	}
	code := inv.Tax.Ext.Get(untdid.ExtKeyDocumentType)
	if code == "" {
		return true
	}
	return code.In(validCreditNoteUNTDIDCodes...)
}

// partialCorrectiveITOnly returns false if the invoice uses UNTDID codes 326
// (partial invoice) or 384 (corrected invoice) but either party is non-Italian.
func partialCorrectiveITOnly(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Tax == nil {
		return true
	}
	code := inv.Tax.Ext.Get(untdid.ExtKeyDocumentType)
	if code != "326" && code != "384" {
		return true
	}
	return partyCountry(inv.Supplier) == l10n.IT && partyCountry(inv.Customer) == l10n.IT
}

// hasPaymentInstructions is true when the invoice carries payment instructions.
// Used by country rules (DE-R-001, NL-R-007) that require suppliers to state
// how they expect to be paid.
func hasPaymentInstructions(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	return inv.Payment != nil && inv.Payment.Instructions != nil
}

// hasOrderingCode is true when the invoice carries a buyer reference via
// `Ordering.Code`. Used by DE-R-015.
func hasOrderingCode(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	return inv.Ordering != nil && inv.Ordering.Code != ""
}
