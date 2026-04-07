package au

import (
	"errors"
	"strconv"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

const abnLength = 11

var abnWeights = []int{10, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19}

// validateTaxIdentity performs validation specific to Australian tax IDs.
func validateTaxIdentity(tID *tax.Identity) error {
	if tID == nil {
		return nil
	}
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code,
			validation.By(validateABN),
			validation.Skip,
		),
	)
}

// validateABN checks Australian Business Numbers (ABNs).
// Reference: https://abr.business.gov.au/Help/AbnFormat
func validateABN(value any) error {
	code, _ := value.(cbc.Code)
	normalized := strings.ReplaceAll(strings.ToUpper(code.String()), " ", "")
	if normalized == "" {
		return nil
	}
	if len(normalized) != abnLength {
		return errors.New("invalid length")
	}
	if _, err := strconv.Atoi(normalized); err != nil {
		return errors.New("invalid characters, expected numeric")
	}

	sum := 0
	for i, r := range normalized {
		digit := int(r - '0')
		if i == 0 {
			digit--
		}
		sum += digit * abnWeights[i]
	}
	if sum%89 != 0 {
		return errors.New("invalid checksum")
	}
	return nil
}
