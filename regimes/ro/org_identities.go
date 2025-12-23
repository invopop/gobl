package ro

import (
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

// Romanian organization identity types:
//
// # CUI (Cod Unic de Înregistrare)
//
// The Unique Registration Code is the main business identifier in Romania.
// It is physically the same number as the CIF (Cod de Identificare Fiscală).
//
// # CNP (Cod Numeric Personal)
//
// Personal Numerical Code is used for individual persons.
// It is a 13-digit code: S YYMMDD JJ NNN C
//
// # References:
//   - ONRC (Trade Registry): https://www.onrc.ro/
//   - CNP Format Info: https://ro.wikipedia.org/wiki/Cod_numeric_personal
const (
	IdentityTypeCUI cbc.Code = "CUI"
	IdentityTypeCNP cbc.Code = "CNP"
)

var identityTypeDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeCUI,
		Name: i18n.String{
			i18n.EN: "Unique Registration Code",
			i18n.RO: "Cod Unic de Înregistrare",
		},
		Desc: i18n.String{
			i18n.EN: "Romanian business registration and tax identification number (CUI/CIF).",
			i18n.RO: "Cod unic de înregistrare și identificare fiscală pentru companii (CUI/CIF).",
		},
	},
	{
		Code: IdentityTypeCNP,
		Name: i18n.String{
			i18n.EN: "Personal Numerical Code",
			i18n.RO: "Cod Numeric Personal",
		},
		Desc: i18n.String{
			i18n.EN: "Romanian personal identification number for individuals (13 digits).",
			i18n.RO: "Cod numeric personal românesc pentru persoane fizice (13 cifre).",
		},
	},
}

// CNP Control Key
var cnpWeights = []int{2, 7, 9, 1, 4, 6, 3, 5, 8, 2, 7, 9}

// normalizeOrgIdentity normalizes Romanian organization identity codes.
func normalizeOrgIdentity(id *org.Identity) {
	if id == nil {
		return
	}

	switch id.Type {
	case IdentityTypeCUI:
		// Use standard alphanumeric normalization (uppercase, remove spaces/dashes/etc.)
		id.Code = cbc.NormalizeAlphanumericalCode(id.Code)
		// Remove "RO" prefix if present
		id.Code = cbc.Code(strings.TrimPrefix(id.Code.String(), "RO"))
	case IdentityTypeCNP:
		// CNP should be numeric only
		id.Code = cbc.NormalizeNumericalCode(id.Code)
	}
}

// validateOrgIdentity validates Romanian organization identities.
func validateOrgIdentity(id *org.Identity) error {
	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.Required,
			validation.By(validateOrgIdentityCode(id.Type)),
		),
	)
}

func validateOrgIdentityCode(idType cbc.Code) validation.RuleFunc {
	return func(value any) error {
		code, ok := value.(cbc.Code)
		if !ok || code == "" {
			return nil
		}

		switch idType {
		case IdentityTypeCUI:
			// CUI validation is identical to tax code validation
			return validateTaxCode(code)
		case IdentityTypeCNP:
			return validateCNP(code)
		default:
			return nil
		}
	}
}

// validateCNP checks the format and checksum of the Personal Numerical Code.
func validateCNP(code cbc.Code) error {
	val := code.String()

	if len(val) != 13 {
		return validation.NewError("validation_cnp_length", "CNP must be exactly 13 digits")
	}

	for _, c := range val {
		if c < '0' || c > '9' {
			return validation.NewError("validation_cnp_format", "CNP must contain only digits")
		}
	}

	// 1. Validate Sex/Century Digit (First digit)
	// 1-8: Residents (M/F for various centuries)
	// 9: Foreign citizens (Rezidenți străini)
	// 0: Invalid
	sexDigit := val[0]
	if sexDigit == '0' {
		return validation.NewError("validation_cnp_sex", "invalid CNP first digit")
	}

	// 2. Validate Checksum
	return validateCNPChecksum(val)
}

// validateCNPChecksum calculates the control digit using the standard algorithm.
// Weights: 279146358279
// Modulo 11 rule: If remainder is 10, control digit is 1. Otherwise, it is the remainder.
func validateCNPChecksum(val string) error {
	var sum int

	// Calculate sum of first 12 digits multiplied by weights
	for i := range 12 {
		digit := int(val[i] - '0')
		sum += digit * cnpWeights[i]
	}

	remainder := sum % 11
	controlDigit := remainder
	// here, using a trick like the modulo 10 optimization from the tax identity validation
	// would require a lookup table that would be less readable for a negligible gain in performance
	if remainder == 10 {
		controlDigit = 1
	}

	// The last digit (13th) is the control digit
	actualControl := int(val[12] - '0')

	if actualControl != controlDigit {
		return validation.NewError("validation_cnp_checksum", "invalid CNP checksum")
	}

	return nil
}
