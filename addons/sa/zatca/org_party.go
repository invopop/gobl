package zatca

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

func orgAddressRules() *rules.Set {
	return rules.For(new(org.Address),
		rules.Field("country",
			rules.Assert("01", "address must have a country code", is.Present),
		),
		rules.When(
			is.Func("country code is SA", countryCodeIsSA),
			rules.Field("street",
				rules.Assert("02", "address in SA must have a street name (BR-KSA-09), (BR-KSA-63)", is.Present),
			),
			rules.Field("num",
				rules.Assert("03", "address in SA must have a 4 digits building number (BR-KSA-09), (BR-KSA-63), (BR-KSA-37)", is.Present, is.Matches(`^\d{4}$`)),
			),
			rules.Field("code",
				rules.Assert("04", "address in SA must have a 5 digits postal code (BR-KSA-09), (BR-KSA-63) (BR-KSA-67), (BR-KSA-66)", is.Present, is.Matches(`^\d{5}$`)),
			),
			rules.Field("locality",
				rules.Assert("05", "address in SA must have a city name (BR-KSA-09), (BR-KSA-63)", is.Present),
			),
			// mapped to district
			rules.Field("street_extra",
				rules.Assert("06", "address in SA must have a district name (BR-KSA-09), (BR-KSA-63)", is.Present),
			),
		),
	)
}

func countryCodeIsSA(val any) bool {
	address, ok := val.(*org.Address)
	return ok && address != nil && address.Country == l10n.ISOCountryCode("SA")
}
