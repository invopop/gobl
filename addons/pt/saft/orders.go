package saft

import (
	"fmt"
	"slices"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateOrder(ord *bill.Order) error {
	dt := orderDocType(ord)

	return validation.ValidateStruct(ord,
		validation.Field(&ord.Tax,
			validation.By(validateOrderTax),
			validation.Skip,
		),
		validation.Field(&ord.Series,
			validateSeriesFormat(dt),
			validation.Skip,
		),
		validation.Field(&ord.Code,
			validateCodeFormat(ord.Series, dt),
			validation.Skip,
		),
		validation.Field(&ord.ValueDate,
			validation.Required,
			validation.Skip,
		),
		validation.Field(&ord.Lines,
			validation.Each(
				bill.RequireLineTaxCategory(tax.CategoryVAT),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}

func orderDocType(ord *bill.Order) cbc.Code {
	if ord.Tax == nil || ord.Tax.Ext == nil {
		return cbc.CodeEmpty
	}
	return ord.Tax.Ext[ExtKeyWorkType]
}

func validateOrderTax(val any) error {
	t, _ := val.(*bill.Tax)
	if t == nil {
		t = new(bill.Tax)
	}

	return validation.ValidateStruct(t,
		validation.Field(&t.Ext,
			validation.By(validateOrderTaxExt),
			validation.Skip,
		),
	)
}

func validateOrderTaxExt(val any) error {
	ext, _ := val.(tax.Extensions)
	if ext == nil {
		ext = make(tax.Extensions)
	}

	if wt, ok := ext[ExtKeyWorkType]; ok {
		if slices.Contains(invoiceWorkTypes, wt) {
			return validation.Errors{
				ExtKeyWorkType.String(): fmt.Errorf("value '%s' invalid", wt),
			}
		}
	}

	return validation.Validate(val,
		tax.ExtensionsRequire(ExtKeyWorkType),
		validation.Skip,
	)
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
