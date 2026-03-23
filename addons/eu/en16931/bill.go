package en16931

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

var discountKeyMap = tax.Extensions{
	bill.DiscountKeyEarlyCompletion:  "41",
	bill.DiscountKeyMilitary:         "62",
	bill.DiscountKeyWorkAccident:     "63",
	bill.DiscountKeySpecialAgreement: "64",
	bill.DiscountKeyProductionError:  "65",
	bill.DiscountKeyNewOutlet:        "66",
	bill.DiscountKeySample:           "67",
	bill.DiscountKeyEndOfRange:       "68",
	bill.DiscountKeyIncoterm:         "70",
	bill.DiscountKeyPOSThreshold:     "71",
	bill.DiscountKeySpecialRebate:    "100",
	bill.DiscountKeyTemporary:        "103",
	bill.DiscountKeyStandard:         "104",
	bill.DiscountKeyYarlyTurnover:    "105",
}

// The following map is useful to get started, but for most users it will make
// sense to use the UNTDID codes directly in the extensions.
var chargeKeyMap = tax.Extensions{
	bill.ChargeKeyStampDuty: "ST",
	bill.ChargeKeyOutlay:    "AAE",
	bill.ChargeKeyTax:       "TX",
	bill.ChargeKeyCustoms:   "ABW",
	bill.ChargeKeyDelivery:  "DL",
	bill.ChargeKeyPacking:   "PC",
	bill.ChargeKeyHandling:  "HD",
	bill.ChargeKeyInsurance: "IN",
	bill.ChargeKeyStorage:   "ABA",
	bill.ChargeKeyAdmin:     "AEM",
	bill.ChargeKeyCleaning:  "CG",
}

func normalizeBillInvoice(m *bill.Invoice) {
	if m.Tax == nil {
		m.Tax = &bill.Tax{}
	}
}

func normalizeBillDiscount(m *bill.Discount) {
	if val, ok := discountKeyMap[m.Key]; ok {
		m.Ext = m.Ext.Merge(tax.Extensions{
			untdid.ExtKeyAllowance: val,
		})
	}
}

func normalizeBillLineDiscount(m *bill.LineDiscount) {
	if val, ok := discountKeyMap[m.Key]; ok {
		m.Ext = m.Ext.Merge(tax.Extensions{
			untdid.ExtKeyAllowance: val,
		})
	}
}

func normalizeBillCharge(m *bill.Charge) {
	if val, ok := chargeKeyMap[m.Key]; ok {
		m.Ext = m.Ext.Merge(tax.Extensions{
			untdid.ExtKeyCharge: val,
		})
	}
}

func normalizeBillLineCharge(m *bill.LineCharge) {
	if val, ok := chargeKeyMap[m.Key]; ok {
		m.Ext = m.Ext.Merge(tax.Extensions{
			untdid.ExtKeyCharge: val,
		})
	}
}

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		// Tax details
		rules.Field("tax",
			rules.Assert("01", "tax details are required with ext and UNTDID document type", is.Present),
			rules.Field("ext",
				rules.Assert("02", "document type extension is required",
					tax.ExtensionsRequire(untdid.ExtKeyDocumentType),
				),
			),
		),
		// Lines: BR-16 requires at least one line
		rules.Field("lines",
			rules.Assert("03", "at least one line is required (BR-16)", is.Present),
		),
		// Supplier: BR-8 requires addresses
		rules.Field("supplier",
			rules.Field("addresses",
				rules.Assert("04", "supplier addresses are required (BR-8)", is.Present),
			),
		),
		// Customer: BR-10 requires addresses when customer is provided
		rules.Field("customer",
			rules.Field("addresses",
				rules.Assert("05", "customer addresses are required (BR-10)", is.Present),
			),
		),
		// Payment: BR-CO-25 requires payment details when amount is due
		rules.When(is.Func("is due standard invoice", isDueStandardInvoice),
			rules.Field("payment",
				rules.Assert("06", "payment details are required when amount is due (BR-CO-25)", is.Present),
				rules.Field("terms",
					rules.Assert("07", "payment terms are required when amount is due (BR-CO-25)", is.Present),
				),
			),
		),
	)
}

func billDiscountRules() *rules.Set {
	return rules.For(new(bill.Discount),
		// BR-32: taxes are required on document-level discounts
		rules.Field("taxes",
			rules.Assert("01", "taxes are required (BR-32)", is.Present),
		),
		// BR-33: either reason or allowance type extension required
		rules.Assert("02", "either a reason or an allowance type extension is required (BR-33)",
			is.Func("reason or allowance", billDiscountHasReasonOrAllowance),
		),
	)
}

func billLineDiscountRules() *rules.Set {
	return rules.For(new(bill.LineDiscount),
		// BR-41: either reason or allowance type extension required
		rules.Assert("01", "either a reason or an allowance type extension is required (BR-41)",
			is.Func("reason or allowance", billLineDiscountHasReasonOrAllowance),
		),
	)
}

func billChargeRules() *rules.Set {
	return rules.For(new(bill.Charge),
		// BR-36: either reason or charge type extension required
		rules.Assert("01", "either a reason or a charge type extension is required (BR-36)",
			is.Func("reason or charge", billChargeHasReasonOrExt),
		),
	)
}

func billLineChargeRules() *rules.Set {
	return rules.For(new(bill.LineCharge),
		// BR-44: either reason or charge type extension required
		rules.Assert("01", "either a reason or a charge type extension is required (BR-44)",
			is.Func("reason or charge", billLineChargeHasReasonOrExt),
		),
	)
}

func isDue(inv *bill.Invoice) bool {
	return inv.Totals != nil &&
		((inv.Totals.Due != nil && !inv.Totals.Due.IsZero()) ||
			(inv.Totals.Due == nil && !inv.Totals.Payable.IsZero()))
}

func isDueStandardInvoice(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && inv.Type.In(bill.InvoiceTypeStandard) && isDue(inv)
}

func billDiscountHasReasonOrAllowance(val any) bool {
	d, ok := val.(*bill.Discount)
	return !ok || d == nil || d.Reason != "" || d.Ext.Has(untdid.ExtKeyAllowance)
}

func billLineDiscountHasReasonOrAllowance(val any) bool {
	d, ok := val.(*bill.LineDiscount)
	return !ok || d == nil || d.Reason != "" || d.Ext.Has(untdid.ExtKeyAllowance)
}

func billChargeHasReasonOrExt(val any) bool {
	c, ok := val.(*bill.Charge)
	return !ok || c == nil || c.Reason != "" || c.Ext.Has(untdid.ExtKeyCharge)
}

func billLineChargeHasReasonOrExt(val any) bool {
	c, ok := val.(*bill.LineCharge)
	return !ok || c == nil || c.Reason != "" || c.Ext.Has(untdid.ExtKeyCharge)
}
