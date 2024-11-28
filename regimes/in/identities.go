package in

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

const (
	// IdentityKeyPAN represents the Indian Permanent Account Number (PAN). It is a unique identifier assigned
	// to individuals, companies, and other entities.
	IdentityKeyPAN cbc.Key = "in-pan"
)

var panRegexPattern = regexp.MustCompile(`^[A-Z]{5}[0-9]{4}[A-Z]$`)

var identityKeyDefinitions = []*cbc.KeyDefinition{
	{
		Key: IdentityKeyPAN,
		Name: i18n.String{
			i18n.EN: "Permanent Account Number",
			i18n.HI: "स्थायी खाता संख्या",
		},
	},
}

func normalizePAN(id *org.Identity) {
	if id == nil || id.Key != IdentityKeyPAN {
		return
	}
	code := cbc.NormalizeAlphanumericalCode(id.Code).String()
	id.Code = cbc.Code(code)
}

func validatePAN(id *org.Identity) error {
	if id == nil || id.Key != IdentityKeyPAN {
		return nil
	}

	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.Match(panRegexPattern),
		),
	)
}
