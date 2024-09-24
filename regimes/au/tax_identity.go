package au

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var taxWeightTableABN = [11]int{10, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19}

// var taxWeightTableACN = [8]int{8, 7, 6, 5, 4, 3, 2, 1}

func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateTaxCode)),
	)
}

func validateTaxCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return errors.New("invalid format")
	}
	val := code.String()
	reABN := regexp.MustCompile(`^\d{11}$`)
	// reACN := regexp.MustCompile(`^\d{9}$`)

	// If 11 digits long check ABN, if 9 digits long ACN
	switch {
	case reABN.MatchString(val):
		return businessNumberCheck(val)
	// case reACN.MatchString(val):
	// 	return companyNumberCheck(val)
	default:
		return errors.New("invalid format")
	}
}

func businessNumberCheck(val string) error {
	if z, _ := strconv.Atoi(val); z == 0 {
		return errors.New("zeros")
	}

	// Source: https://abr.business.gov.au/Help/AbnFormat
	firstDigit, err := strconv.Atoi(string(val[0]))
	if err != nil {
		return errors.New("invalid format")
	}
	firstDigit -= 1
	modifiedABN := strconv.Itoa(firstDigit) + val[1:]
	sum := 0
	for i := 0; i < 11; i++ {
		digit, err := strconv.Atoi(string(modifiedABN[i]))
		if err != nil {
			return errors.New("invalid format")
		}
		sum += digit * taxWeightTableABN[i]
	}
	if sum%89 == 0 {
		return nil
	}
	return errors.New("checksum mismatch")
}

// func companyNumberCheck(val string) error {
// 	if z, _ := strconv.Atoi(val); z == 0 {
// 		return errors.New("zeros")
// 	}

// 	// Source: https://asic.gov.au/for-business/registering-a-company/steps-to-register-a-company/australian-company-numbers/australian-company-number-digit-check
// 	checkDigit, err := strconv.Atoi(string(val[8]))
// 	if err != nil {
// 		return errors.New("invalid format")
// 	}
// 	sum := 0
// 	for i := 0; i < 8; i++ {
// 		digit, err := strconv.Atoi(string(val[i]))
// 		if err != nil {
// 			return errors.New("invalid format")
// 		}
// 		sum += digit * taxWeightTableACN[i]
// 	}

// 	remainder := sum % 10
// 	if (10-remainder == checkDigit) || (remainder == 0 && checkDigit == 0) {
// 		return nil
// 	}
// 	return errors.New("checksum mismatch")
// }
