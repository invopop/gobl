package lv

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// References:
// - https://www.vid.gov.lv/lv/pievienotas-vertibas-nodokla-maksataji
// - https://lookuptax.com/docs/tax-identification-number/latvia-tax-id-guide
// - https://vatdb.com/guides/validate-lv-vat-number/

var (
	taxCodeRegexps = []*regexp.Regexp{
		regexp.MustCompile(`^LV\d{11}$`),
		regexp.MustCompile(`^\d{11}$`),
	}
	// Weights for the Mod-11 checksum calculation (positions 1-10)
	taxCodeWeights = []int{9, 1, 4, 8, 3, 10, 2, 5, 7, 6}
)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid VAT identity code",
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

func validateTaxCode(value any) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	// Check format - must be LV + 11 digits, or just 11 digits
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

	// Strip LV prefix if present
	if len(val) == 13 {
		val = val[2:]
	}

	// Get first two digits
	firstDigit, err := strconv.Atoi(string(val[0]))
	if err != nil {
		return errors.New("invalid first digit")
	}

	// Personal codes (first digit 0-3) - no checksum validation needed
	// These follow the format similar to personas kods without the hyphen
	if firstDigit >= 0 && firstDigit <= 3 {
		return validatePersonalCode(val)
	}

	// Legal entities (first digit 4-9) - validate Mod-11 checksum
	if firstDigit >= 4 && firstDigit <= 9 {
		return validateMod11Checksum(val)
	}

	return errors.New("invalid first digit")
}

// validatePersonalCode validates personal codes (starting with 0-3)
func validatePersonalCode(val string) error {

	firstDigits, err := strconv.Atoi(val[0:2])
	if err != nil {
		return errors.New("invalid first two digits")
	}

	// For modern personal codes, prefix 3X (new format from 1 July 2017)
	// The second digit is random between 2-9, digits 3-11 are random 0-9
	if firstDigits >= 32 && firstDigits <= 39 {
		return nil
	}

	// For legacy personal codes, validate date components
	// Format: DDMMYYCXXXX where digits 1-2 are day, 3-4 are month, 5-6 are year,
	// digit 7 is century (0=18xx, 1=19xx, 2=20xx)
	day, err1 := strconv.Atoi(val[0:2])
	month, err2 := strconv.Atoi(val[2:4])
	year, err3 := strconv.Atoi(val[4:6])
	century, err4 := strconv.Atoi(string(val[6]))

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return errors.New("invalid personal code digits")
	}

	// Determine full year from century digit
	// Century digit: 0 = 18th century (1800-1899), 1 = 19th century (1900-1999), 2 = 21st century (2000-2099)
	var fullYear int
	switch century {
	case 0:
		fullYear = 1800 + year
	case 1:
		fullYear = 1900 + year
	case 2:
		fullYear = 2000 + year
	default:
		return errors.New("invalid century digit in personal code")
	}

	// Validate the date using cal package which handles proper date validation
	// including leap years and month-specific day limits
	d := cal.MakeDate(fullYear, time.Month(month), day)
	if !d.IsValid() {
		return errors.New("invalid birth date in personal code")
	}

	return nil
}

func validateMod11Checksum(val string) error {
	var sum int
	for i := 0; i < 10; i++ {
		digit := int(val[i] - '0')
		sum += taxCodeWeights[i] * digit
	}

	remainder := sum % 11
	checkDigit := 3 - remainder

	if checkDigit < 0 {
		checkDigit += 11
	}

	// If the calculated checksum is 10 or 11, the number is invalid
	if checkDigit >= 10 {
		return errors.New("calculated checksum results in an invalid corporate digit")
	}

	actualCheck := int(val[10] - '0')

	if checkDigit != actualCheck {
		return errors.New("checksum mismatch")
	}

	return nil
}
