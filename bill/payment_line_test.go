package bill

import (
	"context"
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentLineCalculation(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		pl := &PaymentLine{
			Index:  1,
			Amount: num.MakeAmount(6050, 2),
		}
		require.NoError(t, pl.calculate(nil, currency.EUR, tax.RoundingRulePrecise))
		assert.Nil(t, pl.Payable)
		assert.Nil(t, pl.Advances)
		assert.Equal(t, "60.50", pl.Amount.String())
		assert.Nil(t, pl.Tax)
		assert.Nil(t, pl.Due)
	})
	t.Run("basic with payable", func(t *testing.T) {
		pl := &PaymentLine{
			Index:   1,
			Payable: num.NewAmount(6050, 2),
			Amount:  num.MakeAmount(6050, 2),
		}
		require.NoError(t, pl.calculate(nil, currency.EUR, tax.RoundingRulePrecise))
		assert.Equal(t, "60.50", pl.Payable.String())
		assert.Nil(t, pl.Advances)
		assert.Equal(t, "60.50", pl.Amount.String())
		assert.Nil(t, pl.Tax)
		assert.Equal(t, "0.00", pl.Due.String())
	})
	t.Run("basic with payable and advances", func(t *testing.T) {
		pl := &PaymentLine{
			Index:       1,
			Installment: 2,
			Payable:     num.NewAmount(6050, 2),
			Advances:    num.NewAmount(2000, 2), // 20€ already paid
			Amount:      num.MakeAmount(2000, 2),
		}
		require.NoError(t, pl.calculate(nil, currency.EUR, tax.RoundingRulePrecise))
		assert.Equal(t, "60.50", pl.Payable.String())
		assert.Equal(t, "20.00", pl.Advances.String())
		assert.Equal(t, "20.00", pl.Amount.String())
		assert.Nil(t, pl.Tax)
		assert.Equal(t, "20.50", pl.Due.String())
	})
	t.Run("basic with payable and taxes", func(t *testing.T) {
		pl := &PaymentLine{
			Index:       1,
			Installment: 2,
			Payable:     num.NewAmount(6050, 2),
			Advances:    num.NewAmount(2000, 2), // 20€ already paid
			Amount:      num.MakeAmount(2000, 2),
			Tax: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code: "VAT",
						Rates: []*tax.RateTotal{
							{
								Base:    num.MakeAmount(1653, 2),
								Percent: num.NewPercentage(21, 2),
							},
						},
					},
				},
			},
		}
		require.NoError(t, pl.calculate(nil, currency.EUR, tax.RoundingRulePrecise))
		assert.Equal(t, "60.50", pl.Payable.String())
		assert.Equal(t, "20.00", pl.Advances.String())
		assert.Equal(t, "20.00", pl.Amount.String())
		assert.Equal(t, "16.53", pl.Tax.Categories[0].Rates[0].Base.String(), "should be same as original")
		assert.Equal(t, "3.47", pl.Tax.Sum.String())
		assert.Equal(t, "20.50", pl.Due.String())
	})
	t.Run("with partial payments and taxes 50%", func(t *testing.T) {
		pl := &PaymentLine{
			Index: 1,
			Document: &org.DocumentRef{
				Series:    "F1",
				Code:      "01234",
				IssueDate: cal.NewDate(2025, 1, 24),
				Payable:   num.NewAmount(12100, 2),
				Tax: &tax.Total{
					Categories: []*tax.CategoryTotal{
						{
							Code: "VAT",
							Rates: []*tax.RateTotal{
								{
									Base:    num.MakeAmount(10000, 2),
									Percent: num.NewPercentage(21, 2),
								},
							},
						},
					},
				},
			},
			Amount: num.MakeAmount(6050, 2),
		}
		require.NoError(t, pl.calculate(nil, currency.EUR, tax.RoundingRulePrecise))
		assert.Equal(t, "121.00", pl.Payable.String())
		assert.Nil(t, pl.Advances)
		assert.Equal(t, "60.50", pl.Amount.String(), "should be half of the total")
		assert.Equal(t, "10.50", pl.Tax.Sum.String(), "should be half of the tax")
		assert.Equal(t, "60.50", pl.Due.String(), "should be half of the payable amount")
	})

	t.Run("with partial payments and taxes 25%", func(t *testing.T) {
		pl := &PaymentLine{
			Index: 1,
			Document: &org.DocumentRef{
				Series:    "F1",
				Code:      "01234",
				IssueDate: cal.NewDate(2025, 1, 24),
				Payable:   num.NewAmount(12100, 2),
				Tax: &tax.Total{
					Categories: []*tax.CategoryTotal{
						{
							Code: "VAT",
							Rates: []*tax.RateTotal{
								{
									Base:    num.MakeAmount(10000, 2),
									Percent: num.NewPercentage(21, 2),
								},
							},
						},
					},
				},
			},
			Amount: num.MakeAmount(3025, 2),
		}
		require.NoError(t, pl.calculate(nil, currency.EUR, tax.RoundingRulePrecise))
		assert.Equal(t, "121.00", pl.Payable.String())
		assert.Nil(t, pl.Advances)
		assert.Equal(t, "30.25", pl.Amount.String())
		assert.Equal(t, "5.25", pl.Tax.Sum.String())
		assert.Equal(t, "90.75", pl.Due.String())
	})

	t.Run("with partial payments and taxes 25% and advances", func(t *testing.T) {
		pl := &PaymentLine{
			Index:       1,
			Installment: 2,
			Document: &org.DocumentRef{
				Series:    "F1",
				Code:      "01234",
				IssueDate: cal.NewDate(2025, 1, 24),
				Payable:   num.NewAmount(12100, 2),
				Tax: &tax.Total{
					Categories: []*tax.CategoryTotal{
						{
							Code: "VAT",
							Rates: []*tax.RateTotal{
								{
									Base:    num.MakeAmount(10000, 2),
									Percent: num.NewPercentage(21, 2),
								},
							},
						},
					},
				},
			},
			Advances: num.NewAmount(2000, 2), // 20€ already paid
			Amount:   num.MakeAmount(3025, 2),
		}
		require.NoError(t, pl.calculate(nil, currency.EUR, tax.RoundingRulePrecise))
		assert.Equal(t, 2, pl.Installment, "should be the second installment")
		assert.Equal(t, "121.00", pl.Payable.String())
		assert.Equal(t, "20.00", pl.Advances.String())
		assert.Equal(t, "30.25", pl.Amount.String())
		assert.Equal(t, "5.25", pl.Tax.Sum.String())
		assert.Equal(t, "70.75", pl.Due.String())
	})

	t.Run("with taxes and currency", func(t *testing.T) {
		pl := &PaymentLine{
			Index: 1,
			Document: &org.DocumentRef{
				Series:    "F1",
				Code:      "01234",
				IssueDate: cal.NewDate(2025, 1, 24),
				Currency:  currency.USD,
				Payable:   num.NewAmount(12100, 2),
				Tax: &tax.Total{
					Categories: []*tax.CategoryTotal{
						{
							Code: "VAT",
							Rates: []*tax.RateTotal{
								{
									Base:    num.MakeAmount(10000, 2),
									Percent: num.NewPercentage(21, 2),
								},
							},
						},
					},
				},
			},
			Amount: num.MakeAmount(10285, 2),
		}
		rates := []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(85, 2),
			},
		}
		require.NoError(t, pl.calculate(rates, currency.EUR, tax.RoundingRulePrecise))
		assert.Equal(t, "102.85", pl.Payable.String())
		assert.Nil(t, pl.Advances)
		assert.Equal(t, "102.85", pl.Amount.String())
		assert.Equal(t, "17.85", pl.Tax.Sum.String())
		assert.Equal(t, "0.00", pl.Due.String())
	})

	t.Run("with partial payments and taxes 25% and currency", func(t *testing.T) {
		pl := &PaymentLine{
			Index: 1,
			Document: &org.DocumentRef{
				Series:    "F1",
				Code:      "01234",
				IssueDate: cal.NewDate(2025, 1, 24),
				Currency:  currency.USD,
				Payable:   num.NewAmount(12100, 2),
				Tax: &tax.Total{
					Categories: []*tax.CategoryTotal{
						{
							Code: "VAT",
							Rates: []*tax.RateTotal{
								{
									Base:    num.MakeAmount(10000, 2),
									Percent: num.NewPercentage(21, 2),
								},
							},
						},
					},
				},
			},
			Amount: num.MakeAmount(3025, 2),
		}
		rates := []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(85, 2),
			},
		}
		require.NoError(t, pl.calculate(rates, currency.EUR, tax.RoundingRulePrecise))
		assert.Equal(t, "102.85", pl.Payable.String())
		assert.Nil(t, pl.Advances)
		assert.Equal(t, "30.25", pl.Amount.String())
		assert.Equal(t, "5.25", pl.Tax.Sum.String())
		assert.Equal(t, "72.60", pl.Due.String())
	})

	t.Run("missing exchange rate", func(t *testing.T) {
		pl := &PaymentLine{
			Index: 1,
			Document: &org.DocumentRef{
				Series:    "F1",
				Code:      "01234",
				IssueDate: cal.NewDate(2025, 1, 24),
				Currency:  currency.GBP,
				Payable:   num.NewAmount(12100, 2),
				Tax: &tax.Total{
					Categories: []*tax.CategoryTotal{
						{
							Code: "VAT",
							Rates: []*tax.RateTotal{
								{
									Base:    num.MakeAmount(10000, 2),
									Percent: num.NewPercentage(21, 2),
								},
							},
						},
					},
				},
			},
			Amount: num.MakeAmount(10285, 2),
		}
		rates := []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(85, 2),
			},
		}
		err := pl.calculate(rates, currency.EUR, tax.RoundingRulePrecise)
		require.ErrorContains(t, err, "document: (currency: missing exchange rate from GBP to EUR.).")
	})
}

func TestPaymentLineValidation(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		pl := &PaymentLine{
			Index:  1,
			Amount: num.MakeAmount(6050, 2),
		}
		assert.NoError(t, pl.ValidateWithContext(context.Background()))
	})
	t.Run("big installment", func(t *testing.T) {
		pl := &PaymentLine{
			Index:       1,
			Installment: 1000,
			Amount:      num.MakeAmount(6050, 2),
		}
		err := pl.ValidateWithContext(context.Background())
		assert.ErrorContains(t, err, "installment: must be no greater than 999.")
	})
	t.Run("min installment", func(t *testing.T) {
		pl := &PaymentLine{
			Index:       1,
			Installment: -1,
			Amount:      num.MakeAmount(6050, 2),
		}
		err := pl.ValidateWithContext(context.Background())
		assert.ErrorContains(t, err, "installment: must be no less than 1")
	})
	t.Run("zero installment", func(t *testing.T) {
		pl := &PaymentLine{
			Index:       1,
			Installment: 0, // same as empty
			Amount:      num.MakeAmount(6050, 2),
		}
		err := pl.ValidateWithContext(context.Background())
		assert.NoError(t, err)
	})

	t.Run("advances less than payable", func(t *testing.T) {
		pl := &PaymentLine{
			Index:    1,
			Payable:  num.NewAmount(6050, 2),
			Advances: num.NewAmount(7000, 2), // more than payable
			Amount:   num.MakeAmount(6050, 2),
		}
		err := pl.ValidateWithContext(context.Background())
		assert.ErrorContains(t, err, "advances: must be no greater than 60.50")
	})

	t.Run("amount more than payable", func(t *testing.T) {
		pl := &PaymentLine{
			Index:    1,
			Payable:  num.NewAmount(6050, 2),
			Advances: num.NewAmount(2000, 2),
			Amount:   num.MakeAmount(6050, 2),
		}
		err := pl.ValidateWithContext(context.Background())
		assert.ErrorContains(t, err, "amount: must be no greater than 40.50.")
	})
}
