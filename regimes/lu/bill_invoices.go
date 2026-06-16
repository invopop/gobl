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
				rules.When(
					is.Func("supplier has no TVA tax ID code", supplierHasNoTVACode),
					rules.Field("identities",
						rules.Assert("01",
							"invoice LU supplier without TVA tax ID code must have an RCS identity",
							org.IdentityTypeIn(IdentityTypeRCS),
						),
					),
				),
			),
		),
	)
}

func supplierHasNoTVACode(value any) bool {
	party, _ := value.(*org.Party)
	return party == nil || party.TaxID == nil || party.TaxID.Code == ""
}
