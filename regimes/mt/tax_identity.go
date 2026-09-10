package mt

import (
	"regexp"

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
				rules.AssertIfPresent("01", "invalid Maltese VAT number format",
					is.MatchesRegexp(taxCodeRegexp),
				),
				rules.AssertIfPresent("02", "invalid Maltese VAT number checksum",
					is.StringFunc("checksum", validVATChecksum),
				),
			),
		),
	)
}

// validVATChecksum verifies the Maltese VAT number check digits: the first six
// digits are weighted by [3, 4, 6, 7, 8, 9] and the last two digits, taken as a
// number, must equal 37 - (weighted sum mod 37) (which yields 37 when the sum is
// an exact multiple of 37). Malformed input defers to the format assertion so the
// two checks report independently.
//
// NOTE: the MTCA, legislation.mt and the EU VIES service publish the "MT" +
// 8-digit format but not the check-digit algorithm. This mod-37 scheme comes from
// community implementations (vat-validator, VatDB) and is verified in the tests
// against MT12701906, MT12357210 and MT13043536.
func validVATChecksum(code string) bool {
	if !taxCodeRegexp.MatchString(code) {
		return true
	}
	weights := []int{3, 4, 6, 7, 8, 9}
	sum := 0
	for i := range 6 {
		sum += int(code[i]-'0') * weights[i]
	}
	check := 37 - (sum % 37)
	return int(code[6]-'0')*10+int(code[7]-'0') == check
}
