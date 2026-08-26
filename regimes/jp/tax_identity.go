package jp

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Japanese qualified-invoice registration number (登録番号): the letter "T"
// followed by 13 digits. For a corporation this is "T" + its 13-digit Corporate
// Number (法人番号); for a sole proprietor or other registrant without a Corporate
// Number, the NTA issues a distinct 13-digit number.
//
// References:
//   - Registration number format:
//     https://www.invoice-kohyo.nta.go.jp/about-toroku/index.html
//   - Public lookup site (for definitive verification):
//     https://www.invoice-kohyo.nta.go.jp/
var taxCodeRegexp = regexp.MustCompile(`^T[0-9]{13}$`)

var errInvalidFormat = errors.New("invalid format")

// Note on checksum validation: the Corporate Number carries a documented mod-9
// check digit (see https://www.houjin-bangou.nta.go.jp/documents/checkdigit.pdf),
// so "T" + a Corporate Number could in principle be checksum-verified. GOBL does
// NOT enforce it here, for two reasons:
//
//  1. The NTA check-digit specification is scoped to the Corporate Number system.
//     No official source extends it to the qualified-invoice registration-number
//     space as a whole, which also includes numbers the NTA issues to registrants
//     without a Corporate Number. Enforcing a checksum requires positive
//     confirmation that every code in the space satisfies it; that confirmation
//     does not exist. (A shared algorithm would also give no collision-avoidance
//     benefit - disjoint ranges do that - so the shared-algorithm assumption has
//     no design rationale behind it either.)
//  2. A corporate and a non-corporate registration number are indistinguishable by
//     shape (both are "T" + 13 digits), so the checksum cannot be applied
//     conditionally.
//
// In a validation library the failure modes are asymmetric: enforcing a checksum
// that some valid numbers do not satisfy rejects real invoices, whereas
// validating format only merely fails to catch a mistyped digit. We therefore
// validate the format and leave definitive verification to the NTA public lookup
// site above.

func normalizeTaxIdentity(tID *tax.Identity) {
	tax.NormalizeIdentity(tID)
}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Japanese registration number",
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
	if !taxCodeRegexp.MatchString(code.String()) {
		return errInvalidFormat
	}
	return nil
}
