package nz

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var (
	thresholdMid  = num.MakeAmount(200, 0)
	thresholdHigh = num.MakeAmount(1000, 0)
)

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.When(
				invoiceTotalExceeds(inv, thresholdMid),
				validation.By(validateSupplierTaxID),
			),
			validation.When(
				inv.HasTags(TagSecondHandGoods),
				validation.By(validatePartyAddress),
			),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.When(
				invoiceTotalExceeds(inv, thresholdHigh),
				validation.Required,
				validation.By(validateCustomerDetails),
			),
			validation.When(
				inv.HasTags(tax.TagExport),
				validation.Required,
				validation.By(validatePartyAddress),
			),
			validation.Skip,
		),
	)
}

func invoiceTotalExceeds(inv *bill.Invoice, threshold num.Amount) bool {
	if inv.Totals == nil {
		return false
	}
	return inv.Totals.TotalWithTax.Compare(threshold) > 0
}

func validateSupplierTaxID(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
	)
}

func validateCustomerDetails(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Addresses,
			validation.When(
				!hasCustomerIdentifier(p),
				validation.Required,
			),
			validation.Skip,
		),
	)
}

func validatePartyAddress(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Addresses, validation.Required),
	)
}

func hasCustomerIdentifier(p *org.Party) bool {
	return len(p.Addresses) > 0 ||
		len(p.Telephones) > 0 ||
		len(p.Emails) > 0 ||
		len(p.Websites) > 0 ||
		len(p.Identities) > 0
}
