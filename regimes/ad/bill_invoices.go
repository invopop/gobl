package ad

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// billInvoiceRules defines Andorran invoice-level requirements.
// An Andorran supplier must identify themselves with either:
//   - a tax ID code (NRT as IGI registration, for businesses above the
//     €40,000 annual threshold), or
//   - an org identity of type NRT (for entities below the threshold who
//     hold an NRT but are not required to register for IGI).
func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.Field("supplier",
				rules.Assert("01",
					"invoice supplier in Andorra must have an NRT tax ID code or an NRT identity",
					is.Func("has NRT", hasSupplierNRT),
				),
			),
		),
	)
}

// hasSupplierNRT returns true if the supplier has identified themselves
// with either a tax ID code or an explicit NRT org identity.
func hasSupplierNRT(value any) bool {
	party, _ := value.(*org.Party)
	return hasTaxIDCode(party) || hasSupplierNRTIdentity(party)
}

// hasTaxIDCode checks whether the supplier has a tax identity with a
// non-empty code — i.e. they are registered for IGI.
func hasTaxIDCode(party *org.Party) bool {
	return party != nil && party.TaxID != nil && party.TaxID.Code != ""
}

// hasSupplierNRTIdentity checks whether the supplier has declared an
// org identity of type NRT — used for entities below the IGI threshold.
func hasSupplierNRTIdentity(party *org.Party) bool {
	if party == nil || len(party.Identities) == 0 {
		return false
	}
	for _, id := range party.Identities {
		if id.Type == IdentityTypeNRT {
			return true
		}
	}
	return false
}