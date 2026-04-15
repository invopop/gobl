package hu

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// References:
// - python-stdnum: https://github.com/arthurdejong/python-stdnum/blob/master/stdnum/hu/anum.py
// - OECD TIN documentation: https://www.oecd.org/content/dam/oecd/en/topics/policy-issue-focus/aeoi/hungary-tin.pdf
// - NAV Online-Invoice: https://github.com/nav-gov-hu/Online-Invoice

var taxCodeRegexps = []*regexp.Regexp{
	regexp.MustCompile(`^[1-9]\d{7}$`),
}

// Tax code checksum weights for the 8-digit törzsszám.
var taxCodeWeights = []int{9, 7, 3, 1, 9, 7, 3}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Hungarian adószám code",
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

func validateTaxCode(value any) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	match := false
	for _, re := range taxCodeRegexps {
		if re.MatchString(val) {
			match = true
			break
		}
	}
	if !match {
		return errors.New("invalid format")
	}

	return validateTaxCodeChecksum(val)
}

// validateTaxCodeChecksum validates the modulo-10 weighted checksum of the
// 8-digit törzsszám. Weights [9, 7, 3, 1, 9, 7, 3] are applied to the first
// 7 digits. The 8th digit is the check digit: (10 - sum%10) % 10.
func validateTaxCodeChecksum(val string) error {
	sum := 0
	for i := range 7 {
		sum += int(val[i]-'0') * taxCodeWeights[i]
	}

	expected := (10 - (sum % 10)) % 10
	checkDigit := int(val[7] - '0')

	if checkDigit != expected {
		return errors.New("checksum mismatch")
	}

	return nil
}

// normalizeTaxIdentity strips the HU prefix, spaces, dashes, and other
// non-alphanumeric characters, then extracts the 8-digit base number
// if the full 11-digit adószám was provided.
func normalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil || tID.Code == "" {
		return
	}
	// Standard GOBL normalization: uppercase, strip non-alphanumeric, remove country prefix.
	// This handles "HU13895459-2-41" -> "13895459241"
	tax.NormalizeIdentity(tID)

	// If the full 11-digit adószám was provided (after stripping dashes),
	// extract just the 8-digit base number (törzsszám) for EU VAT validation.
	code := tID.Code.String()
	if len(code) == 11 {
		code = code[:8]
	}

	tID.Code = cbc.Code(code)
}
