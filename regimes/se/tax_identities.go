package se

import (
	"errors"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pkg/luhn"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// The full length of a Swedish tax ID, including the check digits.
	taxCodeLength = 12
	// The length of the code before the check digits.
	taxCodeLengthWithoutCheckDigits = 10
	// The check digits of a Swedish tax ID.
	taxCodeCheckDigit = "01"
)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn("SE"),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Swedish VAT identity code",
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

// validateTaxCode validates the tax code for Swedish tax identities.
// Assumes the code has already been normalized, is made of 12 numeric characters,
// retaining the checksum at the end, plus 2 control digits "01".
func validateTaxCode(code cbc.Code) error {
	// Normalised Swedish tax IDs must have a specific length.
	if len(code) != taxCodeLength {
		return errors.New("invalid length")
	}
	// Swedish tax IDs must finish in "01".
	if code[10:] != taxCodeCheckDigit {
		return errors.New("invalid check digit, expected 01")
	}
	// Swedish tax IDs must be exclusively numeric.
	if _, err := strconv.Atoi(string(code)); err != nil {
		return errors.New("invalid characters, expected numeric")
	}
	// The code prior to the check digit must be Luhn-valid.
	if !luhn.Check(code[:10]) {
		return errors.New("invalid identification number checksum")
	}
	return nil
}
