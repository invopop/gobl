package nz

import (
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// irdRegexp matches the basic shape of a New Zealand IRD number: 8 or 9 digits.
var irdRegexp = regexp.MustCompile(`^[0-9]{8,9}$`)

// IRD numbers must fall within this inclusive range to be considered valid.
const (
	irdMin = 10000000
	irdMax = 200000000
)

// irdWeights are the primary weighting factors applied to the IRD base number.
var irdWeights = []int{3, 2, 7, 6, 5, 4, 3, 2}

// irdWeightsSecondary are the weighting factors used when the primary check
// digit calculation yields 10.
var irdWeightsSecondary = []int{7, 4, 3, 2, 5, 2, 7, 6}

// normalizeTaxIdentity strips whitespace, separators and country prefix from a
// New Zealand IRD number.
func normalizeTaxIdentity(tID *tax.Identity) {
	tax.NormalizeIdentity(tID)
}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.Assert("01", "invoice tax id code must be a valid IRD number",
					is.MatchesRegexp(irdRegexp)),
				rules.Assert("02", "invoice tax id code checksum must be valid",
					is.Func("valid", isValidTaxIdentityCode)),
			),
		),
	)
}

// isValidTaxIdentityCode reports whether the value is an IRD number within range
// and with a valid check digit. Empty or non-conforming values are left to the
// format rule above.
func isValidTaxIdentityCode(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return true
	}
	return validIRD(code.String())
}

// validIRD reports whether code is a valid New Zealand IRD number, applying the
// Inland Revenue weighted modulus-11 check-digit algorithm.
func validIRD(code string) bool {
	n, err := strconv.Atoi(code)
	if err != nil || n < irdMin || n > irdMax {
		return false
	}

	// Left-pad to 9 digits so the base is always the first 8 digits and the
	// check digit the 9th.
	if len(code) == 8 {
		code = "0" + code
	}

	check := int(code[8] - '0')

	calc := irdCheckDigit(code, irdWeights)
	if calc == 10 {
		// Re-run with the secondary weights when the first pass is inconclusive.
		calc = irdCheckDigit(code, irdWeightsSecondary)
		if calc == 10 {
			return false
		}
	}
	return calc == check
}

// irdCheckDigit returns the expected check digit for the base (first 8 digits)
// of code given a set of weights.
func irdCheckDigit(code string, weights []int) int {
	sum := 0
	for i, w := range weights {
		sum += int(code[i]-'0') * w
	}
	rem := sum % 11
	if rem == 0 {
		return 0
	}
	return 11 - rem
}
