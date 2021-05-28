package currency

import "github.com/invopop/gobl/num"

// ExchangeRates represents an array of currency exchange rates
type ExchangeRates []ExchangeRate

// ExchangeRate contains data on the rate to be used when converting data.
type ExchangeRate struct {
	Currency string     `json:"currency"`
	Value    num.Amount `json:"value"`
}
