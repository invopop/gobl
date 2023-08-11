package mx

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

// SAT item identity codes (ClaveProdServ) regular expression.
var (
	itemIdentityValidCodeRegexp        = regexp.MustCompile(`^\d{8}$`)
	itemIdentityNormalizableCodeRegexp = regexp.MustCompile(`^\d{6}$`)
)

func validateItem(item *org.Item) error {
	return validation.ValidateStruct(item,
		validation.Field(&item.Identities,
			org.HasIdentityKey(IdentityKeyProductCode),
			validation.By(validItemIdentities),
			validation.Skip,
		),
	)
}

func validItemIdentities(value interface{}) error {
	ids, ok := value.([]*org.Identity)
	if !ok {
		return nil
	}
	for _, id := range ids {
		if id.Key == IdentityKeyProductCode {
			if itemIdentityValidCodeRegexp.MatchString(string(id.Code)) {
				return nil
			}
			return errors.New("SAT code must have 8 digits")
		}
	}
	return nil
}

func normalizeItem(item *org.Item) error {
	for _, id := range item.Identities {
		if id.Type == IdentityTypeSAT && itemIdentityNormalizableCodeRegexp.MatchString(string(id.Code)) {
			id.Code = id.Code + "00"
		}
	}
	return nil
}
