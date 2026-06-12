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

// billInvoiceRules requires the supplier to carry at least one official
// identifier. Note that the RCS number is a company-registry reference, not
// a tax identifier: a VAT-registered supplier is still legally required to
// state its TVA number on invoices. Accepting an RCS identity as an
// alternative is a deliberate baseline so that businesses below the VAT
// registration threshold (e.g. small franchised businesses) can still issue
// valid GOBL invoices; it is not a claim of legal equivalence.
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
