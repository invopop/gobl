package au

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var customerABNThreshold = num.MakeAmount(100000, 2)

func validateBillInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(validateBillInvoiceSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.When(
				requiresCustomerABN(inv),
				validation.Required,
			).Else(
				validation.Skip,
			),
			validation.By(validateBillInvoiceCustomer(inv)),
			validation.Skip,
		),
	)
}

func validateBillInvoiceSupplier(value any) error {
	party, ok := value.(*org.Party)
	if !ok || party == nil {
		return nil
	}
	return validation.ValidateStruct(party,
		validation.Field(&party.Name,
			validation.Required,
			validation.Skip,
		),
		validation.Field(&party.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
	)
}

func validateBillInvoiceCustomer(inv *bill.Invoice) validation.RuleFunc {
	return func(value any) error {
		party, ok := value.(*org.Party)
		if !ok || party == nil || !requiresCustomerABN(inv) {
			return nil
		}
		return validation.ValidateStruct(party,
			validation.Field(&party.TaxID,
				validation.Required,
				tax.RequireIdentityCode,
				validation.Skip,
			),
		)
	}
}

func requiresCustomerABN(inv *bill.Invoice) bool {
	return inv != nil &&
		inv.Totals != nil &&
		inv.Totals.TotalWithTax.Compare(customerABNThreshold) >= 0
}
