package au

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

type invoiceValidator struct {
	inv *bill.Invoice
}

func validateInvoice(inv *bill.Invoice) error {
	v := &invoiceValidator{inv: inv}
	return v.validate()
}

func (v *invoiceValidator) validate() error {
	inv := v.inv
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Type,
			validation.In(
				bill.InvoiceTypeStandard,
				bill.InvoiceTypeCreditNote,
				bill.InvoiceTypeDebitNote,
			),
			validation.Skip,
		),
		validation.Field(&inv.IssueDate,
			validation.Required,
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.By(v.supplier),
			validation.Skip,
		),
		// Customer ID necessary when total over AUD 1000
		validation.Field(&inv.Customer,
			validation.When(
				inv.Totals.Total.Compare(num.MakeAmount(1000, 0)) == 1,
				validation.By(v.customer),
			),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) supplier(val any) error {
	obj, _ := val.(*org.Party)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.By(v.checkSupplierCountry),
			validation.Skip,
		),
		validation.Field(&obj.Name,
			validation.Required,
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) customer(val any) error {
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

// Supplier must have ABN, therefore be Australian
func (v *invoiceValidator) checkSupplierCountry(value interface{}) error {
	obj, _ := value.(*tax.Identity)
	if obj == nil {
		return nil
	}

	return validation.ValidateStruct(obj,
		validation.Field(&obj.Country,
			validation.In(l10n.AU.Tax()),
			validation.Skip,
		),
	)
}
