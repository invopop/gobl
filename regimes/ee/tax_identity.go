package ee

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Estonian VAT numbers (KMKR, käibemaksukohustuslase registreerimisnumber) consist
// of exactly 9 digits. The EE country prefix is stripped by NormalizeIdentity before
// this validation runs, so the rules below only see the 9 digits.
//
// Checksum: weighted sum of the first 8 digits with weights [3,7,1,3,7,1,3,7]. The
// 9th digit must equal (10 - sum%10) % 10.
//
// Reference: https://meta.cdq.com/DataModel:CDQ/Business_Partner/Identifier/EU_VAT_ID_EE

var (
	taxCodeRegexp  = regexp.MustCompile(`^\d{9}$`)
	taxCodeWeights = []int{3, 7, 1, 3, 7, 1, 3, 7}
)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Estonian VAT identity code",
					is.Func("valid", isValidTaxIdentityCode),
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
	return validateTaxCode(code) == nil
}

func validateTaxCode(code cbc.Code) error {
	val := code.String()
	if !taxCodeRegexp.MatchString(val) {
		return errors.New("invalid format")
	}
	return validateChecksum(val)
}

func validateChecksum(val string) error {
	var sum int
	for i, w := range taxCodeWeights {
		sum += w * int(val[i]-'0')
	}
	expected := (10 - sum%10) % 10
	if int(val[8]-'0') != expected {
		return errors.New("checksum mismatch")
	}
	return nil
}
