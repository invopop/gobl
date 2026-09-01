package nz

import (
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	minIRD = 10000000
	maxIRD = 200000000
)

var (
	// IRD numbers are 8 or 9 digits, with no separators or letters.
	irdPattern = regexp.MustCompile(`^[0-9]{8,9}$`)

	// New Zealand IRD uses a modulus-11 check with a secondary fallback when
	// the primary calculation resolves to 10.
	irdWeights    = []int{3, 2, 7, 6, 5, 4, 3, 2}
	irdWeightsAlt = []int{7, 4, 3, 2, 5, 2, 7, 6}
)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.Assert(
					"01",
					"invoice tax id code must be a valid IRD number",
					is.MatchesRegexp(irdPattern),
				),
				rules.Assert(
					"02",
					"invoice tax id code checksum must be valid",
					is.Func("valid", validTaxIdentityCode),
				),
			),
		),
	)
}

func validTaxIdentityCode(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return true
	}
	if !irdPattern.MatchString(code.String()) {
		return false
	}
	return validateIRD(code.String())
}

// validateIRD applies the published IRD modulus-11 validation rules.
func validateIRD(code string) bool {
	number, err := strconv.Atoi(code)
	if err != nil || number < minIRD || number > maxIRD {
		return false
	}

	if len(code) == 8 {
		code = "0" + code
	}

	expected := int(code[8] - '0')
	actual := calculateIRDCheckDigit(code, irdWeights)

	if actual == 10 {
		actual = calculateIRDCheckDigit(code, irdWeightsAlt)
		if actual == 10 {
			return false
		}
	}

	return actual == expected
}

func calculateIRDCheckDigit(code string, weights []int) int {
	total := 0
	for i, weight := range weights {
		total += int(code[i]-'0') * weight
	}

	remainder := total % 11
	if remainder == 0 {
		return 0
	}

	return 11 - remainder
}
