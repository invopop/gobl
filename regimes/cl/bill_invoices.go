package cl

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
			is.InContext(tax.RegimeIn(CountryCode)),
			// Every Chilean invoice is issued under the supplier's RUT.
			// PartyHasTaxIDCode also rejects a tax identity present with an
			// empty code, which a nil-check alone would let through.
			rules.Field("supplier",
				rules.Assert("01", "invoice supplier tax ID code required for Chilean regime",
					org.PartyHasTaxIDCode(),
				),
				// Calculate normally derives the regime from this country; the
				// explicit check also protects documents that declare CL directly.
				rules.Field("tax_id",
					rules.Field("country",
						rules.Assert("02", "invoice supplier tax ID country must be CL for Chilean regime",
							is.In(CountryCode),
						),
					),
				),
			),
			rules.When(
				is.Func("not simplified", invoiceNotSimplified),
				// Chilean "boletas" (simplified invoices, DTE type 39) do not require the
				// customer's real RUT — the SII's own technical spec allows a generic
				// "consumer final" RUT (66.666.666-6) or an internal code instead. The
				// customer's RUT is only required on standard ("factura") invoices, which
				// is why this rule is scoped to non-simplified documents.
				//
				// Source: https://www.sii.cl/factura_electronica/factura_mercado/formato_boleta_electronica.pdf
				rules.Field("customer",
					rules.Assert("03", "invoice customer tax ID code required for Chilean regime",
						org.PartyHasTaxIDCode(),
					),
				),
			),
		),
	)
}

func invoiceNotSimplified(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	return !tax.TagSimplified.In(inv.GetTags()...)
}
