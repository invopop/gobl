package en16931

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
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
		validation.Field(&inv.Payment,
			validation.When(
				//BR-CO-25
				inv.Totals != nil && ((inv.Totals.Due != nil && !inv.Totals.Due.IsZero()) || !inv.Totals.Payable.IsZero()),
				validation.Required,
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
	if charge.Reason == "" && (charge.Ext == nil || charge.Ext[untdid.ExtKeyCharge] == "") {
		return validation.NewError("BR-44", "either a reason or a charge type extension is required")
	}
	return nil
}

func validateBillLineDiscount(value any) error {
	// BR-41
	discount, ok := value.(*bill.LineDiscount)
	if !ok || discount == nil {
		return nil
	}
	if discount.Reason == "" && (discount.Ext == nil || discount.Ext[untdid.ExtKeyAllowance] == "") {
		return validation.NewError("BR-41", "either a reason or an allowance type extension is required")
	}
	return nil
}

func validateBillCharge(charge *bill.Charge) error {
	// BR-36
	if charge.Reason == "" && (charge.Ext == nil || charge.Ext[untdid.ExtKeyCharge] == "") {
		return validation.NewError("BR-36", "either a reason or a charge type extension is required")
	}
	return nil
}

func validateBillDiscount(discount *bill.Discount) error {
	// BR-33
	if discount.Reason == "" && (discount.Ext == nil || discount.Ext[untdid.ExtKeyAllowance] == "") {
		return validation.NewError("BR-33", "either a reason or an allowance type extension is required")
	}
	return nil
}
