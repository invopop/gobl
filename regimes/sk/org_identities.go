package sk

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// IdentityKeyICO is the IČO (Identifikačné číslo organizácie), the number that
	// the Statistical Office assigns to every legal entity and entrepreneur. It
	// goes on commercial documents when assigned (§ 3a Commercial Code) and is
	// the legal registration identifier of the Peppol BIS (BT-30 and BT-47).
	IdentityKeyICO cbc.Key = "sk-ico"
)

// icoRegexp matches a normalized IČO: exactly 8 digits.
var icoRegexp = regexp.MustCompile(`^\d{8}$`)

// nonDigitRegexp matches everything that is not a digit.
var nonDigitRegexp = regexp.MustCompile(`\D`)

var identityDefinitions = []*cbc.Definition{
	{
		Key: IdentityKeyICO,
		Name: i18n.String{
			i18n.EN: "Organisation Identification Number (IČO)",
			i18n.SK: "Identifikačné číslo organizácie (IČO)",
		},
	},
}

// normalizeICO removes spaces and other separators from an IČO, which is often
// written in groups of three digits such as "36 396 567". The IČO is the legal
// registration identifier of the party, so it gets the legal scope that formats
// such as UBL use for BT-30 and BT-47.
func normalizeICO(id *org.Identity) {
	if id == nil || id.Key != IdentityKeyICO {
		return
	}
	id.Code = cbc.Code(nonDigitRegexp.ReplaceAllString(id.Code.String(), ""))
	if id.Scope == cbc.KeyEmpty {
		id.Scope = org.IdentityScopeLegal
	}
}

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				org.IdentitiesKeyIn(IdentityKeyICO),
				rules.Field("code",
					rules.Assert("01", "identity code for key sk-ico (IČO) must be 8 digits",
						is.MatchesRegexp(icoRegexp),
					),
				),
			),
		),
	)
}
