package sg

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.When(
				inv.HasTags(TagInvoiceReceipt),
				validation.By(validateReceiptSupplier),
				validation.Skip,
			).Else(
				validation.By(validateInvoiceSupplier),
				validation.Skip,
			),
		),
		validation.Field(&inv.Customer,
			validation.When(
				inv.HasTags(TagInvoiceReceipt) || inv.HasTags(tax.TagSimplified),
				validation.Skip,
			).Else(validation.Required),
		),
	)
}

func validateInvoiceSupplier(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
		validation.Field(&p.Name,
			validation.Required,
		),
		validation.Field(&p.Addresses,
			validation.Required,
		),
	)
}

func validateReceiptSupplier(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
		validation.Field(&p.Name,
			validation.Required,
		),
	)
}
