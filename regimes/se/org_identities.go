package se

import (
	"errors"
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/luhn"
	"github.com/invopop/validation"
)

const (
	// IdentityTypeOrgNr defines the key for the Swedish Organization Number (Organisationsnummer).
	IdentityTypeOrgNr cbc.Code = "ON" // Officially SE-ON
	// IdentityTypePersonNr defines the key for the Swedish Person Number (Personnummer).
	IdentityTypePersonNr cbc.Code = "PN" // Unofficial
	// IdentityTypeCoordinationNr defines the key for the Swedish Coordination Number (Samordningsnummer).
	IdentityTypeCoordinationNr cbc.Code = "CN" // Unofficial
)

var (
	// ValidOrgIdentityTypes defines the keys for the Swedish organization identities.
	ValidOrgIdentityTypes = []cbc.Code{IdentityTypeOrgNr, IdentityTypePersonNr, IdentityTypeCoordinationNr}

	// Regular expressions for validating Swedish identity codes
	orgNrRegex        = regexp.MustCompile(`^\d{10}$`)              // 0123456789
	individualNrRegex = regexp.MustCompile(`^\d{6}[\-\+]{1}\d{4}$`) // 010101-0101 or 010101+0101
)

var identityTypeDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeOrgNr,
		Name: i18n.String{
			i18n.EN: "Organization Number",
			i18n.SE: "Organisationsnummer",
		},
		Desc: i18n.String{
			i18n.EN: "Swedish company registration number.",
			i18n.SE: "Svenskt företagsregistreringsnummer.",
		},
	},
	{
		Code: IdentityTypePersonNr,
		Name: i18n.String{
			i18n.EN: "Person Number",
			i18n.SE: "Personnummer",
		},
		Desc: i18n.String{
			i18n.EN: "Swedish personal registration number.",
			i18n.SE: "Svenskt personnummer.",
		},
	},
	{
		Code: IdentityTypeCoordinationNr,
		Name: i18n.String{
			i18n.EN: "Coordination Number",
			i18n.SE: "Samordningsnummer",
		},
		Desc: i18n.String{
			i18n.EN: "Swedish coordination number.",
			i18n.SE: "Svenskt samordningsnummer.",
		},
	},
}

// normalizeOrgIdentity performs normalization specific to Swedish identity codes.
//
//   - For organization numbers, it returns a 10 digit number, removing any separators.
//   - For individual numbers, it returns a 10 digit number with the separator. If none are present, a hyphen will be added. If a plus sign (`+`) is present anywhere, it will be used as the separator.
//
// If too many or too few numbers are present, it does nothing.
func normalizeOrgIdentity(id *org.Identity) {
	if id == nil {
		return
	}
	switch id.Type {
	case IdentityTypeOrgNr:
		// Organization numbers should be numeric only, with no separators
		code := cbc.NormalizeNumericalCode(id.Code).String()
		// Only if we have 12 digits, i.e the check digits are present
		// can we safely remove them
		if len(code) == taxCodeLength {
			code = strings.TrimSuffix(code, taxCodeCheckDigit)
		}

		// If we don't have the expected number of digits, it's likely not valid and no safe operation
		// can be performed.
		if len(code) != taxCodeLengthWithoutCheckDigits {
			return
		}

		id.Code = cbc.Code(code)

	case IdentityTypePersonNr, IdentityTypeCoordinationNr:
		// Individual numbers should maintain separator (- or +)
		code := strings.TrimSpace(id.Code.String())

		// If there's no separator but we have the right number of digits,
		// insert a hyphen at the right position, since it's the most
		// statistically likely separator.
		if len(code) == taxCodeLengthWithoutCheckDigits && !strings.ContainsAny(code, "-+") {
			code = code[:6] + "-" + code[6:]
		} else {
			// Extract digits and keep the separator
			digitsOnly := ""
			for _, c := range code {
				if c >= '0' && c <= '9' {
					digitsOnly += string(c)
				}
			}

			// If we don't have the expected number of digits, it's likely not valid and no safe operation
			// can be performed.
			if len(digitsOnly) != taxCodeLengthWithoutCheckDigits {
				return
			}

			// Determine the separator. If a plus sign is present anywhere, preserve it.
			separator := "-"
			if strings.Contains(code, "+") {
				separator = "+"
			}
			code = digitsOnly[:6] + separator + digitsOnly[6:]
		}

		id.Code = cbc.Code(code)

	default:
		return
	}
}

// validateOrgIdentity performs validation for Swedish organization identities.
// Assumes the code has already been normalized.
//
//   - For organization numbers, it checks if the number is 10 digits long.
//   - For individual numbers, it checks if the number is 10 digits long and if the Luhn checksum is valid.
//
// If the number is not valid, it returns an error.
//
// If the organization type is not valid, it returns nil.
func validateOrgIdentity(id *org.Identity) error {
	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.By(func(value any) error {
				return validateOrgIdentityCode(value, id)
			}),
			validation.Skip,
		),
	)
}

func validateOrgIdentityCode(value any, id *org.Identity) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}

	// Normalize to digits only for type check
	digitsOnly := cbc.NormalizeNumericalCode(code).String()

	switch id.Type {
	case IdentityTypeOrgNr:
		if !orgNrRegex.MatchString(digitsOnly) {
			return errors.New("invalid organization number format")
		}
	case IdentityTypePersonNr, IdentityTypeCoordinationNr:
		if !individualNrRegex.MatchString(code.String()) {
			return errors.New("invalid person or coordination number format")
		}
	default:
		return nil
	}
	if !luhn.Check(cbc.Code(digitsOnly)) {
		return errors.New("invalid checksum")
	}
	return nil
}
