package choruspro

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// normalizeInvoice applies Chorus Pro specific normalization rules
func normalizeInvoice(inv *bill.Invoice) {
	if inv == nil {
		return
	}

	// Ensure required extensions are set with default values if not present
	if inv.Tax == nil {
		inv.Tax = &bill.Tax{}
	}
	if inv.Tax.Ext.IsZero() {
		inv.Tax.Ext = tax.MakeExtensions()
	}

	// Set default framework type if not specified. This breaks away from the
	// typical deterministic behavior of assigning extensions in GOBL, due to
	// complexity of trying to apply scenarios.
	if !inv.Tax.Ext.Has(ExtKeyFramework) {
		inv.Tax.Ext = inv.Tax.Ext.Merge(
			tax.ExtensionsOf(tax.ExtMap{
				ExtKeyFramework: ExtFrameworkCodeSupplier,
			}),
		)
	}

}

// normalizeBillLine is necessary as Chorus Pro requires quantity to be rounded to 4 decimals
func normalizeBillLine(line *bill.Line) {
	if line == nil {
		return
	}
	line.Quantity = line.Quantity.RescaleDown(4)
}

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Assert("07", "invoice must be in EUR or provide exchange rate for conversion", currency.CanConvertTo(currency.EUR)),
		// Customer validation (only when customer exists)
		rules.Field("customer",
			rules.Field("ext",
				// Always be set to '1' as this can only be used for B2G operations.
				rules.Assert("01", "customer scheme extension must be '1'",
					tax.ExtensionsHasCodes(ExtKeyScheme, "1"),
				),
			),
			rules.Field("identities",
				// Further assertions are made in the org.Party rules
				rules.Assert("02", "customer identities are required", is.Present),
			),
		),
		// Tax validation
		rules.Field("tax",
			rules.Assert("03", "tax object is required with extensions", is.Present),
			rules.Field("ext",
				rules.Assert("04", "tax extensions are required", is.Present),
				rules.Assert("05", "framework extension is required",
					tax.ExtensionsRequire(ExtKeyFramework),
				),
			),
		),
		// Totals validation for paid framework
		rules.When(
			is.Func("framework is paid", invoiceFrameworkIsPaid),
			rules.Field("totals",
				rules.Assert("06", "must be paid in full for framework 'A2'",
					is.Func("paid", invoiceTotalsPaid),
				),
			),
		),
	)
}

func invoiceFrameworkIsPaid(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && inv.Tax != nil && inv.Tax.Ext.Get(ExtKeyFramework) == ExtFrameworkCodePaid
}

func invoiceTotalsPaid(val any) bool {
	totals, ok := val.(*bill.Totals)
	if !ok || totals == nil {
		return false
	}
	return totals.Paid()
}
