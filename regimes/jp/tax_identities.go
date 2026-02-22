package jp

import (
	"errors"
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// References:
// - Corporate Number: https://www.houjin-bangou.nta.go.jp/en/setsumei/
// - Qualified Invoice: https://www.nta.go.jp/taxes/shiraberu/zeimokubetsu/shohi/keigenzeiritsu1.htm

// taxIdentityPattern matches a Qualified Invoice Issuer Number:
// "T" followed by 13 digits.
var taxIdentityPattern = regexp.MustCompile(`^T\d{13}$`)

// normalizeTaxIdentity performs normalization on tax identities
func normalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil {
		return
	}

	// Clean the code
	code := strings.ToUpper(strings.TrimSpace(tID.Code.String()))
	code = tax.IdentityCodeBadCharsRegexp.ReplaceAllString(code, "")

	// Remove country prefix if present
	code = strings.TrimPrefix(code, string(l10n.JP))

	// For Qualified Invoice Issuer numbers, ensure a proper format
	if strings.HasPrefix(code, "T") {
		// Already has T prefix, keep it
		tID.Code = cbc.Code(code)
		return
	}

	// If it looks like a 13-digit corporate number without a T prefix,
	// store as-is (validation will reject it as a tax identity code).
	tID.Code = cbc.Code(code)
}

// validateTaxIdentity validates a tax identity for Japan
func validateTaxIdentity(tID *tax.Identity) error {
	if tID == nil {
		return nil
	}

	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code,
			validation.By(validateTaxCode),
			validation.Skip,
		),
	)
}

// validateTaxCode validates the tax code format
func validateTaxCode(value any) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}

	val := code.String()

	// Check if it's a Qualified Invoice Issuer Registration Number (T + 13 digits)
	if !taxIdentityPattern.MatchString(val) {
		return errors.New("must be 'T' followed by 13 digits")
	}

	// Validate check digit
	if err := validateTNumberCheckDigit(val); err != nil {
		return errors.New("invalid check digit")
	}

	return nil
}

// validateTNumberCheckDigit validates the check digit of a T-number.
// The first digit of the 13-digit part (immediately after "T") is the check digit,
// calculated from the remaining 12 digits using the NTA corporate number algorithm.
func validateTNumberCheckDigit(tNumber string) error {
	if len(tNumber) != 14 || tNumber[0] != 'T' {
		return errors.New("invalid format")
	}

	digits := tNumber[1:] // Strip the "T" prefix; digits[0] is the check digit

	// NTA checksum algorithm: sum base digits (indices 1–12) with weights
	// alternating 2,1 from left (odd index → weight 2, even index → weight 1).
	// Check digit = 9 - (sum mod 9), or 0 if sum mod 9 == 0.
	var sum int
	for j := 1; j <= 12; j++ {
		d := int(digits[j] - '0')
		w := 1
		if j%2 == 1 {
			w = 2
		}
		sum += d * w
	}

	remainder := sum % 9
	var expectedCheck int
	if remainder == 0 {
		expectedCheck = 0
	} else {
		expectedCheck = 9 - remainder
	}

	actualCheck := int(digits[0] - '0')
	if actualCheck != expectedCheck {
		return errors.New("check digit mismatch")
	}

	return nil
}
