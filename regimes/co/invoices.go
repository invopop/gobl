package co

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

type invoiceValidator struct {
	inv *bill.Invoice
}

func validateInvoice(inv *bill.Invoice) error {
	v := &invoiceValidator{inv: inv}
	return v.validate()
}

func (v *invoiceValidator) validate() error {
	inv := v.inv
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
			validation.By(v.validSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.By(v.validCustomer),
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
			validation.Each(validation.By(v.preceding)),
			validation.Skip,
		),
		validation.Field(&inv.Outlays,
			validation.Empty,
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) validSupplier(value interface{}) error {
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
				tax.ExtensionsRequires(ExtKeyDIANMunicipality),
			),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) validCustomer(value interface{}) error {
	obj, _ := value.(*org.Party)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.When(
				!v.inv.Tax.ContainsTag(tax.TagSimplified),
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
				tax.ExtensionsRequires(ExtKeyDIANMunicipality),
			),
			validation.Skip,
		),
	)
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

func (v *invoiceValidator) preceding(value interface{}) error {
	obj, ok := value.(*bill.Preceding)
	if !ok || obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Ext,
			validation.When(
				v.inv.Type == bill.InvoiceTypeCreditNote,
				tax.ExtensionsRequires(ExtKeyDIANCreditCode),
			),
			validation.When(
				v.inv.Type == bill.InvoiceTypeDebitNote,
				tax.ExtensionsRequires(ExtKeyDIANDebitCode),
			),
		),
		validation.Field(&obj.Reason, validation.Required),
	)
}

func normalizeParty(p *org.Party) {
	if p == nil {
		return
	}
	// 2024-03-14: Migrate Tax ID Zone to extensions "co-dian-municipality"
	if p.TaxID != nil && p.TaxID.Zone != "" { //nolint:staticcheck
		if p.Ext == nil {
			p.Ext = make(tax.Extensions)
		}
		p.Ext[ExtKeyDIANMunicipality] = tax.ExtValue(p.TaxID.Zone) //nolint:staticcheck
		p.TaxID.Zone = ""                                          //nolint:staticcheck
	}
}
