package cl

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Chilean RUT tax identity code length",
					is.Func("valid length", isValidTaxIdentityLength),
				),
				rules.When(is.Func("valid length", isValidTaxIdentityLength),
					rules.AssertIfPresent("02", "Chilean RUT tax identity code body must be numeric",
						is.Func("valid body", isValidTaxIdentityBody),
					),
					rules.When(is.Func("valid body", isValidTaxIdentityBody),
						rules.AssertIfPresent("03", "invalid Chilean RUT tax identity check digit",
							is.Func("valid check digit", isValidTaxIdentityCheckDigit),
						),
						rules.When(is.Func("valid check digit", isValidTaxIdentityCheckDigit),
							rules.AssertIfPresent("04", "invalid Chilean RUT tax identity checksum",
								is.Func("valid checksum", isValidTaxIdentityChecksum),
							),
						),
					),
				),
			),
		),
	)
}

func normalizeTaxIdentity(tID *tax.Identity) {
	tax.NormalizeIdentity(tID)
}

func isValidTaxIdentityLength(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return false
	}
	return isValidTaxCodeLength(code)
}

func isValidTaxIdentityBody(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || code == "" || !isValidTaxCodeLength(code) {
		return false
	}

	s := code.String()
	body := s[:len(s)-1]
	for _, c := range body {
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}

func isValidTaxIdentityCheckDigit(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || code == "" || !isValidTaxCodeLength(code) || !isValidTaxIdentityBody(value) {
		return false
	}

	s := code.String()
	return isValidCheckDigit(s[len(s)-1])
}

func isValidTaxIdentityChecksum(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || code == "" || !isValidTaxCodeLength(code) || !isValidTaxIdentityBody(value) || !isValidTaxIdentityCheckDigit(value) {
		return false
	}

	s := code.String()
	body := s[:len(s)-1]
	check := s[len(s)-1]
	return expectedRUTCheckDigit(body) == check
}

func isValidTaxCodeLength(code cbc.Code) bool {
	s := code.String()
	return len(s) >= 8 && len(s) <= 9
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
