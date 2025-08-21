package tax_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// taxableLine is a very simple implementation of what the totals calculator requires.
type taxableLine struct {
	taxes  tax.Set
	amount num.Amount
}

func (tl *taxableLine) GetTaxes() tax.Set {
	return tl.taxes
}

func (tl *taxableLine) GetTotal() num.Amount {
	return tl.amount
}

func TestTotalClone(t *testing.T) {
	var tt *tax.Total
	assert.NotPanics(t, func() {
		_ = tt.Clone()
	})
	tt = &tax.Total{
		Categories: []*tax.CategoryTotal{
			{
				Code:     tax.CategoryVAT,
				Retained: false,
				Rates: []*tax.RateTotal{
					{
						Key:     tax.RateStandard,
						Base:    num.MakeAmount(10000, 2),
						Percent: num.NewPercentage(210, 3),
						Amount:  num.MakeAmount(2100, 2),
						Surcharge: &tax.RateTotalSurcharge{
							Percent: num.MakePercentage(10, 3),
							Amount:  num.MakeAmount(100, 2),
						},
					},
				},
				Amount:    num.MakeAmount(2100, 2),
				Surcharge: num.NewAmount(100, 2),
			},
		},
		Sum: num.MakeAmount(2200, 2),
	}
	tt2 := tt.Clone()
	d1, err := json.Marshal(tt)
	require.NoError(t, err)
	d2, err := json.Marshal(tt2)
	require.NoError(t, err)

	assert.JSONEq(t, string(d1), string(d2))

	tt.Categories[0].Rates[0].Base = num.MakeAmount(20000, 2)
	assert.NotEqual(t, tt.Categories[0].Rates[0].Base, tt2.Categories[0].Rates[0].Base)
}

func TestTotalNegate(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var tt *tax.Total
		assert.NotPanics(t, func() {
			_ = tt.Negate()
		})
	})

	tt := &tax.Total{
		Categories: []*tax.CategoryTotal{
			{
				Code:     tax.CategoryVAT,
				Retained: false,
				Rates: []*tax.RateTotal{
					{
						Key:     tax.RateStandard,
						Base:    num.MakeAmount(10000, 2),
						Percent: num.NewPercentage(210, 3),
						Amount:  num.MakeAmount(2100, 2),
					},
				},
				Amount: num.MakeAmount(2100, 2),
			},
		},
		Sum: num.MakeAmount(2100, 2),
	}
	tt2 := tt.Negate()
	assert.Equal(t, int64(-2100), tt2.Category("VAT").Rates[0].Amount.Value())
}

func TestTotalCategory(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var tt *tax.Total
		assert.NotPanics(t, func() {
			_ = tt.Category("VAT")
		})
	})

	t.Run("empty", func(t *testing.T) {
		tt := &tax.Total{}
		assert.Nil(t, tt.Category("VAT"))
	})

	t.Run("basic", func(t *testing.T) {
		tt := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code: tax.CategoryVAT,
					Rates: []*tax.RateTotal{
						{
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(210, 3),
						},
					},
				},
			},
		}
		cat := tt.Category("VAT")
		assert.NotNil(t, cat)
		assert.Equal(t, tax.CategoryVAT, cat.Code)
		assert.Len(t, cat.Rates, 1)
	})
}

func TestTotalExchange(t *testing.T) {
	er := &currency.ExchangeRate{
		From:   currency.EUR,
		To:     currency.USD,
		Amount: num.MakeAmount(120, 2), // 1 EUR = 1.20 USD
	}

	t.Run("nil", func(t *testing.T) {
		var tt *tax.Total
		assert.NotPanics(t, func() {
			tt.Exchange(er, tax.RoundingRulePrecise)
		})
	})
	t.Run("basic example", func(t *testing.T) {
		tt := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code:     tax.CategoryVAT,
					Retained: false,
					Rates: []*tax.RateTotal{
						{
							Key:     tax.RateStandard,
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(210, 3),
							Amount:  num.MakeAmount(2100, 2),
						},
					},
					Amount: num.MakeAmount(2100, 2),
				},
			},
			Sum: num.MakeAmount(2100, 2),
		}
		tt.Exchange(er, tax.RoundingRulePrecise)
		assert.Equal(t, int64(2520), tt.Sum.Value())
		assert.Equal(t, int64(2520), tt.Category("VAT").Amount.Value())
		assert.Equal(t, int64(2520), tt.Category("VAT").Rates[0].Amount.Value())
	})
}

func TestTotalScale(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var tt *tax.Total
		assert.NotPanics(t, func() {
			tt.Scale(num.MakeAmount(5, 2), currency.EUR, tax.RoundingRulePrecise)
		})
	})
	t.Run("empty", func(t *testing.T) {
		tt := &tax.Total{}
		assert.NotPanics(t, func() {
			tt.Scale(num.MakeAmount(5, 2), currency.EUR, tax.RoundingRulePrecise)
		})
	})
	t.Run("basic", func(t *testing.T) {
		tt := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code: tax.CategoryVAT,
					Rates: []*tax.RateTotal{
						{
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(210, 3),
						},
					},
				},
			},
		}
		tt.Scale(num.MakeAmount(50, 2), currency.EUR, tax.RoundingRulePrecise)
		assert.Equal(t, int64(1050), tt.Sum.Value())
		assert.Equal(t, int64(1050), tt.Category("VAT").Amount.Value())
		assert.Equal(t, int64(1050), tt.Category("VAT").Rates[0].Amount.Value())
	})
}

func TestTotalMerge(t *testing.T) {
	t.Run("basic merge", func(t *testing.T) {
		tt := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code:     tax.CategoryVAT,
					Retained: false,
					Rates: []*tax.RateTotal{
						{
							Key:     tax.RateStandard,
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(210, 3),
							Amount:  num.MakeAmount(2100, 2),
						},
					},
					Amount: num.MakeAmount(2100, 2),
				},
			},
			Sum: num.MakeAmount(2100, 2),
		}
		tt2 := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code:     tax.CategoryVAT,
					Retained: false,
					Rates: []*tax.RateTotal{
						{
							Key:     tax.RateStandard,
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(210, 3),
							Amount:  num.MakeAmount(2100, 2),
						},
					},
					Amount: num.MakeAmount(2100, 2),
				},
			},
			Sum: num.MakeAmount(2100, 2),
		}
		tt3 := tt.Merge(tt2)
		assert.Equal(t, int64(4200), tt3.Category("VAT").Amount.Value())
	})
	t.Run("invert then merge", func(t *testing.T) {
		tt := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code:     tax.CategoryVAT,
					Retained: false,
					Rates: []*tax.RateTotal{
						{
							Key:     tax.RateStandard,
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(210, 3),
							Amount:  num.MakeAmount(2100, 2),
						},
					},
					Amount: num.MakeAmount(2100, 2),
				},
			},
			Sum: num.MakeAmount(2100, 2),
		}
		tt2 := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code:     tax.CategoryVAT,
					Retained: false,
					Rates: []*tax.RateTotal{
						{
							Key:     tax.RateStandard,
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(210, 3),
							Amount:  num.MakeAmount(2100, 2),
						},
					},
					Amount: num.MakeAmount(2100, 2),
				},
			},
			Sum: num.MakeAmount(2100, 2),
		}
		tt = tt.Negate()
		tt3 := tt.Merge(tt2)
		assert.Equal(t, int64(0), tt3.Category("VAT").Amount.Value())
	})
	t.Run("merge exempt", func(t *testing.T) {
		tt := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code:     tax.CategoryVAT,
					Retained: false,
					Rates: []*tax.RateTotal{
						{
							Base:   num.MakeAmount(10000, 2),
							Amount: num.MakeAmount(0, 2),
						},
					},
					Amount: num.MakeAmount(0, 2),
				},
			},
			Sum: num.MakeAmount(0, 2),
		}
		tt2 := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code:     tax.CategoryVAT,
					Retained: false,
					Rates: []*tax.RateTotal{
						{
							Base:   num.MakeAmount(10000, 2),
							Amount: num.MakeAmount(0, 2),
						},
					},
					Amount: num.MakeAmount(0, 2),
				},
			},
			Sum: num.MakeAmount(0, 2),
		}
		tt3 := tt.Merge(tt2)
		assert.Equal(t, int64(0), tt3.Sum.Value())
		assert.Equal(t, int64(0), tt3.Category("VAT").Amount.Value())
		assert.Equal(t, int64(0), tt3.Category("VAT").Rates[0].Amount.Value())
	})
	t.Run("merge with different rate keys", func(t *testing.T) {
		tt := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code:     tax.CategoryVAT,
					Retained: false,
					Rates: []*tax.RateTotal{
						{
							Key:     tax.RateStandard,
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(210, 3),
							Amount:  num.MakeAmount(2100, 2),
						},
					},
					Amount: num.MakeAmount(2100, 2),
				},
			},
			Sum: num.MakeAmount(2100, 2),
		}
		tt2 := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code:     tax.CategoryVAT,
					Retained: false,
					Rates: []*tax.RateTotal{
						{
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(210, 3),
							Amount:  num.MakeAmount(2100, 2),
						},
					},
					Amount: num.MakeAmount(2100, 2),
				},
			},
			Sum: num.MakeAmount(2100, 2),
		}
		tt3 := tt.Merge(tt2)
		assert.Equal(t, int64(4200), tt3.Sum.Value())
		assert.Equal(t, int64(4200), tt3.Category("VAT").Amount.Value())
		assert.Equal(t, int64(4200), tt3.Category("VAT").Rates[0].Amount.Value())
		assert.Equal(t, tax.RateStandard, tt3.Category("VAT").Rates[0].Key)
	})
	t.Run("merge with different rate percents", func(t *testing.T) {
		tt := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code:     tax.CategoryVAT,
					Retained: false,
					Rates: []*tax.RateTotal{
						{
							Key:     tax.RateStandard,
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(210, 3),
							Amount:  num.MakeAmount(2100, 2),
						},
					},
					Amount: num.MakeAmount(2100, 2),
				},
			},
			Sum: num.MakeAmount(2100, 2),
		}
		tt2 := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code:     tax.CategoryVAT,
					Retained: false,
					Rates: []*tax.RateTotal{
						{
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(200, 3),
							Amount:  num.MakeAmount(2000, 2),
						},
					},
					Amount: num.MakeAmount(2000, 2),
				},
			},
			Sum: num.MakeAmount(2100, 2),
		}
		tt3 := tt.Merge(tt2)
		assert.Equal(t, int64(4200), tt3.Sum.Value())
		assert.Equal(t, int64(4100), tt3.Category("VAT").Amount.Value())
		assert.Equal(t, int64(2100), tt3.Category("VAT").Rates[0].Amount.Value())
		assert.Equal(t, int64(2000), tt3.Category("VAT").Rates[1].Amount.Value())
	})
	t.Run("merge with different categories", func(t *testing.T) {
		tt := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code:     tax.CategoryVAT,
					Retained: false,
					Rates: []*tax.RateTotal{
						{
							Key:     tax.RateStandard,
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(210, 3),
							Amount:  num.MakeAmount(2100, 2),
						},
					},
					Amount: num.MakeAmount(2100, 2),
				},
			},
			Sum: num.MakeAmount(2100, 2),
		}
		tt2 := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code:     "IRPF",
					Retained: true,
					Rates: []*tax.RateTotal{
						{
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(150, 3),
							Amount:  num.MakeAmount(1500, 2),
						},
					},
					Amount: num.MakeAmount(1500, 2),
				},
			},
			Sum: num.MakeAmount(-1500, 2),
		}
		tt3 := tt.Merge(tt2)
		assert.Equal(t, int64(2100), tt3.Category("VAT").Amount.Value())
		assert.Equal(t, int64(2100), tt3.Category("VAT").Rates[0].Amount.Value())
		assert.Equal(t, int64(1500), tt3.Category("IRPF").Rates[0].Amount.Value())
		assert.Equal(t, int64(600), tt3.Sum.Value())
	})
	t.Run("merge with same surcharge", func(t *testing.T) {
		tt := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code:     tax.CategoryVAT,
					Retained: false,
					Rates: []*tax.RateTotal{
						{
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(210, 3),
							Surcharge: &tax.RateTotalSurcharge{
								Percent: num.MakePercentage(10, 3),
								Amount:  num.MakeAmount(100, 2),
							},
							Amount: num.MakeAmount(2100, 2),
						},
					},
					Amount:    num.MakeAmount(2200, 2),
					Surcharge: num.NewAmount(100, 2),
				},
			},
			Sum: num.MakeAmount(2200, 2),
		}
		tt2 := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code:     tax.CategoryVAT,
					Retained: false,
					Rates: []*tax.RateTotal{
						{
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(210, 3),
							Amount:  num.MakeAmount(2000, 2),
							Surcharge: &tax.RateTotalSurcharge{
								Percent: num.MakePercentage(10, 3),
								Amount:  num.MakeAmount(100, 2),
							},
						},
					},
					Amount:    num.MakeAmount(2100, 2),
					Surcharge: num.NewAmount(100, 2),
				},
			},
			Sum: num.MakeAmount(2200, 2),
		}
		tt3 := tt.Merge(tt2)
		assert.Equal(t, int64(4400), tt3.Sum.Value())
		assert.Equal(t, int64(4300), tt3.Category("VAT").Amount.Value())
		assert.Equal(t, int64(4100), tt3.Category("VAT").Rates[0].Amount.Value())
		assert.Equal(t, int64(200), tt3.Category("VAT").Surcharge.Value())
	})
	t.Run("merge with different surcharge", func(t *testing.T) {
		tt := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code:     tax.CategoryVAT,
					Retained: false,
					Rates: []*tax.RateTotal{
						{
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(210, 3),
							Surcharge: &tax.RateTotalSurcharge{
								Percent: num.MakePercentage(10, 3),
								Amount:  num.MakeAmount(100, 2),
							},
							Amount: num.MakeAmount(2100, 2),
						},
					},
					Amount:    num.MakeAmount(2200, 2),
					Surcharge: num.NewAmount(100, 2),
				},
			},
			Sum: num.MakeAmount(2200, 2),
		}
		tt2 := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code:     tax.CategoryVAT,
					Retained: false,
					Rates: []*tax.RateTotal{
						{
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(210, 3),
							Amount:  num.MakeAmount(2000, 2),
							Surcharge: &tax.RateTotalSurcharge{
								Percent: num.MakePercentage(11, 3),
								Amount:  num.MakeAmount(110, 2),
							},
						},
					},
					Amount:    num.MakeAmount(2100, 2),
					Surcharge: num.NewAmount(110, 2),
				},
			},
			Sum: num.MakeAmount(2210, 2),
		}
		tt3 := tt.Merge(tt2)
		assert.Equal(t, int64(4410), tt3.Sum.Value())
		assert.Equal(t, int64(4300), tt3.Category("VAT").Amount.Value())
		assert.Equal(t, int64(2100), tt3.Category("VAT").Rates[0].Amount.Value())
		assert.Equal(t, int64(2000), tt3.Category("VAT").Rates[1].Amount.Value())
	})
}

func TestTotalCalculate(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var tt *tax.Total
		assert.NotPanics(t, func() {
			tt.Calculate(currency.EUR, tax.RoundingRulePrecise)
		})
	})
	t.Run("empty", func(t *testing.T) {
		tt := &tax.Total{}
		tt.Calculate(currency.EUR, tax.RoundingRulePrecise)
		assert.Equal(t, int64(0), tt.Sum.Value())
	})
	t.Run("basic", func(t *testing.T) {
		tt := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code: tax.CategoryVAT,
					Rates: []*tax.RateTotal{
						{
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(210, 3),
						},
					},
				},
			},
		}
		tt.Calculate(currency.EUR, tax.RoundingRulePrecise)
		assert.Equal(t, int64(2100), tt.Sum.Value())
		assert.Equal(t, int64(2100), tt.Category("VAT").Amount.Value())
		assert.Equal(t, int64(2100), tt.Category("VAT").Rates[0].Amount.Value())
	})
	t.Run("basic with surcharge", func(t *testing.T) {
		tt := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code:     tax.CategoryVAT,
					Retained: false,
					Rates: []*tax.RateTotal{
						{
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(210, 3),
							Surcharge: &tax.RateTotalSurcharge{
								Percent: num.MakePercentage(10, 3),
							},
						},
					},
				},
			},
		}
		tt.Calculate(currency.EUR, tax.RoundingRulePrecise)
		data, _ := json.Marshal(tt)
		fmt.Printf("TOTAL: %s\n", string(data))
		assert.Equal(t, int64(2200), tt.Sum.Value())
		assert.Equal(t, int64(2100), tt.Category(tax.CategoryVAT).Amount.Value())
		assert.Equal(t, int64(2100), tt.Category(tax.CategoryVAT).Rates[0].Amount.Value())
		assert.Equal(t, int64(100), tt.Category(tax.CategoryVAT).Surcharge.Value())
	})

	t.Run("basic with retained surcharge", func(t *testing.T) {
		tt := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code:     tax.CategoryVAT,
					Retained: false,
					Rates: []*tax.RateTotal{
						{
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(210, 3),
							Surcharge: &tax.RateTotalSurcharge{
								Percent: num.MakePercentage(10, 3),
							},
						},
					},
				},
				{
					Code:     "IRPF",
					Retained: true,
					Rates: []*tax.RateTotal{
						{
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(150, 3),
							Surcharge: &tax.RateTotalSurcharge{
								Percent: num.MakePercentage(10, 3),
							},
						},
					},
				},
			},
		}
		tt.Calculate(currency.EUR, tax.RoundingRulePrecise)
		data, _ := json.Marshal(tt)
		fmt.Printf("TOTAL: %s\n", string(data))
		assert.Equal(t, "22.00", tt.Sum.String())
		assert.Equal(t, "16.00", tt.Retained.String())
		assert.Equal(t, int64(2100), tt.Category(tax.CategoryVAT).Amount.Value())
		assert.Equal(t, int64(2100), tt.Category(tax.CategoryVAT).Rates[0].Amount.Value())
		assert.Equal(t, "15.00", tt.Category("IRPF").Rates[0].Amount.String())
		assert.Equal(t, int64(100), tt.Category(tax.CategoryVAT).Surcharge.Value())
	})

	t.Run("basic with informative tax", func(t *testing.T) {
		tt := &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code: tax.CategoryVAT,
					Rates: []*tax.RateTotal{
						{
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(210, 3),
						},
					},
				},
				{
					Code:        "ISS",
					Informative: true,
					Rates: []*tax.RateTotal{
						{
							Base:    num.MakeAmount(10000, 2),
							Percent: num.NewPercentage(50, 3),
						},
					},
				},
			},
		}
		tt.Calculate(currency.EUR, tax.RoundingRulePrecise)
		assert.Equal(t, int64(2100), tt.Sum.Value()) // Informative tax with surcharge should not affect Sum
		assert.Equal(t, int64(500), tt.Category("ISS").Amount.Value())
		assert.True(t, tt.Category("ISS").Informative)
	})
}
