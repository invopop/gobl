package be

import (
	"errors"
	"math"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Source: https://github.com/ltns35/go-vat

// Belgium tax codes are split between personal and enterprise tax IDs:
// - Personal IDs are 9 characters long.
// - Enterprise IDs are 10 characters long and must always start with a `0` or a `1`.

var (
	taxCodeRegexps = []*regexp.Regexp{
		regexp.MustCompile(`^[01]?\d{9}$`),
	}
)

// validateTaxIdentity checks to ensure the NIT code looks okay.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateTaxCode)),
	)
}

func validateTaxCode(value interface{}) error {
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

	return commercialCheck(val)
}

func commercialCheck(val string) error {
	// Pad regular citizen IDs with a 0 at the beginning to
	// ensure we can use the same regular checks.
	if len(val) == 9 {
		val = "0" + val
	}

	str := val[:8]
	num, _ := strconv.Atoi(str) //nolint:errcheck

	chk := 97 - math.Mod(float64(num), 97)

	str = val[8:10]
	last, _ := strconv.Atoi(str) //nolint:errcheck

	if float64(last) != chk {
		return errors.New("checksum mismatch")
	}

	return nil
}
