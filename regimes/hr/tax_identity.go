package hr

import (
	"regexp"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// The personal identification number consists of eleven digits,
// ten digits that are determined at random and one control digit.
// The "HR" country prefix, spaces and hyphens are removed during
// normalization (tax.NormalizeIdentity), so only the digits reach
// this validation.
// Reference: https://narodne-novine.nn.hr/clanci/sluzbeni/2008_05_60_2033.html
var oibRegexp = regexp.MustCompile(`^\d{11}$`)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "tax identity code for HR must be 11 digits",
					is.MatchesRegexp(oibRegexp),
				),
				rules.AssertIfPresent("02", "tax identity code for HR failed the ISO 7064 MOD 11,10 checksum",
					is.StringFunc("checksum", validOIBChecksum),
				),
			),
		),
	)
}

// validOIBChecksum verifies the OIB check digit using the ISO 7064 MOD 11,10
// algorithm, the same check as u:ctrlOIB in the Tax Administration's
// HR-CIUS-EXT-EN16931-UBL schematron. Each assertion is evaluated
// independently, so a code that already failed the format rule still reaches
// this function and is rejected here too.
// Reference: https://narodne-novine.nn.hr/clanci/sluzbeni/2009_01_1_6.html
func validOIBChecksum(val string) bool {
	if !oibRegexp.MatchString(val) {
		return false
	}
	r := 10
	for i := range 10 {
		a := (r + int(val[i]-'0')) % 10
		if a == 0 {
			a = 10
		}
		r = (a * 2) % 11
	}
	check := 11 - r
	if check == 10 {
		check = 0
	}
	return check == int(val[10]-'0')
}
