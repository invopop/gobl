package at

import (
	"errors"
	"math"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Source: https://github.com/ltns35/go-vat

var (
	taxCodeMultipliers = []int{
		1,
		2,
		1,
		2,
		1,
		2,
		1,
	}
	taxCodeRegexps = []*regexp.Regexp{
		regexp.MustCompile(`^U\d{8}$`),
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
	var total float64 = 0
	for i, m := range taxCodeMultipliers {
		num := int(val[i+1] - '0')
		x := float64(num * m)
		if x > 9 {
			total += math.Floor(x/10) + math.Mod(x, 10)
		} else {
			total += x
		}
	}

	total = 10 - math.Mod(total+4, 10)
	if total == 10 {
		total = 0
	}

	lastNum := int(val[8] - '0')
	if lastNum != int(total) {
		return errors.New("checksum mismatch")
	}

	return nil
}
