package pt

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
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
		validation.Field(&inv.Currency, validation.In(currency.EUR)),
		validation.Field(&inv.Supplier,
			validation.By(v.supplier),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(
				validation.By(v.line),
				validation.Skip, // Prevents each line's `ValidateWithContext` function from being called again.
			),
			validation.Skip, // Prevents each line's `ValidateWithContext` function from being called again.
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
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) line(value interface{}) error {
	line, _ := value.(*bill.Line)
	if line == nil {
		return nil
	}

	return validation.ValidateStruct(line,
		validation.Field(&line.Quantity,
			num.Min(num.MakeAmount(0, 0)),
		),
		validation.Field(&line.Item,
			validation.By(v.item),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) item(value interface{}) error {
	item, _ := value.(*org.Item)
	if item == nil {
		return nil
	}
	return validation.ValidateStruct(item,
		validation.Field(&item.Price,
			num.Min(num.MakeAmount(0, 0)),
		),
	)
}
