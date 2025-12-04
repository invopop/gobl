package tbai

import (
	"strings"

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

func normalizeInvoice(inv *bill.Invoice) {
	if inv == nil {
		return
	}
	normalizeInvoiceTax(inv)
}

func normalizeInvoiceTax(inv *bill.Invoice) {
	tx := inv.Tax
	if tx == nil {
		tx = &bill.Tax{}
	}
	if tx.Ext == nil {
		tx.Ext = make(tax.Extensions)
	}
	if tx.Ext.Has(ExtKeyRegion) {
		return
	}
	if inv.Supplier == nil || len(inv.Supplier.Addresses) == 0 {
		return
	}
	addr := inv.Supplier.Addresses[0]
	// Take a set of different names for the same region and attempt
	// to use them to set the region code automatically.
	switch strings.ToLower(addr.Region) {
	case "alava", "álava", "araba", "vi":
		tx.Ext[ExtKeyRegion] = "VI"
	case "bizkaia", "vizcaya", "bi":
		tx.Ext[ExtKeyRegion] = "BI"
	case "gipuzkoa", "guipuzcoa", "guipúzcoa", "ss":
		tx.Ext[ExtKeyRegion] = "SS"
	default:
		return
	}
	if len(tx.Ext) > 0 {
		inv.Tax = tx
	}
}

func normalizeBillLine(line *bill.Line) {
	if line == nil || line.Item == nil {
		return
	}
	vt := line.Taxes.Get(tax.CategoryVAT)
	if vt == nil {
		return
	}
	switch line.Item.Key {
	case org.ItemKeyGoods:
		vt.Ext = vt.Ext.SetOneOf(ExtKeyProduct, "goods", "resale")
	case org.ItemKeyServices, cbc.KeyEmpty:
		vt.Ext = vt.Ext.Set(ExtKeyProduct, "services")
	}
}

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Tax,
			validation.Required,
			validation.By(validateInvoiceTax),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.By(validateInvoiceCustomer),
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
		validation.Field(&inv.Lines,
			validation.Each(
				validation.By(validateInvoiceLine),
				validation.Skip,
			),
			validation.Skip,
		),
		validation.Field(&inv.Notes,
			org.ValidateNotesHasKey(org.NoteKeyGeneral),
			validation.Skip,
		),
	)
}

func validateInvoiceTax(val any) error {
	obj, ok := val.(*bill.Tax)
	if obj == nil || !ok {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Ext,
			tax.ExtensionsRequire(ExtKeyRegion),
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

func validateInvoicePreceding(val any) error {
	p, ok := val.(*org.DocumentRef)
	if !ok {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.IssueDate, validation.Required),
		validation.Field(&p.Ext,
			tax.ExtensionsRequire(ExtKeyCorrection),
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
				obj.Key == tax.KeyExempt,
				tax.ExtensionsRequire(ExtKeyExempt),
			),
			validation.Skip,
		),
	)
}
