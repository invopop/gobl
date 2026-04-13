package saft

import (
	"errors"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func billPaymentLineRules() *rules.Set {
	return rules.For(new(bill.PaymentLine),
		rules.Field("document",
			rules.Assert("01", "cannot be blank", is.Present),
			rules.Field("issue_date",
				rules.Assert("02", "cannot be blank", is.Present),
			),
		),
		rules.Field("tax",
			rules.Assert("03", "cannot be blank", is.Present),
			rules.Assert("04", "missing category VAT",
				is.FuncError("has VAT category", paymentLineTaxHasVAT),
			),
			rules.Field("categories",
				rules.Each(
					rules.Field("rates",
						rules.Assert("05", "only one rate allowed per line", is.Length(0, 1)),
					),
				),
			),
		),
		rules.Assert("06", "exemption notes invalid",
			is.FuncError("exemption notes", paymentLineExemptionNotesValid),
		),
	)
}

// paymentLineTaxHasVAT checks that the payment line's tax total has a VAT category.
func paymentLineTaxHasVAT(val any) error {
	lt, ok := val.(*tax.Total)
	if !ok || lt == nil {
		return nil
	}
	if lt.Category(tax.CategoryVAT) == nil {
		return errors.New("missing category VAT")
	}
	return nil
}

// paymentLineExemptionNotesValid validates exemption notes for a payment line.
func paymentLineExemptionNotesValid(val any) error {
	pl, ok := val.(*bill.PaymentLine)
	if !ok || pl == nil {
		return nil
	}
	return validateExemptionNotes(pl.Notes, paymentLineTaxExemptionCode(pl))
}

func paymentLineTaxExemptionCode(pl *bill.PaymentLine) cbc.Code {
	if pl.Tax == nil {
		return ""
	}

	vat := pl.Tax.Category(tax.CategoryVAT)
	if vat == nil || len(vat.Rates) == 0 {
		return ""
	}

	// Since there's a validation that only allows one rate per line,
	// we can safely check the first rate
	rate := vat.Rates[0]
	if rate == nil {
		return ""
	}

	return rate.Ext.Get(ExtKeyExemption)
}
