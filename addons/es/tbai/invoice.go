package tbai

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var invoiceCorrectionDefinitions = tax.CorrectionSet{
	{
		Schema: bill.ShortSchemaInvoice,
		Extensions: []cbc.Key{
			ExtKeyCorrection,
		},
	},
}

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Preceding,
			validation.When(
				inv.Type.In(es.InvoiceCorrectionTypes...),
				validation.Required,
			),
			validation.By(validateInvoicePreceding),
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

func validateInvoicePreceding(val any) error {
	p, ok := val.(*org.DocumentRef)
	if !ok {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.IssueDate, validation.Required),
		validation.Field(&p.Ext,
			tax.ExtensionsRequires(ExtKeyCorrection),
			validation.Skip,
		),
	)
}

func validateInvoiceLine(value any) error {
	obj, _ := value.(*bill.Line)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Taxes,
			validation.Each(
				validation.By(validateInvoiceLineTax),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}

func validateInvoiceLineTax(value any) error {
	obj, ok := value.(*tax.Combo)
	if obj == nil || !ok {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Ext,
			validation.When(
				obj.Rate == tax.RateExempt,
				tax.ExtensionsRequires(ExtKeyExemption),
			),
			validation.Skip,
		),
	)
}
