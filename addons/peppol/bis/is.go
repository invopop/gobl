package bis

import (
	"regexp"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// validISInvoiceDocumentTypes is the IS-R-001 recommended UNTDID 1001 subset.
var validISInvoiceDocumentTypes = []cbc.Code{"380", "381"}

var isAccountRe = regexp.MustCompile(`^\d{12}$`)

// isEINDAGIDateRe matches the YYYY-MM-DD format required by IS-R-008.
var isEINDAGIDateRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func billInvoiceRulesIS() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(supplierCountryIs(l10n.IS),
			// IS-R-001 (warning): invoice type should be 380 or 381.
			rules.Field("tax",
				rules.Field("ext",
					rules.AssertIfPresent("IS-R-001", "Icelandic invoice document type should be 380 or 381 (IS-R-001)",
						is.Func("is doc type", isDocumentTypeValid),
					),
				),
			),
			// IS-R-002: Icelandic supplier must have legal identity.
			rules.Field("supplier",
				rules.Assert("IS-R-002", "Icelandic supplier must have a legal identity (IS-R-002)",
					is.Func("is supplier legal", partyHasLegalIdentity),
				),
				rules.Field("addresses",
					rules.Assert("IS-R-003", "Icelandic supplier address must have street and postcode (IS-R-003)",
						is.Func("is address complete", firstAddressStreetAndCode),
					),
				),
			),
			rules.Field("customer",
				rules.When(is.Func("customer is IS", func(val any) bool { return partyCountry(valAsParty(val)) == l10n.IS }),
					rules.Assert("IS-R-004", "Icelandic customer must have a legal identity (IS-R-004)",
						is.Func("is customer legal", partyHasLegalIdentity),
					),
					rules.Field("addresses",
						rules.Assert("IS-R-005", "Icelandic customer address must have street and postcode (IS-R-005)",
							is.Func("is customer address complete", firstAddressStreetAndCode),
						),
					),
				),
			),
			// IS-R-008/R-009/R-010: EINDAGI note handling.
			rules.Assert("IS-R-008", "Icelandic EINDAGI note must be in YYYY-MM-DD format (IS-R-008)",
				is.Func("is eindagi format", isEINDAGIFormat),
			),
			rules.Assert("IS-R-009", "Icelandic invoice with EINDAGI must have a due date (IS-R-009)",
				is.Func("is eindagi due", isEINDAGIDuePresent),
			),
			rules.Assert("IS-R-010", "Icelandic EINDAGI date must be on or after due date (IS-R-010)",
				is.Func("is eindagi order", isEINDAGIAfterDue),
			),
		),
	)
}

func payInstructionsRulesIS() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(supplierCountryIs(l10n.IS),
			rules.Field("payment",
				rules.Field("instructions",
					rules.Assert("IS-R-006", "Icelandic payment means 9 requires 12-digit account (IS-R-006)",
						is.Func("is 9", isPaymentCode9Account),
					),
					rules.Assert("IS-R-007", "Icelandic payment means 42 requires 12-digit account (IS-R-007)",
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

func isDocumentTypeValid(val any) bool {
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
	return code.In(validISInvoiceDocumentTypes...)
}

func partyHasLegalIdentity(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	for _, id := range p.Identities {
		if id != nil && id.Scope == "legal" {
			return true
		}
	}
	return p.TaxID != nil && p.TaxID.Code != ""
}

func firstAddressStreetAndCode(val any) bool {
	addrs, ok := val.([]*org.Address)
	if !ok || len(addrs) == 0 {
		return true
	}
	a := addrs[0]
	if a == nil {
		return true
	}
	return a.Street != "" && a.Code != ""
}

func isEINDAGIFormat(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	for _, n := range inv.Notes {
		if n == nil {
			continue
		}
		if n.Src == "EINDAGI" {
			if !isEINDAGIDateRe.MatchString(n.Text) {
				return false
			}
		}
	}
	return true
}

func isEINDAGIDuePresent(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	hasEINDAGI := false
	for _, n := range inv.Notes {
		if n != nil && n.Src == "EINDAGI" {
			hasEINDAGI = true
			break
		}
	}
	if !hasEINDAGI {
		return true
	}
	return inv.Payment != nil && inv.Payment.Terms != nil && len(inv.Payment.Terms.DueDates) > 0
}

func isEINDAGIAfterDue(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Payment == nil || inv.Payment.Terms == nil || len(inv.Payment.Terms.DueDates) == 0 {
		return true
	}
	firstDue := inv.Payment.Terms.DueDates[0]
	if firstDue == nil || firstDue.Date == nil {
		return true
	}
	for _, n := range inv.Notes {
		if n == nil || n.Src != "EINDAGI" {
			continue
		}
		if !isEINDAGIDateRe.MatchString(n.Text) {
			continue
		}
		// Compare by string comparison since both are YYYY-MM-DD.
		if n.Text < firstDue.Date.String() {
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
	if len(instr.CreditTransfer) == 0 {
		return false
	}
	for _, ct := range instr.CreditTransfer {
		if ct == nil || !isAccountRe.MatchString(ct.Number) {
			return false
		}
	}
	return true
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
	if len(instr.CreditTransfer) == 0 {
		return false
	}
	for _, ct := range instr.CreditTransfer {
		if ct == nil || !isAccountRe.MatchString(ct.Number) {
			return false
		}
	}
	return true
}
