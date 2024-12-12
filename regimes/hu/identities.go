package hu

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

const (
	// IdentityKeyGroupNumber is the key used when a person belongs to a VAT group
	// The main tax identity field contains the vat number of the group and it is required to include
	// the vat number of the group member, which is included in identites.
	IdentityKeyGroupNumber cbc.Key = "hu-group-number"
)

var identityKeyDefinitions = []*cbc.KeyDefinition{
	{
		Key: IdentityKeyGroupNumber,
		Name: i18n.String{
			i18n.EN: "Group Member Number",
			i18n.HU: "Csoport tag adószáma",
		},
	},
}

func validateIdentity(id *org.Identity) error {
	if id == nil || id.Key != IdentityKeyGroupNumber {
		return nil
	}
	switch id.Key {
	case IdentityKeyGroupNumber:
		return validation.ValidateStruct(id,
			validation.Field(&id.Code, validation.By(validateGroupCode)),
		)
	default:
		return nil
	}
}

func validateGroupCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok {
		return nil
	}
	if code == "" {
		return nil
	}

	for _, v := range code {
		x := v - 48
		if x < 0 || x > 9 {
			return errors.New("contains invalid characters")
		}
	}

	if len(code) != 11 {
		return errors.New("invalid length")
	}

	str := code.String()

	// Calculate check-digit
	result := 9*int(str[0]-'0') + 7*int(str[1]-'0') + 3*int(str[2]-'0') + int(str[3]-'0') + 9*int(str[4]-'0') + 7*int(str[5]-'0') + 3*int(str[6]-'0')
	checkDigit := (10 - result%10) % 10

	compare, err := strconv.Atoi(string(code[7]))
	if err != nil {
		return fmt.Errorf("invalid check digit: %w", err)
	}
	if compare != checkDigit {
		return errors.New("checksum mismatch")
	}

	if len(code) == 11 {
		if !validAreaCodes[code[9:11]] {
			return errors.New("invalid area code")
		}
		if code[8:9] != "4" {
			return errors.New("invalid VAT code")
		}
	}
	return nil
}
