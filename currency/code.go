package currency

import "errors"

// Code is the ISO currency code
type Code string

// Def provides a structure for the currencies
type Def struct {
	Name    string `json:"name"`    // name of the currency
	Country string `json:"country"` // name of the country it belongs to
	Code    Code   `json:"code"`    // three-letter currency code
	Num     string `json:"num"`     // three-digit currency code
	Units   int    `json:"units"`   // how many cents are used for the currency
}

// Validate ensures the currency code is valid according
// to the ISO 4217 three-letter list.
func (c Code) Validate() error {
	if string(c) == "" {
		return nil
	}
	for _, cc := range iso4217Codes {
		if c == cc {
			return nil
		}
	}
	return errors.New("invalid language code")
}
