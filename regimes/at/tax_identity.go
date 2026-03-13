package at

import (
	"errors"
	"math"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// Source: https://github.com/ltns35/go-vat

var (
	taxCodeMultipliers = []int{
		1,
		2,
		1,
		2,
		1,
		2,
		1,
	}
	taxCodeRegexps = []*regexp.Regexp{
		regexp.MustCompile(`^U\d{8}$`),
	}
)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn("AT"),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Austrian VAT identity code",
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

func validateTaxCode(value interface{}) error {
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

	return commercialCheck(val)
}

func commercialCheck(val string) error {
	var total float64
	for i, m := range taxCodeMultipliers {
		num := int(val[i+1] - '0')
		x := float64(num * m)
		if x > 9 {
			total += math.Floor(x/10) + math.Mod(x, 10)
		} else {
			total += x
		}
	}

	total = 10 - math.Mod(total+4, 10)
	if total == 10 {
		total = 0
	}

	lastNum := int(val[8] - '0')
	if lastNum != int(total) {
		return errors.New("checksum mismatch")
	}

	return nil
}
