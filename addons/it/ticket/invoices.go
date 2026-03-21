package ticket

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
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
	if inv.Tax.Ext != nil && inv.Tax.Ext.Has(ExtKeyLottery) {
		inv.Tax.Ext[ExtKeyLottery] = cbc.NormalizeAlphanumericalCode(inv.Tax.Ext[ExtKeyLottery])
	}
}

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
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
		),
		rules.When(bill.InvoiceTypeIn(bill.InvoiceTypeCorrective),
			rules.Field("preceding",
				rules.Assert("05", "preceding documents are required for corrective invoices", is.Present),
			),
		),
		rules.Field("lines",
			rules.Each(
				rules.Assert("06", "line taxes must include VAT category",
					is.FuncError("has VAT", lineHasVATCategory),
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

func lineHasVATCategory(val any) error {
	return bill.RequireLineTaxCategory(tax.CategoryVAT).Validate(val)
}
