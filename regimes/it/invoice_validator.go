package it

import (
	"errors"

	"github.com/invopop/gobl/bill"
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
		validation.Field(&inv.Supplier, validation.By(v.supplier)),
		validation.Field(&inv.Customer, validation.By(v.customer)),
	)
}

func (v *invoiceValidator) supplier(value interface{}) error {
	supplier, _ := value.(*org.Party)
	if supplier == nil {
		return errors.New("missing supplier details")
	}

	return validation.ValidateStruct(supplier,
		validation.Field(&supplier.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			tax.IdentityTypeIn(TaxIdentityTypeBusiness, TaxIdentityTypeGovernment),
		),
		validation.Field(&supplier.Addresses,
			validation.By(validateAddress("supplier")),
		),
		validation.Field(&supplier.Registration,
			validation.By(validateRegistration),
		),
	)
}

func (v *invoiceValidator) customer(value interface{}) error {
	customer, _ := value.(*org.Party)
	if customer == nil {
		return errors.New("missing customer details")
	}

	// Customers must have a tax ID (PartitaIVA) if they are Italian legal entities
	// like government offices and companies.
	return validation.ValidateStruct(customer,
		validation.Field(&customer.TaxID,
			validation.When(
				customer.TaxID.Country.In(l10n.IT),
				validation.Required,
				tax.RequireIdentityCode,
				tax.IdentityTypeIn(
					TaxIdentityTypeBusiness,
					TaxIdentityTypeGovernment,
					TaxIdentityTypeIndividual,
				),
			),
		),
		validation.Field(&customer.Addresses,
			validation.By(validateAddress("customer")),
		),
	)
}

func validateAddress(partyType string) validation.RuleFunc {
	return func(value interface{}) error {
		v, ok := value.([]*org.Address)
		if !ok {
			return errors.New("value is not a slice of Address")
		}

		if len(v) != 1 {
			return errors.New(partyType + " must have exactly one address")
		}

		address := v[0]

		return validation.ValidateStruct(address,
			validation.Field(&address.Country),
			validation.Field(&address.Locality, validation.Required),
			validation.Field(&address.Code, validation.Required),
			validation.Field(&address.Street, validation.Required),
			validation.Field(&address.Number, validation.Required),
		)
	}
}

func validateRegistration(value interface{}) error {
	v, ok := value.(*org.Registration)
	if !ok {
		return errors.New("value is not a valid Registration")
	}

	if v == nil {
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

	rf := ss.Codes[KeyFatturaPARegimeFiscale]
	if rf == "" {
		return errors.New("missing scenario related to RegimeFiscale")
	}

	return nil
}
