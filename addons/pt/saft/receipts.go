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

func validateReceipt(rct *bill.Receipt) error {
	rt := receiptType(rct)

	return validation.ValidateStruct(rct,
		validation.Field(&rct.Series,
			validateSeriesFormat(rt),
			validation.Skip,
		),
		validation.Field(&rct.Code,
			validateCodeFormat(rct.Series, rt),
			validation.Skip,
		),
		validation.Field(&rct.Ext,
			tax.ExtensionsRequire(ExtKeyReceiptType),
			validation.Skip,
		),
		validation.Field(&rct.Supplier,
			validation.By(validateSupplier),
			validation.Skip,
		),
		validation.Field(&rct.Customer,
			validation.By(validateCustomer),
			validation.Skip,
		),
		validation.Field(&rct.Lines,
			validation.Each(
				validation.By(validateReceiptLine),
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

func validateReceiptLine(val any) error {
	rl, _ := val.(*bill.ReceiptLine)
	if rl == nil {
		return nil
	}

	return validation.ValidateStruct(rl,
		validation.Field(&rl.Document,
			validation.Required,
			validation.By(validateLineDocument),
			validation.Skip,
		),
		validation.Field(&rl.Tax,
			validation.Required,
			validation.By(validateLineTax),
			validation.Skip,
		),
		validation.Field(&rl.Debit, num.Min(num.AmountZero)),
		validation.Field(&rl.Credit, num.Min(num.AmountZero)),
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

func receiptType(rct *bill.Receipt) cbc.Code {
	if rct == nil || rct.Ext == nil {
		return cbc.CodeEmpty
	}

	return rct.Ext[ExtKeyReceiptType]
}

func normalizeReceipt(rct *bill.Receipt) {
	if rct.Ext == nil {
		rct.Ext = tax.Extensions{}
	}

	// TODO: This could be done with scenarios when supported
	if rct.HasTags(TagVATCash) {
		rct.Ext[ExtKeyReceiptType] = ReceiptTypeCash
	} else {
		rct.Ext[ExtKeyReceiptType] = ReceiptTypeOther
	}
}
