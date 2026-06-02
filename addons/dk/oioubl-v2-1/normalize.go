package oioubl

import (
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
)

// normalize applies the OIOUBL-specific normalizations during Calculate.
func normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Combo:
		normalizeTaxCombo(obj)
	case *pay.Instructions:
		normalizePayInstructions(obj)
	}
}

// normalizePayInstructions records the OIOUBL paymentchannelcode-1.1 value in the
// dk-oioubl-payment-channel extension, derived from the payment means, so the
// gobl.ubl serializer emits cbc:PaymentChannelCode directly.
func normalizePayInstructions(instr *pay.Instructions) {
	if instr == nil {
		return
	}
	if ch := oioublPaymentChannel(instr.Ext.Get(untdid.ExtKeyPaymentMeans)); ch != "" {
		instr.Ext = instr.Ext.Set(ExtKeyPaymentChannel, ch)
	}
}

// oioublPaymentChannel maps a UNTDID 4461 payment means to its OIOUBL payment
// channel: Giro (50) → DK:GIRO, FIK (93) → DK:FIK, direct debit (49) carries no
// channel, and every other settled means defaults to IBAN.
func oioublPaymentChannel(means cbc.Code) cbc.Code {
	switch means {
	case "":
		return ""
	case "50":
		return ExtValuePaymentChannelGiro
	case "93":
		return ExtValuePaymentChannelFIK
	case "49":
		return ""
	default:
		return ExtValuePaymentChannelIBAN
	}
}

// normalizeTaxCombo records the OIOUBL taxcategoryid-1.1 category for a VAT combo
// in the dk-oioubl-tax-category extension, derived from the EN 16931 UNTDID
// category. This moves the mapping out of the gobl.ubl serializer, which then
// emits the value directly. The GOBL category itself is left untouched — in
// particular VAT-exempt stays "exempt", so EN 16931 keeps requiring the
// exemption reason (and allows the VATEX code), even though OIOUBL reports it as
// ZeroRated (OIOUBL 2.1 has no exempt category).
func normalizeTaxCombo(c *tax.Combo) {
	if c == nil || c.Category != tax.CategoryVAT {
		return
	}
	if oc := oioublTaxCategory(c.Ext.Get(untdid.ExtKeyTaxCategory)); oc != "" {
		c.Ext = c.Ext.Set(ExtKeyTaxCategory, oc)
	}
}

// oioublTaxCategory maps an EN 16931 UNTDID 5305 VAT category to its OIOUBL
// taxcategoryid-1.1 equivalent. Exempt (E) has no OIOUBL counterpart and is
// reported as ZeroRated, as both mean no VAT is charged.
func oioublTaxCategory(untdidCat cbc.Code) cbc.Code {
	switch untdidCat {
	case "S":
		return ExtValueTaxCategoryStandardRated
	case "Z", "E":
		return ExtValueTaxCategoryZeroRated
	case "AE":
		return ExtValueTaxCategoryReverseCharge
	}
	return ""
}
