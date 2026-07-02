package is

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	isrules "github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// kennitalaRegexp matches the 10-digit Icelandic Kennitala. The hyphen commonly
// displayed after the sixth digit is stripped during normalization.
var kennitalaRegexp = regexp.MustCompile(`^[0-9]{10}$`)

// kennitalaWeights are the mod-11 multipliers applied to positions 0-7 of a
// Kennitala.
//
// A Kennitala is a 10-digit identifier of the form DDMMYY-RRCS: DDMMYY is the
// person's date of birth (or the entity's registration date), RR is a random
// pair assigned at issuance, C is the mod-11 check digit and S is a century
// marker (8 = 1800s, 9 = 1900s, 0 = 2000s). Þjóðskrá Íslands (Registers
// Iceland) issues Kennitalas to individuals; Skatturinn's company register
// (Fyrirtækjaskrá) issues them to legal entities, with 4 added to D0 to avoid
// overlap (e.g. a company registered on the 15th begins with 55). The same
// 10-digit identifier serves as the tax identity and, for VAT-registered
// businesses, the VSK number.
//
// See: https://lookuptax.com/docs/tax-identification-number/iceland-tax-id-guide
var kennitalaWeights = []int{3, 2, 7, 6, 5, 4, 3, 2}

// normalizeTaxIdentity strips whitespace, separators and the country prefix.
// The standard tax.NormalizeIdentity is sufficient for Icelandic Kennitalas.
func normalizeTaxIdentity(tID *tax.Identity) {
	tax.NormalizeIdentity(tID)
}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "tax id code must be a 10-digit Kennitala",
					isrules.MatchesRegexp(kennitalaRegexp)),
				rules.AssertIfPresent("02", "tax id code Kennitala checksum is invalid",
					isrules.Func("valid", isValidKennitalaCode)),
			),
		),
	)
}

// isValidKennitalaCode reports whether value is a Kennitala with a valid mod-11
// check digit. Unexpected types return false rather than silently passing.
func isValidKennitalaCode(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok {
		return false
	}
	return validKennitala(code.String())
}

// validKennitala applies the Icelandic mod-11 check-digit algorithm. The
// remainder == 1 case is rejected because the registry never issues
// Kennitalas whose calculated check digit would be a two-digit value.
func validKennitala(s string) bool {
	if len(s) != 10 {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	sum := 0
	for i, w := range kennitalaWeights {
		d := int(s[i] - '0')
		sum += d * w
	}
	rem := sum % 11
	var check int
	switch rem {
	case 0:
		check = 0
	case 1:
		return false
	default:
		check = 11 - rem
	}
	return int(s[8]-'0') == check
}
