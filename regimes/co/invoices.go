package co

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
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
				bill.InvoiceTypeProforma,
			),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.By(v.validParty),
			validation.By(v.validSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.By(v.validParty),
			validation.Skip,
		),
		validation.Field(&inv.Preceding,
			validation.When(
				inv.Type.In(bill.InvoiceTypeCreditNote),
				validation.Required,
			),
			validation.Each(validation.By(v.preceding)),
			validation.Skip,
		),
		validation.Field(&inv.Outlays,
			validation.Empty,
		),
	)
}

func (v *invoiceValidator) validParty(value interface{}) error {
	obj, _ := value.(*org.Party)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.Required,
			validation.By(v.validTaxIdentity),
			validation.Skip,
		),
		validation.Field(&obj.Addresses,
			validation.When(
				obj.TaxID != nil && obj.TaxID.Country.In(l10n.CO),
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

func municipalityCodeRequired(tID *tax.Identity) bool {
	if tID == nil {
		return false
	}
	if !tID.Country.In(l10n.CO) {
		return false
	}
	if tID.Type == TaxIdentityTypeCitizen || tID.Code == TaxCodeFinalCustomer {
		return false
	}
	return true
}

func (v *invoiceValidator) validSupplier(value interface{}) error {
	obj, _ := value.(*org.Party)
	if obj == nil || obj.TaxID == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			tax.RequireIdentityType,
			tax.IdentityTypeIn(TaxIdentityTypeTIN),
			tax.RequireIdentityCode,
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) validTaxIdentity(value interface{}) error {
	obj, _ := value.(*tax.Identity)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Code, validation.Required),
		validation.Field(&obj.Type,
			validation.Required,
			isValidTaxIdentityTypeKey,
			validation.When(!obj.Country.In(l10n.CO),
				// Certain types are exclusive of CO identities
				validation.NotIn(TaxIdentityTypeTIN, TaxIdentityTypeCitizen),
			),
		),
	)
}

func (v *invoiceValidator) preceding(value interface{}) error {
	obj, ok := value.(*bill.Preceding)
	if !ok || obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Ext,
			tax.ExtensionsRequires(ExtKeyDIANCorrection),
		),
		validation.Field(&obj.Reason, validation.Required),
	)
}

func normalizeParty(p *org.Party) error {
	// 2024-03-14: Migrate Tax ID Zone to extensions "co-dian-municipality"
	if p.TaxID != nil && p.TaxID.Zone != "" {
		if p.Ext == nil {
			p.Ext = make(tax.Extensions)
		}
		p.Ext[ExtKeyDIANMunicipality] = tax.ExtValue(p.TaxID.Zone)
		p.TaxID.Zone = ""
	}
	return nil
}
