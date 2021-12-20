package nl

import (
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regions/common"
	"github.com/invopop/gobl/tax"
)

var taxRegion = tax.Region{
	Code: "ES",
	Name: i18n.String{
		i18n.EN: "The Netherlands",
		i18n.NL: "Nederland",
	},
	Categories: []tax.Category{
		//
		// VAT
		//
		{
			Code: common.TaxCategoryVAT,
			Name: i18n.String{
				i18n.EN: "VAT",
				i18n.NL: "BTW",
			},
			Desc: i18n.String{
				i18n.EN: "Value Added Tax",
				i18n.NL: "Belasting Toegevoegde Waarde",
			},
			Retained: false,
			Defs: []tax.Def{
				{
					Code: common.TaxRateVATZero,
					Name: i18n.String{
						i18n.EN: "VAT Zero Rate",
						i18n.NL: `BTW 0%-tarief`,
					},
					Values: []tax.Value{
						{
							Percent: num.MakePercentage(0, 0),
						},
					},
				},
				{
					Code: common.TaxRateVATStandard,
					Name: i18n.String{
						i18n.EN: "VAT Standard Rate",
						i18n.NL: "BTW Standaardtarief",
					},
					Values: []tax.Value{
						{
							Since:   org.NewDate(1900, 1, 1),
							Percent: num.MakePercentage(100, 21),
						},
					},
				},
				{
					Code: common.TaxRateVATReduced,
					Name: i18n.String{
						i18n.EN: "VAT Reduced Rate",
						i18n.NL: "BTW Gereduceerd Tarief",
					},
					Values: []tax.Value{
						{
							Since:   org.NewDate(1900, 1, 1),
							Percent: num.MakePercentage(100, 9),
						},
					},
				},
			},
		},
	},
}
