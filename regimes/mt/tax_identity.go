package mt

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Maltese VAT identification numbers are 8 digits (after removing the "MT"
// prefix). Article 10 and Article 12 registrations carry this format; Article 11
// small-undertaking numbers have no prefix and are not VAT identification
// numbers (Cap. 406 art. 13), so they are not validated here.
var taxCodeRegexp = regexp.MustCompile(`^\d{8}$`)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Maltese VAT identity code",
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

// validateTaxCode validates a normalized Maltese VAT identity code: 8 digits
// with a valid check number.
func validateTaxCode(code cbc.Code) error {
	val := code.String()
	if !taxCodeRegexp.MatchString(val) {
		return errors.New("invalid format")
	}
	if !validateVATChecksum(val) {
		return errors.New("checksum mismatch")
	}
	return nil
}

// validateVATChecksum verifies the Maltese VAT number check digits: the first
// six digits are weighted by [3, 4, 6, 7, 8, 9] and the last two digits, taken
// as a number, must equal 37 - (weighted sum mod 37) (which yields 37 when the
// sum is an exact multiple of 37).
//
// The digit conversion via val[i]-'0' assumes the input is 8 ASCII digits,
// guaranteed by taxCodeRegexp in validateTaxCode.
//
// NOTE: the MTCA, legislation.mt and the EU VIES service publish the "MT" +
// 8-digit format but not the check-digit algorithm. This mod-37 scheme comes
// from community implementations (vat-validator, VatDB) and is verified in the
// tests against MT12701906, MT12357210 and MT13043536.
func validateVATChecksum(val string) bool {
	weights := []int{3, 4, 6, 7, 8, 9}
	sum := 0
	for i := range 6 {
		sum += int(val[i]-'0') * weights[i]
	}
	check := 37 - (sum % 37)
	return int(val[6]-'0')*10+int(val[7]-'0') == check
}
