package gb

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Source: https://github.com/ltns35/go-vat

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
	taxCodeRegexps = []*regexp.Regexp{
		regexp.MustCompile(`^\d{9}$`),
		regexp.MustCompile(`^\d{12}$`),
		regexp.MustCompile(`^GD\d{3}$`),
		regexp.MustCompile(`^HA\d{3}$`),
	}
)

// validateTaxIdentity checks to ensure the NIT code looks okay.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateTaxCode)),
	)
}

func validateTaxCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	match := false
	for _, re := range taxCodeRegexps {
		if re.MatchString(val) {
			match = true
			break
		}
	}
	if !match {
		return errors.New("invalid format")
	}

	if strings.HasPrefix(val, "GD") {
		return governmentDepartmentCheck(val)
	}

	if strings.HasPrefix(val, "HA") {
		return healthAuthorityCheck(val)
	}

	return commercialCheck(val)
}

func governmentDepartmentCheck(val string) error {
	const expect = 499 // from 000
	val = val[2:]
	num, _ := strconv.Atoi(val)
	if num > expect {
		return errors.New("invalid government department number")
	}
	return nil
}

func healthAuthorityCheck(val string) error {
	const expect = 500 // to 999
	val = val[2:]
	num, _ := strconv.Atoi(val)
	if num < expect {
		return errors.New("invalid health authority number")
	}
	return nil
}

// Specific file used as example: https://github.com/ltns35/go-vat/blob/main/countries/united_kingdom.go
func commercialCheck(val string) error {
	// 0 VAT numbers disallowed!
	if z, _ := strconv.Atoi(val); z == 0 {
		return errors.New("cannot be only zeros")
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
		return nil
	}

	// Now try the new method by subtracting 55 from the check digit if we can - else add 42
	if checkDigit >= 55 {
		checkDigit = checkDigit - 55
	} else {
		checkDigit = checkDigit + 42
	}

	if checkDigit == lastDigits && num > 1000000 {
		return nil
	}

	return errors.New("checksum mismatch")
}
