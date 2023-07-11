package mx

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

// SAT item identity codes (ClaveProdServ) regular expression.
var (
	itemIdentityValidCodeRegexp = regexp.MustCompile(`^\d{8}$`)
)

func validateItem(item *org.Item) error {
	return validation.ValidateStruct(item,
		validation.Field(&item.Unit, validation.Required),
		validation.Field(&item.Identities, validation.By(validItemIdentities)),
	)
}

func validItemIdentities(value interface{}) error {
	ids, _ := value.([]*org.Identity)

	for _, id := range ids {
		if id.Type == IdentityTypeSAT {
			if itemIdentityCodeRegexp.MatchString(string(id.Code)) {
				return nil
			}
			return errors.New("SAT code must have 8 digits")
		}
	}

	return errors.New("SAT code must be present")
}
