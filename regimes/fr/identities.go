package fr

import (
	"regexp"

	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

var (
	sirenRegexp = regexp.MustCompile(`^\d{9}$`)
	siretRegexp = regexp.MustCompile(`^\d{14}$`)
)

// validateIdentity validates SIREN and SIRET identity formats
func validateIdentity(id *org.Identity) error {
	if id == nil {
		return nil
	}

	switch id.Type {
	case IdentityTypeSIREN:
		return validation.ValidateStruct(id,
			validation.Field(&id.Code,
				validation.Required,
				validation.Match(sirenRegexp).Error("must be exactly 9 digits"),
			),
		)
	case IdentityTypeSIRET:
		return validation.ValidateStruct(id,
			validation.Field(&id.Code,
				validation.Required,
				validation.Match(siretRegexp).Error("must be exactly 14 digits"),
			),
		)
	}

	return nil
}
