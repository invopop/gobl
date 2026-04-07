package sa

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
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
				rules.Assert("01", "supplier must have a valid tax ID code or identity",
					is.Func("has tax ID code or CRN/MOM/MLS/700/SAG/OTH identity (BR-KSA-08)", hasSupplierTaxIDOrIdentity),
				),
			),

			rules.Field("customer",
				rules.Assert("03", "customer must have a valid tax ID code or identity",
					is.Func("hast tax id code or TIN/CRN/Mom/MLS/700/SAG/National/Gcc/Iqa/Passport/OTH (BR-KSA-14)", hasCustomerTaxIDOrIdentity),
				),
			),
		),
	)
}

func hasSupplierTaxIDOrIdentity(value any) bool {
	party, _ := value.(*org.Party)
	return hasTaxIDCode(party) || partyHasIdentities(party, supplierValidIdentities)
}

func hasCustomerTaxIDOrIdentity(value any) bool {
	party, _ := value.(*org.Party)
	return hasTaxIDCode(party) || partyHasIdentities(party, customerValidIdentities)
}

func hasTaxIDCode(party *org.Party) bool {
	return party != nil && party.TaxID != nil && party.TaxID.Code != ""
}

func partyHasIdentities(party *org.Party, identitites []cbc.Code) bool {
	if party == nil || len(party.Identities) == 0 {
		return false
	}
	return org.IdentitiesTypeIn(identitites...).Check(party.Identities)
}
