package mx

import (
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
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
	taxIdentityRegexpPerson  = regexp.MustCompile(TaxIdentityPatternPerson)
	taxIdentityRegexpCompany = regexp.MustCompile(TaxIdentityPatternCompany)
	taxCodeBadCharsRegexp    = regexp.MustCompile(`[^A-ZÑ\&0-9]+`)
)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn("MX"),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Mexican RFC tax identity code",
					is.Func("valid", isValidTaxIdentityCode),
				),
			),
		),
	)
}

func isValidTaxIdentityCode(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return false
	}
	return validateTaxCode(code) == nil
}

// normalizeTaxIdentity ensures the tax code is good for mexico
func normalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil {
		return
	}
	tID.Code = normalizeTaxCode(tID.Code)
}

// normalizeTaxCode normalizes a tax code according to SAT rules.
// It handles special cases where company tax codes (RFC) may include an "MX" prefix.
// For example, a valid company code could be "MXG70123123Z". The function attempts
// to validate both with and without the "MX" prefix.
func normalizeTaxCode(code cbc.Code) cbc.Code {
	c := strings.ToUpper(code.String())
	c = taxCodeBadCharsRegexp.ReplaceAllString(c, "")

	codeTrimmed := strings.TrimPrefix(c, "MX")

	// If the trimmed code looks valid, return it
	if typ := determineTaxCodeType(cbc.Code(codeTrimmed)); !typ.IsEmpty() {
		return cbc.Code(codeTrimmed)
	}

	// If the trimmed code doesn't look valid, return the original code
	return cbc.Code(c)
}

// validateTaxCode validates a tax code according to the rules
// defined by the Mexican SAT.
func validateTaxCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	if typ := determineTaxCodeType(code); typ.IsEmpty() {
		return tax.ErrIdentityCodeInvalid
	}
	return nil
}

// determineTaxCodeType determines the type of tax code or provides
// an empty key if it looks invalid.
func determineTaxCodeType(code cbc.Code) cbc.Key {
	str := code.String()
	switch {
	case taxIdentityRegexpPerson.MatchString(str):
		return TaxIdentityTypePerson
	case taxIdentityRegexpCompany.MatchString(str):
		return TaxIdentityTypeCompany
	default:
		return cbc.KeyEmpty
	}
}
