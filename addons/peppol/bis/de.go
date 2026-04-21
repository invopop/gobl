package bis

import (
	"regexp"
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// validDEInvoiceDocumentTypes is the UNTDID 1001 subset allowed on German
// Peppol invoices (DE-R-017, copied from de/xrechnung).
var validDEInvoiceDocumentTypes = []cbc.Code{
	"326", "380", "384", "389", "381", "875", "876", "877",
}

// deCreditTransferMeansCodes are the UNTDID 4461 codes that imply credit
// transfer for DE-R-023.
var deCreditTransferMeansCodes = []cbc.Code{"30", "58"}

// deCardMeansCodes are the UNTDID 4461 codes that imply card payment for DE-R-024.
var deCardMeansCodes = []cbc.Code{"48", "54", "55"}

// deDirectDebitMeansCode is the UNTDID 4461 code for SEPA direct debit (DE-R-025).
var deDirectDebitMeansCode cbc.Code = "59"

var ibanRe = regexp.MustCompile(`^[A-Z]{2}\d{2}[A-Z0-9]{1,30}$`)

func billInvoiceRulesDE() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(supplierCountryIs(l10n.DE),
			// DE-R-017 (warning): invoice type code restricted.
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("DE-R-017", "German invoice document type must be one of 326, 380, 384, 389, 381, 875, 876, 877 (DE-R-017)",
						is.Func("de doc type", deInvoiceDocumentTypeValid),
					),
				),
			),
			// DE-R-001: payment instructions required.
			rules.Assert("DE-R-001", "payment instructions are required (DE-R-001)",
				is.Func("has payment instructions", hasPaymentInstructions),
			),
			// DE-R-015: buyer reference (Ordering.Code) required.
			rules.Assert("DE-R-015", "buyer reference is required (DE-R-015)",
				is.Func("has buyer reference", hasOrderingCode),
			),
			// DE-R-002..R011 — supplier, customer, and delivery address constraints.
			rules.Field("supplier",
				rules.Assert("DE-R-002", "German supplier must provide a contact (DE-R-002)",
					is.Func("has seller contact group", partyHasContactGroup),
				),
				rules.Field("addresses",
					rules.Assert("DE-R-003", "supplier first address must have a city (DE-R-003)",
						is.Func("has locality", firstAddressHasLocalityPE),
					),
					rules.Assert("DE-R-004", "supplier first address must have a postal code (DE-R-004)",
						is.Func("has code", firstAddressHasCodePE),
					),
				),
				rules.Assert("DE-R-005", "supplier contact name is required (DE-R-005)",
					is.Func("has contact name", partyHasContactName),
				),
				rules.Assert("DE-R-006", "supplier contact telephone is required (DE-R-006)",
					is.Func("has contact telephone", partyHasContactTelephone),
				),
				rules.Assert("DE-R-007", "supplier contact email is required (DE-R-007)",
					is.Func("has contact email", partyHasContactEmail),
				),
				rules.Assert("DE-R-027", "supplier telephone must have at least 3 digits (DE-R-027)",
					is.Func("telephone min length", partyTelephoneMinLength),
				),
				rules.Assert("DE-R-028", "supplier email must be a valid email format (DE-R-028)",
					is.Func("valid email", partyEmailFormat),
				),
			),
			rules.Field("customer",
				rules.Field("addresses",
					rules.Assert("DE-R-008", "customer first address must have a city (DE-R-008)",
						is.Func("has locality", firstAddressHasLocalityPE),
					),
					rules.Assert("DE-R-009", "customer first address must have a postal code (DE-R-009)",
						is.Func("has code", firstAddressHasCodePE),
					),
				),
			),
			rules.Field("delivery",
				rules.Field("receiver",
					rules.Field("addresses",
						rules.AssertIfPresent("DE-R-010", "delivery address must have a city (DE-R-010)",
							is.Func("has locality", firstAddressHasLocalityPE),
						),
						rules.AssertIfPresent("DE-R-011", "delivery address must have a postal code (DE-R-011)",
							is.Func("has code", firstAddressHasCodePE),
						),
					),
				),
			),
			// DE-R-018 (early-payment #SKONTO# note format) is not enforced here.
			// The note is written by the caller into the payment-terms text
			// (gobl.ubl does not synthesize it from bill.Payment.Terms.DueDates);
			// if the caller includes a #SKONTO#-prefixed line, they must format
			// it correctly themselves or the Peppol access point will reject.
			//
			// DE-R-022 (attachment filename uniqueness) is a UBL-level concern
			// governing the cac:AdditionalDocumentReference elements gobl.ubl
			// emits.
			//
			// DE-R-026 (warning): corrective invoices should reference a preceding invoice.
			rules.Assert("DE-R-026", "corrective invoices should reference a preceding invoice (DE-R-026)",
				is.Func("preceding present for corrections", correctivePrecedingPresent),
			),
			// DE-R-014 (VAT category rate percent) is not enforced here. GOBL's
			// tax model only carries a numeric percent on standard/reduced-rated
			// tax combos; exempt, zero, reverse-charge, and not-subject-to-VAT
			// categories structurally have no percent. gobl.ubl emits cbc:Percent
			// for every UBL tax subtotal (0 for the non-standard categories), so
			// the schematron assertion is satisfied by construction at emit time.
			// DE-R-016: tax categories S/Z/E/AE/K/G/L/M require supplier VAT/tax ID.
			rules.Assert("DE-R-016", "supplier must have a VAT or tax registration identifier for these tax categories (DE-R-016)",
				is.Func("supplier tax id for categories", deSupplierHasTaxIDForCategory),
			),
		),
	)
}

func payInstructionsRulesDE() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(supplierCountryIs(l10n.DE),
			rules.Field("payment",
				rules.Field("instructions",
					// DE-R-023: credit transfer codes require CreditTransfer; exclude card/debit.
					rules.Assert("DE-R-023", "credit transfer code requires credit transfer details and excludes card/direct debit (DE-R-023)",
						is.Func("de credit transfer exclusive", deCreditTransferExclusive),
					),
					// DE-R-024: card codes require Card; exclude credit transfer/direct debit.
					rules.Assert("DE-R-024", "card code requires card details and excludes credit transfer/direct debit (DE-R-024)",
						is.Func("de card exclusive", deCardExclusive),
					),
					// DE-R-025: direct debit code 59 requires DirectDebit; exclude credit transfer/card.
					rules.Assert("DE-R-025", "direct debit code requires direct debit details and excludes credit transfer/card (DE-R-025)",
						is.Func("de direct debit exclusive", deDirectDebitExclusive),
					),
					// DE-R-030/R031: direct debit requires creditor + debited account.
					rules.Assert("DE-R-030-031", "direct debit requires creditor and account (DE-R-030, DE-R-031)",
						is.Func("de direct debit fields", deDirectDebitFieldsComplete),
					),
					// DE-R-019 (warn): SEPA credit transfer (code 58) should use IBAN.
					rules.Assert("DE-R-019", "SEPA credit transfer account should be a valid IBAN (DE-R-019)",
						is.Func("de sepa iban", deSEPAIBANValid),
					),
					// DE-R-020 (warn): SEPA direct debit (code 59) should use IBAN.
					rules.Assert("DE-R-020", "SEPA direct debit account should be a valid IBAN (DE-R-020)",
						is.Func("de sepa debit iban", deSEPADebitIBANValid),
					),
				),
			),
		),
	)
}

// --- helpers ---

func deInvoiceDocumentTypeValid(val any) bool {
	// bill.Tax.Ext is tax.Extensions; check via the Get-accessor interface
	// rather than asserting on the concrete map type so callers can pass
	// either form.
	type getter interface {
		Get(cbc.Key) cbc.Code
	}
	g, ok := val.(getter)
	if !ok {
		return true
	}
	code := g.Get(untdid.ExtKeyDocumentType)
	if code == "" {
		return true
	}
	return code.In(validDEInvoiceDocumentTypes...)
}

func partyHasContactGroup(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	if len(p.People) > 0 || len(p.Telephones) > 0 || len(p.Emails) > 0 {
		return true
	}
	return false
}

func firstAddressHasLocalityPE(val any) bool {
	addrs, ok := val.([]*org.Address)
	if !ok || len(addrs) == 0 {
		return true // presence enforced elsewhere
	}
	return addrs[0] != nil && addrs[0].Locality != ""
}

func firstAddressHasCodePE(val any) bool {
	addrs, ok := val.([]*org.Address)
	if !ok || len(addrs) == 0 {
		return true
	}
	return addrs[0] != nil && addrs[0].Code != ""
}

func partyHasContactName(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	if len(p.People) == 0 {
		return false
	}
	// Person.Name is a *org.Name struct — check name presence.
	first := p.People[0]
	return first != nil && first.Name != nil && first.Name.Given != ""
}

func partyHasContactTelephone(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	if len(p.Telephones) > 0 {
		return true
	}
	if len(p.People) > 0 && p.People[0] != nil && len(p.People[0].Telephones) > 0 {
		return true
	}
	return false
}

func partyHasContactEmail(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	if len(p.Emails) > 0 {
		return true
	}
	if len(p.People) > 0 && p.People[0] != nil && len(p.People[0].Emails) > 0 {
		return true
	}
	return false
}

func partyTelephoneMinLength(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	phones := p.Telephones
	if len(phones) == 0 && len(p.People) > 0 && p.People[0] != nil {
		phones = p.People[0].Telephones
	}
	for _, t := range phones {
		if t == nil {
			continue
		}
		digits := 0
		for _, c := range t.Number {
			if c >= '0' && c <= '9' {
				digits++
			}
		}
		if digits < 3 {
			return false
		}
	}
	return true
}

func partyEmailFormat(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	emails := p.Emails
	if len(emails) == 0 && len(p.People) > 0 && p.People[0] != nil {
		emails = p.People[0].Emails
	}
	for _, e := range emails {
		if e == nil {
			continue
		}
		addr := e.Address
		if strings.Count(addr, "@") != 1 {
			return false
		}
		idx := strings.Index(addr, "@")
		if idx <= 0 || idx >= len(addr)-1 {
			return false
		}
	}
	return true
}

func correctivePrecedingPresent(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Tax == nil {
		return true
	}
	code := inv.Tax.Ext.Get(untdid.ExtKeyDocumentType)
	if code != "384" { // only corrective invoices
		return true
	}
	return len(inv.Preceding) > 0
}

func deSupplierHasTaxIDForCategory(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Totals == nil || inv.Totals.Taxes == nil {
		return true
	}
	categoriesRequiringTaxID := []cbc.Code{"S", "Z", "E", "AE", "K", "G", "L", "M"}
	hasSuch := false
	for _, cat := range inv.Totals.Taxes.Categories {
		if cat == nil {
			continue
		}
		for _, rt := range cat.Rates {
			if rt == nil {
				continue
			}
			if rt.Ext.Get(untdid.ExtKeyTaxCategory).In(categoriesRequiringTaxID...) {
				hasSuch = true
				break
			}
		}
	}
	if !hasSuch {
		return true
	}
	if inv.Supplier == nil {
		return false
	}
	// DE-R-016 requires a VAT identifier (TaxID) or tax registration identifier
	// (legal-scope identity). Any other identity (DUNS, GLN, custom) does not
	// satisfy the rule.
	if inv.Supplier.TaxID != nil && inv.Supplier.TaxID.Code != "" {
		return true
	}
	for _, id := range inv.Supplier.Identities {
		if id == nil {
			continue
		}
		if id.Scope == org.IdentityScopeLegal && id.Code != "" {
			return true
		}
	}
	return false
}

func deCreditTransferExclusive(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return true
	}
	code := instr.Ext.Get(untdid.ExtKeyPaymentMeans)
	if !code.In(deCreditTransferMeansCodes...) {
		return true
	}
	if len(instr.CreditTransfer) == 0 {
		return false
	}
	return instr.Card == nil && instr.DirectDebit == nil
}

func deCardExclusive(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return true
	}
	code := instr.Ext.Get(untdid.ExtKeyPaymentMeans)
	if !code.In(deCardMeansCodes...) {
		return true
	}
	if instr.Card == nil {
		return false
	}
	return len(instr.CreditTransfer) == 0 && instr.DirectDebit == nil
}

func deDirectDebitExclusive(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return true
	}
	code := instr.Ext.Get(untdid.ExtKeyPaymentMeans)
	if code != deDirectDebitMeansCode {
		return true
	}
	if instr.DirectDebit == nil {
		return false
	}
	return len(instr.CreditTransfer) == 0 && instr.Card == nil
}

func deDirectDebitFieldsComplete(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil || instr.DirectDebit == nil {
		return true
	}
	return instr.DirectDebit.Creditor != "" && instr.DirectDebit.Account != ""
}

func deSEPAIBANValid(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return true
	}
	code := instr.Ext.Get(untdid.ExtKeyPaymentMeans)
	if code != "58" {
		return true
	}
	for _, ct := range instr.CreditTransfer {
		if ct == nil {
			continue
		}
		acc := ct.IBAN
		if acc == "" {
			acc = ct.Number
		}
		if acc != "" && !ibanRe.MatchString(strings.ReplaceAll(strings.ToUpper(acc), " ", "")) {
			return false
		}
	}
	return true
}

func deSEPADebitIBANValid(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return true
	}
	code := instr.Ext.Get(untdid.ExtKeyPaymentMeans)
	if code != "59" || instr.DirectDebit == nil {
		return true
	}
	acc := instr.DirectDebit.Account
	if acc == "" {
		return true
	}
	return ibanRe.MatchString(strings.ReplaceAll(strings.ToUpper(acc), " ", ""))
}
