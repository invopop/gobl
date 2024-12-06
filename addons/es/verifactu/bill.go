package verifactu

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
			ExtKeyDocType,
		},
	},
}

func normalizeInvoice(inv *bill.Invoice) {
	// Try to move any preceding choices to the document level
	for _, row := range inv.Preceding {
		if len(row.Ext) == 0 {
			continue
		}
		found := false
		if row.Ext.Has(ExtKeyDocType) {
			if inv.Tax == nil || !found {
				inv.Tax = inv.Tax.MergeExtensions(tax.Extensions{
					ExtKeyDocType: row.Ext[ExtKeyDocType],
				})
				found = true // only assign first one
			}
			delete(row.Ext, ExtKeyDocType)
		}
	}

	// Try to normalize the correction type, which is especially complex for
	// Verifactu implying that scenarios cannot be used.
	switch inv.Type {
	case bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote:
		inv.Tax = inv.Tax.MergeExtensions(tax.Extensions{
			ExtKeyCorrectionType: "I",
		})
	case bill.InvoiceTypeCorrective:
		if inv.Tax == nil || inv.Tax.Ext.Get(ExtKeyDocType) != "F3" {
			inv.Tax = inv.Tax.MergeExtensions(tax.Extensions{
				ExtKeyCorrectionType: "S",
			})
		} else {
			// Substitutions of simplified invoices cannot have a correction type
			delete(inv.Tax.Ext, ExtKeyCorrectionType)
		}
	}
}

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
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
		validation.Field(&inv.Tax,
			validation.Required,
			validation.By(validateInvoiceTax(inv.Type)),
			validation.Skip,
		),
		validation.Field(&inv.Notes,
			cbc.ValidateNotesHasKey(cbc.NoteKeyGeneral),
			validation.Skip,
		),
	)
}

var docTypesStandard = []cbc.Code{ // Standard invoices
	"F1", "F2",
}
var docTypesCreditDebit = []cbc.Code{ // Credit or Debit notes
	"R1", "R2", "R3", "R4", "R5",
}
var docTypesCorrective = []cbc.Code{ // Substitutions
	"F3", "R1", "R2", "R3", "R4", "R5",
}

func validateInvoiceTax(it cbc.Key) validation.RuleFunc {
	return func(val any) error {
		obj := val.(*bill.Tax)
		return validation.ValidateStruct(obj,
			validation.Field(&obj.Ext,
				tax.ExtensionsRequire(ExtKeyDocType),
				validation.When(
					it.In(bill.InvoiceTypeStandard),
					tax.ExtensionsHasCodes(
						ExtKeyDocType,
						docTypesStandard...,
					),
				),
				validation.When(
					it.In(bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote),
					tax.ExtensionsHasCodes(
						ExtKeyDocType,
						docTypesCreditDebit...,
					),
				),
				validation.When(
					it.In(bill.InvoiceTypeCorrective),
					tax.ExtensionsHasCodes(
						ExtKeyDocType,
						docTypesCorrective...,
					),
				),
				validation.When(
					obj.Ext.Get(ExtKeyDocType).In(docTypesCreditDebit...),
					tax.ExtensionsRequire(ExtKeyCorrectionType),
				),
				validation.Skip,
			),
		)
	}
}

func validateInvoiceCustomer(val any) error {
	obj, ok := val.(*org.Party)
	if !ok || obj == nil {
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
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.IssueDate,
			validation.Required,
		),
	)
}
