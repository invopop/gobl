package it

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

const (
	// IdentityKeyFiscalCode is the type of identity that represents the Italian
	// "Codice Fiscale", a fiscal code issued to individuals and other taxable entities
	// that is independent from the "Partita IVA" or VAT number used by businesses.
	IdentityKeyFiscalCode cbc.Key = "it-fiscal-code"
)

// source http://blog.marketto.it/2016/01/regex-validazione-codice-fiscale-con-omocodia/
var taxIDPersonRegexPattern = regexp.MustCompile(`^(?:[A-Z][AEIOU][AEIOUX]|[AEIOU]X{2}|[B-DF-HJ-NP-TV-Z]{2}[A-Z]){2}(?:[\dLMNP-V]{2}(?:[A-EHLMPR-T](?:[04LQ][1-9MNP-V]|[15MR][\dLMNP-V]|[26NS][0-8LMNP-U])|[DHPS][37PT][0L]|[ACELMRT][37PT][01LM]|[AC-EHLMPR-T][26NS][9V])|(?:[02468LNQSU][048LQU]|[13579MPRTV][26NS])B[26NS][9V])(?:[A-MZ][1-9MNP-V][\dLMNP-V]{2}|[A-M][0L](?:[1-9MNP-V][\dLMNP-V]|[0L][1-9MNP-V]))[A-Z]$`)

var identityKeyDefinitions = []*cbc.KeyDefinition{
	{
		Key: IdentityKeyFiscalCode,
		Name: i18n.String{
			i18n.EN: "Fiscal Code",
			i18n.IT: "Codice Fiscale",
		},
	},
}

func normalizeIdentity(id *org.Identity) {
	if id == nil || id.Key != IdentityKeyFiscalCode {
		return
	}
	code := strings.ToUpper(id.Code.String())
	code = tax.IdentityCodeBadCharsRegexp.ReplaceAllString(code, "")
	code = strings.TrimPrefix(code, string(l10n.IT))
	id.Code = cbc.Code(code)
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

	// Codice fiscale can belong to either a person or a company. Companies use
	// the regular VAT code, so we test the length or assume that we're
	// dealing with a physical person's details.
	if len(code) == 11 {
		return validateTaxCode(value)
	}

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
