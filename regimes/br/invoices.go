package br

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/validation"
)

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.Required,
			validation.Skip,
		),
	)
}
