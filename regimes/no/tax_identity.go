// Package no provides the tax identity validation specific to Norway.
package no

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var trnRegex = regexp.MustCompile(`^\d{9}$`)

func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateTRNCode)),
	)
}

func validateTRNCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	s := code.String()
	if !trnRegex.MatchString(s) {
		return errors.New("must be a 9-digit number")
	}
	if !validateChecksum(s) {
		return errors.New("invalid checksum for TRN")
	}
	return nil
}

func validateChecksum(trn string) bool {
	weights := []int{3, 2, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i, r := range trn[:8] {
		d, _ := strconv.Atoi(string(r))
		sum += d * weights[i]
	}
	mod := sum % 11
	chk := 11 - mod
	switch chk {
	case 10:
		return false
	case 11:
		chk = 0
	}
	last, _ := strconv.Atoi(string(trn[8]))
	return chk == last
}
