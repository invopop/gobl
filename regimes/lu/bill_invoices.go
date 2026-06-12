package lu

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
			is.InContext(tax.RegimeIn(l10n.LU.Tax())),
			rules.Field("supplier",
				rules.Assert("01",
					fmt.Sprintf("invoice LU supplier must have either a TVA tax ID code or an identity of type '%s'", IdentityTypeRCS),
					is.Func("has TVA code or RCS identity", hasVATCodeOrRCSIdentity),
				),
			),
		),
	)
}

func hasVATCodeOrRCSIdentity(value any) bool {
	party, _ := value.(*org.Party)
	return hasTaxIDCode(party) || hasRCSIdentity(party)
}

func hasTaxIDCode(party *org.Party) bool {
	return party != nil && party.TaxID != nil && party.TaxID.Code != ""
}

func hasRCSIdentity(party *org.Party) bool {
	if party == nil || len(party.Identities) == 0 {
		return false
	}
	return org.IdentityForType(party.Identities, IdentityTypeRCS) != nil
}
