package es

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regions/common"
)

// TaxCodeType represents the types of tax code which are issued
// in Spain. The same general format with variations is used for
// national individuals, foreigners, and legal organizations.
type TaxCodeType string

// Supported tax code types.
const (
	NationalTaxCode     TaxCodeType = "N"
	ForeignTaxCode      TaxCodeType = "X"
	OrganizationTaxCode TaxCodeType = "B"
	OtherTaxCode        TaxCodeType = "O"
	UnknownTaxCode      TaxCodeType = "NA"
)

// tax ID standard tables
const (
	taxCodeCheckLetters       = "TRWAGMYFPDXBNJZSQVHLCKE"
	taxCodeForeignTypeLetters = "XYZ"
	taxCodeOtherTypeLetters   = "KLM"
	taxCodeOrgTypeLetters     = "ABCDEFGHJNPQRSUVW"
	taxCodeOrgCheckLetters    = "JABCDEFGHI"
)

const (
	taxCodeMatchType   = "type"
	taxCodeMatchNumber = "number"
	taxCodeMatchCheck  = "check"
)

// Standard simplified errors messages
var (
	ErrTaxCodeNoMatch      = errors.New("no match")
	ErrTaxCodeUnknownType  = errors.New("unknown type")
	ErrTaxCodeInvalidCheck = errors.New("check letter is invalid")
)

// Known combinations of codes
var (
	taxCodeCountryRegexp  = regexp.MustCompile(`^ES`)
	taxCodeNationalRegexp = regexp.MustCompile(`^(?P<number>[0-9]{8})(?P<check>[` + taxCodeCheckLetters + `])$`)
	taxCodeForeignRegexp  = regexp.MustCompile(`^(?P<type>[` + taxCodeForeignTypeLetters + `])(?P<number>[0-9]{7})(?P<check>[` + taxCodeCheckLetters + `])$`)
	taxCodeOtherRegexp    = regexp.MustCompile(`^(?P<type>[` + taxCodeOtherTypeLetters + `])(?P<number>[0-9]{7})(?P<check>[0-9` + taxCodeOrgCheckLetters + `])$`)
	taxCodeOrgRegexp      = regexp.MustCompile(`^(?P<type>[` + taxCodeOrgTypeLetters + `])(?P<number>[0-9]{7})(?P<check>[0-9` + taxCodeOrgCheckLetters + `])$`)
)

// ValidateTaxIdentity looks at the provided identity's code,
// determines the type, and performs the calculations
// required to determine if it is valid.
// These methods assume the code has already been normalized
// and thus only contains upper-case letters and numbers with
// no white space.
func ValidateTaxIdentity(tID *org.TaxIdentity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateTaxCode)),
	)
}

func validateTaxCode(value interface{}) error {
	code, ok := value.(string)
	if !ok {
		return nil
	}
	if code == "" {
		return nil
	}
	typ, err := DetermineTaxCodeType(code)
	if typ == UnknownTaxCode {
		return ErrTaxCodeUnknownType
	}
	return err
}

// NormalizeTaxIdentity removes any whitespace or separation characters and ensures all letters are
// uppercase. It'll also remove the "ES" part at beginning if present such as required
// for EU VIES system which is redundant and not used in the validation process.
func NormalizeTaxIdentity(tID *org.TaxIdentity) error {
	if err := common.NormalizeTaxIdentity(tID); err != nil {
		return err
	}
	tID.Code = taxCodeCountryRegexp.ReplaceAllString(tID.Code, "")
	return nil
}

// DetermineTaxCodeType takes a valid code and determines the type. If the code
// is not valid, the `UnknownTaxCode` type will be returned.
func DetermineTaxCodeType(code string) (TaxCodeType, error) {
	switch {
	case taxCodeOrgRegexp.MatchString(code):
		return OrganizationTaxCode, verifyOrgCode(code)
	case taxCodeNationalRegexp.MatchString(code):
		return NationalTaxCode, verifyNationalCode(code)
	case taxCodeForeignRegexp.MatchString(code):
		return ForeignTaxCode, verifyForeignCode(code)
	case taxCodeOtherRegexp.MatchString(code):
		return OtherTaxCode, verifyOtherCode(code)
	default:
		return UnknownTaxCode, nil
	}
}

func verifyNationalCode(code string) error {
	m, err := extractMatches(taxCodeNationalRegexp, code)
	if err != nil {
		return err
	}

	if m[taxCodeMatchNumber] == "00000000" {
		// exception case
		return ErrTaxCodeInvalidCheck
	}

	n, _ := strconv.Atoi(m[taxCodeMatchNumber])
	if []rune(taxCodeCheckLetters)[n%23] != []rune(m[taxCodeMatchCheck])[0] {
		return ErrTaxCodeInvalidCheck
	}

	return nil // success
}

func verifyForeignCode(code string) error {
	m, err := extractMatches(taxCodeForeignRegexp, code)
	if err != nil {
		return err
	}

	// Extract index from type letters
	ti := strings.Index(taxCodeForeignTypeLetters, m[taxCodeMatchType])
	ft := strconv.Itoa(ti)

	fs := ft + m[taxCodeMatchNumber]
	ci, _ := strconv.Atoi(fs)

	if []rune(taxCodeCheckLetters)[ci%23] != []rune(m[taxCodeMatchCheck])[0] {
		return ErrTaxCodeInvalidCheck
	}

	return nil // success
}

func verifyOrgCode(code string) error {
	m, err := extractMatches(taxCodeOrgRegexp, code)
	if err != nil {
		return err
	}
	return verifyOrgCodeMatches(m)
}

func verifyOtherCode(code string) error {
	m, err := extractMatches(taxCodeOtherRegexp, code)
	if err != nil {
		return err
	}
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
		return ErrTaxCodeInvalidCheck
	}

	return nil
}

// regex handling is a bit long winded, this helper makes it easier to extract
// named matches.
func extractMatches(regex *regexp.Regexp, code string) (map[string]string, error) {
	m := regex.FindStringSubmatch(code)
	if len(m) == 0 {
		return nil, ErrTaxCodeNoMatch
	}

	r := make(map[string]string)
	for i, n := range regex.SubexpNames() {
		if i != 0 && n != "" {
			r[n] = m[i]
		}
	}

	return r, nil
}
