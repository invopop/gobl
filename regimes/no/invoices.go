package no

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			// On standard invoices the supplier address and a customer are
			// required. The supplier tax ID is intentionally not required:
			// Norwegian businesses below the NOK 50,000 turnover threshold are
			// not VAT-registered and cannot charge MVA, but may still issue
			// invoices. Simplified invoices relax both of these.
			rules.When(
				is.Func("not simplified", isNotSimplified),
				rules.Field("supplier",
					rules.Field("addresses",
						rules.Assert("01", "supplier address is required on standard invoices", is.Present),
					),
				),
				rules.Field("customer",
					rules.Assert("02", "customer is required on standard invoices", is.Present),
				),
			),
			// Credit and debit notes must reference the document they correct.
			rules.When(
				bill.InvoiceTypeIn(
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				),
				rules.Field("preceding",
					rules.Assert("03", "preceding document is required for credit and debit notes", is.Present),
				),
			),
		),
	)
}

func isNotSimplified(value any) bool {
	inv, ok := value.(*bill.Invoice)
	return ok && inv != nil && !inv.HasTags(tax.TagSimplified)
}
