package bill

import (
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
)

// defaultCurrencyConversionAccuracy is the number of decimal places added to
// an amount before converting it into another currency.
const defaultCurrencyConversionAccuracy uint32 = 2

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
	if line.Item == nil || line.Item.Price == nil {
		return line
	}
	l2 := *line
	l2.Item = convertItemInto(ex, line.Item)
	l2.Discounts = convertLineDiscountsInto(ex, line.Discounts)
	l2.Charges = convertLineChargesInto(ex, line.Charges)
	l2.Breakdown = convertSubLinesInto(ex, line.Breakdown)
	l2.Substituted = convertSubLinesInto(ex, line.Substituted)
	return &l2
}

func convertSubLinesInto(ex *currency.ExchangeRate, sls []*SubLine) []*SubLine {
	if len(sls) == 0 {
		return sls
	}
	rows := make([]*SubLine, len(sls))
	for i, sl := range sls {
		rows[i] = convertSubLineInto(ex, sl)
	}
	return rows
}

func convertSubLineInto(ex *currency.ExchangeRate, sl *SubLine) *SubLine {
	if sl == nil || sl.Item == nil || sl.Item.Price == nil {
		return sl
	}
	sl2 := *sl
	sl2.Item = convertItemInto(ex, sl.Item)
	sl2.Discounts = convertLineDiscountsInto(ex, sl.Discounts)
	sl2.Charges = convertLineChargesInto(ex, sl.Charges)
	return &sl2
}

// convertItemInto provides a copy of the item with its price in the exchange
// rate's target currency, using a matching alternative price if available.
func convertItemInto(ex *currency.ExchangeRate, item *org.Item) *org.Item {
	accuracy := defaultCurrencyConversionAccuracy
	i2 := *item
	price := *item.Price

	// Keep the current price as an alternative, and use an existing
	// alternative price in the target currency if available.
	altFound := false
	alts := make([]*currency.Amount, 0, len(item.AltPrices)+1)
	for _, ap := range item.AltPrices {
		if !altFound && ap.Currency == ex.To {
			price = ap.Value
			altFound = true
			continue
		}
		alts = append(alts, ap)
	}
	i2.AltPrices = append(alts, &currency.Amount{
		Currency: ex.From,
		Value:    *item.Price,
	})

	if altFound {
		// List price and discount are unknown in the alternative currency
		i2.List = nil
		i2.Discount = nil
	} else {
		price = price.Upscale(accuracy).Multiply(ex.Amount)
		if item.List != nil {
			l := item.List.Upscale(accuracy).Multiply(ex.Amount)
			i2.List = &l
		}
		if item.Discount != nil {
			d := item.Discount.Upscale(accuracy).Multiply(ex.Amount)
			i2.Discount = &d
		}
	}

	i2.Price = &price
	if i2.Currency != "" {
		i2.Currency = ex.To
	}
	return &i2
}

func convertLineDiscountsInto(ex *currency.ExchangeRate, discounts []*LineDiscount) []*LineDiscount {
	if len(discounts) == 0 {
		return discounts
	}
	accuracy := defaultCurrencyConversionAccuracy
	rows := make([]*LineDiscount, len(discounts))
	for i, v := range discounts {
		d := *v
		d.Amount = d.Amount.Upscale(accuracy).Multiply(ex.Amount)
		rows[i] = &d
	}
	return rows
}

func convertLineChargesInto(ex *currency.ExchangeRate, charges []*LineCharge) []*LineCharge {
	if len(charges) == 0 {
		return charges
	}
	accuracy := defaultCurrencyConversionAccuracy
	rows := make([]*LineCharge, len(charges))
	for i, v := range charges {
		c := *v
		c.Amount = c.Amount.Upscale(accuracy).Multiply(ex.Amount)
		rows[i] = &c
	}
	return rows
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
	p2.Advances = make([]*pay.Record, len(pd.Advances))
	for i, a := range pd.Advances {
		if a == nil {
			continue
		}
		a2 := *a
		a2.Amount = a2.Amount.
			Upscale(defaultCurrencyConversionAccuracy).
			Multiply(ex.Amount).
			Downscale(defaultCurrencyConversionAccuracy)
		p2.Advances[i] = &a2
	}
	return &p2
}
