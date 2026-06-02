package oioubl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
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
//
// Deliberately NOT enforced here: F-LIB318 (line quantity unitCode must be in
// the OIOUBL codelist). OIOUBL 2.1 ships an older UN/ECE Rec 20 subset (~1100
// codes) that omits common current codes — GOBL's `piece` (H87), `km` (KMT) and
// the packaging units (box/bottle/pallet → X**) are not accepted, and most have
// no OIOUBL equivalent to map to. Enforcing it would mean either maintaining the
// full ~1100-code allowlist in the addon (a codelist-value check that belongs in
// gobl.ubl, not here) or emitting a fabricated fallback (e.g. ZZ "mutually
// defined"). Instead the converter emits the real Rec 20 code and the phive
// schematron rejects an out-of-list unit downstream — the authoritative gate for
// codelist values. So an invoice using e.g. `piece` fails at generation with
// F-LIB318 and the user picks an OIOUBL-valid unit (C62/each).

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
			rules.Field("addresses",
				rules.Assert("16", "supplier address must be a complete OIOUBL StructuredDK address: a postal code (F-LIB033), a street name or PO box (F-LIB034), and a building number or PO box (F-LIB035)",
					is.Func("complete StructuredDK address", addressStructuredDKComplete)),
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
				rules.Assert("20", "the customer contact person requires an identity code for the OIOUBL Contact/ID (F-INV051)",
					is.Func("first person has an identity code", firstPersonHasIdentityCode)),
			),
			rules.Field("addresses",
				rules.Assert("17", "customer address must be a complete OIOUBL StructuredDK address: a postal code (F-LIB033), a street name or PO box (F-LIB034), and a building number or PO box (F-LIB035)",
					is.Func("complete StructuredDK address", addressStructuredDKComplete)),
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
				rules.When(is.Func("iban bank-transfer credit transfer without a BIC", ibanTransferMissingBIC),
					rules.Assert("18", "a BIC is required on the credit transfer for IBAN bank-transfer payment means 30/31 (F-LIB113)", is.Func("never", neverTrue)),
				),
				rules.When(is.Func("giro payment means without a valid OIOUBL payment id", giroPaymentIDInvalid),
					rules.Assert("14", "Giro (payment-means 50) requires a dk-oioubl-payment-id of 01, 04 or 15 (F-LIB144 / F-LIB147)", is.Func("never", neverTrue)),
				),
				rules.When(is.Func("fik payment means without a valid OIOUBL payment id", fikPaymentIDInvalid),
					rules.Assert("15", "FIK (payment-means 93) requires a dk-oioubl-payment-id of 71, 73 or 75 (F-LIB152)", is.Func("never", neverTrue)),
				),
				rules.When(is.Func("structured giro/fik payment id without a valid reference", structuredPaymentRefInvalid),
					rules.Assert("23", "structured Giro/FIK payment id (04/15/71/75) requires a numeric payment reference of the required length (F-LIB145 / F-LIB153 / F-LIB156 / F-LIB157 / F-LIB312 / F-LIB336)", is.Func("never", neverTrue)),
				),
				rules.When(is.Func("fik kortart 73 carrying a payment reference", fik73WithReference),
					rules.Assert("24", "FIK payment id 73 must not carry a payment reference, OIOUBL has no element for it (F-LIB275)", is.Func("never", neverTrue)),
				),
				rules.When(is.Func("giro kortart 01 with an over-long payment reference", giro01ReferenceTooLong),
					rules.Assert("25", "Giro payment id 01 payment reference must be at most 16 characters (F-LIB149)", is.Func("never", neverTrue)),
				),
				rules.When(is.Func("giro payment means without a 7-8 digit payee account", giroAccountInvalid),
					rules.Assert("21", "Giro (payment-means 50) requires a 7 or 8 digit payee account (F-LIB319 / F-LIB320 / F-LIB321)", is.Func("never", neverTrue)),
				),
				rules.When(is.Func("fik payment means without an 8-character creditor account", fikAccountInvalid),
					rules.Assert("22", "FIK (payment-means 93) requires an 8-character creditor account (F-LIB305)", is.Func("never", neverTrue)),
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

// billTaxComboRules returns the OIOUBL 2.1 rule set applied to every tax combo
// (line- and document-level), validated by type the way GOBL validates combos.
func billTaxComboRules() *rules.Set {
	return rules.For(new(tax.Combo),
		rules.Assert("01", "standard-rated VAT must have a percent greater than zero (F-LIB382)",
			is.Func("standard-rated has a positive percent", standardRatedHasPositivePercent)),
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

// firstPersonHasIdentityCode reports whether the first contact person carries an
// identity code. The gobl.ubl converter maps it to the OIOUBL cac:Contact/cbc:ID,
// which is mandatory for the customer (F-INV051); without it the converter would
// have to fabricate a value, so the addon requires real data instead. An empty
// people set passes here since rule 03 governs presence.
func firstPersonHasIdentityCode(val any) bool {
	people, ok := val.([]*org.Person)
	if !ok || len(people) == 0 {
		return true
	}
	p := people[0]
	return p != nil && len(p.Identities) > 0 && !p.Identities[0].Code.IsEmpty()
}

// addressStructuredDKComplete reports whether the first address (the one the
// gobl.ubl converter emits) carries everything OIOUBL's StructuredDK format
// requires: a postal code (F-LIB033), a street name or PO box (F-LIB034), and a
// building number or PO box (F-LIB035). An empty address set passes here since
// EN 16931 already governs address presence.
func addressStructuredDKComplete(val any) bool {
	addrs, ok := val.([]*org.Address)
	if !ok || len(addrs) == 0 {
		return true
	}
	a := addrs[0]
	if a == nil {
		return true
	}
	hasPostbox := a.PostOfficeBox != ""
	hasCode := a.Code != ""
	hasStreet := a.Street != "" || hasPostbox
	hasNumber := a.Number != "" || hasPostbox
	return hasCode && hasStreet && hasNumber
}

// ibanTransferMissingBIC reports whether an IBAN bank-transfer instruction
// (payment-means 30, which the converter maps to 31, or 31 itself) carries a
// credit transfer with no BIC. OIOUBL requires the FinancialInstitution/ID for
// the IBAN channel (F-LIB113), which the converter sources from the BIC.
func ibanTransferMissingBIC(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return false
	}
	if !instr.Ext.Get(untdid.ExtKeyPaymentMeans).In("30", "31") {
		return false
	}
	for _, ct := range instr.CreditTransfer {
		if ct != nil && ct.BIC == "" {
			return true
		}
	}
	return false
}

// standardRatedHasPositivePercent reports whether a tax combo that maps to the
// OIOUBL StandardRated category (UNTDID 5305 "S") carries a percent greater than
// zero. OIOUBL rejects StandardRated with a zero or absent percent (F-LIB382).
func standardRatedHasPositivePercent(val any) bool {
	var combo *tax.Combo
	switch c := val.(type) {
	case *tax.Combo:
		combo = c
	case tax.Combo:
		combo = &c
	default:
		return true
	}
	if combo == nil || combo.Ext.Get(untdid.ExtKeyTaxCategory) != "S" {
		return true
	}
	return combo.Percent != nil && !combo.Percent.Base().IsZero() && !combo.Percent.Base().IsNegative()
}

// bankTransferCodes are the OIOUBL PaymentMeansCode values that settle to a
// payee bank account: 42 (domestic), 31 (IBAN), and 30 (generic credit transfer,
// which the gobl.ubl converter maps to 31). OIOUBL then requires the account
// identifier (F-LIB107 for 30/31, F-LIB126 for 42), which GOBL core leaves optional.
var bankTransferCodes = []cbc.Code{"30", "31", "42"}

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

// giroAccountInvalid reports whether a Giro (payment-means 50) instruction's
// payee account is missing or not 7-8 numeric digits (F-LIB319/320/321).
func giroAccountInvalid(val any) bool {
	return accountLengthInvalid(val, "50", isGiroAccountNumber)
}

// fikAccountInvalid reports whether a FIK (payment-means 93) instruction's
// creditor account is missing or not exactly 8 characters (F-LIB305).
func fikAccountInvalid(val any) bool {
	return accountLengthInvalid(val, "93", func(s string) bool { return len(s) == 8 })
}

// accountLengthInvalid fires when the instruction uses the given payment-means
// code but no credit transfer carries an account number satisfying ok.
func accountLengthInvalid(val any, code cbc.Code, ok func(string) bool) bool {
	instr, isInstr := val.(*pay.Instructions)
	if !isInstr || instr == nil {
		return false
	}
	if instr.Ext.Get(untdid.ExtKeyPaymentMeans) != code {
		return false
	}
	for _, ct := range instr.CreditTransfer {
		if ct != nil && ok(ct.Number) {
			return false
		}
	}
	return true
}

// structuredPaymentRefInvalid reports whether a Giro/FIK instruction using a
// structured kortart (Giro 04/15, FIK 71/75) is missing the numeric payment
// reference that OIOUBL emits as cbc:InstructionID, or carries one of the wrong
// length: mandatory F-LIB145 (Giro) / F-LIB153 (FIK), numeric F-LIB312 (Giro) /
// F-LIB336 (FIK), length F-LIB149 (Giro <=16) / F-LIB156 (FIK 71 = 15) /
// F-LIB157 (FIK 75 = 16). The simple kortart (Giro 01, FIK 73) carry no
// reference and are emitted without an InstructionID, so they need no rule.
func structuredPaymentRefInvalid(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return false
	}
	means := instr.Ext.Get(untdid.ExtKeyPaymentMeans)
	ref := instr.Ref.String()
	switch instr.Ext.Get(ExtKeyPaymentID) {
	case "04", "15":
		return means == "50" && !isNumericOfLen(ref, 1, 16)
	case "71":
		return means == "93" && !isNumericOfLen(ref, 15, 15)
	case "75":
		return means == "93" && !isNumericOfLen(ref, 16, 16)
	}
	return false
}

// fik73WithReference reports whether a FIK (payment-means 93) instruction uses
// the simple kortart 73, which forbids cbc:InstructionID, yet carries a payment
// reference that OIOUBL has nowhere to put (F-LIB275).
func fik73WithReference(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return false
	}
	return instr.Ext.Get(untdid.ExtKeyPaymentMeans) == "93" &&
		instr.Ext.Get(ExtKeyPaymentID) == "73" &&
		instr.Ref != ""
}

// giro01ReferenceTooLong reports whether a Giro (payment-means 50) instruction
// using kortart 01 carries a payment reference longer than the 16 characters
// OIOUBL allows in cbc:InstructionID (F-LIB149).
func giro01ReferenceTooLong(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return false
	}
	return instr.Ext.Get(untdid.ExtKeyPaymentMeans) == "50" &&
		instr.Ext.Get(ExtKeyPaymentID) == "01" &&
		len(instr.Ref.String()) > 16
}

// isNumericOfLen reports whether s consists only of ASCII digits and has a
// length within [minLen, maxLen].
func isNumericOfLen(s string, minLen, maxLen int) bool {
	if len(s) < minLen || len(s) > maxLen {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func isGiroAccountNumber(s string) bool {
	return isNumericOfLen(s, 7, 8)
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
