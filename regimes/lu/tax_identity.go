package lu

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Source: https://arthurdejong.org/python-stdnum/doc/1.20/stdnum.lu.tva
// Check digits: first 6 digits modulo 89 must equal the last 2 digits.

var (
	taxCodeRegexp = regexp.MustCompile(`^\d{8}$`)
)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn("LU"),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Luxembourgish VAT identity code",
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
	if code == "" {
		return nil
	}
	val := code.String()
	if !taxCodeRegexp.MatchString(val) {
		return errors.New("invalid format")
	}
	return checksum(val)
}

func checksum(val string) error {
	base, err := strconv.Atoi(val[:6])
	if err != nil {
		return errors.New("invalid format")
	}
	check, err := strconv.Atoi(val[6:])
	if err != nil {
		return errors.New("invalid format")
	}
	if base%89 != check {
		return errors.New("checksum mismatch")
	}
	return nil
}
