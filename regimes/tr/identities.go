package tr

// Checksum algorithm based on:
// https://github.com/MhmtMutlu/tckn-vkn-validator

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

const (
	// IdentityTypeTCKN represents the Turkish national identity number
	// (Türkiye Cumhuriyeti Kimlik Numarası), an 11-digit number issued by
	// the civil registry to individuals and sole traders. For international
	// trade, a VKN (tax identity) should be used instead.
	IdentityTypeTCKN cbc.Code = "TCKN"
)

var tcknRegexp = regexp.MustCompile(`^\d{11}$`)

var identityDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeTCKN,
		Name: i18n.String{
			i18n.EN: "National Identity Number (TCKN)",
			i18n.TR: "Türkiye Cumhuriyeti Kimlik Numarası (TCKN)",
		},
	},
}

// normalizeIdentity standardizes TCKN identity codes.
func normalizeIdentity(id *org.Identity) {
	if id == nil || id.Type != IdentityTypeTCKN {
		return
	}
	code := cbc.NormalizeNumericalCode(id.Code)
	id.Code = code
}

// validateIdentity checks that the TCKN identity code is valid.
func validateIdentity(id *org.Identity) error {
	if id == nil || id.Type != IdentityTypeTCKN {
		return nil
	}
	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.Required,
			validation.By(validateTCKNCode),
			validation.Skip,
		),
	)
}

func validateTCKNCode(value any) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	s := code.String()
	if !tcknRegexp.MatchString(s) {
		return errInvalidFormat
	}
	return verifyTCKN(s)
}

// verifyTCKN validates the 11-digit Turkish national identity number checksum.
//
// Algorithm:
//  1. First digit must not be 0.
//  2. Digit 10 = (sum of odd-indexed digits * 7 - sum of even-indexed digits) % 10
//     (indices 0,2,4,6,8 are odd positions; 1,3,5,7 are even positions)
//  3. Digit 11 = sum of first 10 digits % 10
func verifyTCKN(s string) error {
	digits := stringToDigits(s)
	if digits[0] == 0 {
		return errInvalidFormat
	}
	oddSum := digits[0] + digits[2] + digits[4] + digits[6] + digits[8]
	evenSum := digits[1] + digits[3] + digits[5] + digits[7]
	d10 := (oddSum*7 - evenSum) % 10
	if d10 < 0 {
		d10 += 10
	}
	if d10 != digits[9] {
		return errInvalidChecksum
	}
	total := 0
	for i := 0; i < 10; i++ {
		total += digits[i]
	}
	if total%10 != digits[10] {
		return errInvalidChecksum
	}
	return nil
}
