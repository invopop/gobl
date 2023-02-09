package currency

import "github.com/invopop/gobl/num"

// ExchangeRate contains data on the rate to be used when converting amounts from
// the document's base currency to whatever is defined.
//
// It should be possible to take any amount in the matching currency and multiply it
// by the amount defined in the exchange rate to determine the value.
//
// For example, our document is in EUR and some amounts are defined in USD. Our
// ExchangeRate instance may be defined and used as:
//
//	  rate := &currency.ExchangeRate{
//		   Currency: currency.USD,
//	    Amount: "0.875967",
//	  }
//
//	  val := "100.00" // USD
//	  val.Multiply(rate.Amount) // EUR: "87.60"
type ExchangeRate struct {
	// ISO currency code this rate represents.
	Currency Code `json:"currency" jsonschema:"title=Currency"`
	// How much is 1.00 of this currency worth in the documents currency.
	Amount num.Amount `json:"amount" jsonschema:"title=Amount"`
}
