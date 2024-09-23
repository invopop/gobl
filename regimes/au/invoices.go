package au

import (
	"github.com/invopop/gobl/bill"
	// "github.com/invopop/gobl/num"
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
		validation.Field(&inv.Supplier,
			validation.By(v.supplier),
			validation.Skip,
		),
		// Customer ID necessary when total over AUD 1000
		// validation.Field(&inv.Customer,
		// 	validation.When(
		// 		inv.Totals.Total.Compare(num.MakeAmount(1000, 0)) == 1,
		// 		validation.Required,
		// 	).Else(
		// 		validation.Skip,
		// 	),
		// 	validation.Skip,
		// ),
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
			validation.Skip,
		),
	)
}
