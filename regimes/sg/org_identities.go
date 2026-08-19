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
	regexpUENBusiness     = regexp.MustCompile(`^\d{8}[A-Z]$`)
	regexpUENLocalCompany = regexp.MustCompile(`^(19\d{2}|20\d{2})\d{5}[A-Z]$`)
	regexpUENOther        = regexp.MustCompile(`^[RST]\d{2}[A-Z]{2}\d{4}[A-Z]$`)

	// Entity types assigned to UENs for bodies other than businesses and
	// local companies. See https://www.uen.gov.sg/ueninternet/faces/pages/admin/aboutUEN.jspx.
	uenOtherEntityTypes = map[string]struct{}{
		"CC": {}, "CD": {}, "CH": {}, "CL": {}, "CM": {}, "CP": {}, "CS": {}, "CX": {},
		"DP": {}, "FB": {}, "FC": {}, "FM": {}, "FN": {}, "GA": {}, "GB": {}, "GS": {},
		"HS": {}, "LL": {}, "LP": {}, "MB": {}, "MC": {}, "MD": {}, "MH": {}, "MM": {},
		"MQ": {}, "NB": {}, "NR": {}, "PA": {}, "PB": {}, "PF": {}, "RF": {}, "RP": {},
		"SM": {}, "SS": {}, "TC": {}, "TU": {}, "VH": {}, "XL": {},
	}
)

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

// validateUENCode verifies a UEN's shape and its check character. The check
// character tables are the community reverse-engineered versions used by
// python-stdnum; ACRA does not publish the algorithm.
func validateUENCode(code string) bool {
	switch {
	case regexpUENBusiness.MatchString(code):
		return code[8] == uenDigitCheck(code, uenBusinessWeights, uenBusinessCheckAlphabet)
	case regexpUENLocalCompany.MatchString(code):
		return code[9] == uenDigitCheck(code, uenLocalCompanyWeights, uenLocalCompanyCheckAlphabet)
	case regexpUENOther.MatchString(code):
		_, knownType := uenOtherEntityTypes[code[3:5]]
		return knownType && code[9] == uenOtherCheck(code)
	default:
		return false
	}
}

func uenDigitCheck(code string, weights []int, alphabet string) byte {
	sum := 0
	for i, weight := range weights {
		sum += int(code[i]-'0') * weight
	}
	return alphabet[sum%len(alphabet)]
}

func uenOtherCheck(code string) byte {
	sum := 0
	for i, weight := range uenOtherWeights {
		sum += strings.IndexByte(uenOtherCheckAlphabet, code[i]) * weight
	}
	return uenOtherCheckAlphabet[(sum-5)%11]
}
