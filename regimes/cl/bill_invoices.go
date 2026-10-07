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
			// Chilean Notas de Crédito/Débito Electrónicas only correct standard
			// "factura" documents, never simplified ones ("boletas") — per the
			// SII's official guides, under Art. 57 of the D.L. 825 and Art. 71 of
			// its Reglamento. GOBL's preceding document reference doesn't carry
			// the original document's tags forward, so this rule checks the next
			// best thing: the correction itself must not be tagged as simplified.
			//
			// Source: https://www.sii.cl/destacados/factura_electronica/guias_ayuda/nota_credito_corrige_monto_fe.htm
			rules.When(
				bill.InvoiceTypeIn(bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote),
				rules.Assert("04", "credit and debit notes cannot be simplified invoices for Chilean regime", is.Func("not simplified", invoiceNotSimplified)),
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
