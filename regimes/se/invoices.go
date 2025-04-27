package se

import (
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
	inv := v.inv
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.Required,
			validation.By(v.validateOrgParty),
			validation.By(v.validateSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.Required,
			validation.By(v.validateOrgParty),
			validation.By(v.validateCustomer),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) validateOrgParty(value any) error {
	party, ok := value.(*org.Party)
	if !ok || party == nil {
		return nil
	}

	return validation.ValidateStruct(party,
		validation.Field(&party.Name,
			validation.Required,
			validation.Skip,
		),
		validation.Field(&party.Addresses,
			validation.Required,
			validation.Skip,
		),
	)
}

// validateSupplier checks the supplier's tax ID and organization number requirements.
// The supplier's VAT number and ID Number are always required.
func (v *invoiceValidator) validateSupplier(value any) error {
	party, ok := value.(*org.Party)
	if !ok || party == nil {
		return nil
	}

	return validation.ValidateStruct(party,
		// Swedish Tax ID (VAT ID) is always required for the supplier,
		// and must have the correct format.
		validation.Field(&party.TaxID,
			validation.Required,
			validation.By(func(value interface{}) error {
				tID, ok := value.(*tax.Identity)
				if !ok || tID == nil {
					return nil
				}
				return validateTaxIdentity(tID)
			}),
			validation.Skip,
		),
		validation.Field(&party.Identities,
			validation.Required,
			validation.In(
				org.RequireIdentityType(IdentityTypeOrgNr),
				org.RequireIdentityType(IdentityTypePersonNr),
				org.RequireIdentityType(IdentityTypeCoordinationNr),
			),
		),
	)
}

// validateCustomer checks the customer's tax ID and organization number requirements.
// However, the customer may not include any identities, just a name and address.
func (v *invoiceValidator) validateCustomer(value any) error {
	party, ok := value.(*org.Party)
	if !ok || party == nil {
		return nil
	}

	return validation.ValidateStruct(party,
		validation.Field(&party.Identities,
			validation.When(
				party.TaxID != nil && party.TaxID.Country == l10n.TaxCountryCode(l10n.SE),
				validation.In(
					org.RequireIdentityType(IdentityTypeOrgNr),
					org.RequireIdentityType(IdentityTypePersonNr),
					org.RequireIdentityType(IdentityTypeCoordinationNr),
				),
			),
		),
	)
}
