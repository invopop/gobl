package ticket

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func normalizeInvoice(inv *bill.Invoice) {
	if inv.Tax == nil {
		inv.Tax = new(bill.Tax)
	}
	if inv.Tax.PricesInclude == "" {
		inv.Tax.PricesInclude = tax.CategoryVAT
	}
}

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Tax,
			validation.Required,
			validation.By(validateInvoiceTax),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.By(validateInvoiceSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(
				bill.RequireLineTaxCategory(tax.CategoryVAT),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}

func validateInvoiceSupplier(value interface{}) error {
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
	)
}

// This done because the format requires tax to be calculated at item level
// By forcing this we can ensure that the price already has the tax included
func validateInvoiceTax(value interface{}) error {
	t, ok := value.(*bill.Tax)
	if !ok || t == nil {
		return nil
	}

	return validation.ValidateStruct(t,
		validation.Field(&t.PricesInclude,
			validation.Required,
			validation.In(tax.CategoryVAT),
			validation.Skip,
		),
	)

}
