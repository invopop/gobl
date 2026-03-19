package sg

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Reference: https://mytax.iras.gov.sg/ESVWeb/default.aspx?target=GSTListingSearch
// Reference: https://www.oecd.org/content/dam/oecd/en/topics/policy-issue-focus/aeoi/singapore-tin.pdf
// Reference:https://www.mof.gov.sg/docs/default-source/default-document-library/news-and-publications/press-releases/annexe060808.pdf?sfvrsn=4ee26b50_2
// Singapore’s tax authority does not publish a public checksum algorithm for UEN or GST numbers.
// Indeed, IRAS directs users to verify UENs via the official portal

// regexpsGSTCode uses the UEN identities as a base and adds the GST format used
// for international companies.
var regexpsGSTCode = append(
	regexpsUENIdentities,
	regexp.MustCompile(`^M[A-Z0-9]\d{7}[A-Z]$`),
)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn("SG"),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Singaporean tax identity code",
					is.Func("valid", isValidTaxIdentityCode),
				),
			),
		),
	)
}

func isValidTaxIdentityCode(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return false
	}
	return validateTaxCode(code) == nil
}

func validateTaxCode(code cbc.Code) error {
	if code == "" {
		return nil
	}
	val := code.String()
	match := false
	for _, re := range regexpsGSTCode {
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
