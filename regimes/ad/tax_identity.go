package ad

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// NRT (Número de Registre Tributari) is the tax identifier for all taxpayers
// in Andorra — individuals and legal entities, resident and non-resident.
//
// Format: one letter indicating entity type, six digits, one check letter.
// Written with or without hyphens: L-132950-X or L132950X (both valid input).
//
// Known entity-type prefix letters:
//   - F: resident natural persons (NIA prefixed with F)
//   - E: non-resident natural persons
//   - A: joint-stock companies (Societat Anònima, SA)
//   - L: limited liability companies (Societat Limitada, SL)
//   - C, D, G, O, P, U: other entity types (cooperatives, foundations,
//     public bodies, special-purpose entities)
//
// The check-letter algorithm is not publicly documented by the Andorran
// tax authority (Departament de Tributs i de Fronteres). Validation here
// is format-only; authoritative verification requires the official portal.
//
// References:
//   - https://www.oecd.org/content/dam/oecd/en/topics/policy-issue-focus/aeoi/andorra-tin.pdf

// reNRT validates the normalised NRT format: one uppercase letter,
// six digits, one uppercase letter.

var reNRT = regexp.MustCompile(`^[A-Z]\d{6}[A-Z]$`) // used [A-Z] for the prefix rather than [FEACDGLOPU] 
// deliberately. The OECD document doesn't claim its list is exhaustive

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Andorran tax identity code (NRT)",
					is.Func("valid NRT", isValidTaxIdentityCode),
				),
			),
		),
	)
}

func isValidTaxIdentityCode(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return false
	}
	return validateNRTCode(code) == nil
}

func validateNRTCode(code cbc.Code) error {
	if code == "" {
		return nil
	}
	if !reNRT.MatchString(code.String()) {
		return errors.New("invalid NRT format")
	}
	return nil
}