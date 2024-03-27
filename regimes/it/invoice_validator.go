package it

import (
	"errors"
	"regexp"

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
	if err := v.validateScenarios(); err != nil {
		return err
	}

	inv := v.inv
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Currency,
			validation.In(currency.EUR),
		),
		validation.Field(&inv.Supplier,
			validation.By(v.supplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.By(v.customer),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(
				validation.By(validateLine),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) supplier(value interface{}) error {
	supplier, ok := value.(*org.Party)
	if !ok {
		return nil
	}

	return validation.ValidateStruct(supplier,
		validation.Field(&supplier.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			tax.IdentityTypeIn(TaxIdentityTypeBusiness, TaxIdentityTypeGovernment),
		),
		validation.Field(&supplier.Addresses,
			validation.Required,
			validation.Each(validation.By(validateAddress)),
		),
		validation.Field(&supplier.Registration,
			validation.By(validateRegistration),
		),
	)
}

func (v *invoiceValidator) customer(value interface{}) error {
	customer, _ := value.(*org.Party)
	if customer == nil {
		return nil
	}

	// Customers must have a tax ID (PartitaIVA)
	return validation.ValidateStruct(customer,
		validation.Field(&customer.TaxID,
			validation.Required,
			validation.When(
				isItalianParty(customer),
				tax.RequireIdentityCode,
				tax.IdentityTypeIn(
					TaxIdentityTypeBusiness,
					TaxIdentityTypeGovernment,
					TaxIdentityTypeIndividual,
				),
			),
		),
		validation.Field(&customer.Addresses,
			validation.When(
				isItalianParty(customer),
				// TODO: address not required for simplified invoices
				validation.Each(validation.By(validateAddress)),
			),
		),
	)
}

func isItalianParty(party *org.Party) bool {
	if party == nil || party.TaxID == nil {
		return false
	}
	return party.TaxID.Country.In(l10n.IT)
}

func validateAddress(value interface{}) error {
	v, ok := value.(*org.Address)
	if v == nil || !ok {
		return nil
	}
	// Post code and street in addition to the locality are required in Italian invoices.
	return validation.ValidateStruct(v,
		validation.Field(&v.Street, validation.Required),
		validation.Field(&v.Code,
			validation.Required,
			validation.Match(regexp.MustCompile(`^\d{5}$`)),
		),
	)
}

func validateLine(value interface{}) error {
	v, ok := value.(*bill.Line)
	if v == nil || !ok {
		return nil
	}
	return validation.ValidateStruct(v,
		validation.Field(&v.Taxes,
			tax.SetHasCategory(tax.CategoryVAT),
			validation.Each(
				validation.By(validateLineTax),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}

func validateLineTax(value interface{}) error {
	v, ok := value.(*tax.Combo)
	if v == nil || !ok {
		return nil
	}
	return validation.ValidateStruct(v,
		validation.Field(&v.Ext,
			validation.When(
				v.Category.In(
					TaxCategoryIRPEF,
					TaxCategoryIRES,
					TaxCategoryINPS,
					TaxCategoryENASARCO,
					TaxCategoryENPAM,
				),
				tax.ExtensionsRequires(
					ExtKeySDIRetainedTax,
				),
			),
			validation.Skip,
		),
	)
}

func validateRegistration(value interface{}) error {
	v, ok := value.(*org.Registration)
	if v == nil || !ok {
		return nil
	}
	return validation.ValidateStruct(v,
		validation.Field(&v.Entry, validation.Required),
		validation.Field(&v.Office, validation.Required),
	)
}

// validateScenarios checks that the invoice includes scenarios that help determine
// TipoDocumento and RegimeFiscale
func (v *invoiceValidator) validateScenarios() error {
	ss := v.inv.ScenarioSummary()

	td := ss.Codes[KeyFatturaPATipoDocumento]
	if td == "" {
		return errors.New("missing scenario related to TipoDocumento")
	}

	return nil
}
