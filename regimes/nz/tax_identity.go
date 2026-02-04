package nz

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var (
	irdPattern          = regexp.MustCompile(`^\d{8,9}$`)
	irdPrimaryWeights   = []int{3, 2, 7, 6, 5, 4, 3, 2}
	irdSecondaryWeights = []int{7, 4, 3, 2, 5, 2, 7, 6}
)

func normalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil || tID.Code == "" {
		return
	}
	code := tID.Code.String()
	code = strings.ReplaceAll(code, "-", "")
	code = strings.ReplaceAll(code, " ", "")
	tID.Code = cbc.Code(code)
}

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

	str := strings.ReplaceAll(cbc.NormalizeString(code.String()), "-", "")

	if !irdPattern.MatchString(str) {
		return errors.New("invalid tax ID format: must be 8-9 digit IRD number")
	}

	return validateIRD(str)
}

func validateIRD(code string) error {
	if len(code) == 8 {
		code = "0" + code
	}

	num, err := strconv.Atoi(code)
	if err != nil {
		return errors.New("IRD number must contain only digits")
	}

	if num < 10000000 || num > 150000000 {
		return errors.New("IRD number out of valid range (10,000,000 to 150,000,000)")
	}

	baseDigits := make([]int, 8)
	for i := 0; i < 8; i++ {
		baseDigits[i], _ = strconv.Atoi(string(code[i]))
	}

	checkDigit := calculateIRDCheckDigit(baseDigits, irdPrimaryWeights)

	if checkDigit == 10 {
		checkDigit = calculateIRDCheckDigit(baseDigits, irdSecondaryWeights)
		if checkDigit == 10 {
			return errors.New("invalid IRD number: check digit calculation failed")
		}
	}

	actualCheckDigit, _ := strconv.Atoi(string(code[8]))
	if checkDigit != actualCheckDigit {
		return errors.New("invalid IRD number: check digit mismatch")
	}

	return nil
}

func calculateIRDCheckDigit(digits []int, weights []int) int {
	sum := 0
	for i := 0; i < 8; i++ {
		sum += digits[i] * weights[i]
	}

	remainder := sum % 11
	if remainder == 0 {
		return 0
	}
	return 11 - remainder
}
