package hu

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// There are also some rules for individuals' tax codes, that can be found in this document: https://web.archive.org/web/20200622115229/https://ceginformaciosszolgalat.kormany.hu/download/4/27/02000/adoszam_adoazonosito_ellenorzes_2018.pdf

// Number 4 is only valid for the group tax subject to VAT (second tax id)
var (
	validVatCodes = map[cbc.Code]bool{
		"1": true, "2": true, "3": true, "5": true,
	}

	validAreaCodes = map[cbc.Code]bool{
		"02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true,
		"10": true, "11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true,
		"18": true, "19": true, "20": true, "22": true, "23": true, "24": true, "25": true, "26": true,
		"27": true, "28": true, "29": true, "30": true, "31": true, "32": true, "33": true, "34": true,
		"35": true, "36": true, "37": true, "38": true, "39": true, "40": true, "41": true, "42": true,
		"43": true, "44": true, "51": true,
	}
)

// validateTaxIdentity checks to ensure the NIT code looks okay.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateTaxCode)),
	)
}

//We can have 2 different lenghts of tax code, 8 and 11 characters (8+1+2)

func validateTaxCode(value interface{}) error {
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

	if len(code) != 8 && len(code) != 11 {
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
		if !validVatCodes[code[8:9]] {
			return errors.New("invalid VAT code")
		}
	}
	return nil
}
