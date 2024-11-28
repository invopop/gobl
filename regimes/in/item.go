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

func validateItem(it *org.Item) error {

	return validation.ValidateStruct(it,
		validation.Field(&it.Identities,
			validation.Each(

				validation.By(validItemIdentities),
			),
			validation.Skip,
		),
	)
}

func validItemIdentities(value interface{}) error {
	id, ok := value.(*org.Identity)
	if !ok {
		return nil
	}

	if id.Type == IdentityTypeHSN {
		val := string(id.Code)

		if !HSNCodeRegexp.MatchString(val) {
			return errors.New("must be a 4, 6, or 8-digit number")
		}

	}

	return nil
}
