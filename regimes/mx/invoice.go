package mx

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

type invoiceValidator struct {
	inv *bill.Invoice
	ss  *tax.ScenarioSummary
}

func validateInvoice(inv *bill.Invoice) error {
	ss := inv.ScenarioSummary()
	v := &invoiceValidator{inv, ss}
	return v.validate()
}

func (v *invoiceValidator) validate() error {
	inv := v.inv
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Currency, validation.In(currency.MXN)),
		validation.Field(&inv.Tax,
			validation.By(v.validTax),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.By(v.validSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.By(v.validCustomer),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(
				validation.By(v.validLine),
				validation.Skip, // Prevents each line's `ValidateWithContext` function from being called again.
			),
			validation.Skip, // Prevents each line's `ValidateWithContext` function from being called again.
		),
		validation.Field(&inv.Payment,
			validation.By(v.validPayment),
			validation.Skip,
		),
		validation.Field(&inv.Preceding,
			validation.By(v.validPrecedingList),
			validation.Each(validation.By(v.validPrecedingEntry)),
			validation.Skip,
		),
		validation.Field(&inv.Discounts,
			validation.Empty.Error("the SAT doesn't allow discounts at invoice level. Use line discounts instead."),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) validTax(value any) error {
	obj, _ := value.(*bill.Tax)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Ext,
			tax.ExtensionsRequires(
				ExtKeyCFDIIssuePlace,
			),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) validCustomer(value interface{}) error {
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
				ExtKeyCFDIPostCode,
				ExtKeyCFDIFiscalRegime,
				ExtKeyCFDIUse,
			),
		),
	)
}

func (v *invoiceValidator) validSupplier(value interface{}) error {
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
				ExtKeyCFDIFiscalRegime,
			),
		),
	)
}

func (v *invoiceValidator) validLine(value interface{}) error {
	line, _ := value.(*bill.Line)
	if line == nil {
		return nil
	}

	return validation.ValidateStruct(line,
		validation.Field(&line.Quantity, num.Positive),
		validation.Field(&line.Total, num.Positive),
	)
}

func (v *invoiceValidator) validPayment(value interface{}) error {
	pay, _ := value.(*bill.Payment)
	if pay == nil {
		return nil
	}

	return validation.ValidateStruct(pay,
		validation.Field(&pay.Instructions, validation.By(v.validatePayInstructions)),
		validation.Field(&pay.Advances, validation.Each(validation.By(v.validateAdvance))),
		validation.Field(&pay.Terms, validation.By(v.validatePayTerms)),
	)
}

func (v *invoiceValidator) validatePayInstructions(value interface{}) error {
	instr, _ := value.(*pay.Instructions)
	if instr == nil {
		return nil
	}

	return validation.ValidateStruct(instr,
		validation.Field(&instr.Key, isValidPaymentMeanKey),
	)
}

func (v *invoiceValidator) validateAdvance(value interface{}) error {
	adv, _ := value.(*pay.Advance)
	if adv == nil {
		return nil
	}

	fields := []*validation.FieldRules{
		validation.Field(&adv.Key, isValidPaymentMeanKey),
	}

	// Temporary hack necessary to help transition users from using the instructions key to use
	// the advance key. TODO: Expect the payment means key always to be present in every
	// advance (and not the instructions) once users have transitioned.
	if v.inv.Payment.Instructions == nil || v.inv.Payment.Instructions.Key == cbc.KeyEmpty {
		fields = append(fields, validation.Field(&adv.Key, validation.Required))
	}

	return validation.ValidateStruct(adv, fields...)
}

func (v *invoiceValidator) validatePayTerms(value interface{}) error {
	terms, _ := value.(*pay.Terms)
	if terms == nil {
		return nil
	}

	return validation.ValidateStruct(terms,
		validation.Field(&terms.Notes, validation.Length(0, 1000)),
	)
}

func (v *invoiceValidator) validPrecedingList(value interface{}) error {
	list, _ := value.([]*bill.Preceding)
	if len(list) == 0 {
		return nil
	}

	if v.ss.Codes[KeySATTipoRelacion] == "" {
		return fmt.Errorf("cannot be mapped from a `%s` type invoice", v.inv.Type)
	}

	return nil
}

func (v *invoiceValidator) validPrecedingEntry(value interface{}) error {
	entry, _ := value.(*bill.Preceding)
	if entry == nil {
		return nil
	}

	for _, s := range entry.Stamps {
		if s.Provider == StampSATUUID {
			return nil
		}
	}

	return fmt.Errorf("must have a `%s` stamp", StampSATUUID)
}

var isValidPaymentMeanKey = validation.In(validPaymentMeanKeys()...)

func validPaymentMeanKeys() []interface{} {
	keys := make([]interface{}, len(paymentMeansKeyDefinitions))
	for i, keyDef := range paymentMeansKeyDefinitions {
		keys[i] = keyDef.Key
	}

	return keys
}

func normalizeInvoice(inv *bill.Invoice) error {
	// 2024-04-26: copy suppliers post code to invoice, if not already
	// set.
	if inv.Tax == nil {
		inv.Tax = new(bill.Tax)
	}
	if inv.Tax.Ext == nil {
		inv.Tax.Ext = make(tax.Extensions)
	}
	if inv.Tax.Ext.Has(ExtKeyCFDIIssuePlace) {
		return nil
	}
	if inv.Supplier.Ext.Has(ExtKeyCFDIPostCode) {
		inv.Tax.Ext[ExtKeyCFDIIssuePlace] = inv.Supplier.Ext[ExtKeyCFDIPostCode]
		return nil
	}
	if len(inv.Supplier.Addresses) > 0 {
		addr := inv.Supplier.Addresses[0]
		if addr.Code != "" {
			inv.Tax.Ext[ExtKeyCFDIIssuePlace] = tax.ExtValue(addr.Code)
		}
	}
	return nil
}
