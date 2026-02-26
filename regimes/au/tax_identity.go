package au

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Reference: https://www.ato.gov.au/businesses-and-organisations/preparing-lodging-and-paying/business-activity-statements-bass/instructions/checking-the-validity-of-an-abn
var regexpABN = regexp.MustCompile(`^[1-9]\d{10}$`)

// validateTaxIdentity checks to ensure the ABN format looks correct.
func validateTaxIdentity(tID *tax.Identity) error {
	if tID == nil {
		return nil
	}
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code,
			validation.By(validateTaxCode),
			validation.Skip,
		),
	)
}

func validateTaxCode(value any) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()
	if !regexpABN.MatchString(val) {
		return errors.New("invalid format")
	}
	if !validABNChecksum(val) {
		return errors.New("invalid checksum")
	}
	return nil
}

func validABNChecksum(abn string) bool {
	// ABN checksum: subtract 1 from first digit, multiply by weights,
	// sum, and ensure divisible by 89.
	weights := []int{10, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19}
	if len(abn) != len(weights) {
		return false
	}
	sum := 0
	for i := 0; i < len(weights); i++ {
		d := int(abn[i] - '0')
		if i == 0 {
			d--
		}
		sum += d * weights[i]
	}
	return sum%89 == 0
}
