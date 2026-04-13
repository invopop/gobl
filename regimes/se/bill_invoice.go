package se

import (
	"fmt"

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
			is.InContext(tax.RegimeIn(l10n.SE.Tax())),
			rules.Field("supplier",
				rules.Assert("01", fmt.Sprintf("invoice SE supplier must have either tax ID code or identity with %s, %s, or %s type", IdentityTypeOrgNr, IdentityTypePersonNr, IdentityTypeCoordinationNr),
					is.Func(
						fmt.Sprintf("has tax ID code or identity with type in [%s, %s, %s]", IdentityTypeOrgNr, IdentityTypePersonNr, IdentityTypeCoordinationNr),
						hasSupplierTaxIDOrIdentity,
					),
				),
			),
		),
	)
}

func hasSupplierTaxIDOrIdentity(value any) bool {
	party, _ := value.(*org.Party)
	return hasTaxIDCode(party) || hasSupplierIdentity(party)
}

func hasTaxIDCode(party *org.Party) bool {
	return party != nil && party.TaxID != nil && party.TaxID.Code != ""
}

func hasSupplierIdentity(party *org.Party) bool {
	if party == nil || len(party.Identities) == 0 {
		return false
	}
	for _, id := range party.Identities {
		switch id.Type {
		case IdentityTypeOrgNr, IdentityTypePersonNr, IdentityTypeCoordinationNr:
			return true
		}
	}
	return false
}
