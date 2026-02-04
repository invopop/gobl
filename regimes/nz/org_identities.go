package nz

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

var nzbnPattern = regexp.MustCompile(`^\d{13}$`)

var orgIdentityDefinitions = []*cbc.Definition{
	{
		Key: org.IdentityKeyGLN,
		Name: i18n.String{
			i18n.EN: "NZ Business Number",
		},
		Desc: i18n.String{
			i18n.EN: "13-digit identifier based on the GS1 Global Location Number (GLN) standard, starting with NZ prefix 94.",
		},
	},
}

func normalizeOrgIdentity(id *org.Identity) {
	if id == nil || id.Key != org.IdentityKeyGLN {
		return
	}
	code := id.Code.String()
	code = strings.ReplaceAll(cbc.NormalizeString(code), "-", "")
	id.Code = cbc.Code(code)
}

func validateOrgIdentity(id *org.Identity) error {
	if id == nil || id.Key != org.IdentityKeyGLN {
		return nil
	}
	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.Required,
			validation.By(validateNZBNCode),
			validation.Skip,
		),
	)
}

func validateNZBNCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	return validateNZBN(code.String())
}

func validateNZBN(code string) error {
	if !nzbnPattern.MatchString(code) {
		return errors.New("NZBN must be exactly 13 digits")
	}

	if code[0:2] != "94" {
		return errors.New("NZBN must start with '94' (New Zealand GS1 prefix)")
	}

	digits := make([]int, 13)
	for i := 0; i < 13; i++ {
		var err error
		digits[i], err = strconv.Atoi(string(code[i]))
		if err != nil {
			return errors.New("NZBN must contain only digits")
		}
	}

	calculatedCheck := calculateGS1CheckDigit(digits[:12])

	if calculatedCheck != digits[12] {
		return errors.New("invalid NZBN: check digit mismatch")
	}

	return nil
}

func calculateGS1CheckDigit(digits []int) int {
	sum := 0
	for i := 0; i < 12; i++ {
		if i%2 == 0 {
			sum += digits[i] * 1
		} else {
			sum += digits[i] * 3
		}
	}

	remainder := sum % 10
	if remainder == 0 {
		return 0
	}
	return 10 - remainder
}
