package saft

import (
	"fmt"
	"slices"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func billOrderRules() *rules.Set {
	return rules.For(new(bill.Order),
		rules.Assert("01",
			fmt.Sprintf("tax requires '%s' extension", ExtKeyWorkType),
			is.Func("has work type", orderHasWorkType),
		),
		rules.Assert("02", "work type must not be an invoice work type",
			is.FuncError("order work type valid", orderWorkTypeValid),
		),
		rules.Assert("03", "series format must be valid",
			is.FuncError("series format", orderSeriesFormatValid),
		),
		rules.Assert("04", "code format must be valid",
			is.FuncError("code format", orderCodeFormatValid),
		),
		rules.Field("value_date",
			rules.Assert("05", "cannot be blank", is.Present),
		),
		rules.Field("lines",
			rules.Each(
				rules.Assert("06", "line taxes must include VAT category",
					is.FuncError("has VAT", orderLineHasVAT),
				),
			),
		),
	)
}

func orderHasWorkType(val any) bool {
	ord, ok := val.(*bill.Order)
	if !ok || ord == nil {
		return true
	}
	if ord.Tax == nil || ord.Tax.Ext == nil {
		return false
	}
	return tax.ExtensionsRequire(ExtKeyWorkType).Check(ord.Tax.Ext)
}

func orderWorkTypeValid(val any) error {
	ord, ok := val.(*bill.Order)
	if !ok || ord == nil {
		return nil
	}
	if ord.Tax == nil || ord.Tax.Ext == nil {
		return nil
	}
	if wt, ok := ord.Tax.Ext[ExtKeyWorkType]; ok {
		if slices.Contains(invoiceWorkTypes, wt) {
			return fmt.Errorf("value '%s' invalid", wt)
		}
	}
	return nil
}

func orderSeriesFormatValid(val any) error {
	ord, ok := val.(*bill.Order)
	if !ok || ord == nil {
		return nil
	}
	return validateSeriesFormat(orderDocType(ord)).Validate(ord.Series)
}

func orderCodeFormatValid(val any) error {
	ord, ok := val.(*bill.Order)
	if !ok || ord == nil {
		return nil
	}
	dt := orderDocType(ord)
	return validateCodeFormat(ord.Series, dt).Validate(ord.Code)
}

func orderLineHasVAT(val any) error {
	return bill.RequireLineTaxCategory(tax.CategoryVAT).Validate(val)
}

func orderDocType(ord *bill.Order) cbc.Code {
	if ord.Tax == nil || ord.Tax.Ext == nil {
		return cbc.CodeEmpty
	}
	return ord.Tax.Ext[ExtKeyWorkType]
}

func normalizeOrder(ord *bill.Order) {
	if ord == nil {
		return
	}

	normalizeOrderTax(ord)
	normalizeOrderValueDate(ord)
}

func normalizeOrderTax(ord *bill.Order) {
	if ord.Tax == nil {
		ord.Tax = new(bill.Tax)
	}

	if ord.Tax.Ext == nil {
		ord.Tax.Ext = make(tax.Extensions)
	}

	if !ord.Tax.Ext.Has(ExtKeyWorkType) {
		// Map order types to work types
		switch ord.Type {
		case bill.OrderTypePurchase:
			ord.Tax.Ext[ExtKeyWorkType] = WorkTypePurchaseOrder
		case bill.OrderTypeQuote:
			ord.Tax.Ext[ExtKeyWorkType] = WorkTypeBudgets
		}
	}
}

func normalizeOrderValueDate(ord *bill.Order) {
	ord.ValueDate = determineValueDate(
		dateOrToday(&ord.IssueDate, ord.Regime),
		ord.OperationDate,
		ord.ValueDate,
	)
}
