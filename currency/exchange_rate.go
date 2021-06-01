package currency

import "github.com/invopop/gobl/num"

// ExchangeRates represents an array of currency exchange rates
type ExchangeRates []ExchangeRate

// ExchangeRate contains data on the rate to be used when converting amounts from
// the document's base currency to whatever is defined. As a stand-alone object, this
// isn't much use. It must be used inside a bigger document.
type ExchangeRate struct {
	Currency Code       `json:"currency"`
	Value    num.Amount `json:"value"`
}
