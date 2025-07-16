package bill

import (
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/pay"
)

func convertLinesInto(ex *currency.ExchangeRate, lines []*Line) []*Line {
	if len(lines) == 0 {
		return nil
	}
	nls := make([]*Line, len(lines))
	for i, l := range lines {
		nls[i] = convertLineInto(ex, l)
	}
	return nls
}

func convertLineInto(ex *currency.ExchangeRate, line *Line) *Line {
	accuracy := defaultCurrencyConversionAccuracy

	if line.Item == nil || line.Item.Price == nil {
		return line
	}

	l2 := *line
	l2i := *line.Item
	price := *l2i.Price

	// Add current price to the list of alternative prices
	l2i.AltPrices = append(l2i.AltPrices, &currency.Amount{
		Currency: ex.From,
		Value:    price,
	})

	// Use alt price if available
	altFound := false
	for i, ap := range l2i.AltPrices {
		if ap.Currency == ex.To {
			price = ap.Value
			// remove this alt price from the list
			l2i.AltPrices = append(l2i.AltPrices[:i], l2i.AltPrices[i+1:]...)
			altFound = true
			break
		}
	}
	if !altFound {
		// Perform exchange
		price = price.Upscale(accuracy).Multiply(ex.Amount)
	}

	if len(l2.Discounts) > 0 {
		rows := make([]*LineDiscount, len(l2.Discounts))
		for i, v := range line.Discounts {
			d := *v
			d.Amount = d.Amount.Upscale(accuracy).Multiply(ex.Amount)
			rows[i] = &d
		}
		l2.Discounts = rows
	}

	if len(l2.Charges) > 0 {
		rows := make([]*LineCharge, len(l2.Charges))
		for i, v := range line.Charges {
			d := *v
			d.Amount = d.Amount.Upscale(accuracy).Multiply(ex.Amount)
			rows[i] = &d
		}
		l2.Charges = rows
	}

	l2i.Price = &price
	l2i.Currency = ex.To
	l2.Item = &l2i
	return &l2
}

func convertDiscountsInto(ex *currency.ExchangeRate, discounts []*Discount) []*Discount {
	if len(discounts) == 0 {
		return nil
	}
	ds := make([]*Discount, len(discounts))
	for i, d := range discounts {
		ds[i] = convertDiscountInto(ex, d)
	}
	return ds
}

func convertDiscountInto(ex *currency.ExchangeRate, m *Discount) *Discount {
	accuracy := defaultCurrencyConversionAccuracy
	m2 := *m
	m2.Amount = m2.Amount.Upscale(accuracy).Multiply(ex.Amount)
	return &m2
}

func convertChargesInto(ex *currency.ExchangeRate, charges []*Charge) []*Charge {
	if len(charges) == 0 {
		return nil
	}
	cs := make([]*Charge, len(charges))
	for i, c := range charges {
		cs[i] = convertChargeInto(ex, c)
	}
	return cs
}

func convertChargeInto(ex *currency.ExchangeRate, m *Charge) *Charge {
	accuracy := defaultCurrencyConversionAccuracy
	m2 := *m
	m2.Amount = m2.Amount.Upscale(accuracy).Multiply(ex.Amount)
	return &m2
}

func convertPaymentDetailsInto(ex *currency.ExchangeRate, pd *PaymentDetails) *PaymentDetails {
	if pd == nil {
		return nil
	}
	p2 := *pd
	if len(pd.Advances) == 0 {
		return &p2
	}
	p2.Advances = make([]*pay.Advance, len(pd.Advances))
	for i, a := range pd.Advances {
		a2 := *a
		a2.Amount = a2.Amount.
			Upscale(defaultCurrencyConversionAccuracy).
			Multiply(ex.Amount).
			Downscale(defaultCurrencyConversionAccuracy)
		a2.Currency = ex.To
		p2.Advances[i] = &a2
	}
	return &p2
}
