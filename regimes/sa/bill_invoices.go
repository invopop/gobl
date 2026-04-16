package sa

import (
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
			is.InContext(tax.RegimeIn(l10n.SA.Tax())),
			rules.Field("supplier",
				rules.Assert("01", "supplier must have a valid tax ID code (BR-KSA-39)",
					is.Func("valid VAT code", hasTaxIDCode),
				),
				rules.Field("identities",
					rules.Assert("02", "supplier can have 0 or 1 identities (BR-KSA-08)",
						is.Func("identity must be one of: CRN/MOM/MLS/700/SAG/OTH", hasAtMostOneIdentity),
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

func hasAtMostOneIdentity(value any) bool {
	identities, _ := value.([]*org.Identity)
	return len(identities) == 0 || (len(identities) == 1 && org.IdentitiesTypeIn(supplierValidIdentities...).Check(identities))
}
