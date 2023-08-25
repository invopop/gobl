package mx

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
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
			cbc.CodeMapHas(ExtKeyCFDIProdServ),
			validation.By(validItemExtensions),
			validation.Skip,
		),
	)
}

func validItemExtensions(value interface{}) error {
	ids, ok := value.(cbc.CodeMap)
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

func normalizeItem(item *org.Item) error {
	// 2023-08-25: Migrate identities to extensions
	// Pending removal after migrations completed.
	idents := make([]*org.Identity, 0)
	for _, v := range item.Identities {
		if v.Key != "" {
			if item.Ext == nil {
				item.Ext = make(cbc.CodeMap)
			}
			item.Ext[v.Key] = v.Code
		} else {
			idents = append(idents, v)
		}
	}
	item.Identities = idents
	// end.
	for k, v := range item.Ext {
		if k == ExtKeyCFDIProdServ {
			if itemExtensionNormalizableCodeRegexp.MatchString(v.String()) {
				item.Ext[k] = cbc.Code(v.String() + "00")
			}
		}
	}
	return nil
}
