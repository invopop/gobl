package mx

import (
	"regexp"

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
	TaxIdentityCodeForeign cbc.Code = "XEXX010101000"
)

// Tax Identity Type
const (
	TaxIdentityTypePerson  cbc.Key = "person"
	TaxIdentityTypeCompany cbc.Key = "company"
)

// Tax Identity Patterns
const (
	TaxIdentityPatternPerson  = `^([A-ZÑ&]{4})([0-9]{6})([A-Z0-9]{3})$`
	TaxIdentityPatternCompany = `^([A-ZÑ&]{3})([0-9]{6})([A-Z0-9]{3})$`
)

// Tax Identity Regexp
var (
	TaxIdentityRegexpPerson  = regexp.MustCompile(TaxIdentityPatternPerson)
	TaxIdentityRegexpCompany = regexp.MustCompile(TaxIdentityPatternCompany)
)

func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code,
			validation.By(validateTaxCode),
		),
	)
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
