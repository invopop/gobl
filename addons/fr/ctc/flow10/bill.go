package flow10

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// vatKeyToUNTDIDCategory maps each supported GOBL VAT rate key to its
// UNTDID 5305 category code. The Canary Islands (IGIC / "L") and
// Ceuta/Melilla (IPSI / "M") categories are intentionally absent since
// they are not applicable to Flow 10.
var vatKeyToUNTDIDCategory = map[cbc.Key]cbc.Code{
	tax.KeyStandard:       "S",
	tax.KeyZero:           "Z",
	tax.KeyExempt:         "E",
	tax.KeyReverseCharge:  "AE",
	tax.KeyIntraCommunity: "K",
	tax.KeyExport:         "G",
	tax.KeyOutsideScope:   "O",
}

func invoiceIsB2C(inv *bill.Invoice) bool {
	return inv != nil && inv.Tags.HasTags(TagB2C)
}

func paymentIsB2C(pmt *bill.Payment) bool {
	return pmt != nil && pmt.Tags.HasTags(TagB2C)
}

func normalizeInvoice(inv *bill.Invoice) {
	if inv == nil {
		return
	}
	normalizeInvoiceTaxCategories(inv)
	if invoiceIsB2C(inv) {
		normalizeB2CCategoryOnInvoice(inv)
		return
	}
	normalizeParty(inv.Supplier)
	normalizeParty(inv.Customer)
	normalizeInvoiceBillingMode(inv)
}

// normalizeB2CCategoryOnInvoice defaults the B2C transaction category to
// TNT1 (not subject to French VAT) when the caller has not supplied one.
// TNT1 is the safest default: it covers B2C sales that would otherwise
// require explicit per-case classification (intra-EU distance sales,
// out-of-scope, etc.), and a user wanting a narrower code must set it
// explicitly.
func normalizeB2CCategoryOnInvoice(inv *bill.Invoice) {
	if inv.Tax != nil && inv.Tax.Ext.Get(ExtKeyB2CCategory) != "" {
		return
	}
	if inv.Tax == nil {
		inv.Tax = &bill.Tax{}
	}
	inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyB2CCategory, B2CCategoryNotTaxable)
}

// normalizeInvoiceTaxCategories sets the UNTDID 5305 category extension
// on each VAT combo based on its rate key. Combos whose key we do not
// map (IGIC / IPSI, or unknown) are left untouched.
func normalizeInvoiceTaxCategories(inv *bill.Invoice) {
	for _, line := range inv.Lines {
		if line == nil {
			continue
		}
		for _, combo := range line.Taxes {
			if combo == nil || combo.Category != tax.CategoryVAT {
				continue
			}
			if code, ok := vatKeyToUNTDIDCategory[combo.Key]; ok {
				combo.Ext = combo.Ext.Set(untdid.ExtKeyTaxCategory, code)
			}
		}
	}
}

// normalizeInvoiceBillingMode picks a sensible default for the Flow 10
// billing-mode extension when the user has not supplied one. We default to
// the Mixed (M) prefix since it is the safest without line-level analysis:
//   - M2 when the invoice is already paid in full
//   - M1 otherwise
//
// The user can override by setting the extension explicitly.
func normalizeInvoiceBillingMode(inv *bill.Invoice) {
	if inv.Tax != nil && !inv.Tax.Ext.IsZero() && inv.Tax.Ext.Get(ExtKeyBillingMode) != "" {
		return
	}
	mode := BillingModeM1
	if inv.Totals != nil && inv.Totals.Paid() {
		mode = BillingModeM2
	}
	if inv.Tax == nil {
		inv.Tax = &bill.Tax{}
	}
	inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyBillingMode, mode)
}

func billPaymentRules() *rules.Set {
	return rules.For(new(bill.Payment),
		// Flow 10 only reports payment receipts, not requests or advices.
		rules.Field("type",
			rules.Assert("01", "payment type must be 'receipt' for Flow 10 reporting",
				is.In(bill.PaymentTypeReceipt),
			),
		),
		// Payment date and at least one line (needed to report the amount
		// per rate) apply to both B2B and B2C payments.
		rules.Field("value_date",
			rules.Assert("02", "payment value_date (settlement date) is required",
				is.Present,
			),
		),
		// VAT rates reported on payment lines are constrained to the same
		// G1.24 whitelist as invoices, applied to both B2B and B2C.
		rules.Assert("07", "every VAT line rate must be one of the Flow 10 permitted percentages (G1.24): 0, 0.9, 1.05, 1.75, 2.1, 5.5, 7, 8.5, 9.2, 9.6, 10, 13, 19.6, 20, 20.6",
			is.Func("allowed Flow 10 VAT rates", paymentVATRatesAllowed),
		),
		// Supplier SIREN identifies the French reporting party on the
		// payment. Required for both B2B and B2C.
		rules.Field("supplier",
			rules.Assert("08", "supplier is required",
				is.Present,
			),
			rules.Assert("09", "supplier must have a SIREN identity (ISO/IEC 6523 scheme 0002)",
				is.Func("party has SIREN", partyHasSIREN),
			),
		),
		// Only B2B payments must carry an invoice reference per line
		// (invoice ID and issue date) so they can be reconciled against
		// the settled invoice.
		rules.When(
			is.Func("B2B payment", paymentIsB2BAny),
			rules.Field("lines",
				rules.Each(
					rules.Field("document",
						rules.Assert("04", "each payment line must reference a document (invoice) on B2B payments",
							is.Present,
						),
						rules.Field("code",
							rules.Assert("05", "payment line document code (invoice ID) is required on B2B payments",
								is.Present,
							),
						),
						rules.Field("issue_date",
							rules.Assert("06", "payment line document issue_date (invoice issue date) is required on B2B payments",
								is.Present,
							),
						),
					),
				),
			),
		),
	)
}

func paymentIsB2BAny(v any) bool {
	pmt, ok := v.(*bill.Payment)
	return ok && !paymentIsB2C(pmt)
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
