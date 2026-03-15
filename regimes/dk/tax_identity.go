package dk

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// Reference: https://github.com/ltns35/go-vat

var (
	taxCodeRegexps = []*regexp.Regexp{
		regexp.MustCompile(`^\d{8}$`),
	}
)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn("DK"),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Danish VAT identity code",
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

func validateTaxCodeChecksum(val string) error {
	// Danish CVR numbers use modulo-11 checksum with multipliers [2, 7, 6, 5, 4, 3, 2, 1]
	multipliers := []int{2, 7, 6, 5, 4, 3, 2, 1}
	total := 0

	for i := range 8 {
		digit, err := strconv.Atoi(string(val[i]))
		if err != nil {
			return errors.New("invalid digit")
		}
		total += digit * multipliers[i]
	}

	if total%11 != 0 {
		return errors.New("checksum mismatch")
	}

	return nil
}
