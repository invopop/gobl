package choruspro

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// validateInvoice ensures that the invoice meets Chorus Pro requirements
func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Customer,
			validation.By(validateCustomer),
		),
		validation.Field(&inv.Tax,
			validation.Required,
			validation.By(validateTax),
		),
		validation.Field(&inv.Totals,
			validation.When(
				// A2 can only exist if invoice has been paid
				inv.Tax != nil && inv.Tax.Ext != nil && inv.Tax.Ext.Get(ExtKeyFrameWork) == "A2",
				validation.By(validatePaid),
			)),
	)
}

func validateCustomer(value interface{}) error {
	customer, ok := value.(*org.Party)
	if !ok || customer == nil {
		return nil
	}

	return validation.ValidateStruct(customer,
		validation.Field(&customer.Identities,
			validation.Required,
			validation.Each(
				validation.By(validateIdentity),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}

func validateIdentity(value interface{}) error {
	identity, ok := value.(*org.Identity)
	if !ok {
		return nil
	}

	return validation.ValidateStruct(identity,
		validation.Field(&identity.Ext,
			validation.When(
				identity.Type == fr.IdentityTypeSIRET,
				tax.ExtensionsRequire(ExtKeyScheme),
				tax.ExtensionsHasCodes(ExtKeyScheme, "1"),
			),
			validation.Skip,
		),
	)
}

func validateTax(value interface{}) error {
	t, ok := value.(*bill.Tax)
	if !ok || t == nil {
		return nil
	}

	return validation.ValidateStruct(t,
		validation.Field(&t.Ext,
			validation.Required,
			tax.ExtensionsRequire(ExtKeyFrameWork),
		),
	)
}

// normalizeInvoice applies Chorus Pro specific normalization rules
func normalizeInvoice(inv *bill.Invoice) {
	if inv == nil {
		return
	}

	// Ensure required extensions are set with default values if not present
	if inv.Tax == nil {
		inv.Tax = &bill.Tax{}
	}
	if inv.Tax.Ext == nil {
		inv.Tax.Ext = make(tax.Extensions)
	}

	// Set default framework type if not specified
	if !inv.Tax.Ext.Has(ExtKeyFrameWork) {
		inv.Tax.Ext = inv.Tax.Ext.Merge(
			tax.Extensions{
				ExtKeyFrameWork: "A1",
			},
		)
	}

	normalizeSupplier(inv.Supplier)

}

func normalizeSupplier(supplier *org.Party) {
	if supplier == nil || supplier.TaxID == nil || supplier.TaxID.Code == "" {
		return
	}
	if supplier.TaxID.Country != "FR" {
		if l10n.Unions().Code(l10n.EU).HasMember(l10n.Code(supplier.TaxID.Country)) {
			supplier.TaxID.Scheme = "2"
		} else {
			supplier.TaxID.Scheme = "3"
		}
	}
}

func normalizeParty(party *org.Party) {
	// This is a workaround to ensure the addon normalizes after the regime
	// has normalized the tax ID.
	normalizeIdentities(party.Identities)
}

func normalizeIdentities(identities []*org.Identity) {
	if identities == nil {
		return
	}

	// If we have a SIRET, we need to set the scheme to 1
	for _, identity := range identities {
		if identity.Type == fr.IdentityTypeSIRET {
			if identity.Ext == nil {
				identity.Ext = make(tax.Extensions)
			}
			identity.Ext[ExtKeyScheme] = "1"
			break
		}
	}
}

func validatePaid(value interface{}) error {
	totals, ok := value.(*bill.Totals)
	if !ok {
		return nil
	}
	if !totals.Paid() {
		return validation.NewError("totals", "If the invoice has type A2, it must be paid in full")
	}
	return nil
}
