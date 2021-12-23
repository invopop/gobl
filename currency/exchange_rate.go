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
	// Rate to apply when converting the document's currency to this one.
	Value num.Amount `json:"value" jsonschema:"title=Value"`
}
