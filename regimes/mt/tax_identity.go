package mt

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Source: https://github.com/yolk/valvat/blob/master/lib/valvat/checksum/mt.rb
// Source (Compared): https://github.com/arthurdejong/python-stdnum/blob/master/stdnum/mt/vat.py

// Maltese VAT identification numbers are 8 digits. The "MT" prefix is removed
// during normalization by tax.NormalizeIdentity.
var taxCodeRegexp = regexp.MustCompile(`^\d{8}$`)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Maltese VAT identity code",
					is.Func("valid", validateTaxCode),
				),
			),
		),
	)
}

func validateTaxCode(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok {
		return false
	}
	val := code.String()
	if !taxCodeRegexp.MatchString(val) {
		return false
	}
	return validateTaxCodeChecksum(val)
}

// Weights the 6 base digits by [3, 4, 6, 7, 8, 9] and
// compares 37 - (sum mod 37) against the final two. Indexing is safe because
// validateTaxCode has already matched 8 ASCII digits.
//
// The algorithm is not published: neither the VAT Act nor any MTCA guidance states a
// digit count, a weighting or a check rule, only the "MT" prefix. It follows valvat.
// python-stdnum instead tests the whole eight-digit weighted sum for divisibility by
// 37, which also accepts check+37 and check+74; the stricter form is used here.
func validateTaxCodeChecksum(val string) bool {
	weights := []int{3, 4, 6, 7, 8, 9}
	sum := 0
	for i, w := range weights {
		sum += int(val[i]-'0') * w
	}

	// An exact multiple of 37 yields a check of "37", never "00". Some
	// implementations reduce modulo 37 again and expect "00"; valvat maps zero to 37
	// as here. The smallest base that reaches it is 111111; see the tests.
	check := 37 - sum%37

	expected := int(val[6]-'0')*10 + int(val[7]-'0')
	return check == expected
}
