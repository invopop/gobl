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
	i2.Lines = inv.convertLines(ex)
	i2.Discounts = inv.convertDiscounts(ex)
	i2.Charges = inv.convertCharges(ex)
	i2.Payment = inv.convertPaymentDetails(ex)
	i2.Currency = cur

	if err := i2.Calculate(); err != nil {
		return nil, err
	}

	return &i2, nil
}

func (inv *Invoice) convertLines(ex *currency.ExchangeRate) []*Line {
	if len(inv.Lines) == 0 {
		return nil
	}
	lines := make([]*Line, len(inv.Lines))
	for i, l := range inv.Lines {
		lines[i] = l.convertInto(ex)
	}
	return lines
}

func (l *Line) convertInto(ex *currency.ExchangeRate) *Line {
	accuracy := defaultCurrencyConversionAccuracy

	l2 := *l
	l2i := *l.Item

	// Add current price to the list of alternative prices
	l2i.AltPrices = append(l2i.AltPrices, &currency.Amount{
		Currency: ex.From,
		Value:    l2i.Price,
	})

	// Use alt price if available
	altFound := false
	for i, ap := range l2i.AltPrices {
		if ap.Currency == ex.To {
			l2i.Price = ap.Value
			// remove this alt price from the list
			l2i.AltPrices = append(l2i.AltPrices[:i], l2i.AltPrices[i+1:]...)
			altFound = true
			break
		}
	}
	if !altFound {
		// Perform exchange
		l2i.Price = l2i.Price.Upscale(accuracy).Multiply(ex.Amount)
	}

	if len(l2.Discounts) > 0 {
		rows := make([]*LineDiscount, len(l2.Discounts))
		for i, v := range l.Discounts {
			d := *v
			d.Amount = d.Amount.Upscale(accuracy).Multiply(ex.Amount)
			rows[i] = &d
		}
		l2.Discounts = rows
	}

	if len(l2.Charges) > 0 {
		rows := make([]*LineCharge, len(l2.Charges))
		for i, v := range l.Charges {
			d := *v
			d.Amount = d.Amount.Upscale(accuracy).Multiply(ex.Amount)
			rows[i] = &d
		}
		l2.Charges = rows
	}

	l2.Item = &l2i
	return &l2
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

func (inv *Invoice) convertPaymentDetails(ex *currency.ExchangeRate) *PaymentDetails {
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
