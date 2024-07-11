package it

import (
	"errors"
	"fmt"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/validation"
)

const (
	// IdentityKeyFiscalCode is the type of identity that represents the Italian
	// "Codice Fiscale", a fiscal code issued to individuals and other taxable entities
	// that is independent from the "Partita IVA" or VAT number used by businesses.
	IdentityKeyFiscalCode cbc.Key = "it-fiscal-code"
)

func normalizeIdentity(id *org.Identity) error {
	if id == nil || id.Key != IdentityKeyFiscalCode {
		return nil
	}
	code := strings.ToUpper(id.Code.String())
	code = common.TaxCodeBadCharsRegexp.ReplaceAllString(code, "")
	code = strings.TrimPrefix(code, string(l10n.IT))
	id.Code = cbc.Code(code)
	return nil
}

// validateIdentities helps confirm that an identity of a specific type is valid.
func validateIdentity(id *org.Identity) error {
	if id == nil || id.Key != IdentityKeyFiscalCode {
		return nil
	}
	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.Required,
			validation.By(validateFiscalCode),
			validation.Skip,
		),
	)
}

// Based on details at https://en.wikipedia.org/wiki/Italian_fiscal_code
func validateFiscalCode(value interface{}) error {
	val, ok := value.(cbc.Code)
	if !ok || val == cbc.CodeEmpty {
		return nil
	}
	code := val.String()

	matched := taxIDPersonRegexPattern.MatchString(code)
	if !matched {
		return errors.New("invalid format")
	}

	var sum int
	for i := 0; i < 15; i++ {
		c := strings.Index(taxIDCharCode, string(code[i]))
		if c < 10 {
			c += 10 // move numbers to letters
		}
		if !(i%2 == 0) { // even as count starts from 1
			sum += strings.Index(taxIDEvenChars, string(taxIDCharCode[c]))
		} else { // odd
			sum += strings.Index(taxIDOddChars, string(taxIDCharCode[c]))
		}
	}

	x := string(taxIDCharCode[(sum%taxIDCRCMod)+10])
	if x != string(code[15]) {
		return fmt.Errorf("invalid check digit, expected '%s'", x)
	}

	return nil
}
