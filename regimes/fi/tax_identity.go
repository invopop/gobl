package fi

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// TaxCodeType represents the types of tax codes issued in Finland.
type TaxCodeType string

// Supported tax code types.
const (
	NationalTaxCode     TaxCodeType = "H" // Finnish Personal Identity Code (PIC) - (Henkilötunnus)
	ForeignTaxCode      TaxCodeType = "V" // Veronumero; temporary workers, especially in construction or shipbuilding
	OrganizationTaxCode TaxCodeType = "Y" // Y-tunnus
	UnknownTaxCode      TaxCodeType = "NA"
)

var (
	fullTaxCodeNationalRegexp = regexp.MustCompile(`^(\d{6})([+\-ABCDEFXYWVU])(.{4})$`) // includes the century mark
	taxCodeNationalRegexp     = regexp.MustCompile(`^(\d{6})(\d{3})(\d|[A-Z])$`)        // PIC pattern: DDMMYYZZZT
	taxCodeOrgRegexp          = regexp.MustCompile(`^(\d{7})(\d)$`)                     // (Y-tunnus) pattern: NNNNNNNC
	taxCodeForeignRegexp      = regexp.MustCompile(`^\d{12}$`)                          // Veronumero pattern: 12 digits
)

// Check characters for Finnish PIC
const taxCodeCheckCharacters = "0123456789ABCDEFHJKLMNPRSTUVWXY"

// Standard simplified errors messages
var (
	ErrTaxCodeNoMatch      = errors.New("no match")
	ErrTaxCodeInvalidDate  = errors.New("invalid date")
	ErrTaxCodeUnknownType  = errors.New("unknown type")
	ErrTaxCodeInvalidCheck = errors.New("check character is invalid")
)

// This method assumes the code has already been normalized
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

// DetermineTaxCodeType takes a valid code and determines the type.
// If the code is not valid, an error will be returned.
func DetermineTaxCodeType(code cbc.Code) (TaxCodeType, error) {
	switch {
	case taxCodeNationalRegexp.MatchString(string(code)):
		return NationalTaxCode, verifyNationalCode(code)
	case taxCodeOrgRegexp.MatchString(string(code)):
		return OrganizationTaxCode, verifyOrgCode(code)
	case taxCodeForeignRegexp.MatchString(string(code)):
		return ForeignTaxCode, nil // return nil, since foreign codes have no validation rules
	default:
		return UnknownTaxCode, nil
	}
}

// sources of truth:
// https://dvv.fi/en/personal-identity-code
// https://www.finlex.fi/fi/lainsaadanto/2009/661#chp_2__sec_11__heading
// https://www.vero.fi/globalassets/tietoa-verohallinnosta/ohjelmistokehittajille/finnish-tax-administration_the-control-character-for-verifying-the-authenticity-of-finnish-business-ids-and-personal-identity-codes.pdf
func verifyNationalCode(code cbc.Code) error {
	matches := taxCodeNationalRegexp.FindStringSubmatch(string(code))

	datePart := matches[1]    // DDMMYY
	personalNum := matches[2] // ZZZ
	checkChar := matches[3]   // Q

	// Validate the date portion
	day, err := strconv.Atoi(datePart[0:2])
	if err != nil || day < 1 || day > 31 {
		return ErrTaxCodeInvalidDate
	}
	month, err := strconv.Atoi(datePart[2:4])
	if err != nil || month < 1 || month > 12 {
		return ErrTaxCodeInvalidDate
	}

	numStr := datePart + personalNum
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return ErrTaxCodeInvalidCheck
	}

	// Calculate check digit
	checkDigit := num % 31
	expectedCheck := taxCodeCheckCharacters[checkDigit]

	// Compare with provided check character
	if string(expectedCheck) != checkChar {
		return ErrTaxCodeInvalidCheck
	}

	return nil
}

// sources of truth:
// https://taxid.pro/
// https://www.vero.fi/globalassets/tietoa-verohallinnosta/ohjelmistokehittajille/finnish-tax-administration_the-control-character-for-verifying-the-authenticity-of-finnish-business-ids-and-personal-identity-codes.pdf
func verifyOrgCode(code cbc.Code) error {
	matches := taxCodeOrgRegexp.FindStringSubmatch(string(code))

	base := matches[1]       // 7-digit base number
	checkDigit := matches[2] // Check digit

	// Weight factors for Finnish Business ID validation
	weights := []int{7, 9, 10, 5, 8, 4, 2}

	sum := 0
	for i, r := range base {
		digit, err := strconv.Atoi(string(r))
		if err != nil {
			return ErrTaxCodeInvalidCheck
		}
		sum += digit * weights[i]
	}

	// Calculate expected check digit
	remainder := sum % 11
	expectedCheck := 11 - remainder

	// Handle special cases
	if expectedCheck == 11 {
		expectedCheck = 0
	} else if expectedCheck == 10 {
		return ErrTaxCodeInvalidCheck
	}

	providedCheck, err := strconv.Atoi(checkDigit)
	if err != nil {
		return ErrTaxCodeInvalidCheck
	}

	if expectedCheck != providedCheck {
		return ErrTaxCodeInvalidCheck
	}

	return nil
}

func normalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil {
		return
	}

	tax.NormalizeIdentity(tID, l10n.FI)

	if fullTaxCodeNationalRegexp.MatchString(string(tID.Code)) {
		matches := fullTaxCodeNationalRegexp.FindStringSubmatch(string(tID.Code))
		// remove the century mark
		tID.Code = cbc.Code(matches[1] + matches[3])
	}
}
