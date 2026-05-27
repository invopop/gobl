package flow10

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// paymentIsB2CDoc reports whether the payment reports a B2C settlement,
// determined by the absence of a Customer party.
func paymentIsB2CDoc(pmt *bill.Payment) bool {
	return pmt != nil && pmt.Customer == nil
}

// paymentIsB2B returns a Test that passes when the payment has a
// customer party (B2B settlement).
func paymentIsB2B() rules.Test {
	return is.Func("payment is B2B (has customer)", func(v any) bool {
		pmt, ok := v.(*bill.Payment)
		return ok && !paymentIsB2CDoc(pmt)
	})
}

func billPaymentRules() *rules.Set {
	return rules.For(new(bill.Payment),
		rules.Field("type",
			rules.Assert("01", "payment type must be 'receipt' for Flow 10 reporting",
				is.In(bill.PaymentTypeReceipt),
			),
		),
		rules.Field("value_date",
			rules.Assert("02", "payment value_date (settlement date) is required",
				is.Present,
			),
		),
		rules.Assert("03", "payment VAT line percent must be one of the Flow 10 permitted values (G1.24)",
			is.Func("allowed Flow 10 VAT percents", paymentVATPercentsAllowed),
		),
		rules.Field("supplier",
			rules.Assert("04", "payment supplier is required",
				is.Present,
			),
			rules.Assert("05", "payment supplier must have a SIREN identity (ISO/IEC 6523 scheme 0002)",
				is.Func("party has SIREN", partyHasSIREN),
			),
		),
		rules.When(
			paymentIsB2B(),
			rules.Field("lines",
				rules.Each(
					rules.Field("document",
						rules.Assert("06", "payment line document is required when a customer is present",
							is.Present,
						),
						rules.Field("code",
							rules.Assert("07", "payment line document code is required when a customer is present",
								is.Present,
							),
						),
						rules.Field("issue_date",
							rules.Assert("08", "payment line document issue_date is required when a customer is present",
								is.Present,
							),
						),
					),
				),
			),
		),
	)
}

func paymentVATPercentsAllowed(v any) bool {
	pmt, ok := v.(*bill.Payment)
	if !ok || pmt == nil {
		return true
	}
	for _, line := range pmt.Lines {
		if line == nil || line.Tax == nil {
			continue
		}
		for _, cat := range line.Tax.Categories {
			if cat == nil || cat.Code != tax.CategoryVAT {
				continue
			}
			for _, rate := range cat.Rates {
				if rate == nil || rate.Percent == nil {
					continue
				}
				if !percentageInList(*rate.Percent, allowedVATPercents) {
					return false
				}
			}
		}
	}
	return true
}
