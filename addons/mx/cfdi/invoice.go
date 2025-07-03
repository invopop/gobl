package cfdi

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/mx"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func normalizeInvoice(inv *bill.Invoice) {
	normalizeInvoiceIssueDateAndTime(inv)
	if inv.Tags.HasTags(TagGlobal) {
		inv.Customer = nil
	}
}

func normalizeInvoiceIssueDateAndTime(inv *bill.Invoice) {
	// Overwrite the issue date and time to align with
	// CFDI requirements for the emission date, unless the
	// issue time is already set.
	if inv.IssueTime != nil && !inv.IssueTime.IsZero() {
		return
	}
	tz := inv.RegimeDef().TimeLocation()
	dn := cal.ThisSecondIn(tz)
	tn := dn.Time()
	inv.IssueDate = dn.Date()
	inv.IssueTime = &tn
}

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Tax,
			validation.By(validateInvoiceTax(inv.Tags, inv.Preceding)),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.By(validateInvoiceSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.When(
				inv.Tags.HasTags(TagGlobal),
				validation.Empty.Error("cannot be set with global tag"),
			).Else(
				validation.By(validateInvoiceCustomer),
			),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(
				validation.By(validateInvoiceLine(inv.Tags)),
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
		validation.Field(&inv.Payment,
			validation.When(
				inv.HasTags(TagGlobal),
				validation.Required,
				validation.By(validateInvoicePaymentDetails(inv.Tags)),
			),
			validation.Skip,
		),
	)
}

func validateInvoicePaymentDetails(tags tax.Tags) validation.RuleFunc {
	return func(value any) error {
		payment, _ := value.(*bill.PaymentDetails)
		if payment == nil {
			return nil
		}
		return validation.ValidateStruct(payment,
			validation.Field(&payment.Advances,
				validation.When(
					tags.HasTags(TagGlobal),
					validation.Required.Error("must be set with global tag"),
				),
				validation.Skip,
			),
		)
	}
}

func validateInvoiceTax(tags tax.Tags, preceding []*org.DocumentRef) validation.RuleFunc {
	return func(value any) error {
		obj, _ := value.(*bill.Tax)
		if obj == nil {
			return nil
		}
		return validation.ValidateStruct(obj,
			validation.Field(&obj.Ext,
				tax.ExtensionsRequire(
					ExtKeyDocType,
					ExtKeyIssuePlace,
					ExtKeyPaymentMethod,
				),
				validation.When(
					tags.HasTags(TagGlobal),
					tax.ExtensionsRequire(
						ExtKeyGlobalPeriod,
						ExtKeyGlobalMonth,
						ExtKeyGlobalYear,
					),
				).Else(
					tax.ExtensionsRequireAllOrNone(
						ExtKeyGlobalPeriod,
						ExtKeyGlobalMonth,
						ExtKeyGlobalYear,
					),
				),
				validation.When(
					len(preceding) > 0,
					tax.ExtensionsRequire(
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
				tax.ExtensionsRequire(
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
			tax.ExtensionsRequire(
				ExtKeyFiscalRegime,
			),
		),
	)
}

func validateInvoiceLine(tags tax.Tags) validation.RuleFunc {
	return func(value any) error {
		line, _ := value.(*bill.Line)
		if line == nil {
			return nil
		}
		return validation.ValidateStruct(line,
			validation.Field(&line.Quantity, num.Positive),
			validation.Field(&line.Item, validation.By(validateInvoiceLineItem(tags))),
			validation.Field(&line.Total, num.Min(num.AmountZero)),
		)
	}
}

func validateInvoiceLineItem(tags tax.Tags) validation.RuleFunc {
	return func(value any) error {
		item, _ := value.(*org.Item)
		if item == nil {
			return nil
		}
		return validation.ValidateStruct(item,
			validation.Field(&item.Price, num.Positive),
			validation.Field(&item.Ref,
				validation.When(
					tags.HasTags(TagGlobal),
					validation.Required.Error("must be set with global tag"),
				),
				validation.Skip,
			),
			validation.Field(&item.Ext,
				// When Global Tag is applied, the prod-serv code is
				// is overridden during the conversion process so is not
				// required here.
				validation.When(
					!tags.HasTags(TagGlobal),
					tax.ExtensionsRequire(ExtKeyProdServ),
					validation.By(validItemExtensions),
				),
				validation.Skip,
			),
		)
	}
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
