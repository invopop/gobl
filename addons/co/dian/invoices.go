package dian

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var invoiceCorrectionDefinitions = []*tax.CorrectionDefinition{
	{
		Schema: bill.ShortSchemaInvoice,
		Types: []cbc.Key{
			bill.InvoiceTypeCreditNote,
			bill.InvoiceTypeDebitNote,
		},
		Extensions: []cbc.Key{
			ExtKeyCreditCode,
			ExtKeyDebitCode,
		},
		ReasonRequired: true,
		Stamps: []cbc.Key{
			StampCUDE,
		},
	},
}

func normalizeInvoice(inv *bill.Invoice) {
	if inv == nil {
		return
	}
	normalizeTaxResponsibility(inv.Supplier)
	normalizeTaxResponsibility(inv.Customer)
}

func normalizeTaxResponsibility(p *org.Party) {
	if p == nil || !isColombian(p.TaxID) {
		return
	}
	def := tax.Extensions{ExtKeyTaxResponsibility: "R-99-PN"}
	p.Ext = def.Merge(p.Ext)
}

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Type,
			validation.In(
				bill.InvoiceTypeStandard,
				bill.InvoiceTypeCreditNote,
				bill.InvoiceTypeDebitNote,
				bill.InvoiceTypeProforma,
			),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.By(validateInvoiceSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.By(validateInvoiceCustomer(inv.GetTags())),
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
			validation.Each(validation.By(validateInvoicePreceding(inv.Type))),
			validation.Skip,
		),
	)
}

func validateInvoiceSupplier(value interface{}) error {
	obj, _ := value.(*org.Party)
	if obj == nil || obj.TaxID == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			tax.RequireIdentityCode,
			validation.Skip,
		),
		validation.Field(&obj.Addresses,
			validation.When(
				isColombian(obj.TaxID),
				validation.Length(1, 0),
			),
			validation.Skip,
		),
		validation.Field(&obj.Ext,
			validation.When(
				municipalityCodeRequired(obj.TaxID),
				validation.Required,
				tax.ExtensionsRequire(ExtKeyMunicipality),
			),
			validation.When(
				isColombian(obj.TaxID),
				tax.ExtensionsRequire(ExtKeyTaxResponsibility),
			),
			validation.Skip,
		),
	)
}

func validateInvoiceCustomer(tags []cbc.Key) func(value any) error {
	return func(value any) error {
		obj, _ := value.(*org.Party)
		if obj == nil {
			return nil
		}
		return validation.ValidateStruct(obj,
			validation.Field(&obj.TaxID,
				validation.When(
					!tax.TagSimplified.In(tags...),
					validation.Required,
					tax.RequireIdentityCode,
				),
				validation.Skip,
			),
			validation.Field(&obj.Identities,
				validation.When(
					len(obj.Identities) > 0,
					org.RequireIdentityKey(identityKeys...),
				),
				validation.Skip,
			),
			validation.Field(&obj.Addresses,
				validation.When(
					isColombian(obj.TaxID),
					validation.Length(1, 0),
				),
				validation.Skip,
			),
			validation.Field(&obj.Ext,
				validation.When(
					municipalityCodeRequired(obj.TaxID),
					validation.Required,
					tax.ExtensionsRequire(ExtKeyMunicipality),
				),
				validation.When(
					isColombian(obj.TaxID),
					tax.ExtensionsRequire(ExtKeyTaxResponsibility),
				),
				validation.Skip,
			),
		)
	}
}

func isColombian(tID *tax.Identity) bool {
	return tID != nil && tID.Country.In("CO")
}

// municipalityCodeRequired checks if the municipality code is required for the given tax
// identity by checking to see if the customer is a Colombian company.
func municipalityCodeRequired(tID *tax.Identity) bool {
	if tID == nil {
		return false
	}
	if !tID.Country.In("CO") {
		return false
	}
	return tID.Code != ""
}

func validateInvoicePreceding(typ cbc.Key) validation.RuleFunc {
	return func(value any) error {
		obj, ok := value.(*org.DocumentRef)
		if !ok || obj == nil {
			return nil
		}
		return validation.ValidateStruct(obj,
			validation.Field(&obj.Ext,
				validation.When(
					typ == bill.InvoiceTypeCreditNote,
					tax.ExtensionsRequire(ExtKeyCreditCode),
				),
				validation.When(
					typ == bill.InvoiceTypeDebitNote,
					tax.ExtensionsRequire(ExtKeyDebitCode),
				),
			),
			validation.Field(&obj.Reason, validation.Required),
		)
	}
}
