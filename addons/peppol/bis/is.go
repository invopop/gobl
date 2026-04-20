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
			// IS-R-008/R-009/R-010 (EINDAGI) deferred to gobl.ubl — see deferred.go.
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
