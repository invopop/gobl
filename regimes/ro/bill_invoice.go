package ro

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Romanian invoice validation based on:
// - Law 227/2015 (Fiscal Code - Art. 319 Invoicing): https://legislatie.just.ro/Public/DetaliiDocument/171282
// - Law 296/2023 (B2B e-Factura Mandate): https://legislatie.just.ro/Public/DetaliiDocumentAfis/275745
// - OUG 69/2024 (B2C e-Factura Mandate effective Jan 1, 2025): https://legislatie.just.ro/Public/DetaliiDocument/283893
// - OUG 138/2024 (Simplified Invoices Mandate effective Jan 1, 2025): https://legislatie.just.ro/Public/DetaliiDocument/293675
// - OUG 120/2021 (e-Factura Framework & RO_CIUS): https://legislatie.just.ro/Public/DetaliiDocument/246848
//
// Key requirements:
// - Supplier must have a valid Romanian tax ID (CIF/CUI).
// - B2B Customers must have a valid tax ID (CUI).
// - B2C Customers (individuals) are mandatory in e-Factura as of Jan 1, 2025 (OUG 69/2024).
//   If a CNP is not provided, the invoice is still valid but reported with a standard placeholder (13 zeros).
// - Simplified invoices (<100 EUR) are mandatory in e-Factura as of Jan 1, 2025 (OUG 138/2024).
// - Reporting deadline is 5 calendar days from issuance.

// validateBillInvoice validates Romanian invoices according to local requirements.
func validateBillInvoice(inv *bill.Invoice) error {
	// Check for the simplified tag (supported since OUG 138/2024)
	simplified := inv.Tags.HasTags(tax.TagSimplified)

	return validation.ValidateStruct(inv,
		validation.Field(&inv.Type,
			validation.In(
				bill.InvoiceTypeStandard,
				bill.InvoiceTypeCreditNote,
				bill.InvoiceTypeDebitNote,
				bill.InvoiceTypeCorrective,
			),
			validation.Skip,
		),
		validation.Field(&inv.Preceding,
			validation.When(
				inv.Type.In(
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
					bill.InvoiceTypeCorrective,
				),
				validation.Required.Error("preceding invoice reference is mandatory for credit/debit notes (Art. 319 Fiscal Code)"),
				validation.Skip,
			),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.Required,
			validation.By(validateSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			// Customer is MANDATORY for Standard, Credit, and Debit notes.
			// However, if the simplified tag is present, we relax this strict requirement.
			validation.When(
				!simplified,
				validation.Required.Error("customer is required for non-simplified invoices"),
			),
			// Even if optional, if the customer struct is provided, it must be valid
			validation.By(validateCustomer),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Required,
			validation.Each(
				validation.By(validateInvoiceLine),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}

// validateSupplier validates the supplier party on Romanian invoices.
// The supplier must have a valid tax ID.
func validateSupplier(value any) error {
	supplier, ok := value.(*org.Party)
	if !ok || supplier == nil {
		return nil
	}

	return validation.ValidateStruct(supplier,
		validation.Field(&supplier.TaxID,
			validation.Required.Error("supplier must have a tax ID"),
			tax.RequireIdentityCode,
			validation.Skip,
		),
	)
}

// validateCustomer validates the customer party on Romanian invoices.
// For B2B transactions, the customer should have identification.
// For B2C transactions (mandatory since Jan 1, 2025), a Tax ID (CNP) is optional;
// if missing, the e-Factura transmission layer handles the required placeholder.
func validateCustomer(value any) error {
	customer, ok := value.(*org.Party)
	if !ok || customer == nil {
		return nil
	}

	return validation.ValidateStruct(customer,
		validation.Field(&customer.TaxID,
			// Tax ID is mandatory for B2B but optional for B2C (individuals).
			// We skip required validation here to support B2C flows where CNP is withheld for privacy.
			validation.Skip,
		),
	)
}

// validateInvoiceLine validates individual invoice lines.
func validateInvoiceLine(value any) error {
	line, ok := value.(*bill.Line)
	if !ok || line == nil {
		return nil
	}

	return validation.ValidateStruct(line,
		validation.Field(&line.Item,
			validation.Required.Error("line item is required"),
			validation.Skip,
		),
		validation.Field(&line.Quantity,
			validation.Required.Error("line quantity is required"),
			validation.Skip,
		),
	)
}
