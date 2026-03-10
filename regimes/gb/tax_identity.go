package gb

import (
	"strconv"
	"strings"

	"github.com/invopop/gobl/rules"
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
	tID := new(tax.Identity)
	return rules.ForStruct(tID,
		rules.Field(&tID.Code,
			rules.Assert("010", "invalid tax identity code format",
				rules.Matches(taxIdentityCodePattern),
			),
			rules.Assert("020", "invalid tax identity code checksum",
				rules.ByString("checksum", invalidTaxIdentityChecksum),
			),
		).When(tax.IdentityIn("GB")))
}

func invalidTaxIdentityChecksum(code string) bool {
	if strings.HasPrefix(code, "GD") {
		return invalidGovernmentDepartmentID(code)
	}

	if strings.HasPrefix(code, "HA") {
		return invalidHealthAuthorityID(code)
	}

	return invalidCommercialID(code)
}

func invalidGovernmentDepartmentID(val string) bool {
	const expect = 499 // from 000
	val = val[2:]
	num, _ := strconv.Atoi(val)
	return num > expect
}

func invalidHealthAuthorityID(val string) bool {
	const expect = 500 // to 999
	val = val[2:]
	num, _ := strconv.Atoi(val)
	return num > expect
}

// Specific file used as example: https://github.com/ltns35/go-vat/blob/main/countries/united_kingdom.go
func invalidCommercialID(val string) bool {
	// 0 VAT numbers disallowed!
	if z, _ := strconv.Atoi(val); z == 0 {
		return true
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
		return false
	}

	// Now try the new method by subtracting 55 from the check digit if we can - else add 42
	if checkDigit >= 55 {
		checkDigit = checkDigit - 55
	} else {
		checkDigit = checkDigit + 42
	}

	if checkDigit == lastDigits && num > 1000000 {
		return false
	}

	return true
}
