package cl

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// The RUT (Rol Único Tributario, e.g. "76.086.428-5") identifies both
// individuals and companies in Chile. Modern bodies have 7-8 digits, but
// shorter historic ones exist (e.g. "1-9"), so 1-8 digits are accepted.
//
// Validation expects a normalized code ("760864285"): tax.NormalizeIdentity
// (registered in cl.go) strips punctuation and uppercases beforehand, and
// codes still containing separators are rejected.
//
// The check digit uses modulo 11: body digits, right to left, are multiplied
// by the repeating sequence 2-7 and summed; the digit is 11 minus (sum mod
// 11), where 11 maps to "0" and 10 to "K" (always uppercase).
//
// Sources:
//   - SII DTE spec (modulo 11, uppercase "K"):
//     https://www.sii.cl/factura_electronica/Webservice_Registro_Reclamo_DTE_V1.2.pdf
//   - Full algorithm (verified against this implementation):
//     https://lookuptax.com/docs/es/numero-identificacion-fiscal/guia-rut-chile

var rutFormatRegexp = regexp.MustCompile(`^\d{1,8}[0-9K]$`)

var rutMultipliers = []int{2, 3, 4, 5, 6, 7}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Chilean RUT",
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
	return validateRUT(code.String()) == nil
}

func validateRUT(val string) error {
	if !rutFormatRegexp.MatchString(val) {
		return errors.New("invalid format")
	}

	body := val[:len(val)-1]
	check := val[len(val)-1]

	if check != rutCheckDigit(body) {
		return errors.New("checksum mismatch")
	}

	return nil
}

// rutCheckDigit calculates the expected check digit for a RUT body using
// the modulo 11 algorithm described above.
func rutCheckDigit(body string) byte {
	sum := 0
	for i, j := len(body)-1, 0; i >= 0; i, j = i-1, j+1 {
		digit := int(body[i] - '0')
		sum += digit * rutMultipliers[j%len(rutMultipliers)]
	}

	switch result := 11 - (sum % 11); result {
	case 11:
		return '0'
	case 10:
		return 'K'
	default:
		return byte('0' + result)
	}
}
