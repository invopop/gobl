package nl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			rules.HasContext(tax.RegimeIn(l10n.NL.Tax())),
			rules.Field("supplier",
				rules.Assert("01", "invoice supplier must have a tax ID code or a KVK/OIN identity",
					rules.By("has tax ID code or KVK/OIN identity", hasSupplierTaxIDOrIdentity),
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
		if id.Type == IdentityTypeKVK || id.Type == IdentityTypeOIN {
			return true
		}
	}
	return false
}
