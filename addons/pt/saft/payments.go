package saft

import (
	"errors"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Add-on custom tags
const (
	TagVATCash cbc.Key = "vat-cash"
)

func validatePayment(pmt *bill.Payment) error {
	pt := paymentType(pmt)

	return validation.ValidateStruct(pmt,
		validation.Field(&pmt.Series,
			validateSeriesFormat(pt),
			validation.Skip,
		),
		validation.Field(&pmt.Code,
			validateCodeFormat(pmt.Series, pt),
			validation.Skip,
		),
		validation.Field(&pmt.Ext,
			tax.ExtensionsRequire(ExtKeyPaymentType),
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
		validation.Field(&pmt.Lines,
			validation.Each(
				validation.By(validatePaymentLine),
				validation.Skip,
			),
			validation.Skip,
		),
	)
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

func validatePaymentLine(val any) error {
	pl, _ := val.(*bill.PaymentLine)
	if pl == nil {
		return nil
	}

	return validation.ValidateStruct(pl,
		validation.Field(&pl.Document,
			validation.By(validateLineDocument),
			validation.Required,
			validation.Skip,
		),
		validation.Field(&pl.Tax,
			validation.By(validateLineTax),
			validation.Required,
			validation.Skip,
		),
		validation.Field(&pl.Debit, num.Min(num.AmountZero)),
		validation.Field(&pl.Credit, num.Min(num.AmountZero)),
	)
}

func validateLineDocument(val any) error {
	ld, _ := val.(*org.DocumentRef)
	if ld == nil {
		return nil
	}

	return validation.ValidateStruct(ld,
		validation.Field(&ld.IssueDate,
			validation.Required,
			validation.Skip,
		),
	)
}

func validateLineTax(val any) error {
	lt, _ := val.(*tax.Total)
	if lt == nil {
		return nil
	}

	c := lt.Category(tax.CategoryVAT)
	if c == nil {
		return errors.New("missing category VAT")
	}

	return validation.ValidateStruct(c,
		validation.Field(&c.Rates,
			validation.Each(
				validation.By(validateLineTaxRate),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}

func validateLineTaxRate(val any) error {
	r, _ := val.(*tax.RateTotal)
	if r == nil {
		return nil
	}

	return validation.ValidateStruct(r,
		validation.Field(&r.Ext,
			tax.ExtensionsRequire(ExtKeyTaxRate, pt.ExtKeyRegion),
			validation.Skip,
		),
	)
}

func paymentType(pmt *bill.Payment) cbc.Code {
	if pmt == nil || pmt.Ext == nil {
		return cbc.CodeEmpty
	}

	return pmt.Ext[ExtKeyPaymentType]
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
}
