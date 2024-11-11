package au

import (
	"errors"
	"strconv"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

const (
	// IdentityCompanyNumber is the key for the Australian Company Number, sometimes used instead of the ABN
	IdentityCompanyNumber cbc.Key = "ACN"
)

// Weights for ACN checksum
var taxWeightTableACN = [8]int{8, 7, 6, 5, 4, 3, 2, 1}

var identityKeyDefinitions = []*cbc.KeyDefinition{
	{
		Key: IdentityCompanyNumber,
		Name: i18n.String{
			i18n.EN: "Australian Company Number",
		},
	},
}

func normalizeIdentity(id *org.Identity) {
	if id == nil || id.Key != IdentityCompanyNumber {
		return
	}
	code := strings.ToUpper(id.Code.String())
	code = tax.IdentityCodeBadCharsRegexp.ReplaceAllString(code, "")
	id.Code = cbc.Code(code)
}

// validateIdentitiy helps confirm that an identity of a specific type is valid.
func validateIdentity(id *org.Identity) error {
	if id == nil || id.Key != IdentityCompanyNumber {
		return nil
	}
	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.Required,
			validation.By(validateCompanyNumber),
			validation.Skip,
		),
	)
}

// Source: https://asic.gov.au/for-business/registering-a-company/steps-to-register-a-company/australian-company-numbers/australian-company-number-digit-check
func validateCompanyNumber(value interface{}) error {
	val, ok := value.(cbc.Code)
	if !ok || val == cbc.CodeEmpty {
		return nil
	}
	code := val.String()
	if z, _ := strconv.Atoi(code); z == 0 {
		return errors.New("invalid format")
	}
	if len(code) != 9 {
		return errors.New("invalid format")
	}
	checkDigit, err := strconv.Atoi(string(val[8]))
	if err != nil {
		return errors.New("invalid format")
	}
	sum := 0
	for i := 0; i < 8; i++ {
		digit, err := strconv.Atoi(string(val[i]))
		if err != nil {
			return errors.New("invalid format")
		}
		sum += digit * taxWeightTableACN[i]
	}
	remainder := sum % 10
	if (10-remainder == checkDigit) || (remainder == 0 && checkDigit == 0) {
		return nil
	}
	return errors.New("checksum mismatch")
}
