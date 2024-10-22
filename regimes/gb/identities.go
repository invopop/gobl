package gb

import (
	"errors"
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

const (
	// IdentityUTR represents the UK Unique Taxpayer Reference (UTR).
	IdentityUTR cbc.Key = "gb-utr"
	// IdentityNINO represents the UK National Insurance Number (NINO).
	IdentityNINO cbc.Key = "gb-nino"
)

var badCharsRegexPattern = regexp.MustCompile(`[^\d]`)
var ninoPattern = `^[A-CEGHJ-PR-TW-Z]{2}\d{6}[A-D]$`
var utrPattern = `^[1-9]\d{9}$`

// https://design.tax.service.gov.uk/hmrc-design-patterns/unique-taxpayer-reference/
// https://www.gov.uk/hmrc-internal-manuals/national-insurance-manual/nim39110

var identityKeyDefinitions = []*cbc.KeyDefinition{
	{
		Key: IdentityUTR,
		Name: i18n.String{
			i18n.EN: "Unique Taxpayer Reference",
		},
	},
	{
		Key: IdentityNINO,
		Name: i18n.String{
			i18n.EN: "National Insurance Number",
		},
	},
}

func normalizeTaxNumber(id *org.Identity) {
	if id == nil || (id.Key != IdentityUTR && id.Key != IdentityNINO) {
		return
	}

	if id.Key == IdentityUTR {
		code := id.Code.String()
		code = badCharsRegexPattern.ReplaceAllString(code, "")
		id.Code = cbc.Code(code)
	} else if id.Key == IdentityNINO {
		code := id.Code.String()
		code = strings.ToUpper(code)
		code = tax.IdentityCodeBadCharsRegexp.ReplaceAllString(code, "")
		id.Code = cbc.Code(code)
	}
}

func validateTaxNumber(id *org.Identity) error {
	if id == nil {
		return nil
	}

	if id.Key == IdentityNINO {
		return validation.ValidateStruct(id,
			validation.Field(&id.Code, validation.By(validateNinoCode)),
		)
	} else if id.Key == IdentityUTR {
		return validation.ValidateStruct(id,
			validation.Field(&id.Code, validation.By(validateUtrCode)),
		)
	}

	return nil
}

// validateUtrCode validates the normalized Unique Taxpayer Reference (UTR).
func validateUtrCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	// UK UTR pattern: 10 digits, first digit cannot be 0

	matched, err := regexp.MatchString(utrPattern, val)
	if err != nil {
		return err
	}

	if !matched {
		return errors.New("invalid UTR format")
	}

	return nil
}

// validateNinoCode validates the normalized National Insurance Number (NINO).
func validateNinoCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	// UK NINO pattern: Two prefix letters (valid), six digits, one suffix letter (A-D)

	matched, err := regexp.MatchString(ninoPattern, val)
	if err != nil {
		return err
	}

	if !matched {
		return errors.New("invalid NINO format")
	}

	// Check prefix letters
	if !isValidPrefix(val[:2]) {
		return errors.New("invalid prefix letters")
	}

	return nil
}

// isValidPrefix checks if the prefix letters are valid according to the specified rules.
func isValidPrefix(prefix string) bool {
	// Disallowed prefixes
	disallowedPrefixes := []string{"BG", "GB", "NK", "KN", "TN", "NT", "ZZ"}
	if contains(disallowedPrefixes, prefix) {
		return false
	}

	// First letter should not be D, F, I, Q, U, or V
	if strings.ContainsAny(string(prefix[0]), "DFIQUV") {
		return false
	}

	// Second letter should not be D, F, I, Q, U, V or O
	if strings.ContainsAny(string(prefix[1]), "DFIQUV") || prefix[1] == 'O' {
		return false
	}

	return true
}

// contains checks if a slice contains a specific string.
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
