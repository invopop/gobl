package pt

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
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

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Type,
			validation.In(
				bill.InvoiceTypeStandard,
				bill.InvoiceTypeCreditNote,
				bill.InvoiceTypeDebitNote,
				bill.InvoiceTypeProforma,
				bill.InvoiceTypeOther,
			),
			validation.Skip,
		),
		validation.Field(&inv.Preceding,
			validation.When(
				inv.Type.In(
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				),
				validation.Required,
			),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.By(validateSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(
				validation.By(validateInvoiceLine),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}

func validateSupplier(val any) error {
	obj, _ := val.(*org.Party)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
	)
}

func validateInvoiceLine(val any) error {
	line, _ := val.(*bill.Line)
	if line == nil {
		return nil
	}

	return validation.ValidateStruct(line,
		validation.Field(&line.Quantity,
			num.Min(num.MakeAmount(0, 0)),
		),
		validation.Field(&line.Item,
			validation.By(validateItem),
			validation.Skip,
		),
	)
}

func validateItem(val any) error {
	item, _ := val.(*org.Item)
	if item == nil {
		return nil
	}
	return validation.ValidateStruct(item,
		validation.Field(&item.Price,
			num.Min(num.MakeAmount(0, 0)),
		),
	)
}
