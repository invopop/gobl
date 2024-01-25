package it

// The tax code here refers to Partita IVA, which is a distinct construct from
// Codice Fiscale. Italy operates with two types of tax identification codes.
// Though not all Italian persons possess Partita IVA, all parties engaged in
// economic activities within Italy are required to have one.

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Italian identity types required for transactions inside the same country,
// but not needed for inter-country sales.
const (
	TaxIdentityTypeBusiness   cbc.Key = "business" // default
	TaxIdentityTypeGovernment cbc.Key = "government"
	TaxIdentityTypeIndividual cbc.Key = "individual"
)

var taxIdentityTypeDefinitions = []*tax.KeyDefinition{
	{
		Key: TaxIdentityTypeBusiness,
		Name: i18n.String{
			i18n.EN: "Legal Person",
		},
	},
	{
		Key: TaxIdentityTypeGovernment,
		Name: i18n.String{
			i18n.EN: "Public Administration",
		},
	},
	{
		Key: TaxIdentityTypeIndividual,
		Name: i18n.String{
			i18n.EN: "Natural Person",
		},
	},
}

const (
	taxIDEvenChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	taxIDOddChars  = "BAKPLCQDREVOSFTGUHMINJWZYX"
	taxIDCharCode  = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	taxIDCRCMod    = 26
)

// source http://blog.marketto.it/2016/01/regex-validazione-codice-fiscale-con-omocodia/
var taxIDPersonRegexPattern = regexp.MustCompile(`^(?:[A-Z][AEIOU][AEIOUX]|[AEIOU]X{2}|[B-DF-HJ-NP-TV-Z]{2}[A-Z]){2}(?:[\dLMNP-V]{2}(?:[A-EHLMPR-T](?:[04LQ][1-9MNP-V]|[15MR][\dLMNP-V]|[26NS][0-8LMNP-U])|[DHPS][37PT][0L]|[ACELMRT][37PT][01LM]|[AC-EHLMPR-T][26NS][9V])|(?:[02468LNQSU][048LQU]|[13579MPRTV][26NS])B[26NS][9V])(?:[A-MZ][1-9MNP-V][\dLMNP-V]{2}|[A-M][0L](?:[1-9MNP-V][\dLMNP-V]|[0L][1-9MNP-V]))[A-Z]$`)

// validateTaxIdentity performs checks on the tax codes according to the type
// that was set. Additional validation is laid out at the invoice layer.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code,
			validation.When(
				!tID.Type.In(TaxIdentityTypeIndividual), // include empty type!
				validation.By(validateTaxCode),
			),
			validation.When(
				tID.Type.In(TaxIdentityTypeIndividual),
				validation.By(validateIndividualTaxCode),
			),
		),
	)
}

// normalizeTaxIdentity removes any whitespace or separation characters and ensures all letters are
// uppercase.
func normalizeTaxIdentity(tID *tax.Identity) error {
	if err := common.NormalizeTaxIdentity(tID); err != nil {
		return err
	}

	// try to determine the type automatically
	if tID.Type == cbc.KeyEmpty {
		if tID.Code == "" {
			return nil
		}
		// note that we don't yet have a way to determine government codes
		if taxIDPersonRegexPattern.MatchString(tID.Code.String()) {
			tID.Type = TaxIdentityTypeIndividual
		} else {
			tID.Type = TaxIdentityTypeBusiness
		}
	}

	return nil
}

// source: https://it.wikipedia.org/wiki/Partita_IVA#Struttura_del_codice_identificativo_di_partita_IVA
func validateTaxCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	str := code.String()

	// Check code is just numbers
	for _, v := range str {
		x := v - 48
		if x < 0 || x > 9 {
			return errors.New("contains invalid characters")
		}
	}

	if len(str) != 11 {
		return errors.New("invalid length")
	}

	chk := common.ComputeLuhnCheckDigit(str[:10])
	if chk != str[10:] {
		return errors.New("invalid check digit")
	}

	return nil
}

// Based on details at https://en.wikipedia.org/wiki/Italian_fiscal_code
func validateIndividualTaxCode(value interface{}) error {
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
