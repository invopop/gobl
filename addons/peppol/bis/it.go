package bis

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

func orgPartyRulesIT() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(supplierCountryIs(l10n.IT),
			rules.Field("supplier",
				// IT-R-001: Tax registration identifier length between 11 and 16.
				rules.Assert("IT-R-001", "Italian tax registration identifier length must be 11-16 (IT-R-001)",
					is.Func("it tax id length", italianTaxIDLength),
				),
				rules.Field("addresses",
					// IT-R-002: address line 1 (street) required.
					rules.Assert("IT-R-002", "Italian supplier address line 1 is required (IT-R-002)",
						is.Func("it street", firstAddressHasStreet),
					),
					// IT-R-003: city required.
					rules.Assert("IT-R-003", "Italian supplier city is required (IT-R-003)",
						is.Func("it locality", firstAddressHasLocalityPE),
					),
					// IT-R-004: post code required.
					rules.Assert("IT-R-004", "Italian supplier post code is required (IT-R-004)",
						is.Func("it code", firstAddressHasCodePE),
					),
				),
			),
		),
	)
}

func italianTaxIDLength(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil || p.TaxID == nil {
		return true
	}
	code := p.TaxID.Code.String()
	if code == "" {
		return true
	}
	l := len(code)
	return l >= 11 && l <= 16
}

func firstAddressHasStreet(val any) bool {
	addrs, ok := val.([]*org.Address)
	if !ok || len(addrs) == 0 {
		return true
	}
	return addrs[0] != nil && addrs[0].Street != ""
}
