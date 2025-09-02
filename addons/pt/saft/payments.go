package saft

import (
	"errors"

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
		validation.Field(&pmt.Lines,
			validation.Each(
				validation.By(validatePaymentLine),
				validation.Skip,
			),
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

func validatePaymentLine(val any) error {
	pl, _ := val.(*bill.PaymentLine)
	if pl == nil {
		return nil
	}

	return validation.ValidateStruct(pl,
		validation.Field(&pl.Document,
			validation.By(validatePaymentLineDocument),
			validation.Required,
			validation.Skip,
		),
		validation.Field(&pl.Tax,
			validation.By(validatePaymentLineTax),
			validation.Required,
			validation.Skip,
		),
	)
}

func validatePaymentLineDocument(val any) error {
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

func validatePaymentLineTax(val any) error {
	lt, _ := val.(*tax.Total)
	if lt == nil {
		return nil
	}

	c := lt.Category(tax.CategoryVAT)
	if c == nil {
		return errors.New("missing category VAT")
	}

	return validation.ValidateStruct(lt,
		validation.Field(&lt.Categories,
			validation.Each(
				validation.By(validateLineTaxCategory),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}

func validateLineTaxCategory(val any) error {
	tc, _ := val.(*tax.CategoryTotal)
	if tc == nil {
		return nil
	}

	return validation.ValidateStruct(tc,
		validation.Field(&tc.Rates,
			// According to point 4.4.4.14.6. of Portaria nª 302/2016,
			// multiple tax rates (even for the same document) must be
			// reported broken down in different payment lines.
			validation.Length(0, 1).Error("only one rate allowed per line"),
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

	return validation.ValidateStruct(r, validateVATExt(&r.Ext))
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
