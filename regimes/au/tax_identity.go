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
	if !regexpABN.MatchString(code.String()) {
		return errors.New("invalid format")
	}
	return nil
}
