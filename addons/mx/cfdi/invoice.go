package cfdi

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/mx/sat"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func normalizeInvoice(inv *bill.Invoice) {
	normalizeParty(inv.Supplier)
	normalizeParty(inv.Customer)
	for _, line := range inv.Lines {
		normalizeItem(line.Item)
	}

	// 2024-04-26: copy suppliers post code to invoice, if not already
	// set.
	if inv.Tax == nil {
		inv.Tax = new(bill.Tax)
	}
	if inv.Tax.Ext == nil {
		inv.Tax.Ext = make(tax.Extensions)
	}
	if inv.Tax.Ext.Has(ExtKeyIssuePlace) {
		return
	}
	if inv.Supplier.Ext.Has(ExtKeyPostCode) {
		inv.Tax.Ext[ExtKeyIssuePlace] = inv.Supplier.Ext[ExtKeyPostCode]
		return
	}
	if len(inv.Supplier.Addresses) > 0 {
		addr := inv.Supplier.Addresses[0]
		if addr.Code != "" {
			inv.Tax.Ext[ExtKeyIssuePlace] = tax.ExtValue(addr.Code)
		}
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

func validateInvoiceTax(preceding []*bill.Preceding) validation.RuleFunc {
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
				obj.TaxID != nil && obj.TaxID.Country.In("MX"),
				tax.ExtensionsRequires(
					ExtKeyPostCode,
					ExtKeyFiscalRegime,
					ExtKeyUse,
				),
			),
		),
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
	entry, _ := value.(*bill.Preceding)
	if entry == nil {
		return nil
	}
	return validation.ValidateStruct(entry,
		validation.Field(
			&entry.Stamps,
			head.StampsHas(sat.StampUUID),
			validation.Skip,
		),
	)
}
