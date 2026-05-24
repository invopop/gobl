package no

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var (
	taxCodeRegexps = []*regexp.Regexp{
		regexp.MustCompile(`^\d{9}$`),
	}
)

func normalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil {
		return
	}
	tax.NormalizeIdentity(tID)
	tID.Code = cbc.Code(strings.TrimSuffix(tID.Code.String(), "MVA"))
}

func normalizeOrgIdentity(id *org.Identity) {
	if id == nil || id.Type != IdentityTypeORG {
		return
	}
	code := strings.ToUpper(id.Code.String())
	code = tax.IdentityCodeBadCharsRegexp.ReplaceAllString(code, "")
	code = strings.TrimPrefix(code, "NO")
	code = strings.TrimSuffix(code, "MVA")
	id.Code = cbc.Code(code)
}

func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateTaxCode)),
	)
}

func validateTaxCode(value any) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	match := false
	for _, re := range taxCodeRegexps {
		if re.MatchString(val) {
			match = true
			break
		}
	}
	if !match {
		return errors.New("invalid format")
	}

	return validateTaxCodeChecksum(val)
}

func validateTaxCodeChecksum(val string) error {
	// Norwegian organization numbers use modulo-11 checksum.
	// The check digit is the last (9th) digit, calculated from the first 8 digits
	// using weights [3, 2, 7, 6, 5, 4, 3, 2].
	weights := []int{3, 2, 7, 6, 5, 4, 3, 2}
	total := 0

	for i := range 8 {
		digit, err := strconv.Atoi(string(val[i]))
		if err != nil {
			return errors.New("invalid digit")
		}
		total += digit * weights[i]
	}

	remainder := total % 11
	if remainder == 1 {
		return errors.New("invalid number: no valid check digit exists")
	}

	checkDigit := 0
	if remainder != 0 {
		checkDigit = 11 - remainder
	}

	lastDigit, err := strconv.Atoi(string(val[8]))
	if err != nil {
		return errors.New("invalid digit")
	}

	if lastDigit != checkDigit {
		return errors.New("checksum mismatch")
	}

	return nil
}
