package sk

import (
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// IdentityTypeICO is the IČO (Identifikačné číslo organizácie), the number that
	// the Statistical Office assigns to every legal entity and entrepreneur. It
	// goes on commercial documents when assigned (§ 3a Commercial Code) and is
	// the legal registration identifier of the Peppol BIS (BT-30 and BT-47).
	IdentityTypeICO cbc.Code = "ICO"
)

// icoRegexp matches a normalized IČO: exactly 8 digits.
var icoRegexp = regexp.MustCompile(`^\d{8}$`)

// icoSeparators are the characters an IČO may be written with, such as the
// spaces of "36 396 567". Any other character stays and fails the format rule.
var icoSeparators = strings.NewReplacer(" ", "", "\u00a0", "", ".", "", "-", "", "/", "")

var identityTypeDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeICO,
		Name: i18n.String{
			i18n.EN: "Organisation Identification Number (IČO)",
			i18n.SK: "Identifikačné číslo organizácie (IČO)",
		},
	},
}

// normalizeICO removes the separators from an IČO, which is often written in
// groups of three digits such as "36 396 567". The IČO is the legal
// registration identifier of the party, so it gets the legal scope that formats
// such as UBL use for BT-30 and BT-47.
func normalizeICO(id *org.Identity) {
	if id == nil || id.Type != IdentityTypeICO {
		return
	}
	id.Code = cbc.Code(icoSeparators.Replace(id.Code.String()))
	if id.Scope == cbc.KeyEmpty {
		id.Scope = org.IdentityScopeLegal
	}
}

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				org.IdentityTypeIn(IdentityTypeICO),
				rules.Field("code",
					rules.Assert("01", "identity code for type ICO (IČO) must be 8 digits",
						is.MatchesRegexp(icoRegexp),
					),
				),
			),
		),
	)
}
