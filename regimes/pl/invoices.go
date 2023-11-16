package pl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/common"
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
		validation.Field(&inv.Currency, validation.In(currency.PLN, currency.EUR)),
		// Only commercial and simplified supported at this time for Poland.
		// Rectification state determined by Preceding value.
		validation.Field(&inv.Type, validation.In(
			bill.InvoiceTypeStandard,
			bill.InvoiceTypeCorrective,
			bill.InvoiceTypeProforma,
		)),
		// validation.Field(&inv.Preceding,
		// 	validation.Each(validation.By(v.preceding)),
		// ),
		validation.Field(&inv.Supplier,
			validation.Required,
			validation.By(v.supplier),
		),
		validation.Field(&inv.Customer,
			validation.When(
				!inv.Tax.ContainsTag(common.TagSimplified),
				validation.Required,
				validation.By(v.commercialCustomer),
			),
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
			validation.By(validatePolishTaxIdentity),
		),
		// TODO check if name exists
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
	// Customers must have a tax ID if a Polish entity
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.Required,
			validation.When(
				obj.TaxID.Country.In(l10n.PL),
				validation.By(validatePolishTaxIdentity),
			), // TODO check if id is valid when other entity
			// TODO check if name exists
		),
	)
}

// func (v *invoiceValidator) preceding(value interface{}) error {
// 	obj, _ := value.(*bill.Preceding)
// 	if obj == nil {
// 		return nil
// 	}
// 	return validation.ValidateStruct(obj,
// 		validation.Field(&obj.Changes,
// 			validation.Required,
// 			validation.Each(isValidCorrectionChangeKey),
// 		),
// 		validation.Field(&obj.CorrectionMethod,
// 			validation.Required,
// 			isValidCorrectionMethodKey,
// 		),
// 	)
// }
