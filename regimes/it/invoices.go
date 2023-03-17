package it

import (
	"errors"

	"github.com/invopop/gobl/bill"
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
		validation.Field(&inv.Customer, validation.By(v.customer)),
		validation.Field(&inv.Supplier, validation.By(v.supplier)),
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

func (v *invoiceValidator) customer(value interface{}) error {
	p, _ := value.(*org.Party)

	if p == nil {
		return nil
	}

	// Customers must have a tax ID (PartitaIVA) if they are legal entities like
	// government offices and companies.
	return validation.ValidateStruct(p,
		validation.Field(&p.Type, validation.Required),
		validation.Field(&p.TaxID,
			validation.Required,
			validation.When(
				p.Type == PartyTypePublicAdministration || p.Type == PartyTypeLegalPerson,
				tax.RequireIdentityCode,
			),
		),
		validation.Field(&p.Identities,
			validation.Required,
			validation.When(
				p.Type == PartyTypeNaturalPerson,
				validation.By(v.customerNaturalPerson),
			),
		),
	)
}

func (v *invoiceValidator) customerNaturalPerson(value interface{}) error {
	p, _ := value.(*org.Party)
	if p == nil {
		return nil
	}

	var cf *org.Identity

	for _, id := range p.Identities {
		if id.Type == IdentityTypeCodiceFiscale {
			cf = id
		}
	}

	if cf == nil {
		return errors.New("natural persons must have a codice fiscale")
	}

	return nil
}
