package currency

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
)

// ExchangeRate contains data on the rate to be used when converting amounts from
// one currency into another.
//
// For reference, naming here is based on the following english grammar examples:
// - Exchange from USD to EUR.
// - Convert from USD into EUR.
//
// If the destination or document's currency is EUR and some amounts
// are defined in USD, the `ExchangeRate` instance may be defined and used
// as follows:
//
//	  rate := &currency.ExchangeRate{
//		From:   currency.USD,
//		To:     currency.EUR,
//		Amount: "0.875967",
//	  }
//
//	  val := MakeAmount(100, 2) // 100.00 USD
//	  rate.Convert(val)         // 87.60 EUR
type ExchangeRate struct {
	// Currency code this will be converted from.
	From Code `json:"from" jsonschema:"title=From"`
	// Currency code this exchange rate will convert into.
	To Code `json:"to" jsonschema:"title=To"`
	// At represents the effective date and time at which the exchange rate
	// is determined by the source. The time may be zero if referring to a
	// specific day only.
	At *cal.DateTime `json:"at,omitempty" jsonschema:"title=At"`
	// Source key provides a reference to the source the exchange rate was
	// obtained from. Typically this will be determined by an application
	// used to update exchange rates automatically.
	Source cbc.Key `json:"source,omitempty" jsonschema:"title=Source"`
	// How much is 1 of the "from" currency worth in the "to" currency.
	Amount num.Amount `json:"amount" jsonschema:"title=Amount"`
}

// Convert performs the currency conversion defined by the exchange rate.
func (er *ExchangeRate) Convert(amount num.Amount) num.Amount {
	a := amount.Multiply(er.Amount)
	z := er.To.Def().Zero()
	return a.Rescale(z.Exp()) // ensure scale always matches destination currency
}

// MatchExchangeRate will attempt to find the matching exchange rate that
// will convert from one currency into another. Will return nil if no
// match is found or the currencies are the same.
func MatchExchangeRate(rates []*ExchangeRate, from, to Code) *ExchangeRate {
	if from == to {
		return nil
	}
	for _, rate := range rates {
		if rate.From == from && rate.To == to {
			return rate
		}
	}
	return nil
}

// Convert will convert the provided amount from one currency into another or return
// nil if no match can be found. If the currencies are the same, the original
// amount will be returned.
func Convert(rates []*ExchangeRate, from, to Code, amount num.Amount) *num.Amount {
	if from == to {
		return &amount
	}
	if rate := MatchExchangeRate(rates, from, to); rate != nil {
		a := rate.Convert(amount)
		return &a
	}
	return nil
}
