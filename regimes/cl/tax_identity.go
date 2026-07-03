package cl

import (
	"errors"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Chilean RUT tax identity code",
					is.Func("valid", isValidTaxIdentityCode),
				),
			),
		),
	)
}

func normalizeTaxIdentity(tID *tax.Identity) {
	tax.NormalizeIdentity(tID)
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

	s := code.String()
	if len(s) < 8 || len(s) > 9 {
		return errors.New("invalid length")
	}

	body := s[:len(s)-1]
	check := s[len(s)-1]
	for _, c := range body {
		if c < '0' || c > '9' {
			return errors.New("body contains invalid characters")
		}
	}
	if !isValidCheckDigit(check) {
		return errors.New("invalid check digit")
	}
	if expectedRUTCheckDigit(body) != check {
		return errors.New("checksum mismatch")
	}

	return nil
}

func isValidCheckDigit(c byte) bool {
	return (c >= '0' && c <= '9') || c == 'K'
}

func expectedRUTCheckDigit(body string) byte {
	sum := 0
	factor := 2
	for i := len(body) - 1; i >= 0; i-- {
		sum += int(body[i]-'0') * factor
		factor++
		if factor > 7 {
			factor = 2
		}
	}

	value := 11 - (sum % 11)
	switch value {
	case 11:
		return '0'
	case 10:
		return 'K'
	default:
		return byte('0' + value)
	}
}
