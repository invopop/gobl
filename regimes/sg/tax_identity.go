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

// RegexGSTCode is the regex pattern for the GST registration number.
var RegexGSTCode = regexp.MustCompile(`^[M][A-Z0-9]\d{7}[A-Z]$`)

// validateTaxIdentity checks to ensure the NIT code looks okay.
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
	if !RegexGSTCode.MatchString(val) {
		return errors.New("invalid format")
	}
	return nil
}
