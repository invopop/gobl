package scontrino

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func normalizeInvoice(inv *bill.Invoice) {
	panic("check if required")
}

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(validateSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(
				validation.By(validateLine),
				validation.Skip,
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
	)
}

func validateLine(value interface{}) error {
	line, ok := value.(*bill.Line)
	if !ok {
		return nil
	}

	return validation.ValidateStruct(line,
		validation.Field(&line.Taxes,
			validation.Required,
			bill.RequireLineTaxCategory(tax.CategoryVAT),
			validation.Skip,
		),
		validation.Field(&line.Quantity, validation.Required),
	)
}
