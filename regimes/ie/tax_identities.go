package ie

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Irish VAT numbers can be in two formats:
// Old format: 1 digit, 1 letter/symbol (+, *), 5 digits, 1 letter (e.g., 1A12345B)
// New format: 7 digits, 1-2 letters (e.g., 1234567FA or 1234567F)

var (
	// Old format: 1 digit, 1 letter/symbol (+, *), 5 digits, 1 letter
	taxCodeOldFormatRegexp = regexp.MustCompile(`^(\d)([A-Z+*])(\d{5})([A-Z])$`)
	// New format: 7 digits, 1 or 2 letters
	taxCodeNewFormatRegexp = regexp.MustCompile(`^(\d{7})([A-W][AH]?)$`)
	// Check characters for modulus 23
	taxCodeCheckChars = "WABCDEFGHIJKLMNOPQRSTUV"
)

// validateTaxIdentity performs validation specific to Irish tax IDs.
// Assumes the code has already been normalized.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code,
			validation.By(validateTaxCode),
		),
	)
}

// validateTaxCode validates the tax code for Irish tax identities.
// Assumes the code has already been normalized.
func validateTaxCode(value any) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}

	val := code.String()

	// Check if it matches the old format
	if matches := taxCodeOldFormatRegexp.FindStringSubmatch(val); matches != nil {
		return validateOldFormat(matches)
	}

	// Check if it matches the new format
	if matches := taxCodeNewFormatRegexp.FindStringSubmatch(val); matches != nil {
		return validateNewFormat(matches)
	}

	return errors.New("invalid format")
}

// validateOldFormat validates Irish VAT numbers in old format: 1 digit, 1 letter/symbol (+, *), 5 digits, 1 letter
// Formula: r = (8*0 + 7*c3 + 6*c4 + 5*c5 + 4*c6 + 3*c7 + 2*c1) % 23
func validateOldFormat(matches []string) error {
	if len(matches) != 5 {
		return errors.New("invalid old format structure")
	}

	firstDigit := matches[1]
	// matches[2] is the letter/symbol (not used in checksum)
	fiveDigits := matches[3]
	checkChar := matches[4]

	// Calculate checksum: 7*c3 + 6*c4 + 5*c5 + 4*c6 + 3*c7 + 2*c1
	// where c1=firstDigit and c3-c7 are the five digits
	c1, _ := strconv.Atoi(firstDigit)
	c3, _ := strconv.Atoi(string(fiveDigits[0]))
	c4, _ := strconv.Atoi(string(fiveDigits[1]))
	c5, _ := strconv.Atoi(string(fiveDigits[2]))
	c6, _ := strconv.Atoi(string(fiveDigits[3]))
	c7, _ := strconv.Atoi(string(fiveDigits[4]))

	sum := 7*c3 + 6*c4 + 5*c5 + 4*c6 + 3*c7 + 2*c1
	r := sum % 23

	expectedChar := string(taxCodeCheckChars[r])
	if checkChar != expectedChar {
		return errors.New("checksum mismatch")
	}

	return nil
}

// validateNewFormat validates Irish VAT numbers in new format: 7 digits, 1-2 letters
// Formula: r = (8*c1 + 7*c2 + 6*c3 + 5*c4 + 4*c5 + 3*c6 + 2*c7 + extra) % 23
// where extra is 9 for 'A' or 72 for 'H' as the 9th character
// The 8th character (first letter) is the check character
func validateNewFormat(matches []string) error {
	if len(matches) != 3 {
		return errors.New("invalid new format structure")
	}

	sevenDigits := matches[1]
	letters := matches[2]

	// The first letter is the check character
	checkChar := string(letters[0])

	// Calculate weighted sum using multipliers [8,7,6,5,4,3,2]
	c1, _ := strconv.Atoi(string(sevenDigits[0]))
	c2, _ := strconv.Atoi(string(sevenDigits[1]))
	c3, _ := strconv.Atoi(string(sevenDigits[2]))
	c4, _ := strconv.Atoi(string(sevenDigits[3]))
	c5, _ := strconv.Atoi(string(sevenDigits[4]))
	c6, _ := strconv.Atoi(string(sevenDigits[5]))
	c7, _ := strconv.Atoi(string(sevenDigits[6]))

	sum := 8*c1 + 7*c2 + 6*c3 + 5*c4 + 4*c5 + 3*c6 + 2*c7

	// For 9-character format, add extra based on the 9th character
	if len(letters) == 2 {
		lastChar := letters[1]
		switch lastChar {
		case 'A':
			sum += 9
		case 'H':
			sum += 72
		}
	}

	r := sum % 23

	expectedChar := string(taxCodeCheckChars[r])
	if checkChar != expectedChar {
		return errors.New("checksum mismatch")
	}

	return nil
}
