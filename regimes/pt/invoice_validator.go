package pt

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/common"
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
		validation.Field(&inv.Supplier, validation.Required),
		validation.Field(&inv.Customer, validation.When(
			!inv.Tax.ContainsTag(common.TagSimplified),
			validation.Required,
		)),
		validation.Field(&inv.Lines,
			validation.Each(
				validation.By(v.validLine),
				validation.Skip, // Prevents each line's `ValidateWithContext` function from being called again.
			),
			validation.Skip, // Prevents each line's `ValidateWithContext` function from being called again.
		),
	)
}

func (v *invoiceValidator) validLine(value interface{}) error {
	line, _ := value.(*bill.Line)
	if line == nil {
		return nil
	}

	return validation.ValidateStruct(line,
		validation.Field(&line.Quantity, num.Min(num.MakeAmount(0, 0))),
		validation.Field(&line.Item, validation.By(v.validItem)),
	)
}

func (v *invoiceValidator) validItem(value interface{}) error {
	item, _ := value.(*org.Item)
	if item == nil {
		return nil
	}

	return validation.ValidateStruct(item,
		validation.Field(&item.Price, num.Min(num.MakeAmount(0, 0))),
	)
}
