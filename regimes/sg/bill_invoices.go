package sg

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.HasContext(tax.RegimeIn(CountryCode)),
			rules.Field("supplier",
				rules.Assert("01", "invoice supplier in Singapore must have a GST tax ID code or a UEN identity",
					is.Func("has GST tax ID code or UEN identity", hasSupplierTaxIDOrIdentity),
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
		case IdentityTypeUEN:
			return true
		}
	}
	return false
}
