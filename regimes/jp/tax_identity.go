package jp

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Japanese qualified-invoice registration number (登録番号): the letter "T"
// followed by 13 digits. For a corporation this is "T" + its 13-digit Corporate
// Number (法人番号); for a sole proprietor or other registrant without a Corporate
// Number, the NTA issues a distinct 13-digit number.
//
// References:
//   - Registration number format:
//     https://www.invoice-kohyo.nta.go.jp/about-toroku/index.html
//   - Public lookup site (for definitive verification):
//     https://www.invoice-kohyo.nta.go.jp/
var taxCodeRegexp = regexp.MustCompile(`^T[0-9]{13}$`)

// Check-digit validation: the 13 digits following "T" are a 1-digit check
// digit followed by the 12-digit base number, per the NTA's published
// algorithm:
// https://www.houjin-bangou.nta.go.jp/documents/checkdigit.pdf
//
//	sum         = Σ (base digit × weight), summed right-to-left, weights
//	              alternating 1, 2, 1, 2, ... starting at the rightmost digit
//	check digit = 9 - (sum mod 9)
//
// Worked example from the NTA PDF: base "700110005901" yields sum 37,
// 37 mod 9 = 1, check digit 9-1 = 8, giving "8700110005901".
//
// Note that the NTA only documents this algorithm for Corporate Numbers, but
// any registration number carrying a "T" is claiming to be a valid qualified
// invoice issuer number, so GOBL checks it regardless. A registrant without
// one can simply omits tax_id.code, which stays valid.
func isValidChecksum(code string) bool {
	// code is "T" + check digit (1) + base number (12).
	if len(code) != 14 {
		return false
	}
	checkDigit := int(code[1] - '0')
	base := code[2:]

	sum := 0
	weight := 1
	for i := len(base) - 1; i >= 0; i-- {
		sum += int(base[i]-'0') * weight
		if weight == 1 {
			weight = 2
		} else {
			weight = 1
		}
	}

	return checkDigit == 9-(sum%9)
}

func normalizeTaxIdentity(tID *tax.Identity) {
	tax.NormalizeIdentity(tID)
}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Japanese registration number",
					is.Func("valid", isValidTaxIdentityCode),
				),
			),
		),
	)
}

func isValidTaxIdentityCode(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok {
		return false
	}
	val := code.String()
	return taxCodeRegexp.MatchString(val) && isValidChecksum(val)
}
