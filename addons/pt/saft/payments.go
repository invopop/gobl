package saft

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Add-on custom tags
const (
	TagVATCash cbc.Key = "vat-cash"
)

func validatePayment(pmt *bill.Payment) error {
	dt := paymentDocType(pmt)

	return validation.ValidateStruct(pmt,
		validation.Field(&pmt.Series,
			validateSeriesFormat(dt),
			validation.Skip,
		),
		validation.Field(&pmt.Code,
			validateCodeFormat(pmt.Series, dt),
			validation.Skip,
		),
		validation.Field(&pmt.Ext,
			tax.ExtensionsRequire(ExtKeyPaymentType),
			tax.ExtensionsRequire(ExtKeySource),
			validation.When(
				pmt.Ext[ExtKeySource] != SourceBillingProduced,
				tax.ExtensionsRequire(ExtKeySourceRef),
			),
			validation.By(validateSourceRef(dt)),
			validation.Skip,
		),
		validation.Field(&pmt.Supplier,
			validation.By(validateSupplier),
			validation.Skip,
		),
		validation.Field(&pmt.Customer,
			validation.By(validateCustomer),
			validation.Skip,
		),
		validation.Field(&pmt.Total, num.ZeroOrPositive),
	)
}

func paymentDocType(pmt *bill.Payment) cbc.Code {
	if pmt.Ext == nil {
		return cbc.CodeEmpty
	}
	return pmt.Ext[ExtKeyPaymentType]
}

func validateSupplier(val any) error {
	sup, _ := val.(*org.Party)
	if sup == nil {
		return nil
	}

	return validation.ValidateStruct(sup,
		validation.Field(&sup.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
	)
}

func validateCustomer(val any) error {
	cus, _ := val.(*org.Party)
	if cus == nil {
		return nil
	}

	return validation.ValidateStruct(cus,
		validation.Field(&cus.Name,
			validation.When(cus.TaxID != nil && cus.TaxID.Code != cbc.CodeEmpty, validation.Required),
		),
	)
}

func normalizePayment(pmt *bill.Payment) {
	if pmt.Ext == nil {
		pmt.Ext = tax.Extensions{}
	}

	// TODO: This could be done with scenarios when supported
	if pmt.HasTags(TagVATCash) {
		pmt.Ext[ExtKeyPaymentType] = PaymentTypeCash
	} else {
		pmt.Ext[ExtKeyPaymentType] = PaymentTypeOther
	}

	if !pmt.Ext.Has(ExtKeySource) {
		pmt.Ext[ExtKeySource] = SourceBillingProduced
	}
}
