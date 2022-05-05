package tax_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regions/common"
	"github.com/invopop/gobl/regions/es"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

// taxableLine is a very simple implementation of what the totals calculator requires.
type taxableLine struct {
	taxes  tax.Map
	amount num.Amount
}

func (tl *taxableLine) GetTaxes() tax.Map {
	return tl.taxes
}

func (tl *taxableLine) GetTotal() num.Amount {
	return tl.amount
}

func TestTotalCalculate(t *testing.T) {
	spain := es.New()
	date := cal.MakeDate(2022, 01, 24)
	zero := num.MakeAmount(0, 2)
	var tests = []struct {
		desc        string
		lines       []tax.TaxableLine
		taxIncluded tax.Code
		want        *tax.Total
		err         error
	}{
		{
			desc: "basic no tax",
			lines: []tax.TaxableLine{
				&taxableLine{taxes: nil, amount: num.MakeAmount(10000, 2)},
			},
			taxIncluded: "",
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{},
				Sum:        zero,
			},
		},
		{
			desc: "with VAT",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: map[tax.Code]tax.Key{
						common.TaxCategoryVAT: common.TaxRateStandard,
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			taxIncluded: "",
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     common.TaxCategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     common.TaxRateStandard,
								Base:    num.MakeAmount(10000, 2),
								Percent: num.MakePercentage(210, 3),
								Amount:  num.MakeAmount(2100, 2),
							},
						},
						Base:   num.MakeAmount(10000, 2),
						Amount: num.MakeAmount(2100, 2),
					},
				},
				Sum: num.MakeAmount(2100, 2),
			},
		},
		{
			desc: "with multiline VAT",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: map[tax.Code]tax.Key{
						common.TaxCategoryVAT: common.TaxRateStandard,
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: map[tax.Code]tax.Key{
						common.TaxCategoryVAT: common.TaxRateStandard,
					},
					amount: num.MakeAmount(15000, 2),
				},
			},
			taxIncluded: "",
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     common.TaxCategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     common.TaxRateStandard,
								Base:    num.MakeAmount(25000, 2),
								Percent: num.MakePercentage(210, 3),
								Amount:  num.MakeAmount(5250, 2),
							},
						},
						Base:   num.MakeAmount(25000, 2),
						Amount: num.MakeAmount(5250, 2),
					},
				},
				Sum: num.MakeAmount(5250, 2),
			},
		},
		{
			desc: "with multirate VAT",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: map[tax.Code]tax.Key{
						common.TaxCategoryVAT: common.TaxRateStandard,
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: map[tax.Code]tax.Key{
						common.TaxCategoryVAT: common.TaxRateReduced,
					},
					amount: num.MakeAmount(15000, 2),
				},
			},
			taxIncluded: "",
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     common.TaxCategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     common.TaxRateStandard,
								Base:    num.MakeAmount(10000, 2),
								Percent: num.MakePercentage(210, 3),
								Amount:  num.MakeAmount(2100, 2),
							},
							{
								Key:     common.TaxRateReduced,
								Base:    num.MakeAmount(15000, 2),
								Percent: num.MakePercentage(100, 3),
								Amount:  num.MakeAmount(1500, 2),
							},
						},
						Base:   num.MakeAmount(25000, 2),
						Amount: num.MakeAmount(3600, 2),
					},
				},
				Sum: num.MakeAmount(3600, 2),
			},
		},
		{
			desc: "with multirate VAT included in price",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: map[tax.Code]tax.Key{
						common.TaxCategoryVAT: common.TaxRateStandard,
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: map[tax.Code]tax.Key{
						common.TaxCategoryVAT: common.TaxRateReduced,
					},
					amount: num.MakeAmount(15000, 2),
				},
			},
			taxIncluded: common.TaxCategoryVAT,
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     common.TaxCategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     common.TaxRateStandard,
								Base:    num.MakeAmount(8264, 2),
								Percent: num.MakePercentage(210, 3),
								Amount:  num.MakeAmount(1736, 2),
							},
							{
								Key:     common.TaxRateReduced,
								Base:    num.MakeAmount(13636, 2),
								Percent: num.MakePercentage(100, 3),
								Amount:  num.MakeAmount(1364, 2),
							},
						},
						Base:   num.MakeAmount(21900, 2),
						Amount: num.MakeAmount(3100, 2),
					},
				},
				Sum: num.MakeAmount(3100, 2),
			},
		},
		{
			desc: "with multirate VAT and retained tax",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: map[tax.Code]tax.Key{
						common.TaxCategoryVAT: common.TaxRateStandard,
						es.TaxCategoryIRPF:    es.TaxRatePro,
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: map[tax.Code]tax.Key{
						common.TaxCategoryVAT: common.TaxRateReduced,
					},
					amount: num.MakeAmount(15000, 2),
				},
			},
			taxIncluded: "",
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     common.TaxCategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     common.TaxRateStandard,
								Base:    num.MakeAmount(10000, 2),
								Percent: num.MakePercentage(210, 3),
								Amount:  num.MakeAmount(2100, 2),
							},
							{
								Key:     common.TaxRateReduced,
								Base:    num.MakeAmount(15000, 2),
								Percent: num.MakePercentage(100, 3),
								Amount:  num.MakeAmount(1500, 2),
							},
						},
						Base:   num.MakeAmount(25000, 2),
						Amount: num.MakeAmount(3600, 2),
					},
					{
						Code:     es.TaxCategoryIRPF,
						Retained: true,
						Rates: []*tax.RateTotal{
							{
								Key:     es.TaxRatePro,
								Base:    num.MakeAmount(10000, 2),
								Percent: num.MakePercentage(150, 3),
								Amount:  num.MakeAmount(1500, 2),
							},
						},
						Base:   num.MakeAmount(10000, 2),
						Amount: num.MakeAmount(1500, 2),
					},
				},
				Sum: num.MakeAmount(2100, 2),
			},
		},

		{
			desc: "with multirate VAT included in price plus retained tax",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: map[tax.Code]tax.Key{
						common.TaxCategoryVAT: common.TaxRateStandard,
						es.TaxCategoryIRPF:    es.TaxRatePro,
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: map[tax.Code]tax.Key{
						common.TaxCategoryVAT: common.TaxRateReduced,
					},
					amount: num.MakeAmount(15000, 2),
				},
			},
			taxIncluded: common.TaxCategoryVAT,
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     common.TaxCategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     common.TaxRateStandard,
								Base:    num.MakeAmount(8264, 2),
								Percent: num.MakePercentage(210, 3),
								Amount:  num.MakeAmount(1736, 2),
							},
							{
								Key:     common.TaxRateReduced,
								Base:    num.MakeAmount(13636, 2),
								Percent: num.MakePercentage(100, 3),
								Amount:  num.MakeAmount(1364, 2),
							},
						},
						Base:   num.MakeAmount(21900, 2),
						Amount: num.MakeAmount(3100, 2),
					},
					{
						Code:     es.TaxCategoryIRPF,
						Retained: true,
						Rates: []*tax.RateTotal{
							{
								Key:     es.TaxRatePro,
								Base:    num.MakeAmount(8264, 2),
								Percent: num.MakePercentage(150, 3),
								Amount:  num.MakeAmount(1240, 2),
							},
						},
						Base:   num.MakeAmount(8264, 2),
						Amount: num.MakeAmount(1240, 2),
					},
				},
				Sum: num.MakeAmount(1860, 2),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			tot := tax.NewTotal(zero)
			err := tot.Calculate(spain, test.lines, test.taxIncluded, date, zero)
			if test.err != nil {
				assert.ErrorIs(t, err, test.err)
			}
			if test.want != nil {
				if !assert.EqualValues(t, test.want, tot) {
					data, _ := json.MarshalIndent(tot, "", "  ")
					t.Logf("data output: %v", string(data))
				}
			}
		})
	}

}
