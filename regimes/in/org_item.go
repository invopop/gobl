package in

import (
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

func validateOrgItem(it *org.Item) error {
	return validation.ValidateStruct(it,
		validation.Field(&it.Identities,
			org.RequireIdentityType(IdentityTypeHSN),
			validation.Skip,
		),
	)
}
