package gb

import (
	"strconv"
	"strings"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Source: https://github.com/ltns35/go-vat

const taxIdentityCodePattern = `^(\d{9}|\d{12}|GD\d{3}|HA\d{3})$`

var (
	taxCodeMultipliers = []int{
		8,
		7,
		6,
		5,
		4,
		3,
		2,
	}
)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn("GB"),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid UK VAT identity format",
					is.Matches(taxIdentityCodePattern),
				),
				rules.AssertIfPresent("02", "all-zero UK VAT identity codes are not allowed",
					is.StringFunc("not-zeros", taxCodeNotZeros),
				),
				rules.AssertIfPresent("03", "UK VAT identity checksum mismatch",
					is.StringFunc("checksum", taxCodeChecksumValid),
				),
			),
		),
	)
}

func taxCodeNotZeros(code string) bool {
	if strings.HasPrefix(code, "GD") || strings.HasPrefix(code, "HA") {
		return true
	}
	z, _ := strconv.Atoi(code)
	return z != 0
}

func taxCodeChecksumValid(code string) bool {
	if strings.HasPrefix(code, "GD") {
		return validGovernmentDepartmentID(code)
	}

	if strings.HasPrefix(code, "HA") {
		return validHealthAuthorityID(code)
	}

	return validCommercialID(code)
}

func validGovernmentDepartmentID(val string) bool {
	const max = 499 // range 000-499
	val = val[2:]
	num, _ := strconv.Atoi(val)
	return num <= max
}

func validHealthAuthorityID(val string) bool {
	const min = 500 // range 500-999
	val = val[2:]
	num, _ := strconv.Atoi(val)
	return num >= min
}

// Specific file used as example: https://github.com/ltns35/go-vat/blob/main/countries/united_kingdom.go
func validCommercialID(val string) bool {
	if len(val) < 9 {
		return false
	}

	// 0 VAT numbers disallowed!
	if z, _ := strconv.Atoi(val); z == 0 {
		return false
	}

	// Check range is OK for modulus 97 calculation
	str := val[:7]
	num, _ := strconv.Atoi(str)

	// Extract the next digit and multiply by the counter.
	sum := 0
	for i, m := range taxCodeMultipliers {
		x := int(val[i] - '0')
		sum += x * m
	}

	// Old numbers use a simple 97 modulus, but new numbers use an adaptation of that (less 55). Our
	// tax number could use either system, so we check it against both.

	// Establish check digits by subtracting 97 from total until negative.
	checkDigit := sum
	for checkDigit > 0 {
		checkDigit = checkDigit - 97
	}

	// Get the absolute value and compare it with the last two characters of the VAT number. If the
	// same, then it is a valid traditional check digit. However, even then the number must fit within
	// certain specified ranges.
	if checkDigit < 0 {
		checkDigit = 0 - checkDigit
	}

	lastDigitsStr := val[7:9]
	lastDigits, _ := strconv.Atoi(lastDigitsStr)

	if checkDigit == lastDigits && num < 9990001 && (num < 100000 || num > 999999) && (num < 9490001 || num > 9700000) {
		return true
	}

	// Now try the new method by subtracting 55 from the check digit if we can - else add 42
	if checkDigit >= 55 {
		checkDigit = checkDigit - 55
	} else {
		checkDigit = checkDigit + 42
	}

	if checkDigit == lastDigits && num > 1000000 {
		return true
	}

	return false // invalid
}
