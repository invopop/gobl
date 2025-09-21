package au

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var (
	abnRegex = regexp.MustCompile(`^\d{11}$`)
)

func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateABN)),
	)
}

func validateABN(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	if !abnRegex.MatchString(val) {
		return errors.New("must be 11 digits")
	}

	return validateABNChecksum(val)
}

// Source: https://www.ato.gov.au/businesses-and-organisations/hiring-and-paying-your-workers/payg-withholding/payments-you-need-to-withhold-from/withholding-from-suppliers/checking-an-abn
// Source: https://abr.business.gov.au/Help/AbnFormat
func validateABNChecksum(abn string) error {
	weights := []int{10, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19}

	digits := make([]int, 11)
	for i, char := range abn {
		digit, _ := strconv.Atoi(string(char))
		if i == 0 {
			digit--
		}
		digits[i] = digit
	}

	sum := 0
	for i := 0; i < 11; i++ {
		sum += digits[i] * weights[i]
	}

	if sum%89 != 0 {
		return errors.New("invalid ABN checksum")
	}

	return nil
}
