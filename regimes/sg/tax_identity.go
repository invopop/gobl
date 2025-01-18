package sg

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Reference: https://lookuptax.com/docs/tax-identification-number/singapore-tax-id-guide#nric-number
// Reference: https://mytax.iras.gov.sg/ESVWeb/default.aspx?target=GSTListingSearch

var (
	taxCodeRegexps = []*regexp.Regexp{
		regexp.MustCompile(`^(19[0-9]{2}|20[0-9]{2})\d{5}[A-Z]$`), // UEN (ROC)
		regexp.MustCompile(`^\d{9}[A-Z]$`),                        // UEN (ROB)
		regexp.MustCompile(`^[TS]\d{2}[A-Z]\w{1}\d{4}[A-Z]$`),     // UEN (Others)
		regexp.MustCompile(`^[STFGM]\d{7}[A-Z]$`),                 // NIRC/FIN
		regexp.MustCompile(`^\w{2}\d{7}[A-Z]$`),                   // GST

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

	return nil
}
