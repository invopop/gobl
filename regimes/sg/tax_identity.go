package sg

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Reference: https://mytax.iras.gov.sg/ESVWeb/default.aspx?target=GSTListingSearch
// Reference: https://www.oecd.org/content/dam/oecd/en/topics/policy-issue-focus/aeoi/singapore-tin.pdf
// Reference:https://www.mof.gov.sg/docs/default-source/default-document-library/news-and-publications/press-releases/annexe060808.pdf?sfvrsn=4ee26b50_2
// IRAS-assigned GST registration numbers have no public checksum algorithm.
// The final character may be either a letter or a digit.
var regexpGSTNumber = regexp.MustCompile(`^M[A-Z0-9]\d{7}[A-Z0-9]$`)

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
	text := code.String()
	return regexpGSTNumber.MatchString(text) || validateUENCode(text)
}
