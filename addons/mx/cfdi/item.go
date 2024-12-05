package cfdi

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// SAT item identity codes (ClaveProdServ) regular expression.
var (
	itemExtensionValidCodeRegexp        = regexp.MustCompile(`^\d{8}$`)
	itemExtensionNormalizableCodeRegexp = regexp.MustCompile(`^\d{6}$`)
)

func validateItem(value any) error {
	item, _ := value.(*org.Item)
	if item == nil {
		return nil
	}
	return validation.ValidateStruct(item,
		validation.Field(&item.Ext,
			tax.ExtensionsRequire(ExtKeyProdServ),
			validation.By(validItemExtensions),
			validation.Skip,
		),
	)
}

func validItemExtensions(value interface{}) error {
	ids, ok := value.(tax.Extensions)
	if !ok {
		return nil
	}
	for k, v := range ids {
		if k == ExtKeyProdServ {
			if itemExtensionValidCodeRegexp.MatchString(string(v)) {
				return nil
			}
			return validation.Errors{
				k.String(): validation.NewError("invalid", "must have 8 digits"),
			}
		}
	}
	return nil
}

// extension keys that have been migrated from identities to
// extensions.
var migratedExtensionKeys = []cbc.Key{
	ExtKeyProdServ,
	ExtKeyFiscalRegime,
	ExtKeyUse,
}

func normalizeItem(item *org.Item) {
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
