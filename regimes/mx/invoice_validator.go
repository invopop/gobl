package mx

import (
	"github.com/invopop/gobl/bill"
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
		// validation.Field(&inv.Lines,
		// 	validation.Each(validation.By(v.validLine)),
		// ),
	)
}

// func (v *invoiceValidator) validLine(value interface{}) error {
// 	line, _ := value.(*bill.Line)
// 	if line == nil {
// 		return nil
// 	}
// 	return validation.ValidateStruct(line,
// 		validation.Field(&line.Quantity,
// 			validation.By(validLineQuantity),
// 		),
// 	)
// }

// func validLineQuantity(value interface{}) error {
// 	quantity, ok := value.(num.Amount)
// 	if !ok {
// 		return nil
// 	}

// 	if quantity.Compare(num.MakeAmount(0, 0)) != 1 {
// 		return validation.NewError("quantity", "must be positive")
// 	}
// 	return nil
// }

func (v *invoiceValidator) validCustomer(value interface{}) error {
	obj, _ := value.(*org.Party)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID, validation.Required),
	)
}
