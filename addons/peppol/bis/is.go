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

// NoteSrcEINDAGI marks an org.Note as the Icelandic EINDAGI ("due date")
// supporting document description. gobl.ubl emits notes with this source as
// a cac:AdditionalDocumentReference with cbc:DocumentDescription="EINDAGI";
// IS-R-008/R-009/R-010 validate the hand-authored YYYY-MM-DD value the
// caller writes into Note.Text.
const NoteSrcEINDAGI cbc.Key = "eindagi"

// isEINDAGIDateRe matches the YYYY-MM-DD format required by IS-R-008.
var isEINDAGIDateRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func billInvoiceRulesIS() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(supplierCountryIs(l10n.IS),
			rules.Field("supplier",
				rules.Assert("IS-01", "Icelandic supplier must have a legal identity (IS-R-002)",
					is.Func("is supplier legal", partyHasLegalIdentity),
				),
				rules.Field("addresses",
					rules.Assert("IS-02", "Icelandic supplier address must have street and postcode (IS-R-003)",
						is.Func("is address complete", firstAddressStreetAndCode),
					),
				),
			),
			rules.Field("customer",
				rules.When(is.Func("customer is IS", func(val any) bool { return partyCountry(valAsParty(val)) == l10n.IS }),
					rules.Assert("IS-03", "Icelandic customer must have a legal identity (IS-R-004)",
						is.Func("is customer legal", partyHasLegalIdentity),
					),
					rules.Field("addresses",
						rules.Assert("IS-04", "Icelandic customer address must have street and postcode (IS-R-005)",
							is.Func("is customer address complete", firstAddressStreetAndCode),
						),
					),
				),
			),
			// IS-R-008/R-009/R-010: EINDAGI is an org.Note with Src="eindagi"
			// (see NoteSrcEINDAGI). Callers author the note by hand; we catch
			// format, due-date-presence, and date-order mistakes here.
			rules.Assert("IS-05", "EINDAGI note text must use YYYY-MM-DD format (IS-R-008)",
				is.Func("is eindagi format", isEINDAGIFormatValid),
			),
			rules.Assert("IS-06", "invoice with EINDAGI note must include a payment due date (IS-R-009)",
				is.Func("is eindagi due-date present", isEINDAGIDueDatePresent),
			),
			rules.Assert("IS-07", "EINDAGI date must be on or after the first due date (IS-R-010)",
				is.Func("is eindagi after due", isEINDAGIAfterFirstDue),
			),
		),
	)
}

func payInstructionsRulesIS() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(supplierCountryIs(l10n.IS),
			rules.Field("payment",
				rules.Field("instructions",
					rules.Assert("IS-08", "Icelandic payment means 9 requires 12-digit account (IS-R-006)",
						is.Func("is 9", isPaymentCode9Account),
					),
					rules.Assert("IS-09", "Icelandic payment means 42 requires 12-digit account (IS-R-007)",
						is.Func("is 42", isPaymentCode42Account),
					),
				),
			),
		),
	)
}

func valAsParty(v any) *org.Party {
	p, ok := v.(*org.Party)
	if !ok {
		return nil
	}
	return p
}

func partyHasLegalIdentity(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	for _, id := range p.Identities {
		if id != nil && id.Scope == org.IdentityScopeLegal {
			return true
		}
	}
	return p.TaxID != nil && p.TaxID.Code != ""
}

func firstAddressStreetAndCode(val any) bool {
	addrs, ok := val.([]*org.Address)
	if !ok || len(addrs) == 0 {
		return true // presence is enforced elsewhere
	}
	a := addrs[0]
	if a == nil {
		return false
	}
	return a.Street != "" && a.Code != ""
}

// validISAccount accepts either a 12-digit Icelandic domestic account or an
// IS-prefix IBAN (IS + 24 alphanumeric chars, 26 total).
func validISAccount(s string) bool {
	if s == "" {
		return false
	}
	if len(s) == 12 && onlyDigits(s) {
		return true
	}
	upper := strings.ToUpper(strings.ReplaceAll(s, " ", ""))
	if len(upper) == 26 && strings.HasPrefix(upper, "IS") {
		return true
	}
	return false
}

// paymentCreditTransferHasValidAccount returns true when every credit transfer
// entry carries a valid account (IBAN preferred, Number as fallback).
func paymentCreditTransferHasValidAccount(instr *pay.Instructions) bool {
	if len(instr.CreditTransfer) == 0 {
		return false
	}
	for _, ct := range instr.CreditTransfer {
		if ct == nil {
			return false
		}
		if !validISAccount(ct.IBAN) && !validISAccount(ct.Number) {
			return false
		}
	}
	return true
}

func isPaymentCode9Account(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return true
	}
	code := instr.Ext.Get(untdid.ExtKeyPaymentMeans)
	if code != "9" {
		return true
	}
	return paymentCreditTransferHasValidAccount(instr)
}

func isPaymentCode42Account(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return true
	}
	code := instr.Ext.Get(untdid.ExtKeyPaymentMeans)
	if code != "42" {
		return true
	}
	return paymentCreditTransferHasValidAccount(instr)
}

// eindagiNotes returns all notes on the invoice marked as EINDAGI.
func eindagiNotes(inv *bill.Invoice) []*org.Note {
	if inv == nil {
		return nil
	}
	var out []*org.Note
	for _, n := range inv.Notes {
		if n != nil && n.Src == NoteSrcEINDAGI {
			out = append(out, n)
		}
	}
	return out
}

// isEINDAGIFormatValid checks IS-R-008: any EINDAGI note must carry a
// YYYY-MM-DD date in its Text field.
func isEINDAGIFormatValid(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	for _, n := range eindagiNotes(inv) {
		if !isEINDAGIDateRe.MatchString(strings.TrimSpace(n.Text)) {
			return false
		}
	}
	return true
}

// isEINDAGIDueDatePresent checks IS-R-009: an EINDAGI note requires at
// least one Payment.Terms.DueDates entry with a Date set (BT-9).
func isEINDAGIDueDatePresent(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	if len(eindagiNotes(inv)) == 0 {
		return true
	}
	if inv.Payment == nil || inv.Payment.Terms == nil {
		return false
	}
	for _, d := range inv.Payment.Terms.DueDates {
		if d != nil && d.Date != nil {
			return true
		}
	}
	return false
}

// isEINDAGIAfterFirstDue checks IS-R-010: the EINDAGI date must be on or
// after the first Payment.Terms.DueDates date. YYYY-MM-DD strings compare
// correctly lexicographically.
func isEINDAGIAfterFirstDue(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	notes := eindagiNotes(inv)
	if len(notes) == 0 {
		return true
	}
	if inv.Payment == nil || inv.Payment.Terms == nil || len(inv.Payment.Terms.DueDates) == 0 {
		return true // IS-R-009 already catches the missing-due-date case
	}
	var firstDue string
	for _, d := range inv.Payment.Terms.DueDates {
		if d != nil && d.Date != nil {
			firstDue = d.Date.String()
			break
		}
	}
	if firstDue == "" {
		return true
	}
	for _, n := range notes {
		text := strings.TrimSpace(n.Text)
		if !isEINDAGIDateRe.MatchString(text) {
			continue // IS-R-008 catches format errors
		}
		if text < firstDue {
			return false
		}
	}
	return true
}
