package saft

import (
	"errors"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validatePaymentLine(pl *bill.PaymentLine) error {
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
