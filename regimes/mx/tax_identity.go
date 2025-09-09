package mx

import (
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

/*
 * Sources of data:
 *
 *  - https://tramites.aguascalientes.gob.mx/download/documentos/D20230407194800_Estructura%20RFC.pdf
 *  - http://omawww.sat.gob.mx/fichas_tematicas/reforma_fiscal/Documents/Especificaciones_tecnicas_EF_SOCAP.PDF
 *
 */

// Constants used to specific tax identity codes.
const (
	TaxIdentityCodeGeneric cbc.Code = "XAXX010101000"
	TaxIdentityCodeForeign cbc.Code = "XEXX010101000"
)

// Tax Identity Type
const (
	TaxIdentityTypePerson  cbc.Key = "person"
	TaxIdentityTypeCompany cbc.Key = "company"
)

// Tax Identity Patterns
const (
	TaxIdentityPatternPerson  = `^([A-ZÑ\&]{4})([0-9]{6})([A-Z0-9]{3})$`
	TaxIdentityPatternCompany = `^([A-ZÑ\&]{3})([0-9]{6})([A-Z0-9]{3})$`
)

// Tax Identity Regexp
var (
	TaxIdentityRegexpPerson  = regexp.MustCompile(TaxIdentityPatternPerson)
	TaxIdentityRegexpCompany = regexp.MustCompile(TaxIdentityPatternCompany)
	TaxCodeBadCharsRegexp    = regexp.MustCompile(`[^A-ZÑ\&0-9]+`)
)

// ValidateTaxIdentity validates a tax identity for SAT.
func ValidateTaxIdentity(tID *tax.Identity) error {
	if tID == nil {
		return nil
	}
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code,
			validation.By(ValidateTaxCode),
			validation.Skip, // don't apply regular code validation
		),
	)
}

// NormalizeTaxIdentity ensures the tax code is good for mexico
func NormalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil {
		return
	}
	tID.Code = NormalizeTaxCode(tID.Code)
}

// NormalizeTaxCode normalizes a tax code according to SAT rules.
// It handles special cases where company tax codes (RFC) may include an "MX" prefix.
// For example, a valid company code could be "MXG70123123Z". The function attempts
// to validate both with and without the "MX" prefix.
func NormalizeTaxCode(code cbc.Code) cbc.Code {
	c := strings.ToUpper(code.String())
	c = TaxCodeBadCharsRegexp.ReplaceAllString(c, "")

	codeTrimmed := strings.TrimPrefix(c, "MX")

	// If the trimmed code looks valid, return it
	if typ := DetermineTaxCodeType(cbc.Code(codeTrimmed)); !typ.IsEmpty() {
		return cbc.Code(codeTrimmed)
	}

	// If the trimmed code doesn't look valid, return the original code
	return cbc.Code(c)
}

// ValidateTaxCode validates a tax code according to the rules
// defined by the Mexican SAT.
func ValidateTaxCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	if typ := DetermineTaxCodeType(code); typ.IsEmpty() {
		return tax.ErrIdentityCodeInvalid
	}
	return nil
}

// DetermineTaxCodeType determines the type of tax code or provides
// an empty key if it looks invalid.
func DetermineTaxCodeType(code cbc.Code) cbc.Key {
	str := code.String()
	switch {
	case TaxIdentityRegexpPerson.MatchString(str):
		return TaxIdentityTypePerson
	case TaxIdentityRegexpCompany.MatchString(str):
		return TaxIdentityTypeCompany
	default:
		return cbc.KeyEmpty
	}
}
