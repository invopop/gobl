package bis

import (
	"strings"

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
		rules.Assert("01", "tax category must match VATEX code (PEPPOL-EN16931-P0104..P0111)",
			is.Func("vatex-category coherence", vatexCategoryCoherent),
		),
	)
}

// vatexCategoryCoherent returns true when the combo's VATEX code (if any) is
// paired with the matching UNTDID 5305 tax category. For known VATEX codes
// the mapping is explicit; for any other `VATEX-EU-*` code we fail closed
// and require category `E` (the exemption default under the CEF codelist).
func vatexCategoryCoherent(val any) bool {
	combo, ok := val.(*tax.Combo)
	if !ok || combo == nil {
		return true
	}
	vatex := combo.Ext.Get(cef.ExtKeyVATEX)
	if vatex == "" {
		return true
	}
	current := combo.Ext.Get(untdid.ExtKeyTaxCategory)
	if required, ok := vatexCategoryMap[vatex]; ok {
		return current == required
	}
	// Unknown VATEX-EU-* codes default to category E so we fail closed on new
	// or mistyped codes instead of silently accepting any pairing.
	if strings.HasPrefix(vatex.String(), "VATEX-EU-") {
		return current == "E"
	}
	return true
}
