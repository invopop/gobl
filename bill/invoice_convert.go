package bill

import (
	"fmt"

	"github.com/invopop/gobl/currency"
)

// ConvertInto will use the defined exchange rates in the invoice to convert all the prices
// into the given currency.
//
// The intent of this method is help convert the invoice amounts when the destination is
// unable or unwilling to handle the current currency. This is typically the case
// with tax related reports or declarations.
//
// The method will return a new invoice with all the amounts converted into the given
// currency or an error if the conversion is not possible.
//
// Conversion is done by first exchanging the lowest common amounts to the destination
// currency, then recalculating the totals.
func (inv *Invoice) ConvertInto(cur currency.Code) (*Invoice, error) {
	// Calculate ensures that all the totals and amounts have been prepared
	// so we can make assumptions about the data that will be available,
	// including the original currency!
	if err := inv.Calculate(); err != nil {
		return nil, err
	}

	if inv.Currency == cur {
		return inv, nil
	}
	ex := currency.MatchExchangeRate(inv.ExchangeRates, inv.Currency, cur)
	if ex == nil {
		return nil, fmt.Errorf("no exchange rate defined for '%v' to '%v'", inv.Currency, cur)
	}

	i2 := *inv
	i2.Totals = nil
	i2.Lines = convertLinesInto(ex, inv.Lines)
	i2.Discounts = convertDiscountsInto(ex, inv.Discounts)
	i2.Charges = convertChargesInto(ex, inv.Charges)
	i2.Payment = convertPaymentDetailsInto(ex, inv.Payment)
	i2.Currency = cur

	if err := i2.Calculate(); err != nil {
		return nil, err
	}

	return &i2, nil
}
