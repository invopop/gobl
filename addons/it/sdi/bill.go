package sdi

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
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

	// Normalize Italian supplier telephone numbers by stripping '+39' prefix
	if isItalianParty(party) && len(party.Telephones) > 0 {
		for _, tel := range party.Telephones {
			if tel != nil && len(tel.Number) >= 3 && tel.Number[:3] == "+39" {
				tel.Number = tel.Number[3:]
			}
		}
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
		validation.Field(&inv.Ordering,
			// Need to access tagas so we pass the invoice directly
			validation.By(validateInvoiceOrdering(inv)),
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
		validation.Field(&supplier.Name,
			validation.By(validateLatin1String),
			validation.Skip,
		),
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
		validation.Field(&supplier.Telephones,
			validation.When(
				isItalianParty(supplier) && len(supplier.Telephones) > 0,
				validation.Each(validation.By(validateTelephone)),
			),
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
		validation.Field(&customer.Name,
			validation.By(validateLatin1String),
			validation.When(
				(customer.TaxID != nil && customer.TaxID.Code != cbc.CodeEmpty) || customer.People == nil,
				validation.Required,
			),
			validation.Skip,
		),
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
		validation.Field(&customer.People,
			validation.When(
				(customer.Name == ""),
				validation.Required,
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
			validation.By(validateLatin1String),
			validation.Skip,
		),
	)
}

func validateCharge(val any) error {
	charge, _ := val.(*bill.Charge)
	if charge == nil || !charge.Key.Has(KeyFundContribution) {
		return nil
	}

	return validation.ValidateStruct(charge,
		validation.Field(&charge.Percent,
			validation.Required,
			validation.Skip,
		),
		validation.Field(&charge.Ext,
			tax.ExtensionsRequire(ExtKeyFundType),
			validation.Skip,
		),
		validation.Field(&charge.Taxes,
			tax.SetHasCategory(tax.CategoryVAT),
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

func validateInvoiceOrdering(inv *bill.Invoice) validation.RuleFunc {
	return func(value any) error {
		o, _ := value.(*bill.Ordering)
		if o == nil {
			return nil
		}

		return validation.ValidateStruct(o,
			validation.Field(&o.Despatch,
				validation.When(
					inv.HasTags(TagDeferred),
					validation.Each(validation.By(validateDespatch)),
				).Else(
					validation.Nil.Error("can only be set when invoice has deferred tag")),
				validation.Skip,
			),
		)
	}
}

func validateDespatch(value any) error {
	d, ok := value.(*org.DocumentRef)
	if !ok || d == nil {
		return nil
	}
	return validation.ValidateStruct(d,
		validation.Field(&d.IssueDate,
			validation.Required,
			validation.Skip,
		),
	)
}

func validateTelephone(value any) error {
	t, ok := value.(*org.Telephone)
	if !ok {
		return nil
	}
	return validation.ValidateStruct(t,
		validation.Field(&t.Number,
			validation.Required,
			validation.Length(5, 12),
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
		validation.Field(&v.Street,
			validation.When(v.PostOfficeBox == "",
				validation.Required.Error("either street or post office box must be set"),
			),
			validation.By(validateLatin1String),
			validation.Skip,
		),
		validation.Field(&v.PostOfficeBox,
			validation.When(v.Street == "",
				validation.Required.Error("either street or post office box must be set"),
			),
			validation.By(validateLatin1String),
			validation.Skip,
		),
		validation.Field(&v.Country,
			validation.Required,
			validation.Skip,
		),
		validation.Field(&v.Locality,
			validation.Required,
			validation.By(validateLatin1String),
			validation.Skip,
		),
		validation.Field(&v.Code,
			validation.When(v.Country.In("IT"),
				validation.Required,
				validation.Match(regexp.MustCompile(`^\d{5}$`)),
			),
			validation.Skip,
		),
	)
}

// validateLatin1String ensures that the item name only contains characters
// from Latin and Latin-1 range (ASCII 0-127 and extended Latin-1 128-255).
func validateLatin1String(val any) error {
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

func validateInvoiceSupplierRegistration(value interface{}) error {
	v, ok := value.(*org.Registration)
	if v == nil || !ok {
		return nil
	}
	return validation.ValidateStruct(v,
		validation.Field(&v.Entry,
			validation.Required,
			validation.Skip,
		),
		validation.Field(&v.Office,
			validation.Required,
			validation.Skip,
		),
	)
}
