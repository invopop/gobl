package oioubl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// validPaymentMeansCodes are the UNTDID 4461 payment-means values accepted for
// OIOUBL (F-LIB100). It includes "30" (generic credit transfer) because the
// gobl.ubl converter maps it to OIOUBL's "31"; the remaining values are emitted
// as-is. Codes outside this set can't produce valid OIOUBL output.
var validPaymentMeansCodes = []cbc.Code{
	"1", "10", "20", "30", "31", "42", "48", "49", "50", "58", "59", "93", "97",
}

// Rule citations reference the OIOUBL Invoice schematron (F-INV) first and
// the CreditNote equivalent (F-CRN) second where one exists. F-INV142 is
// invoice-only because OIOUBL CreditNote uses BillingReference rather than
// OrderLineReference.

var (
	roundingMin = num.MakeAmount(-1000, 2)
	roundingMax = num.MakeAmount(1000, 2)
)

// billInvoiceRules returns the OIOUBL 2.1 rule set for bill.Invoice
// (covers both invoices and credit notes).
func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Field("code",
			rules.Assert("05", "invoice code is required (F-INV009 / F-CRN006)", is.Present),
		),
		rules.Field("supplier",
			rules.Field("inboxes",
				rules.Assert("01", "supplier inboxes are required (F-INV031 / F-CRN028)", is.Present),
			),
		),
		rules.Field("customer",
			rules.Field("inboxes",
				rules.Assert("02", "customer inboxes are required (F-INV044 / F-CRN040)", is.Present),
			),
			// F-INV046 requires exactly one Contact in OIOUBL output;
			// gobl.ubl picks one Person at emit time, so the addon asserts presence only.
			rules.Field("people",
				rules.Assert("03", "customer people are required (F-INV046 / F-CRN042)", is.Present),
			),
		),
		rules.When(is.Func("non-credit-note invoice with line order ref", invoiceWithLineOrderRef),
			rules.Field("ordering",
				rules.Assert("07", "ordering is required when any invoice line has an order reference (F-INV142)", is.Present),
			),
		),
		rules.Field("totals",
			rules.Field("rounding",
				rules.AssertIfPresent("08", "rounding must be between -10.00 and 10.00 (F-INV338 / F-CRN208)", is.Func("in rounding range", roundingInRange)),
			),
		),
		// F-INV239 / F-CRN158: gobl.ubl emits cac:DeliveryLocation whenever
		// delivery.receiver is set; the schematron then requires either an ID
		// (sourced from delivery.identities[0].code) or an Address (sourced
		// from receiver.addresses).
		rules.Field("delivery",
			rules.When(is.Func("receiver set without identities or addresses", deliveryReceiverWithoutLocationData),
				rules.Assert("11", "delivery requires either identities or receiver.addresses (F-INV239 / F-CRN158)", is.Func("never", neverTrue)),
			),
		),
		rules.Field("payment",
			rules.Field("instructions",
				rules.Field("ext",
					rules.AssertIfPresent("12", "payment-means code must be one of the OIOUBL allowed values (F-LIB100)",
						tax.ExtensionsHasCodes(untdid.ExtKeyPaymentMeans, validPaymentMeansCodes...)),
				),
				rules.When(is.Func("bank-transfer payment means without a payee account", bankTransferMissingAccount),
					rules.Assert("13", "a credit transfer account (IBAN or number) is required for bank-transfer payment means (F-LIB107 / F-LIB126)", is.Func("never", neverTrue)),
				),
				rules.When(is.Func("giro payment means without a valid OIOUBL payment id", giroPaymentIDInvalid),
					rules.Assert("14", "Giro (payment-means 50) requires a dk-oioubl-payment-id of 01, 04 or 15 (F-LIB144 / F-LIB147)", is.Func("never", neverTrue)),
				),
				rules.When(is.Func("fik payment means without a valid OIOUBL payment id", fikPaymentIDInvalid),
					rules.Assert("15", "FIK (payment-means 93) requires a dk-oioubl-payment-id of 71, 73 or 75 (F-LIB152)", is.Func("never", neverTrue)),
				),
			),
		),
		rules.Field("lines",
			rules.Each(
				rules.Field("quantity",
					rules.Assert("06", "line quantity must not be zero (F-INV147 / F-CRN088)", is.Func("non-zero amount", quantityNonZero)),
				),
				rules.Field("discounts",
					rules.Each(
						rules.Field("amount",
							rules.Assert("09", "line discount amount must not be negative (F-INV335 / F-CRN203)", is.Func("non-negative amount", amountNonNegative)),
						),
					),
				),
				rules.Field("charges",
					rules.Each(
						rules.Field("amount",
							rules.Assert("10", "line charge amount must not be negative (F-INV335 / F-CRN203)", is.Func("non-negative amount", amountNonNegative)),
						),
					),
				),
			),
		),
	)
}

func invoiceWithLineOrderRef(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	if inv.Type == bill.InvoiceTypeCreditNote {
		return false
	}
	for _, line := range inv.Lines {
		if !line.Order.IsEmpty() {
			return true
		}
	}
	return false
}

func quantityNonZero(val any) bool {
	switch a := val.(type) {
	case num.Amount:
		return !a.IsZero()
	case *num.Amount:
		return a == nil || !a.IsZero()
	}
	return true
}

func amountNonNegative(val any) bool {
	switch a := val.(type) {
	case num.Amount:
		return !a.IsNegative()
	case *num.Amount:
		return a == nil || !a.IsNegative()
	}
	return true
}

func deliveryReceiverWithoutLocationData(val any) bool {
	del, ok := val.(*bill.DeliveryDetails)
	if !ok || del == nil || del.Receiver == nil {
		return false
	}
	for _, id := range del.Identities {
		if !id.Code.IsEmpty() {
			return false
		}
	}
	return len(del.Receiver.Addresses) == 0
}

func neverTrue(any) bool {
	return false
}

// bankTransferCodes are the OIOUBL PaymentMeansCode values that settle to a
// payee bank account (42 domestic, 31 IBAN). OIOUBL then requires the account
// identifier (F-LIB107 for 31, F-LIB126 for 42), which GOBL core leaves optional.
var bankTransferCodes = []cbc.Code{"31", "42"}

func bankTransferMissingAccount(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return false
	}
	code := instr.Ext.Get(untdid.ExtKeyPaymentMeans)
	if !code.In(bankTransferCodes...) {
		return false
	}
	for _, ct := range instr.CreditTransfer {
		if ct != nil && (ct.IBAN != "" || ct.Number != "") {
			return false
		}
	}
	return true
}

func giroPaymentIDInvalid(val any) bool {
	return paymentIDInvalidFor(val, "50", giroPaymentIDs)
}

func fikPaymentIDInvalid(val any) bool {
	return paymentIDInvalidFor(val, "93", fikPaymentIDs)
}

// paymentIDInvalidFor reports whether the instruction uses the given OIOUBL
// payment-means code but lacks a dk-oioubl-payment-id from the allowed set
// (covering both the mandatory-presence and the codelist checks).
func paymentIDInvalidFor(val any, code cbc.Code, allowed []cbc.Code) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return false
	}
	if instr.Ext.Get(untdid.ExtKeyPaymentMeans) != code {
		return false
	}
	return !instr.Ext.Get(ExtKeyPaymentID).In(allowed...)
}

func roundingInRange(val any) bool {
	var a num.Amount
	switch v := val.(type) {
	case num.Amount:
		a = v
	case *num.Amount:
		if v == nil {
			return true
		}
		a = *v
	default:
		return true
	}
	return a.Compare(roundingMin) >= 0 && a.Compare(roundingMax) <= 0
}
