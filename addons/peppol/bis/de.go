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

// deCreditTransferMeansCodes are the UNTDID 4461 codes that imply credit
// transfer for DE-R-023.
var deCreditTransferMeansCodes = []cbc.Code{"30", "58"}

// deCardMeansCodes are the UNTDID 4461 codes that imply card payment for DE-R-024.
var deCardMeansCodes = []cbc.Code{"48", "54", "55"}

// deDirectDebitMeansCode is the UNTDID 4461 code for SEPA direct debit (DE-R-025).
var deDirectDebitMeansCode cbc.Code = "59"

// skontoRe matches the DE-R-018 early-payment discount line format:
//
//	#SKONTO#TAGE=N#PROZENT=N(.NN)#[BASISBETRAG=-N(.NN)#]
var skontoRe = regexp.MustCompile(`^#SKONTO#TAGE=\d+#PROZENT=\d+(\.\d{1,2})?#(BASISBETRAG=-?\d+(\.\d{1,2})?#)?$`)

func billInvoiceRulesDE() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(supplierCountryIs(l10n.DE),
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
			// DE-R-018: any #SKONTO# line the caller writes into the payment
			// terms text must match the fixed format. gobl.ubl does not
			// synthesize these — they're hand-authored — so we catch format
			// slips at the GOBL layer before they reach a Peppol access point.
			rules.Assert("DE-R-018", "early-payment discount note must follow the #SKONTO# format (DE-R-018)",
				is.Func("skonto format", skontoFormatValid),
			),
			// DE-R-014 (VAT category rate percent, fatal) and DE-R-022
			// (attachment filename uniqueness, fatal) are UBL-level concerns
			// owned by gobl.ubl: percent is emitted for every tax subtotal
			// (0 for non-standard categories), and attachment filenames live
			// inside the cac:AdditionalDocumentReference elements the
			// converter produces.
			//
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
				),
			),
		),
	)
}

// --- helpers ---

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

// skontoFormatValid enforces DE-R-018: each `#`-prefixed line in the payment
// terms text must match the fixed Skonto format. The text is sourced from
// pay.Terms.Notes (single due-date case) and any pay.DueDate.Notes (multiple
// due-date case) — both feed cac:PaymentTerms/cbc:Note in the UBL output.
func skontoFormatValid(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Payment == nil || inv.Payment.Terms == nil {
		return true
	}
	if !skontoLinesValid(inv.Payment.Terms.Notes) {
		return false
	}
	for _, d := range inv.Payment.Terms.DueDates {
		if d == nil {
			continue
		}
		if !skontoLinesValid(d.Notes) {
			return false
		}
	}
	return true
}

// skontoLinesValid returns true when every `#`-prefixed line in s matches the
// Skonto regex. Lines not starting with `#` are ignored so callers can keep
// an intro paragraph (e.g. "Payment within 30 days net.") before the block.
func skontoLinesValid(s string) bool {
	if s == "" {
		return true
	}
	for _, line := range strings.Split(s, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || trimmed[0] != '#' {
			continue
		}
		if !skontoRe.MatchString(trimmed) {
			return false
		}
	}
	return true
}
