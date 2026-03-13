package de

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// Reference: https://github.com/ltns35/go-vat/blob/main/countries/germany.go

var (
	taxCodeRegexps = []*regexp.Regexp{
		regexp.MustCompile(`^[1-9]\d{8}$`),
	}
)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn("DE"),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid German VAT identity code",
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
	p := 10
	sum := 0
	cd := 0
	for i := 0; i < 8; i++ {
		digit, err := strconv.Atoi(string(val[i]))
		if err != nil {
			return errors.New("invalid digit")
		}
		sum = (digit + p) % 10
		if sum == 0 {
			sum = 10
		}
		p = (2 * sum) % 11
	}

	if 11-p == 10 {
		cd = 0
	} else {
		cd = 11 - p
	}

	ecd, err := strconv.Atoi(string(val[8]))
	if err != nil {
		return errors.New("invalid checksum")
	}
	if cd != ecd {
		return errors.New("checksum mismatch")
	}

	return nil
}
