package pl

import (
	"regexp"
	"strconv"
	"unicode"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

/*
 * Sources of data:
 *
 *  - https://tramites.aguascalientes.gob.mx/download/documentos/D20230407194800_Estructura%20RFC.pdf
 *  - https://pl.wikipedia.org/wiki/Numer_identyfikacji_podatkowej
 *
 */

// Tax Identity Type
const (
	TaxIdentityTypePolish cbc.Key = "polish"
	TaxIdentityTypeOther  cbc.Key = "other"
)

// Tax Identity Patterns
const (
	TaxIdentityPatternPolish = `^[1-9]((\d[1-9])|([1-9]\d))\d{7}$`
	TaxIdentityPatternOther  = `^.{1,50}$`
)

// Tax Identity Regexp
var (
	TaxIdentityRegexpPolish = regexp.MustCompile(TaxIdentityPatternPolish)
	TaxIdentityRegexpOther  = regexp.MustCompile(TaxIdentityPatternOther)
)

func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code,
			validation.Required,
			validation.By(validateTaxCode),
		),
		validation.Field(&tID.Zone, validation.Required),
	)
}

func validatePolishTaxIdentity(value interface{}) error {
	code, ok := value.(cbc.Code)
	str := code.String()
	if !ok {
		return nil
	}
	if TaxIdentityRegexpPolish.MatchString(str) && validateNIPChecksum(code) {
		return nil
	}
	return tax.ErrIdentityCodeInvalid
}

func validateTaxCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok {
		return nil
	}
	if code == "" {
		return nil
	}
	typ := DetermineTaxCodeType(code)
	if typ.IsEmpty() {
		return tax.ErrIdentityCodeInvalid
	}
	if typ == TaxIdentityTypePolish {
		if validateNIPChecksum(code) {
			return nil
		}
		return tax.ErrIdentityCodeInvalid
	}
	return nil
}

// DetermineTaxCodeType determines the type of tax code or provides
// an empty key if it looks invalid.
func DetermineTaxCodeType(code cbc.Code) cbc.Key {
	str := code.String()
	switch {
	case TaxIdentityRegexpPolish.MatchString(str):
		return TaxIdentityTypePolish
	case TaxIdentityRegexpOther.MatchString(str):
		return TaxIdentityTypeOther
	default:
		return cbc.KeyEmpty
	}
}

func validateNIPChecksum(code cbc.Code) bool {
	nipStr := code.String()
	if len(nipStr) != 10 {
		return false
	}

	for _, char := range nipStr {
		if !unicode.IsDigit(char) {
			return false
		}
	}

	digits := make([]int, 10)
	for i, char := range nipStr {
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return false
		}
		digits[i] = digit
	}

	weights := [9]int{6, 5, 7, 2, 3, 4, 5, 6, 7}
	checkSum := 0
	for i, digit := range digits[:9] {
		checkSum += digit * weights[i]
	}
	checkSum %= 11

	return checkSum == digits[9]
}
