package tax_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regions/common"
	"github.com/invopop/gobl/regions/es"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

// taxableLine is a very simple implementation of what the totals calculator requires.
type taxableLine struct {
	rates  tax.Rates
	amount num.Amount
}

func (tl *taxableLine) GetTaxRates() tax.Rates {
	return tl.rates
}

func (tl *taxableLine) GetTotal() num.Amount {
	return tl.amount
}

func TestTotalCalculate(t *testing.T) {
	spain := es.New()
	date := org.MakeDate(2022, 01, 24)
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
				&taxableLine{rates: nil, amount: num.MakeAmount(10000, 2)},
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
					rates: []*tax.Rate{
						{Category: common.TaxCategoryVAT, Code: common.TaxRateVATStandard},
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
								Code:    common.TaxRateVATStandard,
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
					rates: []*tax.Rate{
						{Category: common.TaxCategoryVAT, Code: common.TaxRateVATStandard},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					rates: []*tax.Rate{
						{Category: common.TaxCategoryVAT, Code: common.TaxRateVATStandard},
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
								Code:    common.TaxRateVATStandard,
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
					rates: []*tax.Rate{
						{Category: common.TaxCategoryVAT, Code: common.TaxRateVATStandard},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					rates: []*tax.Rate{
						{Category: common.TaxCategoryVAT, Code: common.TaxRateVATReduced},
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
								Code:    common.TaxRateVATStandard,
								Base:    num.MakeAmount(10000, 2),
								Percent: num.MakePercentage(210, 3),
								Amount:  num.MakeAmount(2100, 2),
							},
							{
								Code:    common.TaxRateVATReduced,
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
					rates: []*tax.Rate{
						{Category: common.TaxCategoryVAT, Code: common.TaxRateVATStandard},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					rates: []*tax.Rate{
						{Category: common.TaxCategoryVAT, Code: common.TaxRateVATReduced},
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
								Code:    common.TaxRateVATStandard,
								Base:    num.MakeAmount(8264, 2),
								Percent: num.MakePercentage(210, 3),
								Amount:  num.MakeAmount(1736, 2),
							},
							{
								Code:    common.TaxRateVATReduced,
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
					rates: []*tax.Rate{
						{Category: common.TaxCategoryVAT, Code: common.TaxRateVATStandard},
						{Category: es.TaxCategoryIRPF, Code: es.TaxRateIRPFStandard},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					rates: []*tax.Rate{
						{Category: common.TaxCategoryVAT, Code: common.TaxRateVATReduced},
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
								Code:    common.TaxRateVATStandard,
								Base:    num.MakeAmount(10000, 2),
								Percent: num.MakePercentage(210, 3),
								Amount:  num.MakeAmount(2100, 2),
							},
							{
								Code:    common.TaxRateVATReduced,
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
								Code:    es.TaxRateIRPFStandard,
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
					rates: []*tax.Rate{
						{Category: common.TaxCategoryVAT, Code: common.TaxRateVATStandard},
						{Category: es.TaxCategoryIRPF, Code: es.TaxRateIRPFStandard},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					rates: []*tax.Rate{
						{Category: common.TaxCategoryVAT, Code: common.TaxRateVATReduced},
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
								Code:    common.TaxRateVATStandard,
								Base:    num.MakeAmount(8264, 2),
								Percent: num.MakePercentage(210, 3),
								Amount:  num.MakeAmount(1736, 2),
							},
							{
								Code:    common.TaxRateVATReduced,
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
								Code:    es.TaxRateIRPFStandard,
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
			err := tot.Calculate(spain.Taxes(), test.lines, test.taxIncluded, date, zero)
			if test.err != nil {
				assert.ErrorIs(t, err, test.err)
			}
			if test.want != nil {
				if !assert.Equal(t, test.want, tot) {
					data, _ := json.MarshalIndent(tot, "", "  ")
					t.Logf("data output: %v", string(data))
				}
			}
		})
	}

}
