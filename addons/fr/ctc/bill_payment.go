package ctc

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// paymentIsB2C reports whether the payment reports a B2C settlement,
// determined by the absence of a Customer party. Payments themselves
// are not routed by residency — every payment runs the Flow 10
// e-reporting ruleset regardless of where the parties are based — so
// "B2C" here is only used to mean "no customer present".
func paymentIsB2C(pmt *bill.Payment) bool {
	return pmt != nil && pmt.Customer == nil
}

// paymentHasCustomerAny is the "has Customer party" predicate used to
// gate the per-line invoice-reference rules. Despite the historical
// "B2B" labelling, it does not imply cross-border: a domestic FR-FR
// payment receipt has a customer and goes through this branch.
func paymentHasCustomerAny(v any) bool {
	pmt, ok := v.(*bill.Payment)
	return ok && !paymentIsB2C(pmt)
}

func billPaymentRules() *rules.Set {
	return rules.For(new(bill.Payment),
		// Flow 10 only reports payment receipts, not requests or advices.
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
		rules.Assert("03", "every VAT line rate must be one of the Flow 10 permitted percentages (G1.24): 0, 0.9, 1.05, 1.75, 2.1, 5.5, 7, 8.5, 9.2, 9.6, 10, 13, 19.6, 20, 20.6",
			is.Func("allowed Flow 10 VAT rates", paymentVATRatesAllowed),
		),
		rules.Field("supplier",
			rules.Assert("04", "supplier is required",
				is.Present,
			),
			rules.Assert("05", "supplier must have a SIREN identity (ISO/IEC 6523 scheme 0002)",
				is.Func("party has SIREN", partyHasSIREN),
			),
		),
		// Per-line invoice references are required when the payment
		// carries a Customer (cleared invoice receipts), not B2C settlements.
		rules.When(
			is.Func("payment has customer", paymentHasCustomerAny),
			rules.Field("lines",
				rules.Each(
					rules.Field("document",
						rules.Assert("06", "each payment line must reference a document (invoice) when a customer is present",
							is.Present,
						),
						rules.Field("code",
							rules.Assert("07", "payment line document code (invoice ID) is required when a customer is present",
								is.Present,
							),
						),
						rules.Field("issue_date",
							rules.Assert("08", "payment line document issue_date (invoice issue date) is required when a customer is present",
								is.Present,
							),
						),
					),
				),
			),
		),
	)
}

// paymentVATRatesAllowed reports whether every VAT rate total on the
// payment's lines matches one of the G1.24 whitelist percentages.
func paymentVATRatesAllowed(v any) bool {
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
				if !percentageInList(*rate.Percent, allowedVATRates) {
					return false
				}
			}
		}
	}
	return true
}
