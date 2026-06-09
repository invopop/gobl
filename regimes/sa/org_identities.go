package sa

import (
	"regexp"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

var (
	identitiesRegexp = regexp.MustCompile(`^[a-zA-Z0-9]*$`)
)

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.Field("code",
				rules.Assert("01", "SA identity code must be valid",
					is.MatchesRegexp(identitiesRegexp),
				),
			),
		),
	)
}
