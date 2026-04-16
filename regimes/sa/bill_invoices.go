package sa

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
			is.InContext(tax.RegimeIn(l10n.SA.Tax())),
			rules.Field("supplier",
				rules.Assert("01", "supplier must have a valid tax ID code. An additional identity is optional but must be at most 1",
					is.Func("has tax ID code or CRN/MOM/MLS/700/SAG/OTH identity (BR-KSA-39), (BR-KSA-08)", hasValidTaxIDAndIdentities),
				),
			),
		),
	)
}

func hasValidTaxIDAndIdentities(value any) bool {
	party, _ := value.(*org.Party)
	if party == nil {
		return false
	}
	return hasTaxIDCode(party) && hasAtMostOneIdentity(party)
}

func hasTaxIDCode(party *org.Party) bool {
	return party.TaxID != nil && party.TaxID.Code != ""
}

func hasAtMostOneIdentity(party *org.Party) bool {
	if len(party.Identities) == 0 {
		return true
	}
	if len(party.Identities) == 1 && org.IdentitiesTypeIn(supplierValidIdentities...).Check(party.Identities) {
		return true
	}
	return false
}
