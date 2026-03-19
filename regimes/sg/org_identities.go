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
	regexpsUENIdentities = []*regexp.Regexp{
		regexp.MustCompile(`^(19[0-9]{2}|20[0-9]{2})\d{5}[A-Z]$`), // UEN (ROC)
		regexp.MustCompile(`^\d{8}[A-Z]$`),                        // UEN (ROB)
		regexp.MustCompile(`^[TS]\d{2}[A-Z]{2}\d{4}[A-Z]$`),       // UEN (Others)
	}
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
			is.HasContext(tax.RegimeIn(CountryCode)),
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
	val := code.String()
	match := false
	for _, re := range regexpsUENIdentities {
		if re.MatchString(val) {
			match = true
			break
		}
	}
	return match
}
