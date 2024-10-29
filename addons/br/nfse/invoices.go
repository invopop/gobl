package nfse

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/validation"
)

func validateInvoice(inv *bill.Invoice) error {
	if inv == nil {
		return nil
	}

	return validation.ValidateStruct(inv,
		validation.Field(&inv.Charges,
			validation.Empty.Error("not supported by nfse"),
			validation.Skip,
		),
		validation.Field(&inv.Discounts,
			validation.Empty.Error("not supported by nfse"),
			validation.Skip,
		),
	)
}
