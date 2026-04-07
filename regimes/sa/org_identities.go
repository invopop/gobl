package sa

import (
	"regexp"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

var (
	tinRegexp   = regexp.MustCompile(`^\d{15}$`)
	otherRegexp = regexp.MustCompile(`^\d{10}$`)
)

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(countryCode)),
			rules.When(
				is.Or(
					org.IdentitiesTypeIn(IdentityTypeCRN),
					org.IdentitiesTypeIn(IdentityTypeMom),
					org.IdentitiesTypeIn(IdentityTypeMLS),
					org.IdentitiesTypeIn(IdentityType700),
					org.IdentitiesTypeIn(IdentityTypeSAG),
					org.IdentitiesTypeIn(IdentityTypeNational),
					org.IdentitiesTypeIn(IdentityTypeGcc),
					org.IdentitiesTypeIn(IdentityTypeIqa),
					org.IdentitiesTypeIn(IdentityTypePassport),
					org.IdentitiesTypeIn(IdentityTypeOTH),
				),
				rules.Field("code",
					rules.Assert("01", "identity code must be valid",
						is.MatchesRegexp(otherRegexp),
					),
				),
			),
			rules.When(
				org.IdentitiesTypeIn(IdentityTypeTIN),
				rules.Field("code",
					rules.Assert("02", "identity code for type TIN must be valid",
						is.MatchesRegexp(tinRegexp),
					),
				),
			),
		),
	)
}
