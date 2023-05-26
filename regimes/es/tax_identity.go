package es

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// TaxCodeType represents the types of tax code which are issued
// in Spain. The same general format with variations is used for
// national individuals, foreigners, and legal organizations.
type TaxCodeType string

// Supported tax code types, from code itself.
const (
	NationalTaxCode     TaxCodeType = "N"
	ForeignTaxCode      TaxCodeType = "X"
	OrganizationTaxCode TaxCodeType = "B"
	OtherTaxCode        TaxCodeType = "O"
	UnknownTaxCode      TaxCodeType = "NA"
)

// The tax identity type is required for TicketBAI documents
// in the Basque Country.
const (
	TaxIdentityTypeFiscal   cbc.Key = "fiscal"
	TaxIdentityTypePassport cbc.Key = "passport"
	TaxIdentityTypeForeign  cbc.Key = "foreign"
	TaxIdentityTypeResident cbc.Key = "resident"
	TaxIdentityTypeOther    cbc.Key = "other"
)

var taxIdentityTypeDefinitions = []*tax.KeyDefinition{
	{
		Key: TaxIdentityTypeFiscal,
		Name: i18n.String{
			i18n.EN: "National Tax Identity",
			i18n.ES: "Número de Identificación Fiscal",
		},
		Codes: cbc.CodeSet{
			KeyTicketBAIIDType: "02",
		},
	},
	{
		Key: TaxIdentityTypePassport,
		Name: i18n.String{
			i18n.EN: "Passport",
			i18n.ES: "Pasaporte",
		},
		Codes: cbc.CodeSet{
			KeyTicketBAIIDType: "03",
		},
	},
	{
		Key: TaxIdentityTypeForeign,
		Name: i18n.String{
			i18n.EN: "National ID Card or similar from a foreign country",
			i18n.ES: "Documento oficial de identificación expedido por el país o territorio de residencia",
		},
		Codes: cbc.CodeSet{
			KeyTicketBAIIDType: "04",
		},
	},
	{
		Key: TaxIdentityTypeResident,
		Name: i18n.String{
			i18n.EN: "Residential permit",
			i18n.ES: "Certificado de residencia",
		},
		Codes: cbc.CodeSet{
			KeyTicketBAIIDType: "05",
		},
	},
	{
		Key: TaxIdentityTypeOther,
		Name: i18n.String{
			i18n.EN: "An other type of source not listed",
			i18n.ES: "Otro documento probatorio",
		},
		Codes: cbc.CodeSet{
			KeyTicketBAIIDType: "06",
		},
	},
}

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
	taxCodeNationalRegexp = regexp.MustCompile(`^(?P<number>[0-9]{8})(?P<check>[` + taxCodeCheckLetters + `])$`)
	taxCodeForeignRegexp  = regexp.MustCompile(`^(?P<type>[` + taxCodeForeignTypeLetters + `])(?P<number>[0-9]{7})(?P<check>[` + taxCodeCheckLetters + `])$`)
	taxCodeOtherRegexp    = regexp.MustCompile(`^(?P<type>[` + taxCodeOtherTypeLetters + `])(?P<number>[0-9]{7})(?P<check>[0-9` + taxCodeOrgCheckLetters + `])$`)
	taxCodeOrgRegexp      = regexp.MustCompile(`^(?P<type>[` + taxCodeOrgTypeLetters + `])(?P<number>[0-9]{7})(?P<check>[0-9` + taxCodeOrgCheckLetters + `])$`)
)

// validateTaxIdentity looks at the provided identity's code,
// determines the type, and performs the calculations
// required to determine if it is valid.
// These methods assume the code has already been normalized
// and thus only contains upper-case letters and numbers with
// no white space.
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
	typ, err := DetermineTaxCodeType(code)
	if typ == UnknownTaxCode {
		return ErrTaxCodeUnknownType
	}
	return err
}

// normalizeTaxIdentity removes any whitespace or separation characters and ensures all letters are
// uppercase. It'll also remove the "ES" part at beginning if present such as required
// for EU VIES system which is redundant and not used in the validation process.
func normalizeTaxIdentity(tID *tax.Identity) error {
	return common.NormalizeTaxIdentity(tID)
}

// DetermineTaxCodeType takes a valid code and determines the type. If the code
// is not valid, the `UnknownTaxCode` type will be returned.
func DetermineTaxCodeType(code cbc.Code) (TaxCodeType, error) {
	switch {
	case taxCodeOrgRegexp.MatchString(string(code)):
		return OrganizationTaxCode, verifyOrgCode(code)
	case taxCodeNationalRegexp.MatchString(string(code)):
		return NationalTaxCode, verifyNationalCode(code)
	case taxCodeForeignRegexp.MatchString(string(code)):
		return ForeignTaxCode, verifyForeignCode(code)
	case taxCodeOtherRegexp.MatchString(string(code)):
		return OtherTaxCode, verifyOtherCode(code)
	default:
		return UnknownTaxCode, nil
	}
}

func verifyNationalCode(code cbc.Code) error {
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

func verifyForeignCode(code cbc.Code) error {
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

func verifyOrgCode(code cbc.Code) error {
	m, err := extractMatches(taxCodeOrgRegexp, code)
	if err != nil {
		return err
	}
	return verifyOrgCodeMatches(m)
}

func verifyOtherCode(code cbc.Code) error {
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
func extractMatches(regex *regexp.Regexp, code cbc.Code) (map[string]string, error) {
	m := regex.FindStringSubmatch(code.String())
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
