package bill

import (
	"fmt"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/pay"
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
	i2.Totals = new(Totals)
	i2.Lines = inv.converLines(ex)
	i2.Discounts = inv.convertDiscounts(ex)
	i2.Charges = inv.convertCharges(ex)
	i2.Outlays = inv.convertOutlays(ex)
	i2.Payment = inv.convertPayment(ex)
	i2.Currency = cur

	if err := i2.Calculate(); err != nil {
		return nil, err
	}

	return &i2, nil
}

func (inv *Invoice) converLines(ex *currency.ExchangeRate) []*Line {
	if len(inv.Lines) == 0 {
		return nil
	}
	lines := make([]*Line, len(inv.Lines))
	for i, l := range inv.Lines {
		lines[i] = l.convertInto(ex)
	}
	return lines
}

func (inv *Invoice) convertDiscounts(ex *currency.ExchangeRate) []*Discount {
	if len(inv.Discounts) == 0 {
		return nil
	}
	discounts := make([]*Discount, len(inv.Discounts))
	for i, d := range inv.Discounts {
		discounts[i] = d.convertInto(ex)
	}
	return discounts
}

func (inv *Invoice) convertCharges(ex *currency.ExchangeRate) []*Charge {
	if len(inv.Charges) == 0 {
		return nil
	}
	charges := make([]*Charge, len(inv.Charges))
	for i, c := range inv.Charges {
		charges[i] = c.convertInto(ex)
	}
	return charges
}

func (inv *Invoice) convertOutlays(ex *currency.ExchangeRate) []*Outlay {
	if len(inv.Outlays) == 0 {
		return nil
	}
	outlays := make([]*Outlay, len(inv.Outlays))
	for i, o := range inv.Outlays {
		o2 := *o
		o2.Amount = o2.Amount.
			Upscale(defaultCurrencyConversionAccuracy).
			Multiply(ex.Amount).
			Downscale(defaultCurrencyConversionAccuracy)
		outlays[i] = &o2
	}
	return outlays
}

func (inv *Invoice) convertPayment(ex *currency.ExchangeRate) *Payment {
	if inv.Payment == nil {
		return nil
	}
	p2 := *inv.Payment
	if len(inv.Payment.Advances) == 0 {
		return &p2
	}
	p2.Advances = make([]*pay.Advance, len(inv.Payment.Advances))
	for i, a := range inv.Payment.Advances {
		a2 := *a
		a2.Amount = a2.Amount.
			Upscale(defaultCurrencyConversionAccuracy).
			Multiply(ex.Amount).
			Downscale(defaultCurrencyConversionAccuracy)
		p2.Advances[i] = &a2
	}
	return &p2
}
