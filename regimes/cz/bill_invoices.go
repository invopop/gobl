package cz

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// billInvoiceRules defines Czech invoice validation rules.
// Supplier must have either a DIČ (tax ID) or IČO (business registration).
func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.Field("supplier",
				rules.Assert("01", fmt.Sprintf("invoice CZ supplier must have either tax ID code or identity with '%s' key", IdentityKeyICO),
					is.Func(
						fmt.Sprintf("has tax ID code or identity with '%s' key", IdentityKeyICO),
						hasTaxIDOrIdentity,
					),
				),
			),
		),
	)
}

func hasTaxIDOrIdentity(value any) bool {
	party, _ := value.(*org.Party)
	return hasTaxIDCode(party) || hasIdentityICO(party)
}

func hasTaxIDCode(party *org.Party) bool {
	return party != nil && party.TaxID != nil && party.TaxID.Code != ""
}

func hasIdentityICO(party *org.Party) bool {
	if party == nil || len(party.Identities) == 0 {
		return false
	}
	return org.IdentityForKey(party.Identities, IdentityKeyICO) != nil
}
