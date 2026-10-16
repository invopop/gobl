package is

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	rulesis "github.com/invopop/gobl/rules/is" // aliased: rules/is collides with this package's name
	"github.com/invopop/gobl/tax"
)

// IdentityTypeKennitala identifies the Icelandic national identification number
// (kennitala), issued by Registers Iceland (Þjóðskrá Íslands) to both natural
// persons and companies.
const IdentityTypeKennitala cbc.Code = "KT"

var identityTypeDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeKennitala,
		Name: i18n.String{
			i18n.EN: "Kennitala",
			i18n.IS: "Kennitala",
		},
		Desc: i18n.String{
			i18n.EN: "Icelandic national identification number for persons and companies.",
			i18n.IS: "Íslensk kennitala fyrir einstaklinga og fyrirtæki.",
		},
	},
}

// normalizeOrgIdentity strips separators and whitespace from a kennitala, leaving
// the bare 10-digit form. If the result is not 10 digits it leaves the value
// untouched, so validation can report the original malformed code.
func normalizeOrgIdentity(id *org.Identity) {
	if id == nil || id.Type != IdentityTypeKennitala {
		return
	}
	code := cbc.NormalizeNumericalCode(id.Code)
	if len(code) != kennitalaLength {
		return
	}
	id.Code = code
}

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			rulesis.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				org.IdentityTypeIn(IdentityTypeKennitala),
				rules.Field("code",
					rules.Assert("01", "invalid kennitala format",
						rulesis.Func("well-formed", orgKennitalaWellFormed),
					),
					rules.Assert("02", "invalid checksum for kennitala",
						rulesis.Func("checksum", orgKennitalaChecksumValid),
					),
				),
			),
		),
	)
}

func orgKennitalaWellFormed(value any) bool {
	code, ok := value.(cbc.Code)
	return ok && isCanonicalKennitala(code)
}

func orgKennitalaChecksumValid(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || !isCanonicalKennitala(code) {
		return true // skip if not well-formed; the format assertion reports that
	}
	return ValidKennitala(code)
}
