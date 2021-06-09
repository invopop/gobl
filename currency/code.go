package currency

import (
	"errors"

	"github.com/invopop/gobl/num"
)

// Code is the ISO currency code
type Code string

// Def provides a structure for the currencies
type Def struct {
	Name  string `json:"name"`  // name of the currency
	Code  Code   `json:"code"`  // three-letter currency code
	Num   string `json:"num"`   // three-digit currency code
	Units uint32 `json:"units"` // how many cents are used for the currency
}

// Validate ensures the currency code is valid according
// to the ISO 4217 three-letter list.
func (c Code) Validate() error {
	if string(c) == "" {
		return nil
	}
	if _, ok := iso4217Defs[c]; ok {
		return nil
	}
	return errors.New("invalid language code")
}

// Get provides the code's currency definition, or
// false if none is found.
func Get(c Code) (Def, bool) {
	d, ok := iso4217Defs[c]
	return d, ok
}

// BaseAmount provides a definition's zero amount with the correct decimal
// places so that it can be used as a base for calculating totals.
func (d Def) BaseAmount() num.Amount {
	return num.MakeAmount(0, d.Units)
}
