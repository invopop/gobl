package be

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.InContext(tax.RegimeIn(l10n.BE.Tax())),
			rules.Field("supplier",
				rules.When(
					is.Func("no BCE identity", supplierNoBCEIdentity),
					rules.Field("tax_id",
						rules.Assert("01", "supplier tax ID required for Belgian regime", is.Present),
						rules.Field("code",
							rules.Assert("02", "supplier tax ID code required for Belgian regime", is.Present),
						),
					),
				),
				rules.When(
					is.Func("no tax ID code", supplierNoTaxIDCode),
					rules.Field("identities",
						rules.Assert("03", "supplier identities must include BCE type",
							is.Func("has BCE type", identitiesIncludeBCE)),
					),
				),
			),
		),
	)
}

func supplierNoBCEIdentity(val any) bool {
	p, _ := val.(*org.Party)
	return !hasIdentityBCE(p)
}

func supplierNoTaxIDCode(val any) bool {
	p, _ := val.(*org.Party)
	return !hasTaxIDCode(p)
}

func identitiesIncludeBCE(val any) bool {
	idents, _ := val.([]*org.Identity)
	return org.IdentityForType(idents, IdentityTypeBCE) != nil
}

func hasTaxIDCode(party *org.Party) bool {
	return party != nil && party.TaxID != nil && party.TaxID.Code != ""
}

func hasIdentityBCE(party *org.Party) bool {
	if party == nil || len(party.Identities) == 0 {
		return false
	}
	return org.IdentityForType(party.Identities, IdentityTypeBCE) != nil
}
