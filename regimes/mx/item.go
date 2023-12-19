package mx

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

func validateItem(item *org.Item) error {
	return validation.ValidateStruct(item,
		validation.Field(&item.Ext,
			tax.ExtMapRequires(ExtKeyCFDIProdServ),
			validation.By(validItemExtensions),
			validation.Skip,
		),
	)
}

func validItemExtensions(value interface{}) error {
	ids, ok := value.(tax.ExtMap)
	if !ok {
		return nil
	}
	for k, v := range ids {
		if k == ExtKeyCFDIProdServ {
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
	ExtKeyCFDIProdServ,
	ExtKeyCFDIFiscalRegime,
	ExtKeyCFDIUse,
}

func normalizeItem(item *org.Item) error {
	// 2023-08-25: Migrate identities to extensions
	// Pending removal after migrations completed.
	idents := make([]*org.Identity, 0)
	for _, v := range item.Identities {
		if v.Key.In(migratedExtensionKeys...) {
			if item.Ext == nil {
				item.Ext = make(tax.ExtMap)
			}
			item.Ext[v.Key] = cbc.KeyOrCode(v.Code)
		} else {
			idents = append(idents, v)
		}
	}
	item.Identities = idents
	// end.
	for k, v := range item.Ext {
		if k == ExtKeyCFDIProdServ {
			if itemExtensionNormalizableCodeRegexp.MatchString(v.String()) {
				item.Ext[k] = cbc.KeyOrCode(v.String() + "00")
			}
		}
	}
	return nil
}
