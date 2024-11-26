package in

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

// Identity type used in India to classify products and services.
const (
	IdentityTypeHSN = cbc.Code("HSN")
)

// HSNCodeRegexp defines the regular expression to validate HSN codes.
var HSNCodeRegexp = regexp.MustCompile(`^(?:\d{4}|\d{6}|\d{8})$`)

func validateItem(value any) error {
	item, ok := value.(*org.Item)
	if !ok || item == nil {
		return nil
	}

	return validation.ValidateStruct(item,
		validation.Field(&item.Identities,
			validation.By(validItemIdentities),
			validation.Skip,
		),
	)
}

func validItemIdentities(value interface{}) error {
	identities, ok := value.([]*org.Identity)
	if !ok {
		return nil
	}

	for _, identity := range identities {
		if identity == nil {
			continue
		}

		if identity.Type == IdentityTypeHSN {
			val := string(identity.Code)

			if !HSNCodeRegexp.MatchString(val) {
				return errors.New("must be a 4, 6, or 8-digit number")
			}
			break
		}
	}

	return nil
}
