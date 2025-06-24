package choruspro

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
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
				inv.Tax != nil && inv.Tax.Ext != nil && inv.Tax.Ext.Get(ExtKeyFramework) == "A2",
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
		validation.Field(&customer.Ext,
			tax.ExtensionsHasCodes(ExtKeyScheme, "1"),
			validation.Skip,
		),
		validation.Field(&customer.Identities,
			validation.Required,
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
			tax.ExtensionsRequire(ExtKeyFramework),
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
	if !inv.Tax.Ext.Has(ExtKeyFramework) {
		inv.Tax.Ext = inv.Tax.Ext.Merge(
			tax.Extensions{
				ExtKeyFramework: "A1",
			},
		)
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
