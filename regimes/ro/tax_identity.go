package ro

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

/*
	The Romanian CIF (Cod de Identificare Fiscală, historically also known as
the CUI, Cod Unic de Înregistrare) is 2 to 10 digits long. The last digit
is a check digit calculated using a fixed weighting key and a MOD 11
algorithm.

	Source: https://fintp.org/wp-content/uploads/2020/03/FinTPc-SRS.pdf
	(The exact algorithm can be found at page 76)
*/

// We use a regex to reject anything outside 2-10 plain digits. No leading 0s.
var taxCodeRegexp = regexp.MustCompile(`^[1-9]\d{1,9}$`)

/*
	taxCodeWeights are applied right-to-left against the code digits
(excluding the check digit), once the code has been left-padded with
zeros to 9 digits.
*/

// Tax code weights from the document
var taxCodeWeights = []int{7, 5, 3, 2, 1, 7, 5, 3, 2}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Romanian CIF/CUI",
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

func validateTaxCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	if !taxCodeRegexp.MatchString(val) {
		return errors.New("invalid format")
	}

	return checksumCheck(val)
}

func checksumCheck(val string) error {

	// Split off the check digit from the payload
	check := int(val[len(val)-1] - '0')
	digits := val[:len(val)-1]

	// Left-pad to 9 digits (regardless of how short the code remains) to reach the standard
	for len(digits) < len(taxCodeWeights) {
		digits = "0" + digits
	}

	// Weighted sum based on digit[i] x taxCodeWeights[i]
	total := 0
	for i, w := range taxCodeWeights {
		total += int(digits[i]-'0') * w
	}

	// MOD 11 reduction
	control := (total * 10) % 11
	if control == 10 {
		control = 0
	}

	// Final check
	if control != check {
		return errors.New("checksum mismatch")
	}

	// Everything good, the code matches
	return nil
}
