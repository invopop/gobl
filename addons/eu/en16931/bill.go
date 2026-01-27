package en16931

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
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

func validateBillInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Tax,
			validation.Required,
			validation.By(validateBillInvoiceTax),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Required, // BR-16 - at least one line
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.By(validateBillInvoiceParty),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.By(validateBillInvoiceParty),
			validation.Skip,
		),
		validation.Field(&inv.Payment,
			validation.When(
				isDue(inv) && inv.Type.In(bill.InvoiceTypeStandard),
				validation.Required.Error("payment details are required when amount is due (BR-CO-25)"), // BR-CO-25
				validation.By(validateBillPayment),
			),
			validation.Skip,
		),
	)
}

func validateBillInvoiceTax(value any) error {
	tx, ok := value.(*bill.Tax)
	if !ok || tx == nil {
		return nil
	}
	return validation.ValidateStruct(tx,
		validation.Field(&tx.Ext,
			tax.ExtensionsRequire(untdid.ExtKeyDocumentType),
			validation.Skip,
		),
	)
}

func validateBillInvoiceParty(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}

	//BR-8 & BR-10
	return validation.ValidateStruct(p,
		validation.Field(&p.Addresses,
			validation.Required,
			validation.Skip,
		),
	)
}

func validateBillLine(line *bill.Line) error {

	return validation.ValidateStruct(line,
		validation.Field(&line.Discounts,
			validation.Each(
				validation.By(validateBillLineDiscount),
			),
			validation.Skip,
		),
		validation.Field(&line.Charges,
			validation.Each(
				validation.By(validateBillLineCharge),
			),
			validation.Skip,
		),
	)
}

func validateBillLineCharge(value any) error {
	// BR-44
	charge, ok := value.(*bill.LineCharge)
	if !ok || charge == nil {
		return nil
	}

	return validation.ValidateStruct(charge,
		validation.Field(&charge.Reason,
			validation.When(
				!charge.Ext.Has(untdid.ExtKeyCharge),
				validation.Required.Error("either a reason or a charge type extension is required"),
			),
			validation.Skip,
		),
		validation.Field(&charge.Ext,
			validation.When(
				charge.Reason == "",
				validation.Required.Error("either a reason or a charge type extension is required"),
			),
			validation.Skip,
		),
	)
}

func validateBillLineDiscount(value any) error {
	// BR-41
	discount, ok := value.(*bill.LineDiscount)
	if !ok || discount == nil {
		return nil
	}

	return validation.ValidateStruct(discount,
		// BR-41
		validation.Field(&discount.Reason,
			validation.When(
				!discount.Ext.Has(untdid.ExtKeyAllowance),
				validation.Required.Error("either a reason or an allowance type extension is required (BR-41)"),
			),
			validation.Skip,
		),
		validation.Field(&discount.Ext,
			validation.When(
				discount.Reason == "",
				validation.Required.Error("either a reason or an allowance type extension is required (BR-41)"),
			),
			validation.Skip,
		),
	)
}

func validateBillCharge(charge *bill.Charge) error {
	// BR-36
	return validation.ValidateStruct(charge,
		validation.Field(&charge.Reason,
			validation.When(
				!charge.Ext.Has(untdid.ExtKeyCharge),
				validation.Required.Error("either a reason or a charge type extension is required (BR-36)"),
			),
			validation.Skip,
		),
		validation.Field(&charge.Ext,
			validation.When(
				charge.Reason == "",
				validation.Required.Error("either a reason or a charge type extension is required (BR-36)"),
			),
			validation.Skip,
		),
	)
}

func validateBillDiscount(discount *bill.Discount) error {
	// BR-32 & BR-33
	return validation.ValidateStruct(discount,
		validation.Field(&discount.Taxes,
			validation.Required.Error("taxes are required (BR-32)"), // BR-32
			validation.Skip,
		),
		validation.Field(&discount.Reason,
			validation.When(
				!discount.Ext.Has(untdid.ExtKeyAllowance),
				validation.Required.Error("either a reason or an allowance type extension is required (BR-33)"),
			),
			validation.Skip,
		),
		validation.Field(&discount.Ext,
			validation.When(
				discount.Reason == "",
				validation.Required.Error("either a reason or an allowance type extension is required (BR-33)"),
			),
			validation.Skip,
		),
	)
}

func validateBillPayment(value any) error {
	payment, ok := value.(*bill.PaymentDetails)
	if !ok || payment == nil {
		return nil
	}
	return validation.ValidateStruct(payment,
		validation.Field(&payment.Terms,
			validation.Required.Error("payment terms are required when amount is due (BR-CO-25)"),
		),
	)
}

func isDue(inv *bill.Invoice) bool {
	return inv.Totals != nil &&
		((inv.Totals.Due != nil && !inv.Totals.Due.IsZero()) ||
			(inv.Totals.Due == nil && !inv.Totals.Payable.IsZero()))
}
