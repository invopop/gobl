package ar

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// invoiceValidator validates Argentine invoices according to AFIP regulations
func invoiceValidator(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(validateInvoiceSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.By(validateInvoiceCustomer(inv)),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(validation.By(validateInvoiceLine)),
			validation.Skip,
		),
	)
}

// validateInvoiceSupplier ensures supplier has required data for Argentine invoices
func validateInvoiceSupplier(value interface{}) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}

	return validation.ValidateStruct(p,
		validation.Field(&p.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
		),
	)
}

// validateInvoiceCustomer validates customer based on invoice type and tags
// Returns a validation function that has access to the invoice context
func validateInvoiceCustomer(inv *bill.Invoice) validation.RuleFunc {
	return func(value interface{}) error {
		p, ok := value.(*org.Party)
		if !ok || p == nil {
			// Customer is optional for simplified invoices and Type B to final consumers
			if inv.HasTags(tax.TagSimplified, TagInvoiceTypeB, TagInvoiceTypeC) {
				return nil
			}
			return validation.ErrRequired
		}

		// For Type A invoices, customer must have a CUIT/CUIL (tax ID)
		if inv.HasTags(TagInvoiceTypeA) {
			return validation.ValidateStruct(p,
				validation.Field(&p.TaxID,
					validation.Required.Error("customer tax ID required for Type A invoices"),
					tax.RequireIdentityCode,
				),
			)
		}

		// For Type E (export) invoices, validate foreign customer
		if inv.HasTags(TagInvoiceTypeE, TagExportServices, TagExportGoods) {
			// Export invoices should have a customer with either a foreign tax ID or identities
			if p.TaxID != nil && p.TaxID.Country.Code() == "AR" {
				return validation.NewError("validation_invalid", "export invoices should not have Argentine customers")
			}
		}

		return nil
	}
}

// validateInvoiceLine ensures invoice lines comply with Argentine requirements
func validateInvoiceLine(value interface{}) error {
	line, ok := value.(*bill.Line)
	if !ok || line == nil {
		return nil
	}

	// Basic line validation
	return validation.ValidateStruct(line,
		validation.Field(&line.Item,
			validation.Required,
		),
	)
}

// normalizeInvoice applies Argentine-specific normalization to invoices
func normalizeInvoice(inv *bill.Invoice) {
	// Normalize supplier and customer using common normalization
	if inv.Supplier != nil {
		normalizeParty(inv.Supplier)
	}
	if inv.Customer != nil {
		normalizeParty(inv.Customer)
	}

	// Handle simplified invoices - remove customer tax ID if it's a final consumer
	// This normalization is currently commented out as it depends on business requirements
	// Uncomment and implement if needed for your use case
	/*
		if inv.HasTags(tax.TagSimplified) {
			if inv.Customer != nil && inv.Customer.TaxID != nil {
				// For simplified invoices to final consumers, clear the tax ID
				inv.Customer.TaxID = nil
			}
		}
	*/

	// Ensure export invoices have zero-rated VAT
	if inv.HasTags(TagInvoiceTypeE, TagExportServices, TagExportGoods) {
		// Export invoices should use zero rate
		// This is typically handled by the tax calculation, but we can validate here
		ensureExportTaxRates(inv)
	}
}

// ensureExportTaxRates verifies that export invoices use appropriate tax rates
// This is currently a placeholder for future enhancement
func ensureExportTaxRates(_ *bill.Invoice) {
	// The actual tax calculation is handled by the regime's tax calculation logic
	// In a complete implementation, you might want to:
	// 1. Verify all lines use zero or exempt VAT rates
	// 2. Add automatic zero-rate application if not specified
	// 3. Validate that export documentation is present
}
