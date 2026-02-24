package no

import (
	"errors"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

const (
	// IdentityTypeOrgNr defines the key for the Norwegian Organization Number (Organisasjonsnummer).
	IdentityTypeOrgNr cbc.Code = "ON"
)

var identityTypeDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeOrgNr,
		Name: i18n.String{
			i18n.EN: "Organization Number",
			i18n.NO: "Organisasjonsnummer",
		},
		Desc: i18n.String{
			i18n.EN: "Norwegian organization number (9 digits).",
			i18n.NO: "Norsk organisasjonsnummer (9 siffer).",
		},
	},
}

// normalizeOrgIdentity normalizes Norwegian org identity codes.
func normalizeOrgIdentity(id *org.Identity) {
	if id == nil {
		return
	}
	if id.Type != IdentityTypeOrgNr {
		return
	}

	digits, ok := cleanNorwayOrgNr(id.Code)
	if ok {
		id.Code = cbc.Code(digits)
	}
}

func validateOrgIdentity(id *org.Identity) error {
	if id == nil {
		return nil
	}

	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.By(func(value any) error {
				code, ok := value.(cbc.Code)
				if !ok || code == "" {
					return nil
				}
				if id.Type != IdentityTypeOrgNr {
					return nil
				}

				digits, ok := cleanNorwayOrgNr(code)
				if !ok {
					return errors.New("invalid organization number format")
				}
				if !validOrgNrMod11(digits) {
					return errors.New("invalid organization number checksum")
				}
				return nil
			}),
			validation.Skip,
		),
	)
}

// MOD11 checksum for Norwegian organization number.
// Weights: 3,2,7,6,5,4,3,2 (for first 8 digits). Check digit is 9th.
func validOrgNrMod11(s string) bool {
	if len(s) != 9 {
		return false
	}
	weights := []int{3, 2, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i := 0; i < 8; i++ {
		d := int(s[i] - '0')
		if d < 0 || d > 9 {
			return false
		}
		sum += d * weights[i]
	}
	rem := sum % 11
	cd := 11 - rem
	if cd == 11 {
		cd = 0
	}
	if cd == 10 {
		return false
	}
	check := int(s[8] - '0')
	return check == cd
}
