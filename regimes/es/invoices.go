package es

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// invoiceValidator adds validation checks to invoices which are relevant
// for the region.
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
		validation.Field(&inv.Currency, validation.In(currency.EUR)),
		// Only commercial and simplified supported at this time for spain.
		// Rectification state determined by Preceding value.
		validation.Field(&inv.Type, validation.In(
			bill.InvoiceTypeStandard,
			bill.InvoiceTypeCorrective,
			bill.InvoiceTypeProforma,
		)),
		validation.Field(&inv.Preceding,
			validation.Each(validation.By(v.preceding)),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.Required,
			validation.By(v.supplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.By(v.commercialCustomer),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(
				validation.By(v.validateLine),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) supplier(value interface{}) error {
	obj, _ := value.(*org.Party)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
		),
	)
}

func (v *invoiceValidator) commercialCustomer(value interface{}) error {
	obj, _ := value.(*org.Party)
	if obj == nil {
		return nil
	}
	if obj.TaxID == nil {
		return nil // validation already handled, this prevents panics
	}
	// Customers must have a tax ID if a Spanish entity
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.Required,
			validation.When(
				obj.TaxID.Country.In(l10n.ES),
				tax.RequireIdentityCode,
			),
		),
	)
}

func (v *invoiceValidator) preceding(value interface{}) error {
	obj, _ := value.(*bill.Preceding)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Changes,
			validation.Required,
			validation.Each(isValidCorrectionChangeKey),
		),
	)
}

func (v *invoiceValidator) validateLine(value interface{}) error {
	obj, _ := value.(*bill.Line)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Taxes,
			validation.Each(
				validation.By(v.validateLineTax),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) validateLineTax(value interface{}) error {
	obj, ok := value.(*tax.Combo)
	if obj == nil || !ok {
		return nil
	}
	zone := l10n.CodeEmpty
	if v.inv.Supplier != nil && v.inv.Supplier.TaxID != nil {
		zone = v.inv.Supplier.TaxID.Zone
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Ext,
			validation.When(
				zone.In(ZonesBasqueCountry...) &&
					obj.Rate == tax.RateExempt,
				tax.ExtMapRequires(ExtKeyTBAIExemption),
			),
			validation.Skip,
		),
	)
}
