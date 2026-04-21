package bis

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// ISO 6523 ICD scheme identifiers we know how to validate. Each maps to a
// checksum/format check below.
const (
	schemeGLN     cbc.Code = "0088" // GLN — PEPPOL-COMMON-R040
	schemeNOOrg   cbc.Code = "0192" // Norwegian organization number — PEPPOL-COMMON-R041
	schemeDKCVR   cbc.Code = "0184" // Danish CVR — PEPPOL-COMMON-R042
	schemeBEEnt   cbc.Code = "0208" // Belgian enterprise — PEPPOL-COMMON-R043
	schemeITIPA   cbc.Code = "0201" // Italian IPA — PEPPOL-COMMON-R044
	schemeITCF    cbc.Code = "9907" // Italian Codice Fiscale — PEPPOL-COMMON-R045/R046
	schemeITPIva  cbc.Code = "9906" // Italian Partita IVA — PEPPOL-COMMON-R047
	schemeSEOrg   cbc.Code = "0007" // Swedish organization number — PEPPOL-COMMON-R049
	schemeAUABN   cbc.Code = "0151" // Australian Business Number — PEPPOL-COMMON-R050
	schemeDKPNum  cbc.Code = "0200" // Danish P-number — PEPPOL-COMMON-R052
	schemeDKSENum cbc.Code = "0198" // Danish SE-number (ERSTORG) — PEPPOL-COMMON-R053
)

func orgPartyRules() *rules.Set {
	return rules.For(new(org.Party),
		// PEPPOL-EN16931-R010/R020: buyer and seller MUST have an electronic
		// address. Applied to any org.Party — suppliers and customers alike.
		rules.Field("inboxes",
			rules.Assert("R010-R020", "party electronic address (inbox) is required (PEPPOL-EN16931-R010/R020)",
				is.Present,
			),
		),
	)
}

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.Assert("COMMON", "identifier format invalid",
			is.FuncError("identity format", identityFormatValid),
		),
	)
}

func orgInboxRules() *rules.Set {
	return rules.For(new(org.Inbox),
		rules.Assert("COMMON", "inbox code format invalid",
			is.FuncError("inbox format", inboxFormatValid),
		),
	)
}

// identityFormatValid checks the identity's code against the checksum/format
// rule implied by its ISO 6523 scheme ID (taken from `Ext[iso-scheme-id]`).
func identityFormatValid(val any) error {
	id, ok := val.(*org.Identity)
	if !ok || id == nil {
		return nil
	}
	scheme := id.Ext.Get(iso.ExtKeySchemeID)
	return checkSchemeFormat(scheme, id.Code)
}

// inboxFormatValid checks an inbox's code against the scheme-implied format.
// Uses Inbox.Scheme directly rather than an extension.
func inboxFormatValid(val any) error {
	ib, ok := val.(*org.Inbox)
	if !ok || ib == nil || ib.Scheme == "" || ib.Code == "" {
		return nil
	}
	return checkSchemeFormat(ib.Scheme, ib.Code)
}

// checkSchemeFormat dispatches to the format validator for the given ISO 6523
// scheme ID. Unknown schemes return nil (no validation performed).
func checkSchemeFormat(scheme cbc.Code, code cbc.Code) error {
	if code == "" {
		return nil
	}
	switch scheme {
	case schemeGLN:
		if !validGLN(code.String()) {
			return errors.New("invalid GLN: must be 13 digits with valid Mod 10 checksum (PEPPOL-COMMON-R040)")
		}
	case schemeNOOrg:
		if !validNorwegianOrg(code.String()) {
			return errors.New("invalid Norwegian organization number: must be 9 digits with valid Mod 11 checksum (PEPPOL-COMMON-R041)")
		}
	case schemeDKCVR:
		if !validDanishCVR(code.String()) {
			return errors.New("invalid Danish CVR: must be 8 digits with valid Mod 11 checksum (PEPPOL-COMMON-R042)")
		}
	case schemeBEEnt:
		if !validBelgianEnterprise(code.String()) {
			return errors.New("invalid Belgian enterprise number: must be 10 digits with valid Mod 97 checksum (PEPPOL-COMMON-R043)")
		}
	case schemeITIPA:
		if !validITIPA(code.String()) {
			return errors.New("invalid Italian IPA code: must be 6 alphanumeric characters (PEPPOL-COMMON-R044)")
		}
	case schemeITCF:
		if !validITCodiceFiscale(code.String()) {
			return errors.New("invalid Italian Codice Fiscale: must be 11 digits or 16 alphanumerics (PEPPOL-COMMON-R045/R046)")
		}
	case schemeITPIva:
		if !validITPartitaIVA(code.String()) {
			return errors.New("invalid Italian Partita IVA: must be 11 digits (PEPPOL-COMMON-R047)")
		}
	case schemeSEOrg:
		if !validSwedishOrg(code.String()) {
			return errors.New("invalid Swedish organization number: must be 10 digits with valid Luhn checksum (PEPPOL-COMMON-R049)")
		}
	case schemeAUABN:
		if !validAustralianABN(code.String()) {
			return errors.New("invalid Australian ABN: must be 11 digits with valid weighted checksum (PEPPOL-COMMON-R050)")
		}
	case schemeDKPNum:
		if !validDanishPNumber(code.String()) {
			return errors.New("invalid Danish P-number: must be 10 digits (PEPPOL-COMMON-R052)")
		}
	case schemeDKSENum:
		if !validDanishSENumber(code.String()) {
			return errors.New("invalid Danish SE-number: must be 8 digits (PEPPOL-COMMON-R053)")
		}
	}
	return nil
}

// Format validators

var digitsOnlyRe = regexp.MustCompile(`^\d+$`)

func onlyDigits(s string) bool {
	return digitsOnlyRe.MatchString(s)
}

// validGLN — 13 digits, GS1 Mod 10 checksum. Last digit is the check digit.
func validGLN(code string) bool {
	if len(code) != 13 || !onlyDigits(code) {
		return false
	}
	sum := 0
	for i := 0; i < 12; i++ {
		n := int(code[i] - '0')
		if i%2 == 0 {
			sum += n
		} else {
			sum += n * 3
		}
	}
	check := (10 - sum%10) % 10
	return check == int(code[12]-'0')
}

// validNorwegianOrg — 9 digits, Mod 11 checksum with weights 3,2,7,6,5,4,3,2.
func validNorwegianOrg(code string) bool {
	if len(code) != 9 || !onlyDigits(code) {
		return false
	}
	weights := []int{3, 2, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i, w := range weights {
		sum += int(code[i]-'0') * w
	}
	check := 11 - sum%11
	switch check {
	case 11:
		check = 0
	case 10:
		return false // not a valid NO org number
	}
	return check == int(code[8]-'0')
}

// validDanishCVR — 8 digits, Mod 11 checksum with weights 2,7,6,5,4,3,2,1.
func validDanishCVR(code string) bool {
	if len(code) != 8 || !onlyDigits(code) {
		return false
	}
	weights := []int{2, 7, 6, 5, 4, 3, 2, 1}
	sum := 0
	for i, w := range weights {
		sum += int(code[i]-'0') * w
	}
	return sum%11 == 0
}

// validBelgianEnterprise — 10 digits, first 8 as number N, N mod 97 equals last 2.
func validBelgianEnterprise(code string) bool {
	if len(code) != 10 || !onlyDigits(code) {
		return false
	}
	n := 0
	for i := 0; i < 8; i++ {
		n = n*10 + int(code[i]-'0')
	}
	check := 97 - (n % 97)
	expected := int(code[8]-'0')*10 + int(code[9]-'0')
	return check == expected
}

// validITIPA — Italian IPA: 6 alphanumeric characters.
var itIPARe = regexp.MustCompile(`^[A-Z0-9]{6}$`)

func validITIPA(code string) bool {
	return itIPARe.MatchString(code)
}

// validITCodiceFiscale — 11 digits (legal entities) or 16 alphanumerics (people).
var itCFPersonRe = regexp.MustCompile(`^[A-Z]{6}\d{2}[A-Z]\d{2}[A-Z]\d{3}[A-Z]$`)

func validITCodiceFiscale(code string) bool {
	if len(code) == 11 && onlyDigits(code) {
		return true
	}
	return len(code) == 16 && itCFPersonRe.MatchString(code)
}

// validITPartitaIVA — 11 digits.
func validITPartitaIVA(code string) bool {
	return len(code) == 11 && onlyDigits(code)
}

// validSwedishOrg — 10 digits, last digit is Luhn checksum.
func validSwedishOrg(code string) bool {
	if len(code) != 10 || !onlyDigits(code) {
		return false
	}
	return luhnValid(code)
}

// luhnValid implements the standard Luhn Mod-10 algorithm over a numeric string.
func luhnValid(code string) bool {
	sum := 0
	double := false
	for i := len(code) - 1; i >= 0; i-- {
		n := int(code[i] - '0')
		if double {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		double = !double
	}
	return sum%10 == 0
}

// validAustralianABN — 11 digits. Subtract 1 from leading digit, apply
// weights 10,1,3,5,7,9,11,13,15,17,19, sum must be divisible by 89.
func validAustralianABN(code string) bool {
	if len(code) != 11 || !onlyDigits(code) {
		return false
	}
	weights := []int{10, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19}
	digits := make([]int, 11)
	for i := range digits {
		digits[i] = int(code[i] - '0')
	}
	digits[0]-- // subtract 1 from leading digit
	sum := 0
	for i, w := range weights {
		sum += digits[i] * w
	}
	return sum%89 == 0
}

// validDanishPNumber — 10 digits.
func validDanishPNumber(code string) bool {
	return len(code) == 10 && onlyDigits(code)
}

// validDanishSENumber — 8 digits.
func validDanishSENumber(code string) bool {
	return len(code) == 8 && onlyDigits(code)
}
