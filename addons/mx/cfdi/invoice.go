package cfdi

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/mx/sat"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

type invoiceValidator struct {
	inv *bill.Invoice
}

func validateInvoice(inv *bill.Invoice) error {
	v := &invoiceValidator{inv}
	return v.validate()
}

func (v *invoiceValidator) validate() error {
	inv := v.inv
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Tax,
			validation.By(v.tax),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.By(v.supplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.By(v.customer),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(
				validation.By(v.line),
				validation.Skip,
			),
			validation.Skip,
		),
		validation.Field(&inv.Payment,
			validation.By(v.payment),
			validation.Skip,
		),
		validation.Field(&inv.Preceding,
			validation.Each(validation.By(v.precedingEntry)),
			validation.Skip,
		),
		validation.Field(&inv.Discounts,
			validation.Empty.Error("the SAT doesn't allow discounts at invoice level. Use line discounts instead."),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) tax(value any) error {
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
				len(v.inv.Preceding) > 0,
				tax.ExtensionsRequires(
					ExtKeyRelType,
				),
			),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) customer(value interface{}) error {
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

func (v *invoiceValidator) supplier(value interface{}) error {
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

func (v *invoiceValidator) line(value interface{}) error {
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

func (v *invoiceValidator) payment(value interface{}) error {
	pay, _ := value.(*bill.Payment)
	if pay == nil {
		return nil
	}

	return validation.ValidateStruct(pay,
		validation.Field(&pay.Instructions, validation.By(v.payInstructions)),
		validation.Field(&pay.Advances, validation.Each(validation.By(v.payAdvance))),
		validation.Field(&pay.Terms, validation.By(v.payTerms)),
	)
}

func (v *invoiceValidator) payInstructions(value interface{}) error {
	instr, _ := value.(*pay.Instructions)
	if instr == nil {
		return nil
	}

	return validation.ValidateStruct(instr,
		validation.Field(&instr.Ext,
			tax.ExtensionsRequires(ExtKeyPaymentMeans),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) payAdvance(value interface{}) error {
	adv, _ := value.(*pay.Advance)
	if adv == nil {
		return nil
	}

	return validation.ValidateStruct(adv,
		validation.Field(&adv.Ext,
			tax.ExtensionsRequires(ExtKeyPaymentMeans),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) payTerms(value interface{}) error {
	terms, _ := value.(*pay.Terms)
	if terms == nil {
		return nil
	}

	return validation.ValidateStruct(terms,
		validation.Field(&terms.Notes, validation.Length(0, 1000)),
	)
}

func (v *invoiceValidator) precedingEntry(value interface{}) error {
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

func normalizeInvoice(inv *bill.Invoice) {
	normalizeParty(inv.Supplier)
	normalizeParty(inv.Customer)
	normalizeInvoiceLines(inv)
	normalizeInvoicePaymentInstructions(inv)
	normalizeInvoicePaymentAdvances(inv)

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

func normalizeInvoiceLines(inv *bill.Invoice) {
	for _, line := range inv.Lines {
		normalizeItem(line.Item)
		for _, combo := range line.Taxes {
			var k tax.ExtValue
			switch combo.Category {
			case sat.TaxCategoryISR:
				k = "001"
			case tax.CategoryVAT, sat.TaxCategoryRVAT:
				k = "002"
			case sat.TaxCategoryIEPS, sat.TaxCategoryRIEPS:
				k = "003"
			}
			if combo.Ext == nil {
				combo.Ext = make(tax.Extensions)
			}
			combo.Ext[ExtKeyTaxType] = k
		}
	}
}
