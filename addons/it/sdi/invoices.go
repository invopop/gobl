package sdi

import (
	"errors"
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
			validation.Required,
			validation.By(validateTax),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.By(validateSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.Required,
			validation.By(validateCustomer),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(
				validation.By(validateLine),
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
	obj, ok := value.(*bill.Tax)
	if !ok {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Ext,
			tax.ExtensionsRequire(
				ExtKeyFormat,
				ExtKeyDocumentType,
			),
			validation.Skip,
		),
	)
}

func validateSupplier(value interface{}) error {
	supplier, ok := value.(*org.Party)
	if !ok || supplier == nil {
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
	customer, ok := value.(*org.Party)
	if !ok {
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
			validation.Required,
			validation.Each(validation.By(validateAddress)),
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

func validateLine(val any) error {
	line, _ := val.(*bill.Line)
	if line == nil {
		return nil
	}

	return validation.ValidateStruct(line,
		validation.Field(&line.Item,
			validation.By(validateItem),
			validation.Skip,
		),
	)
}

func validateItem(val any) error {
	item, _ := val.(*org.Item)
	if item == nil {
		return nil
	}

	return validation.ValidateStruct(item,
		validation.Field(&item.Name,
			validation.By(validateItemName),
			validation.Skip,
		),
	)
}

// validateItemName ensures that the item name only contains characters
// from Latin and Latin-1 range (ASCII 0-127 and extended Latin-1 128-255).
func validateItemName(val any) error {
	name, _ := val.(string)

	for _, r := range name {
		// Check if the character is outside Latin and Latin-1 range
		// Latin and Latin-1 includes ASCII (0-127) and extended Latin-1 (128-255)
		if r > 255 {
			return errors.New("contains characters outside of Latin and Latin-1 range")
		}
	}
	return nil
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
	if party == nil {
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
	if !ok {
		return nil
	}
	// Post code and street in addition to the locality are required in Italian invoices.
	return validation.ValidateStruct(v,
		validation.Field(&v.Street, validation.Required),
		validation.Field(&v.Country, validation.Required),
		validation.Field(&v.Code,
			validation.When(v.Country.In("IT"),
				validation.Required,
				validation.Match(regexp.MustCompile(`^\d{5}$`)),
			),
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
