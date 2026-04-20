package bis

import (
	"github.com/invopop/gobl/catalogues/cef"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// vatexCategoryMap captures the VATEX-to-UNTDID-5305 coherence required by
// PEPPOL-EN16931-P0104 through P0111. A VATEX code constrains which tax
// category can be used alongside it.
var vatexCategoryMap = map[cbc.Code]cbc.Code{
	"VATEX-EU-G":  "G",
	"VATEX-EU-O":  "O",
	"VATEX-EU-IC": "K",
	"VATEX-EU-AE": "AE",
	"VATEX-EU-D":  "E",
	"VATEX-EU-F":  "E",
	"VATEX-EU-I":  "E",
	"VATEX-EU-J":  "E",
}

func taxComboRules() *rules.Set {
	return rules.For(new(tax.Combo),
		rules.Assert("P0104-P0111", "tax category must match VATEX code (PEPPOL-EN16931-P0104..P0111)",
			is.Func("vatex-category coherence", vatexCategoryCoherent),
		),
	)
}

// vatexCategoryCoherent returns true when the combo's VATEX code (if any)
// is paired with the matching UNTDID 5305 tax category.
func vatexCategoryCoherent(val any) bool {
	combo, ok := val.(*tax.Combo)
	if !ok || combo == nil {
		return true
	}
	vatex := combo.Ext.Get(cef.ExtKeyVATEX)
	if vatex == "" {
		return true
	}
	required, ok := vatexCategoryMap[vatex]
	if !ok {
		return true
	}
	current := combo.Ext.Get(untdid.ExtKeyTaxCategory)
	return current == required
}
