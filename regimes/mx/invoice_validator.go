package mx

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
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
		validation.Field(&inv.Customer,
			validation.Required,
			validation.By(v.validCustomer),
		),
		validation.Field(&inv.Lines,
			validation.Each(
				validation.By(v.validLine),
				validation.Skip, // Prevents each line's `ValidateWithContext` function from being called again.
			),
			validation.Skip, // Prevents each line's `ValidateWithContext` function from being called again.
		),
	)
}

func (v *invoiceValidator) validCustomer(value interface{}) error {
	obj, _ := value.(*org.Party)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID, validation.Required),
	)
}

func (v *invoiceValidator) validLine(value interface{}) error {
	line, _ := value.(*bill.Line)
	if line == nil {
		return nil
	}

	return validation.ValidateStruct(line,
		validation.Field(&line.Quantity, validation.By(v.validatePositiveAmount)),
		validation.Field(&line.Total, validation.By(v.validatePositiveAmount)),
		validation.Field(&line.Taxes,
			validation.Required,
			validation.Skip, // Prevents each tax's `ValidateWithContext` function from being called again.
		),
	)
}

func (v *invoiceValidator) validatePositiveAmount(value interface{}) error {
	amount, ok := value.(num.Amount)
	if !ok {
		return nil
	}

	if amount.Compare(num.MakeAmount(0, 0)) != 1 {
		return validation.NewError("validation_positive_amount", "must be positive")
	}
	return nil
}
