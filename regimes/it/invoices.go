package it

import (
	"regexp"

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

// normalizeInvoice is used to ensure the invoice data is correct.
func normalizeInvoice(inv *bill.Invoice) error {
	return normalizeCustomer(inv.Customer)
}

func normalizeCustomer(party *org.Party) error {
	if party == nil {
		return nil
	}
	if !isItalianParty(party) {
		return nil
	}
	// If the party is an individual, move the fiscal code to the identities.
	if party.TaxID.Type == "individual" { //nolint:staticcheck
		id := &org.Identity{
			Key:  IdentityKeyFiscalCode,
			Code: party.TaxID.Code,
		}
		party.TaxID.Code = ""
		party.TaxID.Type = "" //nolint:staticcheck
		party.Identities = org.AddIdentity(party.Identities, id)
	}
	return nil
}

func validateInvoice(inv *bill.Invoice) error {
	v := &invoiceValidator{inv: inv}
	return v.validate()
}

func (v *invoiceValidator) validate() error {
	inv := v.inv
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Tax,
			validation.By(v.tax),
			validation.Skip,
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

func (v *invoiceValidator) tax(value any) error {
	obj, _ := value.(*bill.Tax)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Ext,
			tax.ExtensionsHas(
				ExtKeySDIFormat,
				ExtKeySDIDocumentType,
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
			validation.Skip,
		),
		validation.Field(&supplier.Addresses,
			validation.Required,
			validation.Each(validation.By(validateAddress)),
			validation.Skip,
		),
		validation.Field(&supplier.Registration,
			validation.By(validateRegistration),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) customer(value interface{}) error {
	customer, _ := value.(*org.Party)
	if customer == nil {
		return nil
	}

	// Customers must have either a Tax ID (PartitaIVA)
	// or fiscal identity (codice fiscale)
	return validation.ValidateStruct(customer,
		validation.Field(&customer.TaxID,
			validation.Required,
			validation.When(
				isItalianParty(customer) && !hasFiscalCode(customer),
				tax.RequireIdentityCode,
			),
			validation.Skip,
		),
		validation.Field(&customer.Addresses,
			validation.When(
				isItalianParty(customer),
				// TODO: address not required for simplified invoices
				validation.Each(validation.By(validateAddress)),
			),
			validation.Skip,
		),
		validation.Field(&customer.Identities,
			validation.When(
				isItalianParty(customer) && !hasTaxIDCode(customer),
				org.RequireIdentityKey(IdentityKeyFiscalCode),
			),
			validation.Skip,
		),
	)
}

func hasTaxIDCode(party *org.Party) bool {
	return party != nil && party.TaxID != nil && party.TaxID.Code != ""
}

func hasFiscalCode(party *org.Party) bool {
	if party == nil || party.TaxID == nil {
		return false
	}
	return org.IdentityForKey(party.Identities, IdentityKeyFiscalCode) != nil

}

func isItalianParty(party *org.Party) bool {
	if party == nil || party.TaxID == nil {
		return false
	}
	return party.TaxID.Country.In("IT")
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
