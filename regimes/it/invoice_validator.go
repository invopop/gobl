package it

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
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
		// Currently only supporting invoices and credit notes for Italy
		validation.Field(&inv.Type, validation.In(
			bill.InvoiceTypeNone,
			bill.InvoiceTypeCreditNote,
		)),
		validation.Field(&inv.Supplier, validation.Required, validation.By(v.supplier)),
		validation.Field(&inv.Customer, validation.Required, validation.By(v.customer)),
	)
}

func (v *invoiceValidator) supplier(value interface{}) error {
	sup, _ := value.(*org.Party)
	if sup == nil {
		return nil
	}
	// Suppliers must have a VAT ID (Partita IVA)
	return validation.ValidateStruct(sup,
		validation.Field(&sup.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
		),
	)
}

func (v *invoiceValidator) customer(value interface{}) error {
	cus, _ := value.(*org.Party)
	if cus == nil {
		return nil
	}
	if cus.TaxID == nil {
		return nil // validation already handled, this prevents panics
	}
	// Customers must have a VAT ID (Partita IVA) OR a fiscal code (Codice Fiscale)
	return validation.ValidateStruct(cus,
		validation.Field(&cus.TaxID, validation.When(
			fiscalCode(cus) != nil,
			validation.Required,
			tax.RequireIdentityCode,
		)),
		validation.Field(fiscalCode(cus), validation.When(
			cus.TaxID == nil,
			validation.Required,
		)),
	)
}

func fiscalCode(party *org.Party) *org.Identity {
	for _, identity := range party.Identities {
		if identity.Type == IdentityTypeFiscalCode {
			return identity
		}
	}
	return nil
}
