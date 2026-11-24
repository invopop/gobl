package sk

import (
	"regexp"
	"strconv"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

/*
	The Slovak DIČ (Daňové identifikačné číslo) is 10 digits long. The first
digit gives the type of person, from 1 (natural person) to 7 (VAT group),
and the whole number is divisible by 11. The IČ DPH of a VAT payer is the
prefix SK followed by the same DIČ.

	Source: https://docs.peppol.eu/edelivery/codelists/v9.7/Peppol%20Code%20Lists%20-%20Participant%20identifier%20schemes%20v9.7.html
	(Participant identifier scheme 0245, SK:DIC)
*/

// We use a regex to reject anything outside 10 plain digits starting with 1-7.
// The SK prefix, spaces and hyphens are removed during normalization.
var taxCodeRegexp = regexp.MustCompile(`^[1-7]\d{9}$`)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "tax identity code for SK must be 10 digits starting with 1-7",
					is.MatchesRegexp(taxCodeRegexp),
				),
				rules.AssertIfPresent("02", "tax identity code for SK failed the checksum: the 10 digit number must be divisible by 11",
					is.StringFunc("checksum", validateTaxCodeChecksum),
				),
			),
		),
	)
}

func validateTaxCodeChecksum(code string) bool {
	// Each rule runs on its own, so a malformed code also reaches this check
	n, err := strconv.ParseInt(code, 10, 64)
	if err != nil {
		return false
	}

	// MOD 11 over the whole number
	return n%11 == 0
}
