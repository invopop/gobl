package dk

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
			is.InContext(tax.RegimeIn(l10n.DK.Tax())),
			rules.Field("supplier",
				rules.Assert("01", fmt.Sprintf(
					"invoice DK supplier must have either tax ID code or identity with '%s' or '%s' type",
					IdentityTypeCVR, IdentityTypeCPR,
				),
					is.Func(
						fmt.Sprintf("has tax ID code or identity with '%s' or '%s' type", IdentityTypeCVR, IdentityTypeCPR),
						hasTaxIDOrIdentity,
					),
				),
			),
		),
	)
}

func hasTaxIDOrIdentity(value any) bool {
	party, _ := value.(*org.Party)
	if hasTaxIDCode(party) {
		return true
	}
	if hasIdentityCVR(party) {
		return true
	}
	return hasIdentityCPR(party)
}

func hasTaxIDCode(party *org.Party) bool {
	return party != nil && party.TaxID != nil && party.TaxID.Code != ""
}

func hasIdentityCVR(party *org.Party) bool {
	if party == nil || len(party.Identities) == 0 {
		return false
	}
	return org.IdentityForType(party.Identities, IdentityTypeCVR) != nil
}

func hasIdentityCPR(party *org.Party) bool {
	if party == nil || len(party.Identities) == 0 {
		return false
	}
	return org.IdentityForType(party.Identities, IdentityTypeCPR) != nil
}
