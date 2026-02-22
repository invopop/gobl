package jp

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// validateInvoice enforces NTA Qualified Invoice requirements:
//   - Supplier must provide name and tax ID (T-number).
//   - Customer is required unless the invoice is tagged "simplified" (簡易適格請求書).
//   - Export-tagged invoices must have all VAT lines zero-rated.
//   - Addresses must include at least a street or locality.
func validateInvoice(inv *bill.Invoice) error {
	simplified := inv.HasTags(TagSimplified)
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(validateSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.When(!simplified,
				validation.Required,
				validation.By(validateCustomer),
			),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.By(validateInvoiceLines),
			validation.By(checkExportRates(inv.HasTags(TagExport))),
			validation.Skip,
		),
	)
}

// checkExportRates returns a validator that ensures all VAT lines use the zero rate when the invoice carries the export
// tag.
func checkExportRates(isExport bool) validation.RuleFunc {
	return func(value any) error {
		if !isExport {
			return nil
		}
		return validateExportLines(value)
	}
}

// validateSupplier validates the supplier party on an invoice
func validateSupplier(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}

	return validation.ValidateStruct(p,
		// In Japan, invoices should have a tax ID (Corporate Number or QII Registration Number)
		validation.Field(&p.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
		// Name is required for Japanese invoices
		validation.Field(&p.Name,
			validation.Required,
		),
		// Addresses should be validated
		validation.Field(&p.Addresses,
			validation.By(validateAddresses),
			validation.Skip,
		),
	)
}

// validateCustomer validates the customer party on an invoice
func validateCustomer(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}

	// Customer validation is less strict than supplier, but still requires basic information.
	return validation.ValidateStruct(p,
		validation.Field(&p.Name,
			validation.Required,
		),
	)
}

// validateAddresses validates address information
func validateAddresses(value any) error {
	addresses, ok := value.([]*org.Address)
	if !ok || addresses == nil {
		return nil
	}

	for _, addr := range addresses {
		if err := validateAddress(addr); err != nil {
			return err
		}
	}
	return nil
}

// validateAddress validates a single address
func validateAddress(addr *org.Address) error {
	if addr == nil {
		return nil
	}

	return validation.ValidateStruct(addr,
		// Japanese addresses typically require a street address.
		validation.Field(&addr.Street,
			validation.When(addr.Locality == "", validation.Required),
		),
		// Locality (city/ward) is typically required in Japan
		validation.Field(&addr.Locality,
			validation.When(addr.Street == "", validation.Required),
		),
	)
}

// validateInvoiceLines validates invoice line items
func validateInvoiceLines(value any) error {
	lines, ok := value.([]*bill.Line)
	if !ok || lines == nil {
		return nil
	}

	for _, line := range lines {
		if err := validateInvoiceLine(line); err != nil {
			return err
		}
	}
	return nil
}

// validateExportLines ensures all VAT lines use the zero rate for export invoices.
func validateExportLines(value any) error {
	lines, ok := value.([]*bill.Line)
	if !ok || lines == nil {
		return nil
	}
	for _, line := range lines {
		if line == nil {
			continue
		}
		for _, combo := range line.Taxes {
			if combo == nil {
				continue
			}
			if combo.Category != tax.CategoryVAT {
				continue
			}
			if combo.Key != tax.KeyZero {
				return validation.NewError(
					"validation_export_zero_rate",
					"export invoices must use the zero rate for all VAT lines",
				)
			}
		}
	}
	return nil
}

// validateInvoiceLine validates a single invoice line
func validateInvoiceLine(line *bill.Line) error {
	if line == nil {
		return nil
	}

	return validation.ValidateStruct(line,
		// Each line should have a quantity
		validation.Field(&line.Quantity,
			validation.Required,
		),
		// Each line should have an item
		validation.Field(&line.Item,
			validation.Required,
		),
	)
}
