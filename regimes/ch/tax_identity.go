package ch

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
		5,
		4,
		3,
		2,
		7,
		6,
		5,
		4,
	}
	taxCodeRegexps = []*regexp.Regexp{
		regexp.MustCompile(`^E\d{9}$`),
	}
	taxCodeSuffixes = regexp.MustCompile(`(MWST|TVA|IVA)$`)
)

// normalizeTaxIdentity will remove any whitespace or separation characters from
// the tax code and also make sure the default type is set.
func normalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil {
		return
	}
	tax.NormalizeIdentity(tID)
	// CH has some strange suffixes, remove them.
	tID.Code = cbc.Code(taxCodeSuffixes.ReplaceAllString(tID.Code.String(), ""))
}

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
	var total float64
	for i, m := range taxCodeMultipliers {
		num := int(val[i+1] - '0')
		x := float64(num * m)
		total = total + x
	}

	total = 11 - math.Mod(total, 11)
	if total == 10 {
		return errors.New("invalid code")
	}
	if total == 11 {
		total = 0
	}

	lastNum := int(val[9] - '0')
	if lastNum != int(total) {
		return errors.New("checksum mismatch")
	}

	return nil
}
