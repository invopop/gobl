package pl

import (
	"errors"
	"regexp"
	"strconv"
	"unicode"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

/*
 * Sources of data:
 *
 *  - https://pl.wikipedia.org/wiki/Numer_identyfikacji_podatkowej
 *
 */

const (
	taxIdentityPattern = `^[1-9]((\d[1-9])|([1-9]\d))\d{7}$`
)

var (
	taxIdentityRegexp = regexp.MustCompile(taxIdentityPattern)
)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn("PL"),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Polish VAT identity code",
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

func validateTaxCode(code cbc.Code) error {
	if code == "" {
		return nil
	}

	if taxIdentityRegexp.MatchString(code.String()) {
		if validateNIPChecksum(code) {
			return nil
		}
		return errors.New("checksum mismatch")
	}

	return errors.New("invalid format")
}

func validateNIPChecksum(code cbc.Code) bool {
	nipStr := code.String()
	if len(nipStr) != 10 {
		return false
	}

	for _, char := range nipStr {
		if !unicode.IsDigit(char) {
			return false
		}
	}

	digits := make([]int, 10)
	for i, char := range nipStr {
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return false
		}
		digits[i] = digit
	}

	weights := [9]int{6, 5, 7, 2, 3, 4, 5, 6, 7}
	checkSum := 0
	for i, digit := range digits[:9] {
		checkSum += digit * weights[i]
	}
	checkSum %= 11

	return checkSum == digits[9]
}
