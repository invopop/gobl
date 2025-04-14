package facturx

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

func validateInvoice(inv *bill.Invoice) error {

	return validation.ValidateStruct(inv,
		validation.Field(&inv.Payment,
			validation.Required,
			validation.By(validatePayment),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.Required,
			validation.By(validateParty),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.Required,
			validation.By(validateParty),
			validation.Skip,
		),
	)
}
func validatePayment(val any) error {
	payment, ok := val.(*bill.PaymentDetails)
	if !ok || payment == nil {
		return nil
	}
	// Rest of the validation is handled by en16931 addon
	return validation.ValidateStruct(payment,
		validation.Field(&payment.Instructions,
			validation.Required,
			validation.Skip,
		),
	)
}

func validateParty(val any) error {
	party, ok := val.(*org.Party)
	if !ok || party == nil {
		return nil
	}

	return validation.ValidateStruct(party,
		validation.Field(&party.Addresses,
			validation.Required,
			validation.By(validateAddresses),
			validation.Skip,
		),
	)
}

func validateAddresses(val any) error {
	addresses, ok := val.([]*org.Address)
	if !ok || addresses == nil {
		return nil
	}

	// gobl.cii looks at the first address. Should we validate all?
	return validation.ValidateStruct(addresses[0],
		validation.Field(&addresses[0].Country,
			validation.Required,
			validation.Skip,
		),
	)
}
