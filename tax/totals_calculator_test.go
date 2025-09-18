package tax_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/addons/es/tbai"
	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/addons/pt/saft"
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
							Ext: tax.Extensions{
								tbai.ExtKeyExempt: "E1",
							},
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
								Ext: tax.Extensions{
									tbai.ExtKeyExempt: "E1",
								},
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
							Ext: tax.Extensions{
								pt.ExtKeyRegion:    "PT-AC",
								saft.ExtKeyTaxRate: "NOR",
							},
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
								Ext: tax.Extensions{
									pt.ExtKeyRegion:    "PT-AC",
									saft.ExtKeyTaxRate: "NOR",
								},
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
							Ext: tax.Extensions{
								tbai.ExtKeyExempt: "E1",
							},
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
								Ext: tax.Extensions{
									tbai.ExtKeyExempt: "E1",
								},
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
							Ext: tax.Extensions{
								tbai.ExtKeyExempt: "E1",
							},
						},
					},
					amount: num.MakeAmount(10000, 2),
				},
				&taxableLine{
					taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Ext: tax.Extensions{
								tbai.ExtKeyExempt: "E1",
							},
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
								Ext: tax.Extensions{
									tbai.ExtKeyExempt: "E1",
								},
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
							Ext: tax.Extensions{
								tbai.ExtKeyExempt: "E2",
							},
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
								Ext: tax.Extensions{
									tbai.ExtKeyExempt: "E2",
								},
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
							Ext: tax.Extensions{
								sdi.ExtKeyRetained: "A",
							},
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
							Ext: tax.Extensions{
								sdi.ExtKeyRetained: "J", // truffles!
							},
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
								Ext: tax.Extensions{
									sdi.ExtKeyRetained: "A",
								},
								Base:    num.MakeAmount(10000, 2),
								Percent: num.NewPercentage(20, 2),
								Amount:  num.MakeAmount(2000, 2),
							},
							{
								Ext: tax.Extensions{
									sdi.ExtKeyRetained: "J",
								},
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
