package it

import (
	"errors"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
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

	err := v.validateRetainedTaxes()
	if err != nil {
		return err
	}

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
				tax.IdentityTypeIn(TaxIdentityTypeBusiness, TaxIdentityTypeGovernment, TaxIdentityTypeIndividual),
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

// validateRetainedTaxes checks that the invoices with retained taxes has a valid
// retained tax tag included
func (v *invoiceValidator) validateRetainedTaxes() error {
	inv := v.inv

	for _, line := range inv.Lines {
		for _, combo := range line.Taxes {
			if !isRetainedTax(combo.Category) {
				continue
			}

			if !retainedTaxTagPresent(combo.Tags) {
				return errors.New(
					"invoice with retained taxes must include a valid retained tax tag. " +
						"List of tags are found under retainedTaxTags in regimes/it/tags")
			}
		}
	}

	return nil
}

func retainedTaxTagPresent(tags []cbc.Key) bool {
	if len(tags) == 0 {
		return false
	}

	for _, tag := range tags {
		for _, retainedTag := range retainedTaxTags {
			if tag == retainedTag.Key {
				return true
			}
		}
	}

	return false
}

func isRetainedTax(category cbc.Code) bool {
	for _, c := range categories {
		if c.Code == category {
			return c.Retained
		}
	}

	return false
}
