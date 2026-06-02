package oioubl

import (
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

// normalize applies the OIOUBL-specific normalizations during Calculate.
func normalize(doc any) {
	if c, ok := doc.(*tax.Combo); ok {
		normalizeTaxCombo(c)
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
