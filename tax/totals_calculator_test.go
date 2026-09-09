package tax_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/addons/es/tbai"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTotalBySumCalculate(t *testing.T) {
	date := cal.MakeDate(2022, 01, 24)
	zero := num.MakeAmount(0, 2)
	var tests = []struct {
		desc        string
		country     l10n.TaxCountryCode // default "ES"
		rounding    cbc.Key
		tags        []cbc.Key      // default empty
		ext         tax.Extensions // default empty
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
							Category: tax.CategoryVAT,
							Rate:     tax.RateGeneral,
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			taxIncluded: "",
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(10000, 2),
								Percent: num.NewPercentage(210, 3),
								Amount:  num.MakeAmount(2100, 2),
							},
						},
						Amount: num.MakeAmount(2100, 2),
					},
				},
				Sum: num.MakeAmount(2100, 2),
			},
		},
		{
			desc: "rate from same country",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Country:  "ES",
							Category: tax.CategoryVAT,
							Rate:     tax.RateGeneral,
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			taxIncluded: "",
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(10000, 2),
								Percent: num.NewPercentage(210, 3),
								Amount:  num.MakeAmount(2100, 2),
							},
						},
						Amount: num.MakeAmount(2100, 2),
					},
				},
				Sum: num.MakeAmount(2100, 2),
			},
		},
		{
			desc:    "from unknown tax regime",
			country: "XX", // this will fail validation!
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
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
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(10000, 2),
								Percent: num.NewPercentage(210, 3),
								Amount:  num.MakeAmount(2100, 2),
							},
						},
						Amount: num.MakeAmount(2100, 2),
					},
				},
				Sum: num.MakeAmount(2100, 2),
			},
		},
		{
			desc: "export with local VAT of known regime",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Country:  "PT",
							Rate:     tax.RateGeneral,
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			taxIncluded: "",
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Country: "PT",
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(10000, 2),
								Percent: num.NewPercentage(230, 3),
								Amount:  num.MakeAmount(2300, 2),
							},
						},
						Amount: num.MakeAmount(2300, 2),
					},
				},
				Sum: num.MakeAmount(2300, 2),
			},
		},
		{
			desc: "rate not defined for key in category",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Country:  "ES",
							Key:      tax.KeyStandard,
							Rate:     "foo",
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			taxIncluded: "",
			err:         tax.ErrInvalid,
			errContent:  "invalid: 'foo' rate not defined for key 'standard' in category 'VAT'",
		},
		{
			desc: "remove percent and surcharge if no percent",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Key:      tax.KeyExempt,
							Percent:  num.NewPercentage(0, 2),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:    tax.KeyExempt,
								Base:   num.MakeAmount(10000, 2),
								Amount: num.MakeAmount(0, 2),
							},
						},
						Amount: num.MakeAmount(0, 2),
					},
				},
				Sum: num.MakeAmount(0, 2),
			},
		},
		{
			desc: "export with local VAT of unknown regime",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Country:  "JP",
							Percent:  num.NewPercentage(190, 3),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			taxIncluded: "",
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Country: "JP",
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(10000, 2),
								Percent: num.NewPercentage(190, 3),
								Amount:  num.MakeAmount(1900, 2),
							},
						},
						Amount: num.MakeAmount(1900, 2),
					},
				},
				Sum: num.MakeAmount(1900, 2),
			},
		},
		{
			desc: "with exemption",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Key:      tax.KeyExempt,
							Ext: tax.ExtensionsOf(cbc.CodeMap{
								tbai.ExtKeyExempt: "E1",
							}),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			taxIncluded: "",
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key: tax.KeyExempt,
								Ext: tax.ExtensionsOf(cbc.CodeMap{
									tbai.ExtKeyExempt: "E1",
								}),
								Base:    num.MakeAmount(10000, 2),
								Percent: nil,
								Amount:  num.MakeAmount(0, 2),
							},
						},
						Amount: num.MakeAmount(0, 2),
					},
				},
				Sum: num.MakeAmount(0, 2),
			},
		},
		{
			desc: "with exemption and empty ext",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Key:      tax.KeyExempt,
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Rate:     tax.RateReduced,
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			taxIncluded: "",
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     tax.KeyExempt,
								Base:    num.MakeAmount(10000, 2),
								Percent: nil,
								Amount:  num.MakeAmount(0, 2),
							},
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(10000, 2),
								Percent: num.NewPercentage(100, 3),
								Amount:  num.MakeAmount(1000, 2),
							},
						},
						Amount: num.MakeAmount(1000, 2),
					},
				},
				Sum: num.MakeAmount(1000, 2),
			},
		},
		{
			desc: "with no percents and matching rate keys",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Key:      tax.KeyExempt,
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Key:      tax.KeyExempt,
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			taxIncluded: "",
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     tax.KeyExempt,
								Base:    num.MakeAmount(20000, 2),
								Percent: nil,
								Amount:  num.MakeAmount(0, 2),
							},
						},
						Amount: num.MakeAmount(0, 2),
					},
				},
				Sum: num.MakeAmount(0, 2),
			},
		},
		{
			desc:    "with VAT in Azores",
			country: "PT",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Key:      tax.KeyStandard,
							Rate:     tax.RateGeneral,
							Ext: tax.ExtensionsOf(cbc.CodeMap{
								pt.ExtKeyRegion:    "PT-AC",
								"pt-saft-tax-rate": "NOR",
							}),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			taxIncluded: "",
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(10000, 2),
								Percent: num.NewPercentage(160, 3),
								Amount:  num.MakeAmount(1600, 2),
								Ext: tax.ExtensionsOf(cbc.CodeMap{
									pt.ExtKeyRegion:    "PT-AC",
									"pt-saft-tax-rate": "NOR",
								}),
							},
						},
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
							Category: tax.CategoryVAT,
							Percent:  num.NewPercentage(210, 3),
						},
					},
					amount: num.MakeAmount(100000, 3),
				},
			},
			taxIncluded: "",
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(10000, 2),
								Percent: num.NewPercentage(210, 3),
								Amount:  num.MakeAmount(2100, 2),
							},
						},
						Amount: num.MakeAmount(2100, 2),
					},
				},
				Sum: num.MakeAmount(2100, 2),
			},
		},
		{
			desc: "with VAT percents defined, replace for rate",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Rate:     tax.RateGeneral,
							Percent:  num.NewPercentage(20, 2),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			taxIncluded: "",
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(10000, 2),
								Percent: num.NewPercentage(210, 3),
								Amount:  num.MakeAmount(2100, 2),
							},
						},
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
							Category: tax.CategoryVAT,
							Rate:     tax.RateGeneral,
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Rate:     tax.RateGeneral,
						},
					},
					amount: num.MakeAmount(15000, 2),
				},
			},
			taxIncluded: "",
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(25000, 2),
								Percent: num.NewPercentage(210, 3),
								Amount:  num.MakeAmount(5250, 2),
							},
						},
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
							Category: tax.CategoryVAT,
							Rate:     tax.RateGeneral.With(es.TaxRateEquivalence),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Rate:     tax.RateGeneral.With(es.TaxRateEquivalence),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Rate:     tax.RateGeneral,
						},
					},
					amount: num.MakeAmount(15000, 2),
				},
			},
			taxIncluded: "",
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(20000, 2),
								Percent: num.NewPercentage(210, 3),
								Amount:  num.MakeAmount(4200, 2),
								Surcharge: &tax.RateTotalSurcharge{
									Percent: num.MakePercentage(52, 3),
									Amount:  num.MakeAmount(1040, 2),
								},
							},
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(15000, 2),
								Percent: num.NewPercentage(210, 3),
								Amount:  num.MakeAmount(3150, 2),
							},
						},
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
							Category: tax.CategoryVAT,
							Percent:  num.NewPercentage(210, 3),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
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
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(25000, 2),
								Percent: num.NewPercentage(210, 3),
								Amount:  num.MakeAmount(5250, 2),
							},
						},
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
							Category: tax.CategoryVAT,
							Rate:     tax.RateGeneral,
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Rate:     tax.RateReduced,
						},
					},
					amount: num.MakeAmount(15000, 2),
				},
			},
			taxIncluded: "",
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(10000, 2),
								Percent: num.NewPercentage(210, 3),
								Amount:  num.MakeAmount(2100, 2),
							},
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(15000, 2),
								Percent: num.NewPercentage(100, 3),
								Amount:  num.MakeAmount(1500, 2),
							},
						},
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
							Category: tax.CategoryVAT,
							Percent:  num.NewPercentage(210, 3),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
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
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(10000, 2),
								Percent: num.NewPercentage(210, 3),
								Amount:  num.MakeAmount(2100, 2),
							},
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(15000, 2),
								Percent: num.NewPercentage(100, 3),
								Amount:  num.MakeAmount(1500, 2),
							},
						},
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
							Category: tax.CategoryVAT,
							Rate:     tax.RateGeneral,
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Rate:     tax.RateReduced,
						},
					},
					amount: num.MakeAmount(15000, 2),
				},
			},
			taxIncluded: tax.CategoryVAT,
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(8264, 2),
								Percent: num.NewPercentage(210, 3),
								Amount:  num.MakeAmount(1736, 2),
							},
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(13636, 2),
								Percent: num.NewPercentage(100, 3),
								Amount:  num.MakeAmount(1364, 2),
							},
						},
						Amount: num.MakeAmount(3099, 2),
					},
				},
				Sum: num.MakeAmount(3099, 2),
			},
		},
		{
			desc: "with multirate VAT as percentages, and included in price",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Percent:  num.NewPercentage(21, 2),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Percent:  num.NewPercentage(10, 2),
						},
					},
					amount: num.MakeAmount(15000, 2),
				},
			},
			taxIncluded: tax.CategoryVAT,
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(8264, 2),
								Percent: num.NewPercentage(21, 2),
								Amount:  num.MakeAmount(1736, 2),
							},
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(13636, 2),
								Percent: num.NewPercentage(10, 2),
								Amount:  num.MakeAmount(1364, 2),
							},
						},
						Amount: num.MakeAmount(3099, 2),
					},
				},
				Sum: num.MakeAmount(3099, 2),
			},
		},
		{
			desc: "with multirate VAT and retained tax",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Rate:     tax.RateGeneral,
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
							Category: tax.CategoryVAT,
							Rate:     tax.RateReduced,
						},
					},
					amount: num.MakeAmount(15000, 2),
				},
			},
			taxIncluded: "",
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(10000, 2),
								Percent: num.NewPercentage(210, 3),
								Amount:  num.MakeAmount(2100, 2),
							},
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(15000, 2),
								Percent: num.NewPercentage(100, 3),
								Amount:  num.MakeAmount(1500, 2),
							},
						},
						Amount: num.MakeAmount(3600, 2),
					},
					{
						Code:     es.TaxCategoryIRPF,
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
				Sum:      num.MakeAmount(3600, 2),
				Retained: num.NewAmount(1500, 2),
			},
		},

		{
			desc: "with multirate VAT included in price plus retained tax",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Rate:     tax.RateGeneral,
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
							Category: tax.CategoryVAT,
							Rate:     tax.RateReduced,
						},
					},
					amount: num.MakeAmount(15000, 2),
				},
			},
			taxIncluded: tax.CategoryVAT,
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(8264, 2),
								Percent: num.NewPercentage(210, 3),
								Amount:  num.MakeAmount(1736, 2),
							},
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(13636, 2),
								Percent: num.NewPercentage(100, 3),
								Amount:  num.MakeAmount(1364, 2),
							},
						},
						Amount: num.MakeAmount(3099, 2),
					},
					{
						Code:     es.TaxCategoryIRPF,
						Retained: true,
						Rates: []*tax.RateTotal{
							{
								Base:    num.MakeAmount(8264, 2),
								Percent: num.NewPercentage(150, 3),
								Amount:  num.MakeAmount(1240, 2),
							},
						},
						Amount: num.MakeAmount(1240, 2),
					},
				},
				Sum:      num.MakeAmount(3099, 2),
				Retained: num.NewAmount(1240, 2),
			},
		},
		{
			desc: "with invalid category",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: cbc.Code("FOO"),
							Rate:     tax.RateGeneral,
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
							Rate:     tax.RateGeneral,
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			err:        tax.ErrInvalid,
			errContent: "invalid: 'general' rate not defined in category 'IRPF'",
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
			desc: "with retained tax included",
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
			errContent:  "cannot include retained category 'IRPF'",
		},
		{
			desc:    "with informative tax included",
			country: "BR",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: br.TaxCategoryISS,
						},
					},
				},
			},
			taxIncluded: br.TaxCategoryISS,
			err:         tax.ErrInvalidPricesInclude,
			errContent:  "cannot include informative category 'ISS'",
		},
		{
			desc: "tax included with exempt key",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Key:      tax.KeyExempt,
							Ext: tax.ExtensionsOf(cbc.CodeMap{
								tbai.ExtKeyExempt: "E1",
							}),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			taxIncluded: tax.CategoryVAT,
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code: tax.CategoryVAT,
						Rates: []*tax.RateTotal{
							{
								Key: tax.KeyExempt,
								Ext: tax.ExtensionsOf(cbc.CodeMap{
									tbai.ExtKeyExempt: "E1",
								}),
								Base:   num.MakeAmount(10000, 2),
								Amount: num.MakeAmount(0, 2),
							},
						},
						Amount: num.MakeAmount(0, 2),
					},
				},
				Sum: num.MakeAmount(0, 2),
			},
		},
		{
			desc: "tax included with exempt rate and no key",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Ext: tax.ExtensionsOf(cbc.CodeMap{
								tbai.ExtKeyExempt: "E1",
							}),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Ext: tax.ExtensionsOf(cbc.CodeMap{
								tbai.ExtKeyExempt: "E1",
							}),
						},
					},
					amount: num.MakeAmount(2000, 2),
				},
			},
			taxIncluded: tax.CategoryVAT,
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code: tax.CategoryVAT,
						Rates: []*tax.RateTotal{
							{
								Key: tax.KeyStandard,
								Ext: tax.ExtensionsOf(cbc.CodeMap{
									tbai.ExtKeyExempt: "E1",
								}),
								Base:   num.MakeAmount(12000, 2),
								Amount: num.MakeAmount(0, 2),
							},
						},
						Amount: num.MakeAmount(0, 2),
					},
				},
				Sum: num.MakeAmount(0, 2),
			},
		},
		{
			desc: "tax included with regular and exempt rate",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Percent:  num.NewPercentage(21, 2),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Key:      tax.KeyExempt,
							Ext: tax.ExtensionsOf(cbc.CodeMap{
								tbai.ExtKeyExempt: "E2",
							}),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			taxIncluded: tax.CategoryVAT,
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code: tax.CategoryVAT,
						Rates: []*tax.RateTotal{
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(8264, 2),
								Percent: num.NewPercentage(21, 2),
								Amount:  num.MakeAmount(1736, 2),
							},
							{
								Key: tax.KeyExempt,
								Ext: tax.ExtensionsOf(cbc.CodeMap{
									tbai.ExtKeyExempt: "E2",
								}),
								Base:   num.MakeAmount(10000, 2),
								Amount: num.MakeAmount(0, 2),
							},
						},
						Amount: num.MakeAmount(1736, 2),
					},
				},
				Sum: num.MakeAmount(1736, 2),
			},
		},
		{
			desc:    "multiple different retained rates",
			country: "IT",
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Key:      tax.KeyStandard,
							Percent:  num.NewPercentage(220, 3),
						},
						{
							Category: it.TaxCategoryIRPEF,
							Ext: tax.ExtensionsOf(cbc.CodeMap{
								"it-sdi-retained": "A",
							}),
							Percent: num.NewPercentage(20, 2),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Percent:  num.NewPercentage(220, 3),
						},
						{
							Category: it.TaxCategoryIRPEF,
							Ext: tax.ExtensionsOf(cbc.CodeMap{
								"it-sdi-retained": "J", // truffles!
							}),
							Percent: num.NewPercentage(20, 2),
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code: tax.CategoryVAT,
						Rates: []*tax.RateTotal{
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(20000, 2),
								Percent: num.NewPercentage(220, 3),
								Amount:  num.MakeAmount(4400, 2),
							},
						},
						Amount: num.MakeAmount(4400, 2),
					},
					{
						Code:     it.TaxCategoryIRPEF,
						Retained: true,
						Rates: []*tax.RateTotal{
							{
								Ext: tax.ExtensionsOf(cbc.CodeMap{
									"it-sdi-retained": "A",
								}),
								Base:    num.MakeAmount(10000, 2),
								Percent: num.NewPercentage(20, 2),
								Amount:  num.MakeAmount(2000, 2),
							},
							{
								Ext: tax.ExtensionsOf(cbc.CodeMap{
									"it-sdi-retained": "J",
								}),
								Base:    num.MakeAmount(10000, 2),
								Percent: num.NewPercentage(20, 2),
								Amount:  num.MakeAmount(2000, 2),
							},
						},
						Amount: num.MakeAmount(4000, 2),
					},
				},
				Sum:      num.MakeAmount(4400, 2),
				Retained: num.NewAmount(4000, 2),
			},
		},
		{
			desc:     "currency rounding calculation",
			country:  "GR", // Greece uses currency rounding
			rounding: tax.RoundingRuleCurrency,
			lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Rate:     tax.RateGeneral,
						},
					},
					amount: num.MakeAmount(942, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Rate:     tax.RateReduced,
						},
					},
					amount: num.MakeAmount(942, 2),
				},
			},
			want: &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:     tax.CategoryVAT,
						Retained: false,
						Rates: []*tax.RateTotal{
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(942, 2),
								Percent: num.NewPercentage(24, 2),
								Amount:  num.MakeAmount(226, 2),
							},
							{
								Key:     tax.KeyStandard,
								Base:    num.MakeAmount(942, 2),
								Percent: num.NewPercentage(13, 2),
								Amount:  num.MakeAmount(122, 2),
							},
						},
						Amount: num.MakeAmount(348, 2), // with sum-then-round this would be 3.49
					},
				},
				Sum: num.MakeAmount(348, 2), // with sum-then-round this would be 3.49
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			d := date
			if test.date != nil {
				d = *test.date
			}
			country := l10n.ES.Tax()
			if test.country != "" {
				country = test.country
			}
			tc := &tax.TotalCalculator{
				Country:  country,
				Currency: currency.EUR,
				Rounding: test.rounding,
				Date:     d,
				Lines:    test.lines,
				Includes: test.taxIncluded,
			}
			tot := new(tax.Total)
			err := tc.Calculate(tot)
			if test.err != nil && assert.Error(t, err) {
				assert.ErrorIs(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
			if test.errContent != "" {
				assert.ErrorContains(t, err, test.errContent)
			}
			tot.Round(currency.EUR.Def().Zero())
			if test.want != nil {
				want, err := json.Marshal(test.want)
				require.NoError(t, err)
				got, err := json.Marshal(tot)
				require.NoError(t, err)
				if !assert.JSONEq(t, string(want), string(got)) {
					data, _ := json.MarshalIndent(tot, "", "  ")
					t.Logf("data output: %v", string(data))
				}
			}
		})
	}

}

func TestTotalCalculatorCurrencyRoundingIncludedTaxes(t *testing.T) {
	date := cal.MakeDate(2026, 7, 9)
	zero := currency.EUR.Def().Zero()
	calculate := func(t *testing.T, country l10n.TaxCountryCode, lines ...tax.TaxableLine) *tax.Total {
		t.Helper()
		tc := &tax.TotalCalculator{
			Country:  country,
			Currency: currency.EUR,
			Rounding: tax.RoundingRuleCurrency,
			Date:     date,
			Lines:    lines,
			Includes: tax.CategoryVAT,
		}
		tot := new(tax.Total)
		require.NoError(t, tc.Calculate(tot))
		tot.Round(zero)
		return tot
	}

	t.Run("single rate over multiple lines", func(t *testing.T) {
		// 12 lines of 125.00 with 6% VAT included
		lines := make([]tax.TaxableLine, 12)
		for i := range lines {
			lines[i] = &taxableLine{
				taxes:  tax.Set{{Category: tax.CategoryVAT, Rate: tax.RateReduced}},
				amount: num.MakeAmount(12500, 2),
			}
		}
		tot := calculate(t, "PT", lines...)
		rt := tot.Categories[0].Rates[0]
		assert.Equal(t, "1415.09", rt.Base.String())
		assert.Equal(t, "84.91", rt.Amount.String())
		assert.Equal(t, "84.91", tot.Sum.String())
		// the base and amount add up to the original tax-inclusive sum
		assert.Equal(t, "1500.00", rt.Base.Add(rt.Amount).String())
	})

	t.Run("lines with individual rounding errors", func(t *testing.T) {
		// each line's price of 1.93 with 21% VAT included has a base of
		// 1.5950 that individually rounds up to 1.60
		lines := make([]tax.TaxableLine, 10)
		for i := range lines {
			lines[i] = &taxableLine{
				taxes:  tax.Set{{Category: tax.CategoryVAT, Rate: tax.RateGeneral}},
				amount: num.MakeAmount(193, 2),
			}
		}
		tot := calculate(t, "ES", lines...)
		rt := tot.Categories[0].Rates[0]
		assert.Equal(t, "15.95", rt.Base.String())
		assert.Equal(t, "3.35", rt.Amount.String())
		assert.Equal(t, "19.30", rt.Base.Add(rt.Amount).String())
	})

	t.Run("with retained taxes", func(t *testing.T) {
		lines := make([]tax.TaxableLine, 10)
		for i := range lines {
			lines[i] = &taxableLine{
				taxes: tax.Set{
					{Category: tax.CategoryVAT, Rate: tax.RateGeneral},
					{Category: es.TaxCategoryIRPF, Rate: "pro"},
				},
				amount: num.MakeAmount(193, 2),
			}
		}
		tot := calculate(t, "ES", lines...)
		vat := tot.Categories[0].Rates[0]
		irpf := tot.Categories[1].Rates[0]
		// the retained tax base matches the VAT base exactly
		assert.Equal(t, "15.95", vat.Base.String())
		assert.Equal(t, "15.95", irpf.Base.String())
		assert.Equal(t, "2.39", irpf.Amount.String())
		assert.Equal(t, "3.35", tot.Sum.String())
		assert.Equal(t, "2.39", tot.Retained.String())
	})

	t.Run("retained tax over multiple rate groups", func(t *testing.T) {
		mk := func(cents int64, rate cbc.Key) tax.TaxableLine {
			return &taxableLine{
				taxes: tax.Set{
					{Category: tax.CategoryVAT, Rate: rate},
					{Category: es.TaxCategoryIRPF, Rate: "pro"},
				},
				amount: num.MakeAmount(cents, 2),
			}
		}
		tot := calculate(t, "ES",
			mk(193, tax.RateGeneral),
			mk(314, tax.RateGeneral),
			mk(677, tax.RateGeneral),
			mk(187, tax.RateReduced),
			mk(407, tax.RateReduced),
			mk(297, tax.RateReduced),
		)
		general := tot.Categories[0].Rates[0]
		reduced := tot.Categories[0].Rates[1]
		irpf := tot.Categories[1].Rates[0]
		// each rate group's base and amount add up to its tax-inclusive sum,
		// with the amount maintained from the group total even when it
		// differs from a recalculation over the rounded base
		assert.Equal(t, "9.79", general.Base.String())
		assert.Equal(t, "2.05", general.Amount.String())
		assert.Equal(t, "8.10", reduced.Base.String())
		assert.Equal(t, "0.81", reduced.Amount.String())
		// the retained base spans both groups without rounding drift
		assert.Equal(t, "17.89", irpf.Base.String())
		assert.Equal(t, "2.68", irpf.Amount.String())
	})

	t.Run("with exempt lines", func(t *testing.T) {
		tot := calculate(t, "ES",
			&taxableLine{
				taxes:  tax.Set{{Category: tax.CategoryVAT, Rate: tax.RateGeneral}},
				amount: num.MakeAmount(193, 2),
			},
			&taxableLine{
				taxes:  tax.Set{{Category: tax.CategoryVAT, Key: tax.KeyExempt}},
				amount: num.MakeAmount(1000, 2),
			},
		)
		general := tot.Categories[0].Rates[0]
		exempt := tot.Categories[0].Rates[1]
		assert.Equal(t, "1.60", general.Base.String())
		assert.Equal(t, "0.33", general.Amount.String())
		assert.Equal(t, "10.00", exempt.Base.String())
		assert.Equal(t, "0.00", exempt.Amount.String())
	})
}

func TestTotalCalculatorCurrencyRounding(t *testing.T) {
	date := cal.MakeDate(2026, 7, 9)
	zero := currency.EUR.Def().Zero()
	calculate := func(t *testing.T, country l10n.TaxCountryCode, includes cbc.Code, lines ...tax.TaxableLine) (*tax.Total, error) {
		t.Helper()
		tc := &tax.TotalCalculator{
			Country:  country,
			Currency: currency.EUR,
			Rounding: tax.RoundingRuleCurrency,
			Date:     date,
			Lines:    lines,
			Includes: includes,
		}
		tot := new(tax.Total)
		err := tc.Calculate(tot)
		tot.Round(zero)
		return tot, err
	}

	t.Run("with surcharges", func(t *testing.T) {
		lines := make([]tax.TaxableLine, 2)
		for i := range lines {
			lines[i] = &taxableLine{
				taxes: tax.Set{{
					Category:  tax.CategoryVAT,
					Percent:   num.NewPercentage(21, 2),
					Surcharge: num.NewPercentage(52, 3),
				}},
				amount: num.MakeAmount(10000, 2),
			}
		}
		lines = append(lines, &taxableLine{
			taxes: tax.Set{{
				Category:  tax.CategoryVAT,
				Percent:   num.NewPercentage(10, 2),
				Surcharge: num.NewPercentage(14, 3),
			}},
			amount: num.MakeAmount(10000, 2),
		})
		tot, err := calculate(t, "ES", "", lines...)
		require.NoError(t, err)
		ct := tot.Categories[0]
		rt := ct.Rates[0]
		assert.Equal(t, "200.00", rt.Base.String())
		assert.Equal(t, "42.00", rt.Amount.String())
		assert.Equal(t, "10.40", rt.Surcharge.Amount.String())
		assert.Equal(t, "11.80", ct.Surcharge.String())
		assert.Equal(t, "63.80", tot.Sum.String())
	})

	t.Run("included taxes with surcharges", func(t *testing.T) {
		lines := make([]tax.TaxableLine, 2)
		for i := range lines {
			lines[i] = &taxableLine{
				taxes: tax.Set{{
					Category:  tax.CategoryVAT,
					Percent:   num.NewPercentage(21, 2),
					Surcharge: num.NewPercentage(52, 3),
				}},
				amount: num.MakeAmount(12100, 2),
			}
		}
		tot, err := calculate(t, "ES", tax.CategoryVAT, lines...)
		require.NoError(t, err)
		rt := tot.Categories[0].Rates[0]
		assert.Equal(t, "200.00", rt.Base.String())
		assert.Equal(t, "42.00", rt.Amount.String())
		assert.Equal(t, "10.40", rt.Surcharge.Amount.String())
		assert.Equal(t, "52.40", tot.Sum.String())
	})

	t.Run("with informative categories", func(t *testing.T) {
		tot, err := calculate(t, "BR", "",
			&taxableLine{
				taxes: tax.Set{{
					Category: br.TaxCategoryISS,
					Percent:  num.NewPercentage(5, 2),
				}},
				amount: num.MakeAmount(20000, 2),
			},
		)
		require.NoError(t, err)
		ct := tot.Categories[0]
		assert.True(t, ct.Informative)
		assert.Equal(t, "10.00", ct.Amount.String())
		assert.Equal(t, "0.00", tot.Sum.String())
		assert.Nil(t, tot.Retained)
	})

	t.Run("with multiple retained categories", func(t *testing.T) {
		tot, err := calculate(t, "BR", "",
			&taxableLine{
				taxes: tax.Set{
					{
						Category: br.TaxCategoryPISRet,
						Percent:  num.NewPercentage(65, 4),
					},
					{
						Category: br.TaxCategoryCOFINSRet,
						Percent:  num.NewPercentage(3, 2),
					},
				},
				amount: num.MakeAmount(100000, 2),
			},
		)
		require.NoError(t, err)
		require.NotNil(t, tot.Retained)
		assert.Equal(t, "36.50", tot.Retained.String())
		assert.Equal(t, "0.00", tot.Sum.String())
	})

	t.Run("error including retained category", func(t *testing.T) {
		_, err := calculate(t, "ES", "IRPF",
			&taxableLine{
				taxes:  tax.Set{{Category: "IRPF", Percent: num.NewPercentage(15, 2)}},
				amount: num.MakeAmount(10000, 2),
			},
		)
		require.ErrorContains(t, err, "cannot include retained category 'IRPF'")
	})

	t.Run("error including informative category", func(t *testing.T) {
		_, err := calculate(t, "BR", br.TaxCategoryISS,
			&taxableLine{
				taxes: tax.Set{{
					Category: br.TaxCategoryISS,
					Percent:  num.NewPercentage(5, 2),
				}},
				amount: num.MakeAmount(10000, 2),
			},
		)
		require.ErrorContains(t, err, "cannot include informative category 'ISS'")
	})

	t.Run("lines without the included category", func(t *testing.T) {
		tot, err := calculate(t, "ES", tax.CategoryVAT,
			&taxableLine{
				taxes:  tax.Set{{Category: tax.CategoryVAT, Rate: tax.RateGeneral}},
				amount: num.MakeAmount(193, 2),
			},
			&taxableLine{
				taxes:  tax.Set{{Category: es.TaxCategoryIRPF, Percent: num.NewPercentage(15, 2)}},
				amount: num.MakeAmount(1000, 2),
			},
		)
		require.NoError(t, err)
		vat := tot.Categories[0].Rates[0]
		irpf := tot.Categories[1].Rates[0]
		assert.Equal(t, "1.60", vat.Base.String())
		assert.Equal(t, "0.33", vat.Amount.String())
		assert.Equal(t, "10.00", irpf.Base.String())
		assert.Equal(t, "1.50", irpf.Amount.String())
	})

	t.Run("line totals with extra precision", func(t *testing.T) {
		tot, err := calculate(t, "ES", "",
			&taxableLine{
				taxes:  tax.Set{{Category: tax.CategoryVAT, Rate: tax.RateGeneral}},
				amount: num.MakeAmount(1000049, 4), // 100.0049
			},
		)
		require.NoError(t, err)
		rt := tot.Categories[0].Rates[0]
		assert.Equal(t, "100.00", rt.Base.String())
		assert.Equal(t, "21.00", rt.Amount.String())
	})

	t.Run("repeated calculations", func(t *testing.T) {
		tc := &tax.TotalCalculator{
			Country:  "ES",
			Currency: currency.EUR,
			Rounding: tax.RoundingRuleCurrency,
			Date:     date,
			Lines: []tax.TaxableLine{
				&taxableLine{
					taxes: tax.Set{
						{Category: tax.CategoryVAT, Rate: tax.RateGeneral},
						{Category: es.TaxCategoryIRPF, Rate: "pro"},
					},
					amount: num.MakeAmount(10000, 2),
				},
			},
		}
		tot := new(tax.Total)
		require.NoError(t, tc.Calculate(tot))
		require.NoError(t, tc.Calculate(tot))
		require.NotNil(t, tot.Retained)
		assert.Equal(t, "15.00", tot.Retained.String())
		assert.Equal(t, "21.00", tot.Sum.String())
	})
}
