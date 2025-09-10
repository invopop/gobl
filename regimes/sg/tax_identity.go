package sg

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Reference: https://mytax.iras.gov.sg/ESVWeb/default.aspx?target=GSTListingSearch
// Reference: https://www.oecd.org/content/dam/oecd/en/topics/policy-issue-focus/aeoi/singapore-tin.pdf
// Reference:https://www.mof.gov.sg/docs/default-source/default-document-library/news-and-publications/press-releases/annexe060808.pdf?sfvrsn=4ee26b50_2
// Singapore’s tax authority does not publish a public checksum algorithm for UEN or GST numbers.
// Indeed, IRAS directs users to verify UENs via the official portal

var (
	taxCodeRegexps = []*regexp.Regexp{
		regexp.MustCompile(`^(19[0-9]{2}|20[0-9]{2})\d{5}[A-Z]$`), // UEN (ROC)
		regexp.MustCompile(`^\d{8}[A-Z]$`),                        // UEN (ROB)
		regexp.MustCompile(`^[TS]\d{2}[A-Z]{2}\d{4}[A-Z]$`),       // UEN (Others)
	}
)

// validateTaxIdentity checks to ensure the NIT code looks okay.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code,
			validation.By(validateTaxCode),
			validation.Skip,
		),
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
