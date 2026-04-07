package au

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var customerIdentificationThreshold = num.MakeAmount(100000, 2)

func validateBillInvoice(inv *bill.Invoice) error {
	if inv == nil {
		return nil
	}
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(validateBillInvoiceSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.When(
				requiresCustomerIdentification(inv),
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
			validation.By(validateBillInvoiceSupplierTaxID),
			validation.Skip,
		),
	)
}

func validateBillInvoiceSupplierTaxID(value any) error {
	tID := value.(*tax.Identity)
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Country,
			validation.In(l10n.AU.Tax()),
			validation.Skip,
		),
	)
}

func validateBillInvoiceCustomer(inv *bill.Invoice) validation.RuleFunc {
	return func(value any) error {
		party, ok := value.(*org.Party)
		if !ok || party == nil || !requiresCustomerIdentification(inv) {
			return nil
		}
		return validation.ValidateStruct(party,
			validation.Field(&party.Name,
				validation.Required,
				validation.Skip,
			),
		)
	}
}

func requiresCustomerIdentification(inv *bill.Invoice) bool {
	if inv.HasTags(tax.TagSelfBilled) {
		return true
	}
	return inv.Totals != nil &&
		inv.Totals.TotalWithTax.Compare(customerIdentificationThreshold) >= 0
}
