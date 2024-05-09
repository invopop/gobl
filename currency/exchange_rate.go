package currency

import (
	"fmt"

	"github.com/invopop/gobl/num"
	"github.com/invopop/validation"
)

// ExchangeRate contains data on the rate to be used when converting amounts from
// one currency into another.
//
// For reference, naming here is based on the following english grammar examples:
// - Exchange from USD to EUR.
// - Convert from USD into EUR.
//
// It should be possible to take any amount in the matching currency and multiply it
// by the amount defined in the exchange rate to determine the value.
//
// For example, our document is in EUR and some amounts are defined in USD. Our
// ExchangeRate instance may be defined and used as:
//
//	  rate := &currency.ExchangeRate{
//		From: currency.USD,
//		To: currency.EUR,
//		Amount: "0.875967",
//	  }
//
//	  val := MakeAmount(100, 2) // 100.00 USD
//	  val.Multiply(rate.Amount) // 87.60 EUR
type ExchangeRate struct {
	// Currency code this will be converted from.
	From Code `json:"from" jsonschema:"title=From"`
	// Currency code this exchange rate will convert into.
	To Code `json:"to" jsonschema:"title=To"`
	// How much is 1 of the "from" currency worth in the "to" currency.
	Amount num.Amount `json:"amount" jsonschema:"title=Amount"`
}

// Validate ensures the content of the exchange rate looks good.
func (er *ExchangeRate) Validate() error {
	return validation.ValidateStruct(er,
		validation.Field(&er.From, validation.Required),
		validation.Field(&er.To, validation.Required),
		validation.Field(&er.Amount, num.Positive),
	)
}

// MatchExchangeRate will attempt to find the matching exchange rate that
// will convert from one currency into another. If no match is found,
// nil is returned. If the "from" and "to" parameters are the same, `1` will
// be provided instead.
func MatchExchangeRate(rates []*ExchangeRate, from, to Code) *num.Amount {
	if from == to {
		return num.NewAmount(1, 0)
	}
	for _, rate := range rates {
		if rate.From == from && rate.To == to {
			a := rate.Amount
			return &a
		}
	}
	return nil
}

// Exchange will convert the provided amount from one currency into another or return
// nil if no match can be found.
func Exchange(rates []*ExchangeRate, from, to Code, amount num.Amount) *num.Amount {
	if from == to {
		return &amount
	}
	rate := MatchExchangeRate(rates, from, to)
	if rate == nil {
		return nil
	}
	a := amount.Multiply(*rate)
	z := to.Def().Zero()
	a = a.Rescale(z.Exp()) // ensure scale always matches destination currency
	return &a
}

type exchangeRateValidation struct {
	rates []*ExchangeRate
	to    Code
}

// Validate performs validation on the provided value to see if it
// is present in the exchange rates.
func (erv *exchangeRateValidation) Validate(val any) error {
	cur, ok := val.(Code)
	if !ok || cur == CodeEmpty {
		return nil
	}
	if cur == erv.to {
		return nil
	}
	for _, r := range erv.rates {
		if r.From == cur && r.To == erv.to {
			return nil
		}
	}
	return fmt.Errorf("no exchange rate defined for '%v' to '%v'", cur, erv.to)
}

// CanExchangeTo will check to see if the currency to be validated can
// be converted into one of the provided rates.
func CanExchangeTo(rates []*ExchangeRate, to Code) validation.Rule {
	return &exchangeRateValidation{
		rates: rates,
		to:    to,
	}
}
