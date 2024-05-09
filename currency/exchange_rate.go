package currency

import (
	"github.com/invopop/gobl/num"
	"github.com/invopop/validation"
)

// ExchangeRate contains data on the rate to be used when converting amounts from
// one currency into another.
//
// It should be possible to take any amount in the matching currency and multiply it
// by the amount defined in the exchange rate to determine the value.
//
// For example, our document is in EUR and some amounts are defined in USD. Our
// ExchangeRate instance may be defined and used as:
//
//	  rate := &currency.ExchangeRate{
//		From: currency.USD,
//		Into: currency.EUR,
//		Amount: "0.875967",
//	  }
//
//	  val := "100.00" // USD
//	  val.Multiply(rate.Amount) // EUR: "87.60"
type ExchangeRate struct {
	// Currency code this will be converted from.
	From Code `json:"from" jsonschema:"title=From"`
	// Currency code this exchange rate will convert into.
	Into Code `json:"into" jsonschema:"title=Into"`
	// How much is 1 of the "from" currency worth in the "into" currency.
	Amount num.Amount `json:"amount" jsonschema:"title=Amount"`
}

// Validate ensures the content of the exchange rate looks good.
func (er *ExchangeRate) Validate() error {
	return validation.ValidateStruct(er,
		validation.Field(&er.From, validation.Required),
		validation.Field(&er.Into, validation.Required),
		validation.Field(&er.Amount, num.NotZero),
	)
}

// MatchExchangeRate will attempt to find the matching exchange rate that
// will convert from one currency into another. If no match is found,
// nil is returned. If the from and into parameters are the same, `1` will
// be provided instead.
func MatchExchangeRate(rates []*ExchangeRate, from, into Code) *num.Amount {
	if from == into {
		return num.NewAmount(1, 0)
	}
	for _, rate := range rates {
		if rate.From == from && rate.Into == into {
			a := rate.Amount
			return &a
		}
	}
	return nil
}

// Exchange will convert the provided amount from one currency into another or return
// nil if no match can be found.
func Exchange(rates []*ExchangeRate, from, into Code, amount num.Amount) *num.Amount {
	if from == into {
		return &amount
	}
	rate := MatchExchangeRate(rates, from, into)
	if rate == nil {
		return nil
	}
	a := amount.Multiply(*rate)
	z := into.Def().Zero()
	a = a.Rescale(z.Exp()) // ensure scale always matches destination currency
	return &a
}
