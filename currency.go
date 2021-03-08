package gobl

import "github.com/invopop/gobl/num"

// ExchangeRates represents an array of currency exchange rates
type ExchangeRates []ExchangeRate

// ExchangeRate contains data on the rate to be used when converting data.
// The rate is always multipled by the
type ExchangeRate struct {
	Currency string     `json:"currency"`
	Rate     num.Amount `json:"rate"`
}
