package mydata

import (
	"fmt"
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/gr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Field("series",
			rules.Assert("01", "series is required", is.Present),
		),
		rules.Field("tax",
			rules.Assert("02", "tax is required", is.Present),
			rules.Field("ext",
				rules.Assert("03",
					fmt.Sprintf("tax requires '%s' extension", ExtKeyInvoiceType),
					tax.ExtensionsRequire(ExtKeyInvoiceType),
				),
			),
		),
		rules.Field("supplier",
			rules.Field("tax_id",
				rules.Assert("04", "supplier tax ID is required", is.Present),
				rules.Field("code",
					rules.Assert("18", "supplier tax ID code is required", is.Present),
				),
			),
		),
		rules.When(is.Func("requires customer", invoiceRequiresValidCustomer),
			rules.Field("customer",
				rules.Assert("05", "customer is required", is.Present),
				rules.Field("tax_id",
					rules.Assert("06", "customer tax ID is required", is.Present),
				),
				rules.Field("addresses",
					rules.Assert("07", "customer addresses are required", is.Present),
					rules.Each(
						rules.Field("locality",
							rules.Assert("08", "customer address locality is required", is.Present),
						),
						rules.Field("code",
							rules.Assert("09", "customer address code is required", is.Present),
						),
					),
				),
			),
		),
		rules.Field("lines",
			rules.Each(
				rules.Field("total",
					rules.Assert("10", "line total must be positive", num.Positive),
					rules.Assert("11", "line total must not be zero", num.NotZero),
				),
				rules.Field("item",
					rules.Field("ext",
						rules.Assert("12",
							fmt.Sprintf("item income extensions '%s' and '%s' must both be present",
								ExtKeyIncomeCat, ExtKeyIncomeType),
							is.Func("income ext pair", itemIncomeExtPairValid),
						),
					),
				),
			),
		),
		rules.Field("discounts",
			rules.Assert("13", "discounts are not supported by mydata", is.Empty),
		),
		rules.Field("payment",
			rules.Assert("14", "payment is required", is.Present),
			rules.Assert("15", "payment instructions are required when no advances",
				is.Func("has instructions or advances", paymentHasInstructionsOrAdvances),
			),
		),
		rules.When(bill.InvoiceTypeIn(bill.InvoiceTypeCreditNote),
			rules.Field("preceding",
				rules.Assert("16", "preceding documents are required for credit notes", is.Present),
			),
		),
		rules.Field("preceding",
			rules.Each(
				rules.Field("stamps",
					rules.Assert("17",
						fmt.Sprintf("preceding document requires '%s' stamp", gr.StampIAPRMark),
						head.StampsHas(gr.StampIAPRMark),
					),
				),
			),
		),
	)
}

func invoiceRequiresValidCustomer(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	typeCats := []string{"1", "2", "5"}
	it := inv.Tax.Ext[ExtKeyInvoiceType].String()
	for _, prefix := range typeCats {
		if strings.HasPrefix(it, prefix+".") {
			return true
		}
	}
	return false
}

func itemIncomeExtPairValid(val any) bool {
	ext, ok := val.(tax.Extensions)
	if !ok {
		return true
	}
	hasCat := ext.Has(ExtKeyIncomeCat)
	hasType := ext.Has(ExtKeyIncomeType)
	if !hasCat && !hasType {
		return true // neither set, valid
	}
	return hasCat && hasType // both must be set
}

func paymentHasInstructionsOrAdvances(val any) bool {
	p, ok := val.(*bill.PaymentDetails)
	if !ok || p == nil {
		return true
	}
	if len(p.Advances) > 0 {
		return true
	}
	return p.Instructions != nil
}
