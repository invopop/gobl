package es

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Tax Identity keys that may be determined from the code.
const (
	TaxIdentityNational  cbc.Key = "national"
	TaxIdentityForeigner cbc.Key = "foreigner"
	TaxIdentityOrg       cbc.Key = "org"
	TaxIdentityOther     cbc.Key = "other"
)

// tax ID standard tables
const (
	taxCodeCheckLetters         = "TRWAGMYFPDXBNJZSQVHLCKE"
	taxCodeForeignerTypeLetters = "XYZ"
	taxCodeOtherTypeLetters     = "KLM"
	taxCodeOrgTypeLetters       = "ABCDEFGHJNPQRSUVW"
	taxCodeOrgCheckLetters      = "JABCDEFGHI"
)

const (
	taxCodeMatchType   = "type"
	taxCodeMatchNumber = "number"
	taxCodeMatchCheck  = "check"
)

// Known combinations of codes
var (
	taxCodeNumber          = regexp.MustCompile(`^[0-9]`)
	taxCodeNationalRegexp  = regexp.MustCompile(`^(?P<number>[0-9]{8})(?P<check>[` + taxCodeCheckLetters + `])$`)                                                  //nolint:goconst
	taxCodeForeignerRegexp = regexp.MustCompile(`^(?P<type>[` + taxCodeForeignerTypeLetters + `])(?P<number>[0-9]{7})(?P<check>[` + taxCodeCheckLetters + `])$`)   //nolint:goconst
	taxCodeOtherRegexp     = regexp.MustCompile(`^(?P<type>[` + taxCodeOtherTypeLetters + `])(?P<number>[0-9]{7})(?P<check>[0-9` + taxCodeOrgCheckLetters + `])$`) //nolint:goconst
	taxCodeOrgRegexp       = regexp.MustCompile(`^(?P<type>[` + taxCodeOrgTypeLetters + `])(?P<number>[0-9]{7})(?P<check>[0-9` + taxCodeOrgCheckLetters + `])$`)   //nolint:goconst
)

var (
	errTaxIdentityCodeInvalidFormat   = errors.New("invalid format")
	errTaxIdentityCodeInvalidChecksum = errors.New("invalid check digit")
)

// normalizeTaxIdentity removes any whitespace or separation characters and ensures all letters are
// uppercase. It'll also remove the "ES" part at beginning if present such as required
// for EU VIES system which is redundant and not used in the validation process.
func normalizeTaxIdentity(tID *tax.Identity) {
	tax.NormalizeIdentity(tID)
	if taxCodeNumber.MatchString(tID.Code.String()) {
		tID.Code = cbc.Code(fmt.Sprintf("%09s", tID.Code.String()))
	}
}

// TaxIdentityKey determines the type of tax code and returns the appropriate key.
// An empty key will be returned if the code is not recognized. This will only
// work correctly if the tax identity has been normalized.
func TaxIdentityKey(tID *tax.Identity) cbc.Key {
	if tID == nil || tID.Code == "" || tID.Country != l10n.ES.Tax() {
		return cbc.KeyEmpty
	}
	return taxIdentityKey(tID.Code.String())
}

func taxIdentityKey(code string) cbc.Key {
	switch {
	case taxCodeOrgRegexp.MatchString(code):
		return TaxIdentityOrg
	case taxCodeNationalRegexp.MatchString(code):
		return TaxIdentityNational
	case taxCodeForeignerRegexp.MatchString(code):
		return TaxIdentityForeigner
	case taxCodeOtherRegexp.MatchString(code):
		return TaxIdentityOther
	default:
		return cbc.KeyEmpty
	}
}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn("ES"),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Spanish VAT identity code format or checksum",
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
	return validateTaxIdentityCode(code) == nil
}

func validateTaxIdentityCode(code cbc.Code) error {
	switch taxIdentityKey(code.String()) {
	case TaxIdentityNational:
		return verifyNationalCode(code)
	case TaxIdentityForeigner:
		return verifyForeignCode(code)
	case TaxIdentityOrg:
		return verifyOrgCode(code)
	case TaxIdentityOther:
		return verifyOtherCode(code)
	default:
		return errTaxIdentityCodeInvalidFormat
	}
}

func verifyNationalCode(code cbc.Code) error {
	m := extractMatches(taxCodeNationalRegexp, code)

	if m[taxCodeMatchNumber] == "00000000" {
		// exception case
		return errTaxIdentityCodeInvalidFormat
	}

	n, _ := strconv.Atoi(m[taxCodeMatchNumber])
	if []rune(taxCodeCheckLetters)[n%23] != []rune(m[taxCodeMatchCheck])[0] {
		return errTaxIdentityCodeInvalidChecksum
	}

	return nil // success
}

func verifyForeignCode(code cbc.Code) error {
	m := extractMatches(taxCodeForeignerRegexp, code)

	// Extract index from type letters
	ti := strings.Index(taxCodeForeignerTypeLetters, m[taxCodeMatchType])
	ft := strconv.Itoa(ti)

	fs := ft + m[taxCodeMatchNumber]
	ci, _ := strconv.Atoi(fs)

	if []rune(taxCodeCheckLetters)[ci%23] != []rune(m[taxCodeMatchCheck])[0] {
		return errTaxIdentityCodeInvalidChecksum
	}

	return nil // success
}

func verifyOrgCode(code cbc.Code) error {
	m := extractMatches(taxCodeOrgRegexp, code)
	return verifyOrgCodeMatches(m)
}

func verifyOtherCode(code cbc.Code) error {
	m := extractMatches(taxCodeOtherRegexp, code)
	return verifyOrgCodeMatches(m)
}

func verifyOrgCodeMatches(m map[string]string) error {
	num := []rune(m[taxCodeMatchNumber])
	p := make([]int, len(num))
	for i, v := range num {
		p[i], _ = strconv.Atoi(string(v))
	}

	sumEven := 0
	sumOdd := 0
	for k, v := range p {
		switch k & 1 {
		case 1:
			sumEven += v
		case 0:
			v = v * 2
			if v > 9 {
				v = v - 9
			}
			sumOdd += v
		}
	}

	// Calculate check digit
	cdc := (10 - (sumEven+sumOdd)%10) % 10

	// Extract digit to compare against
	cds := m[taxCodeMatchCheck]
	var cdi int
	if i := strings.Index(taxCodeOrgCheckLetters, cds); i != -1 {
		cdi = i
	} else {
		cdi, _ = strconv.Atoi(cds)
	}

	// compare
	if cdc != cdi {
		return errTaxIdentityCodeInvalidChecksum
	}

	return nil
}

// regex handling is a bit long winded, this helper makes it easier to extract
// named matches.
func extractMatches(regex *regexp.Regexp, code cbc.Code) map[string]string {
	m := regex.FindStringSubmatch(code.String())
	r := make(map[string]string)
	for i, n := range regex.SubexpNames() {
		if i != 0 && n != "" {
			r[n] = m[i]
		}
	}
	return r
}
