package facturae

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

func normalizeInvoice(_ *bill.Invoice) {
	// todo
}

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Customer,
			validation.By(validateInvoiceCustomer),
			validation.Skip,
		),
		validation.Field(&inv.Tax,
			validation.By(validateInvoiceTax),
			validation.Skip,
		),
		validation.Field(&inv.Preceding,
			validation.When(
				inv.Type.In(es.InvoiceCorrectionTypes...),
				validation.Required,
			),
			validation.Each(
				validation.By(validateInvoicePreceding),
			),
			validation.Skip,
		),
	)
}

func validateInvoiceCustomer(val any) error {
	obj, _ := val.(*org.Party)
	if obj == nil {
		return nil
	}
	// Customers must have a tax ID to at least set the country,
	// and Spanish ones should also have an ID. There are more complex
	// rules for exports.
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.Required,
			validation.When(
				obj.TaxID != nil && obj.TaxID.Country.In("ES"),
				tax.RequireIdentityCode,
			),
			validation.Skip,
		),
	)
}

func validateInvoiceTax(val any) error {
	t, ok := val.(*bill.Tax)
	if !ok {
		return nil
	}
	return validation.ValidateStruct(t,
		validation.Field(&t.Ext,
			tax.ExtensionsRequires(
				ExtKeyDocType,
				ExtKeyInvoiceClass,
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
