package nz

import (
	"strconv"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Reference: https://www.ird.govt.nz/digital-service-providers/services-catalogue/customer-and-account/ird-number-validation
// Reference: https://www.oecd.org/content/dam/oecd/en/topics/policy-issue-focus/aeoi/new-zealand-tin.pdf

const (
	taxIdentityCodePattern = `^\d{8,9}$`

	irdMinValue = 10_000_000
	irdMaxValue = 150_000_000
)

var (
	// weightsPass1 are the primary weights used for the IRD modulo-11 check digit algorithm.
	weightsPass1 = []int{3, 2, 7, 6, 5, 4, 3, 2}
	// weightsPass2 are the fallback weights used when pass 1 produces a check digit of 10.
	weightsPass2 = []int{7, 4, 3, 2, 5, 2, 7, 6}
)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid NZ IRD number format",
					is.Matches(taxIdentityCodePattern),
				),
				rules.AssertIfPresent("02", "IRD number out of valid range",
					is.StringFunc("range", irdInValidRange),
				),
				rules.AssertIfPresent("03", "IRD number checksum mismatch",
					is.StringFunc("checksum", irdChecksumValid),
				),
			),
		),
	)
}

func irdInValidRange(code string) bool {
	n, err := strconv.Atoi(code)
	if err != nil {
		return false
	}
	return n >= irdMinValue && n <= irdMaxValue
}

// irdChecksumValid validates the IRD number checksum using a two-pass modulo-11 algorithm.
// Pass 1 uses weights [3,2,7,6,5,4,3,2]. When pass 1 produces an impossible check digit of 10,
// pass 2 falls back to weights [7,4,3,2,5,2,7,6]. If both passes yield 10, the number is invalid.
func irdChecksumValid(code string) bool {
	// Guard: checksum requires exactly 8 or 9 digits; let the format rule handle other lengths.
	if len(code) < 8 || len(code) > 9 {
		return true
	}

	// Split into base digits and check digit.
	base := code[:len(code)-1]
	checkDigit := int(code[len(code)-1] - '0')

	// If the base is 7 digits (8-digit IRD), pad left with a zero.
	if len(base) == 7 {
		base = "0" + base
	}

	computed := irdComputeCheckDigit(base, weightsPass1)
	if computed == 10 {
		// Pass 1 produced an impossible digit — fall back to pass 2 weights.
		computed = irdComputeCheckDigit(base, weightsPass2)
		if computed == 10 {
			// Both passes failed; the number is structurally invalid.
			return false
		}
	}

	return computed == checkDigit
}

// irdComputeCheckDigit applies a set of weights to the 8-digit base string
// and returns the modulo-11 check digit (may return 10 on failure).
func irdComputeCheckDigit(base string, weights []int) int {
	sum := 0
	for i, w := range weights {
		sum += int(base[i]-'0') * w
	}
	rem := sum % 11
	if rem == 0 {
		return 0
	}
	return 11 - rem
}
