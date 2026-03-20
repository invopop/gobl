package pt

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Invoice type tags
const (
	TagInvoiceReceipt cbc.Key = "invoice-receipt"
)

var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*cbc.Definition{
		{
			Key: TagInvoiceReceipt,
			Name: i18n.String{
				i18n.EN: "Invoice-receipt",
				i18n.PT: "Fatura-recibo",
			},
		},
	},
}

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.HasContext(tax.RegimeIn(CountryCode)),
			rules.Field("type",
				rules.Assert("01", "invoice type is not valid for Portugal", is.In(
					bill.InvoiceTypeStandard,
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
					bill.InvoiceTypeProforma,
					bill.InvoiceTypeOther,
				)),
			),
			rules.When(
				bill.InvoiceTypeIn(bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote),
				rules.Field("preceding",
					rules.Assert("02", "preceding is required for credit and debit notes", is.Present),
				),
			),
			rules.Field("supplier",
				rules.Field("tax_id",
					rules.Assert("03", "supplier tax ID is required", is.Present),
					rules.Field("code",
						rules.Assert("04", "supplier tax ID code is required", is.Present),
					),
				),
			),
			rules.Field("lines",
				rules.Each(
					rules.Field("quantity",
						rules.Assert("05", "line quantity must be zero or positive", num.ZeroOrPositive),
					),
					rules.Field("item",
						rules.Field("price",
							rules.Assert("06", "item price must be zero or positive", num.ZeroOrPositive),
						),
					),
				),
			),
			rules.Field("totals",
				rules.Field("due",
					rules.Assert("07", "due amount must be zero or positive", num.ZeroOrPositive),
				),
			),
			rules.Assert("11", "value date must not be after issue date",
				is.Func("value date not after issue date", valueDateNotAfterIssueDate),
			),
			rules.Assert("12", "operation date must not be after issue date",
				is.Func("operation date not after issue date", operationDateNotAfterIssueDate),
			),
			rules.Assert("13", "preceding document issue dates must not be after invoice issue date",
				is.Func("preceding dates not after issue date", precedingDatesNotAfterIssueDate),
			),
			rules.Assert("14", "advance payment dates must not be after invoice issue date",
				is.Func("advance dates not after issue date", advanceDatesNotAfterIssueDate),
			),
			rules.Assert("15", "payment due dates must not be before invoice issue date",
				is.Func("due dates not before issue date", dueDatesNotBeforeIssueDate),
			),
		),
	)
}

func valueDateNotAfterIssueDate(val any) bool {
	inv := val.(*bill.Invoice)
	if inv == nil || inv.ValueDate == nil {
		return true
	}
	return cal.DateBefore(inv.IssueDate).Check(inv.ValueDate)
}

func operationDateNotAfterIssueDate(val any) bool {
	inv := val.(*bill.Invoice)
	if inv == nil || inv.OperationDate == nil {
		return true
	}
	return cal.DateBefore(inv.IssueDate).Check(inv.OperationDate)
}

func precedingDatesNotAfterIssueDate(val any) bool {
	inv := val.(*bill.Invoice)
	if inv == nil {
		return true
	}
	for _, ref := range inv.Preceding {
		if ref == nil || ref.IssueDate == nil {
			continue
		}
		if !cal.DateBefore(inv.IssueDate).Check(ref.IssueDate) {
			return false
		}
	}
	return true
}

func advanceDatesNotAfterIssueDate(val any) bool {
	inv := val.(*bill.Invoice)
	if inv == nil || inv.Payment == nil {
		return true
	}
	for _, adv := range inv.Payment.Advances {
		if adv == nil || adv.Date == nil {
			continue
		}
		if !cal.DateBefore(inv.IssueDate).Check(adv.Date) {
			return false
		}
	}
	return true
}

func dueDatesNotBeforeIssueDate(val any) bool {
	inv := val.(*bill.Invoice)
	if inv == nil || inv.Payment == nil || inv.Payment.Terms == nil {
		return true
	}
	for _, dd := range inv.Payment.Terms.DueDates {
		if dd == nil || dd.Date == nil {
			continue
		}
		if !cal.DateAfter(inv.IssueDate).Check(dd.Date) {
			return false
		}
	}
	return true
}
