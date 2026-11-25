package uy

import (
	"errors"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// RUT (Registro Único Tributario) is a 12-digit tax identification number
// used in Uruguay for all taxpayers.
//
// Format: XX-XXXXXX-001-X (typically displayed with hyphens, but stored without)
//
//   - Positions 1-2:  Registration type (01-22)
//   - Positions 3-8:  Sequence number (cannot be all zeros)
//   - Positions 9-11: Fixed value "001" (always represents the main taxpayer;
//     branch/sucursal identification is handled separately via the CdgDGISucur
//     field in the CFE XML, not within the RUT itself)
//   - Position 12:    Check digit
//
// The check digit is calculated using a modulo 11 algorithm with the
// following weights applied to the first 11 digits:
//
//	Weights: [4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2]
//
// Calculation:
//  1. Multiply each of the first 11 digits by its corresponding weight
//  2. Sum the products
//  3. Check digit = (-sum) mod 11
//  4. If the check digit is 10 or 11, the RUT is invalid
//
// References:
//   - https://arthurdejong.org/python-stdnum/doc/1.20/stdnum.uy.rut
//   - https://github.com/alfius/uy-rut
//   - https://www.oecd.org/content/dam/oecd/en/topics/policy-issue-focus/aeoi/uruguay-tin.pdf
//   - https://www.gub.uy/direccion-general-impositiva/comunicacion/noticias/nueva-numeracion-del-rut

var rutWeights = []int{4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

// normalizeTaxIdentity removes whitespace, hyphens, and the "UY" prefix
// from the tax identity code.
func normalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil {
		return
	}
	tax.NormalizeIdentity(tID)
}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn("UY"),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Uruguay RUT identity code",
					is.Func("valid", isValidTaxIdentityCode),
				),
			),
		),
	)
}

func isValidTaxIdentityCode(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return false
	}
	return validateTaxCode(code) == nil
}

func validateTaxCode(code cbc.Code) error {
	val := code.String()

	// RUT must be exactly 12 digits
	if len(val) != 12 {
		return errors.New("must have 12 digits")
	}

	// Verify all characters are digits
	for _, c := range val {
		if c < '0' || c > '9' {
			return errors.New("must contain only digits")
		}
	}

	// Validate registration type (first two digits must be 01-22)
	prefix := (int(val[0]-'0') * 10) + int(val[1]-'0')
	if prefix < 1 || prefix > 22 {
		return errors.New("invalid registration type")
	}

	// Sequence number (positions 3-8) cannot be all zeros
	if val[2:8] == "000000" {
		return errors.New("invalid sequence number")
	}

	// Positions 9-11 must be "001"
	if val[8:11] != "001" {
		return errors.New("invalid fixed field")
	}

	// Validate check digit using modulo 11 algorithm
	return validateCheckDigit(val)
}

func validateCheckDigit(val string) error {
	sum := 0
	for i := 0; i < 11; i++ {
		sum += int(val[i]-'0') * rutWeights[i]
	}

	// The check digit is (-sum) mod 11.
	// In Go, the modulo of a negative number can be negative, so we
	// normalize it: ((11 - (sum % 11)) % 11).
	check := (11 - (sum % 11)) % 11
	if check >= 10 {
		return errors.New("invalid check digit")
	}

	actual := int(val[11] - '0')
	if actual != check {
		return errors.New("checksum mismatch")
	}

	return nil
}
