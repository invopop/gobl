package dk

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			tax.RegimeIn(l10n.DK.Tax()),
			rules.Field("supplier",
				rules.Assert("01", fmt.Sprintf("invoice DK supplier must have either tax ID code or identity with '%s' type", IdentityTypeCVR),
					rules.By(
						fmt.Sprintf("has tax ID code or identity with '%s' type", IdentityTypeCVR),
						hasTaxIDOrIdentity,
					),
				),
			),
		),
	)
}

func hasTaxIDOrIdentity(value any) bool {
	party, _ := value.(*org.Party)
	return hasTaxIDCode(party) || hasIdentityCVR(party)
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
