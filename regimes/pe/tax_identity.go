package pe

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// taxCodeWeights are the mod-11 multipliers used by the established Peruvian
// RUC check-digit convention. SUNAT's PVS T-Registro manual confirms that RUCs
// are numeric, eleven digits long, and validated with modulus 11, but does not
// publish the internal weights. The edge-case mapping is therefore pinned by
// tests against known RUCs instead of being presented as a statutory formula.
//
// Source: https://www2.sunat.gob.pe/orientacion/pvs/registro/manual/mu-0441-pvs-T-Registro_RM_170_2023_TR.pdf
var taxCodeWeights = []int{5, 4, 3, 2, 7, 6, 5, 4, 3, 2}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Peruvian RUC",
					is.Func("valid mod-11 RUC", isValidRUC),
				),
			),
		),
	)
}

// normalizeTaxIdentity performs the standard tax identity normalization,
// which removes punctuation and the "PE" country prefix, leaving the plain
// 11-digit RUC.
func normalizeTaxIdentity(tID *tax.Identity) {
	tax.NormalizeIdentity(tID)
}

// isValidRUC reports whether the value is a valid Peruvian RUC: eleven
// digits, starting with a taxpayer type prefix assigned by SUNAT, and ending
// with a mod-11 check digit.
func isValidRUC(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok {
		return false
	}
	s := code.String()
	if len(s) != 11 {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	// SUNAT's current RUC guidance documents prefixes 10 and 15 for natural
	// persons and 20 for legal entities. Its PVS T-Registro manual additionally
	// recognizes prefix 17. Prefix 16 appears in third-party validators, but is
	// deliberately not accepted without an authoritative public source.
	//
	// Sources:
	// https://centrovirtual.sunat.gob.pe/tramites/inscribete-ruc
	// https://www2.sunat.gob.pe/orientacion/pvs/registro/manual/mu-0441-pvs-T-Registro_RM_170_2023_TR.pdf
	switch s[:2] {
	case "10", "15", "17", "20":
	default:
		return false
	}

	sum := 0
	for i, w := range taxCodeWeights {
		sum += int(s[i]-'0') * w
	}
	check := 11 - (sum % 11)
	switch check {
	case 10:
		check = 0
	case 11:
		check = 1
	}
	return int(s[10]-'0') == check
}
