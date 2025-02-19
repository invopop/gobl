package sdi

import (
	"regexp"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func normalizeInvoice(inv *bill.Invoice) {
	normalizeSupplier(inv.Supplier)
}

func normalizeSupplier(party *org.Party) {
	if party == nil {
		return
	}
	if party.Ext == nil || party.Ext[ExtKeyFiscalRegime] == "" {
		if party.Ext == nil {
			party.Ext = make(tax.Extensions)
		}
		party.Ext[ExtKeyFiscalRegime] = "RF01" // Ordinary regime is default
	}
}

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Tax,
			validation.By(validateTax),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.By(validateSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.By(validateCustomer),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(
				bill.RequireLineTaxCategory(tax.CategoryVAT),
				validation.Skip,
			),
			validation.Skip,
		),
		validation.Field(&inv.Payment,
			validation.By(validateInvoicePaymentDetails),
			validation.Skip,
		),
	)
}

func validateTax(value any) error {
	obj, _ := value.(*bill.Tax)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Ext,
			tax.ExtensionsHas(
				ExtKeyFormat,
				ExtKeyDocumentType,
			),
			validation.Skip,
		),
	)
}

func validateSupplier(value interface{}) error {
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
			validation.By(validateInvoiceSupplierRegistration),
			validation.Skip,
		),
		validation.Field(&supplier.Ext,
			tax.ExtensionsRequire(ExtKeyFiscalRegime),
			validation.Skip,
		),
	)
}

func validateCustomer(value interface{}) error {
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
				org.RequireIdentityKey(it.IdentityKeyFiscalCode),
			),
			validation.Skip,
		),
	)
}

func validateInvoicePaymentDetails(val any) error {
	p, _ := val.(*bill.PaymentDetails)
	if p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.Instructions,
			validation.When(
				(p.Terms != nil && len(p.Terms.DueDates) > 0),
				validation.Required.Error("cannot be blank when terms with due dates are present"),
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
	return org.IdentityForKey(party.Identities, it.IdentityKeyFiscalCode) != nil

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

func validateInvoiceSupplierRegistration(value interface{}) error {
	v, ok := value.(*org.Registration)
	if v == nil || !ok {
		return nil
	}
	return validation.ValidateStruct(v,
		validation.Field(&v.Entry, validation.Required),
		validation.Field(&v.Office, validation.Required),
	)
}
