package sg

import (
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Reference: https://mytax.iras.gov.sg/ESVWeb/default.aspx?target=GSTListingSearch

const (
	// IdentityTypeUEN represents the Unique Entity Number used in Singapore.
	IdentityTypeUEN cbc.Code = "UEN"
)

var (
	regexpUENLocalCompany = regexp.MustCompile(`^(19[0-9]{2}|20[0-9]{2})\d{5}[A-Z]$`) // UEN (ROC)
	regexpUENBusiness     = regexp.MustCompile(`^\d{8}[A-Z]$`)                        // UEN (ROB)
	regexpUENOther        = regexp.MustCompile(`^[TS]\d{2}[A-Z]{2}\d{4}[A-Z]$`)       // UEN (Others)
)

// UEN check character tables. ACRA does not publish the algorithms
// officially; these are the community reverse-engineered versions also used
// by python-stdnum and verified against real registered entities. The final
// character of every UEN is a check character over the preceding ones.
const (
	uenBusinessCheckAlphabet     = "XMKECAWLJDB"
	uenLocalCompanyCheckAlphabet = "ZKCMDNERGWH"
	uenOtherCheckAlphabet        = "ABCDEFGHJKLMNPQRSTUVWX0123456789"
)

var (
	uenBusinessWeights     = []int{10, 4, 9, 3, 8, 2, 7, 1}
	uenLocalCompanyWeights = []int{10, 8, 6, 4, 9, 7, 5, 3, 1}
	uenOtherWeights        = []int{4, 3, 5, 3, 10, 2, 2, 5, 7}
)

var identityDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeUEN,
		Name: i18n.String{
			i18n.EN: "Unique Entity Number (UEN)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The Unique Entity Number (UEN) is a standard identification number
				issued to entities (businesses and local companies) in Singapore.
				It is used for all interactions with government agencies. The same
				UEN is normally used as the GST registration number, but not always.
			`),
		},
	},
}

func normalizeIdentity(id *org.Identity) {
	if id == nil || id.Type != IdentityTypeUEN {
		return
	}
	code := strings.ToUpper(id.Code.String())
	code = tax.IdentityCodeBadCharsRegexp.ReplaceAllString(code, "")
	id.Code = cbc.Code(strings.TrimPrefix(code, string(l10n.SG)))
}

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				org.IdentityTypeIn(IdentityTypeUEN),
				rules.Field("code",
					rules.Assert("01", "identity code for type UEN must be valid",
						is.Func("valid UEN", orgIdentityCheckUEN),
					),
				),
			),
		),
	)
}

func orgIdentityCheckUEN(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return false
	}
	return validateUENCode(code.String())
}

// validateUENCode determines the UEN sub-format from the code's shape and
// verifies both the format and the trailing check character.
func validateUENCode(val string) bool {
	switch {
	case regexpUENLocalCompany.MatchString(val):
		return uenDigitChecksum(val, uenLocalCompanyWeights, uenLocalCompanyCheckAlphabet)
	case regexpUENBusiness.MatchString(val):
		return uenDigitChecksum(val, uenBusinessWeights, uenBusinessCheckAlphabet)
	case regexpUENOther.MatchString(val):
		return uenOtherChecksum(val)
	}
	return false
}

// uenDigitChecksum verifies the check letter of the all-digit UEN formats
// (ROB and ROC): a weighted sum of the digits, mod 11, indexes into the
// format's check alphabet.
func uenDigitChecksum(val string, weights []int, alphabet string) bool {
	sum := 0
	for i, w := range weights {
		sum += int(val[i]-'0') * w
	}
	return val[len(val)-1] == alphabet[sum%11]
}

// uenOtherChecksum verifies the check character of the "Others" UEN format
// ([TS]yyPQnnnnX). Characters map to positions in a 32-character alphabet,
// are combined in a weighted sum, and (sum - 5) mod 11 indexes back into
// the same alphabet.
func uenOtherChecksum(val string) bool {
	sum := 0
	for i, w := range uenOtherWeights {
		sum += strings.IndexByte(uenOtherCheckAlphabet, val[i]) * w
	}
	idx := ((sum-5)%11 + 11) % 11
	return val[len(val)-1] == uenOtherCheckAlphabet[idx]
}
