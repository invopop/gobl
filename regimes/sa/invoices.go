package sa

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// validateInvoice checks the SA invoice for compliance with VAT
// Implementing Regulations (Article 53).
// Reference: https://zatca.gov.sa/en/RulesRegulations/Taxes/Documents/Implmenting%20Regulations%20of%20the%20VAT%20Law_EN.pdf
func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		// Article 53.5.c, 53.8.b, 57: Supplier TIN required on all invoices
		validation.Field(&inv.Supplier,
			validation.By(validateInvoiceSupplier),
			validation.Skip,
		),
		// Article 53.1.a, 53.5.e: Customer required on standard (non-simplified) invoices
		// Article 53.5.d: Customer TIN required when reverse charge applies
		validation.Field(&inv.Customer,
			validation.When(
				!inv.HasTags(tax.TagSimplified),
				validation.Required,
				validation.By(validateInvoiceCustomer(inv)),
			),
			validation.Skip,
		),
	)
}

func validateInvoiceSupplier(val any) error {
	obj, _ := val.(*org.Party)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
	)
}

func validateInvoiceCustomer(inv *bill.Invoice) validation.RuleFunc {
	return func(val any) error {
		obj, _ := val.(*org.Party)
		if obj == nil {
			return nil
		}
		return validation.ValidateStruct(obj,
			// Article 53.5.d: Customer TIN required on reverse-charge invoices
			validation.Field(&obj.TaxID,
				validation.When(
					inv.HasTags(tax.TagReverseCharge),
					validation.Required,
					tax.RequireIdentityCode,
				),
				validation.Skip,
			),
		)
	}
}
