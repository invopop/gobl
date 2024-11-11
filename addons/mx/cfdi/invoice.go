package cfdi

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/mx"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func normalizeInvoice(inv *bill.Invoice) {
	normalizeParty(inv.Supplier)
	normalizeParty(inv.Customer)
	for _, line := range inv.Lines {
		normalizeItem(line.Item)
	}

}

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Tax,
			validation.By(validateInvoiceTax(inv.Preceding)),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.By(validateInvoiceSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.By(validateInvoiceCustomer),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(
				validation.By(validateInvoiceLine),
				validation.Skip,
			),
			validation.Skip,
		),
		validation.Field(&inv.Preceding,
			validation.Each(validation.By(validateInvoicePreceding)),
			validation.Skip,
		),
		validation.Field(&inv.Discounts,
			validation.Empty.Error("not supported, use line discounts instead"),
			validation.Skip,
		),
		validation.Field(&inv.Charges,
			validation.Empty.Error("not supported"),
			validation.Skip,
		),
	)
}

func validateInvoiceTax(preceding []*org.DocumentRef) validation.RuleFunc {
	return func(value any) error {
		obj, _ := value.(*bill.Tax)
		if obj == nil {
			return nil
		}
		return validation.ValidateStruct(obj,
			validation.Field(&obj.Ext,
				tax.ExtensionsRequires(
					ExtKeyDocType,
					ExtKeyIssuePlace,
				),
				validation.When(
					len(preceding) > 0,
					tax.ExtensionsRequires(
						ExtKeyRelType,
					),
				),
				validation.Skip,
			),
		)
	}
}

func validateInvoiceCustomer(value any) error {
	obj, _ := value.(*org.Party)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
		validation.Field(&obj.Ext,
			validation.When(
				isMexican(obj),
				tax.ExtensionsRequires(
					ExtKeyFiscalRegime,
					ExtKeyUse,
				),
			),
		),
		validation.Field(&obj.Addresses,
			validation.When(
				isMexican(obj),
				validation.Required,
				validation.Each(
					validation.By(validateMexicanCustomerAddress),
				),
			),
		),
	)
}

func validateMexicanCustomerAddress(value any) error {
	obj, _ := value.(*org.Address)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Code,
			validation.Required,
			validation.Match(PostCodeRegexp)),
	)
}

func validateInvoiceSupplier(value any) error {
	obj, _ := value.(*org.Party)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
		validation.Field(&obj.Ext,
			tax.ExtensionsRequires(
				ExtKeyFiscalRegime,
			),
		),
	)
}

func validateInvoiceLine(value any) error {
	line, _ := value.(*bill.Line)
	if line == nil {
		return nil
	}
	return validation.ValidateStruct(line,
		validation.Field(&line.Quantity, num.Positive),
		validation.Field(&line.Item, validation.By(validateItem)),
		validation.Field(&line.Total, num.Min(num.AmountZero)),
	)
}

func validateInvoicePreceding(value interface{}) error {
	entry, _ := value.(*org.DocumentRef)
	if entry == nil {
		return nil
	}
	return validation.ValidateStruct(entry,
		validation.Field(
			&entry.Stamps,
			head.StampsHas(mx.StampSATUUID),
			validation.Skip,
		),
	)
}

func isMexican(party *org.Party) bool {
	if party == nil || party.TaxID == nil {
		return false
	}
	return party.TaxID.Country.In("MX")
}
