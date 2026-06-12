package lu

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// billInvoiceRules requires the supplier to have a TVA tax ID code or an RCS
// identity; the latter covers businesses below the VAT registration threshold.
func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.InContext(tax.RegimeIn(l10n.LU.Tax())),
			rules.Field("supplier",
				rules.Assert("01",
					"invoice LU supplier must have a TVA tax ID code or an RCS identity",
					is.AnyOf(
						is.Func("has TVA tax ID code", hasTaxIDCode),
						is.Func("has RCS identity", hasRCSIdentity),
					),
				),
			),
		),
	)
}

func hasTaxIDCode(value any) bool {
	party, _ := value.(*org.Party)
	return party != nil && party.TaxID != nil && party.TaxID.Code != ""
}

func hasRCSIdentity(value any) bool {
	party, _ := value.(*org.Party)
	if party == nil || len(party.Identities) == 0 {
		return false
	}
	return org.IdentityForType(party.Identities, IdentityTypeRCS) != nil
}
