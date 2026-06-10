package fr

import (
	"regexp"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

var (
	sirenRegexp = regexp.MustCompile(`^\d{9}$`)
	siretRegexp = regexp.MustCompile(`^\d{14}$`)
)

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				org.IdentitiesTypeIn(IdentityTypeSIREN),
				rules.Field("code",
					rules.Assert("01", "identity code for type SIREN must be valid",
						is.MatchesRegexp(sirenRegexp),
					),
				),
			),
			rules.When(
				org.IdentitiesTypeIn(IdentityTypeSIRET),
				rules.Field("code",
					rules.Assert("02", "identity code for type SIRET must be valid",
						is.MatchesRegexp(siretRegexp),
					),
				),
			),
		),
	)
}
