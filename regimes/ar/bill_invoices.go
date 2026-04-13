package ar

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.InContext(tax.RegimeIn(l10n.AR.Tax())),
			rules.Field("supplier",
				rules.Field("tax_id",
					rules.Assert("01", "invoice supplier tax ID required for Argentine regime", is.Present),
				),
			),
		),
	)
}
