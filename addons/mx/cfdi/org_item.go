package cfdi

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// SAT item identity codes (ClaveProdServ) regular expression.
var (
	itemExtensionValidCodeRegexp        = regexp.MustCompile(`^\d{8}$`)
	itemExtensionNormalizableCodeRegexp = regexp.MustCompile(`^\d{6}$`)
)

// extension keys that have been migrated from identities to
// extensions.
var migratedExtensionKeys = []cbc.Key{
	ExtKeyProdServ,
	ExtKeyFiscalRegime,
	ExtKeyUse,
}

func normalizeItem(item *org.Item) {
	if item == nil {
		return
	}
	// 2023-08-25: Migrate identities to extensions
	// Pending removal after migrations completed.
	idents := make([]*org.Identity, 0)
	for _, v := range item.Identities {
		if v.Key.In(migratedExtensionKeys...) {
			if item.Ext == nil {
				item.Ext = make(tax.Extensions)
			}
			item.Ext[v.Key] = v.Code
		} else {
			idents = append(idents, v)
		}
	}
	item.Identities = idents
	// end.
	for k, v := range item.Ext {
		if k == ExtKeyProdServ {
			if itemExtensionNormalizableCodeRegexp.MatchString(v.String()) {
				item.Ext[k] = cbc.Code(v.String() + "00")
			}
		}
	}
}
