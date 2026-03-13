package gr

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// Reference: https://lytrax.io/blog/projects/greek-tin-validator-generator

var (
	taxCodeRegexp = regexp.MustCompile(`^\d{9}$`)
)

// normalizeTaxIdentity requires additional steps for Greece as the language code
// is used in the tax code.
func normalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil {
		return
	}
	// also allow for usage of "GR" which may be used in the tax code
	// by accident.
	tax.NormalizeIdentity(tID, l10n.GR)
	tID.Country = "EL" // always override for greece
}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn("EL"),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Greek VAT identity code",
					rules.By("valid", isValidTaxIdentityCode),
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

func validateTaxCode(value any) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	if !taxCodeRegexp.MatchString(val) {
		return errors.New("invalid format")
	}

	if !hasValidChecksum(val) {
		return errors.New("checksum mismatch")
	}

	return nil
}

func hasValidChecksum(val string) bool {
	digits := make([]int, 9)
	for i, char := range val {
		num, _ := strconv.Atoi(string(char)) // ignore errors, we already validated the format
		digits[i] = num
	}

	var sum int
	for i := 0; i < 8; i++ {
		sum += digits[i] * (1 << uint(8-i))
	}
	checkDigit := sum % 11 % 10

	return checkDigit == digits[8]
}
