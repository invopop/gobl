package tax_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
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

func TestTotalCalculate(t *testing.T) {
	spain := es.New()
	portugal := pt.New()
	date := cal.MakeDate(2022, 01, 24)
	zero := num.MakeAmount(0, 2)
	var tests = []struct {
		desc        string
		regime      *tax.Regime // default, spain
		zone        l10n.Code   // default empty
		lines       []tax.TaxableLine
		date        *cal.Date
		taxIncluded cbc.Code
		want        *tax.Total
		err         error
		errContent  string
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
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Rate:     common.TaxRateStandard,
						},
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
			desc: "with exemption",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Tags:     []cbc.Key{es.TagExempt},
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			taxIncluded: "",
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{},
				Sum:        num.MakeAmount(0, 2),
			},
		},
		{
			desc:   "with VAT in Azores",
			regime: portugal,
			zone:   pt.ZoneAzores,
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Rate:     common.TaxRateStandard,
						},
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
								Percent: num.MakePercentage(160, 3),
								Amount:  num.MakeAmount(1600, 2),
							},
						},
						Base:   num.MakeAmount(10000, 2),
						Amount: num.MakeAmount(1600, 2),
					},
				},
				Sum: num.MakeAmount(1600, 2),
			},
		},
		{
			desc: "with VAT percents defined",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Percent:  num.NewPercentage(210, 3),
						},
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
								// Key:     common.TaxRateStandard,
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
			desc: "with VAT percents defined, rate override",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Rate:     common.TaxRateStandard,
							Percent:  num.NewPercentage(200, 3),
						},
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
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Rate:     common.TaxRateStandard,
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Rate:     common.TaxRateStandard,
						},
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
			desc: "with multiline VAT and Surcharge",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Rate:     common.TaxRateStandard.With(es.TaxRateEquivalence),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Rate:     common.TaxRateStandard.With(es.TaxRateEquivalence),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Rate:     common.TaxRateStandard,
						},
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
								Key:     common.TaxRateStandard.With(es.TaxRateEquivalence),
								Base:    num.MakeAmount(20000, 2),
								Percent: num.MakePercentage(210, 3),
								Amount:  num.MakeAmount(4200, 2),
								Surcharge: &tax.RateTotalSurcharge{
									Percent: num.MakePercentage(52, 3),
									Amount:  num.MakeAmount(1040, 2),
								},
							},
							{
								Key:     common.TaxRateStandard,
								Base:    num.MakeAmount(15000, 2),
								Percent: num.MakePercentage(210, 3),
								Amount:  num.MakeAmount(3150, 2),
							},
						},
						Base:      num.MakeAmount(35000, 2),
						Amount:    num.MakeAmount(7350, 2),
						Surcharge: num.NewAmount(1040, 2),
					},
				},
				Sum: num.MakeAmount(8390, 2),
			},
		},
		{
			desc: "with multiline VAT as percentages",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Percent:  num.NewPercentage(210, 3),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Percent:  num.NewPercentage(2100, 4), // different exp.
						},
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
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Rate:     common.TaxRateStandard,
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Rate:     common.TaxRateReduced,
						},
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
			desc: "with multirate VAT as percentages",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Percent:  num.NewPercentage(210, 3),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Percent:  num.NewPercentage(100, 3),
						},
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
								// Key:     common.TaxRateStandard,
								Base:    num.MakeAmount(10000, 2),
								Percent: num.MakePercentage(210, 3),
								Amount:  num.MakeAmount(2100, 2),
							},
							{
								// Key:     common.TaxRateReduced,
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
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Rate:     common.TaxRateStandard,
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Rate:     common.TaxRateReduced,
						},
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
			desc: "with multirate VAT as percentages, and included in price",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Percent:  num.NewPercentage(210, 3),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Percent:  num.NewPercentage(100, 3),
						},
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
								Base:    num.MakeAmount(8264, 2),
								Percent: num.MakePercentage(210, 3),
								Amount:  num.MakeAmount(1736, 2),
							},
							{
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
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Rate:     common.TaxRateStandard,
						},
						{
							Category: es.TaxCategoryIRPF,
							Rate:     es.TaxRatePro,
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Rate:     common.TaxRateReduced,
						},
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
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Rate:     common.TaxRateStandard,
						},
						{
							Category: es.TaxCategoryIRPF,
							Rate:     es.TaxRatePro,
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Rate:     common.TaxRateReduced,
						},
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
		{
			desc: "with invalid category",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: cbc.Code("FOO"),
							Rate:     common.TaxRateStandard,
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			err:        tax.ErrInvalidCategory,
			errContent: "invalid-category: 'FOO'",
		},
		{
			desc: "with invalid rate",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: es.TaxCategoryIRPF,
							Rate:     common.TaxRateStandard,
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			err:        tax.ErrInvalidRate,
			errContent: "invalid-rate: 'standard' rate not defined in category 'IRPF'",
		},

		{
			desc: "with invalid rate on date",
			date: cal.NewDate(2005, 1, 1),
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: es.TaxCategoryIRPF,
							Rate:     es.TaxRatePro,
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			err:        tax.ErrInvalidDate,
			errContent: "invalid-date: rate value unavailable for 'pro' in 'IRPF' on '2005-01-01'",
		},
		{
			desc: "with invalid tax included",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: es.TaxCategoryIRPF,
							Rate:     es.TaxRatePro,
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			taxIncluded: es.TaxCategoryIRPF,
			err:         tax.ErrInvalidPricesInclude,
			errContent:  "cannot include retained",
		},
		{
			desc: "tax included with missing rate",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: common.TaxCategoryVAT,
							Tags:     []cbc.Key{"random"},
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			taxIncluded: common.TaxCategoryVAT,
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{},
				Sum:        num.MakeAmount(0, 2),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			d := date
			if test.date != nil {
				d = *test.date
			}
			reg := spain
			if test.regime != nil {
				reg = test.regime
			}
			zone := l10n.CodeEmpty
			if test.zone != l10n.CodeEmpty {
				zone = test.zone
			}
			tc := &tax.TotalCalculator{
				Regime:   reg,
				Zone:     zone,
				Zero:     zero,
				Date:     d,
				Includes: test.taxIncluded,
				Lines:    test.lines,
			}
			tot := new(tax.Total)
			err := tc.Calculate(tot)
			if test.err != nil && assert.Error(t, err) {
				assert.ErrorIs(t, err, test.err)
			}
			if test.errContent != "" && assert.Error(t, err) {
				assert.Contains(t, err.Error(), test.errContent)
			}
			if test.want != nil {
				if !assert.EqualValues(t, tot, test.want) {
					data, _ := json.MarshalIndent(tot, "", "  ")
					t.Logf("data output: %v", string(data))
				}
			}
		})
	}

}
