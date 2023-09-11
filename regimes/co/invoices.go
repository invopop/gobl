package co

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
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
		validation.Field(&inv.Currency, validation.In(currency.COP)),
		validation.Field(&inv.Type,
			validation.In(bill.InvoiceTypeStandard, bill.InvoiceTypeCreditNote),
		),
		validation.Field(&inv.Supplier,
			validation.Required,
			validation.By(v.validParty),
			validation.By(v.validSupplier),
		),
		validation.Field(&inv.Customer,
			validation.Required,
			validation.By(v.validParty),
		),
		validation.Field(&inv.Preceding, validation.When(
			inv.Type.In(bill.InvoiceTypeCreditNote),
			validation.Required,
			validation.Each(validation.By(v.preceding)),
		)),
		validation.Field(&inv.Outlays, validation.Empty),
	)
}

func (v *invoiceValidator) validParty(value interface{}) error {
	obj, _ := value.(*org.Party)
	if obj == nil || obj.TaxID == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.Required,
			validation.When(
				obj.TaxID.Country.In(l10n.CO),
				tax.RequireIdentityCode,
				validation.By(v.validTaxIdentity),
			),
		),
		validation.Field(&obj.Addresses,
			validation.When(
				obj.TaxID.Country.In(l10n.CO),
				validation.Length(1, 0),
			),
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
			tax.IdentityTypeIn(TaxIdentityTypeTIN),
		),
	)
}

func (v *invoiceValidator) validTaxIdentity(value interface{}) error {
	obj, _ := value.(*tax.Identity)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Zone,
			validation.In(zoneCodes...),
			validation.When(
				obj.Type.In(TaxIdentityTypeTIN),
				validation.Required,
			),
		),
	)
}

func (v *invoiceValidator) preceding(value interface{}) error {
	obj, ok := value.(*bill.Preceding)
	if !ok {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Method, validation.Required, isValidCorrectionMethodKey),
		validation.Field(&obj.Reason, validation.Required),
	)
}
