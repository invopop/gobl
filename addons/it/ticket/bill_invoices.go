package ticket

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

var invoiceCorrectionDefinitions = tax.CorrectionSet{
	{
		Schema: bill.ShortSchemaInvoice,
		Types:  []cbc.Key{bill.InvoiceTypeCorrective},
		Stamps: []cbc.Key{
			StampRef,
		},
	},
}

func normalizeInvoice(inv *bill.Invoice) {
	if inv.Tax == nil {
		inv.Tax = new(bill.Tax)
	}
	if inv.Tax.PricesInclude == "" {
		inv.Tax.PricesInclude = tax.CategoryVAT
	}
	if !inv.Tax.Ext.IsZero() && inv.Tax.Ext.Has(ExtKeyLottery) {
		inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyLottery, cbc.NormalizeAlphanumericalCode(inv.Tax.Ext.Get(ExtKeyLottery)))
	}
}

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Assert("08", "invoice must be in EUR or provide exchange rate for conversion", currency.CanConvertTo(currency.EUR)),
		rules.Field("tax",
			rules.Assert("01", "tax is required", is.Present),
			rules.Field("prices_include",
				rules.Assert("02", "prices_include is required", is.Present),
				rules.Assert("03", "prices_include must be VAT",
					is.In(tax.CategoryVAT),
				),
			),
		),
		rules.Field("supplier",
			rules.Field("tax_id",
				rules.Assert("04", "supplier tax ID is required", is.Present),
			),
			// A supplier address names the store that issued the receipt, so
			// that merchants with several locations under one tax ID can tell
			// receipts apart. The AdE prints it instead of the registered
			// address and needs every part of it.
			rules.Field("addresses",
				rules.Each(
					rules.Field("street",
						rules.Assert("09", "supplier address street is required", is.Present),
					),
					rules.Field("locality",
						rules.Assert("10", "supplier address locality is required", is.Present),
					),
					rules.Field("code",
						rules.Assert("11", "supplier address code is required", is.Present),
						rules.Assert("12", "supplier address code must be 5 digits", is.Matches(`^\d{5}$`)),
					),
					rules.Field("region",
						rules.Assert("13", "supplier address region is required", is.Present),
						rules.Assert("14", "supplier address region must be a 2-letter province code", is.Matches(`^[A-Z]{2}$`)),
					),
				),
			),
		),
		rules.When(bill.InvoiceTypeIn(bill.InvoiceTypeCorrective),
			rules.Field("preceding",
				rules.Assert("05", "preceding documents are required for corrective invoices", is.Present),
			),
		),
		rules.Field("lines",
			rules.Each(
				rules.Assert("06", "line taxes must include VAT category",
					bill.RequireLineTaxCategory(tax.CategoryVAT),
				),
			),
		),
		rules.When(bill.InvoiceTypeIn(bill.InvoiceTypeCorrective),
			rules.Field("lines",
				rules.Each(
					rules.Field("ext",
						rules.Assert("07",
							fmt.Sprintf("corrective invoice lines require '%s' extension", ExtKeyLine),
							tax.ExtensionsRequire(ExtKeyLine),
						),
					),
				),
			),
		),
	)
}
