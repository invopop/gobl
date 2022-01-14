package currency

import "github.com/invopop/gobl/num"

// ExchangeRates represents an array of currency exchange rates
type ExchangeRates []*ExchangeRate

// ExchangeRate contains data on the rate to be used when converting amounts from
// the document's base currency to whatever is defined. As a stand-alone object, this
// isn't much use. It must be used inside a bigger document.
type ExchangeRate struct {
	// ISO currency code this rate represents.
	Currency Code `json:"currency" jsonschema:"title=Currency"`
	// How much is 1.00 of the document's currency worth for this exchange rate.
	Amount num.Amount `json:"amount" jsonschema:"title=Amount"`
}
