package zatca

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func normalizeInvoice(inv *bill.Invoice) {
	if inv == nil {
		return
	}

	// Ensure Tax object exists
	if inv.Tax == nil {
		inv.Tax = &bill.Tax{}
	}

	// Always set rounding to currency for SA ZATCA
	inv.Tax.Rounding = tax.RoundingRuleCurrency

	// Ensure issue date exists
	if inv.IssueTime == nil {
		inv.IssueTime = &cal.Time{}
	}

	// BR-KSA-O-01: "Not subject to VAT" lines must have a 0% rate.
	for _, line := range inv.Lines {
		vat := line.Taxes.Get(tax.CategoryVAT)
		if vat == nil {
			continue
		}
		if vat.Key == tax.KeyOutsideScope {
			vat.Percent = &num.PercentageZero
		}
	}
}

func billDiscountRules() *rules.Set {
	return rules.For(new(bill.Discount),
		rules.Field("taxes",
			rules.Assert("01", "taxes are required (BR-32)", is.Present),
		),
	)
}
