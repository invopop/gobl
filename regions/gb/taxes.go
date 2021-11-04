package gb

import (
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regions/common"
	"github.com/invopop/gobl/tax"
)

var taxRegion = tax.Region{
	Code: "GB",
	Name: i18n.String{
		i18n.EN: "Great Britain",
	},
	Categories: []tax.Category{
		//
		// VAT
		//
		{
			Code: common.TaxCategoryVAT,
			Name: i18n.String{
				i18n.EN: "VAT",
			},
			Desc: i18n.String{
				i18n.EN: "Value Added Tax",
			},
			Retained: false,
			Defs: []tax.Def{
				{
					Code: common.TaxRateVATZero,
					Name: i18n.String{
						i18n.EN: "VAT Zero Rate",
					},
					Values: []tax.Value{
						{
							Percent: num.MakePercentage(0, 3),
						},
					},
				},
				{
					Code: common.TaxRateVATStandard,
					Name: i18n.String{
						i18n.EN: "VAT Standard Rate",
					},
					Values: []tax.Value{
						{
							Since:   org.NewDate(2011, 1, 4),
							Percent: num.MakePercentage(200, 3),
						},
					},
				},
				{
					Code: common.TaxRateVATReduced,
					Name: i18n.String{
						i18n.EN: "VAT Reduced Rate",
					},
					Values: []tax.Value{
						{
							Since:   org.NewDate(2011, 1, 4),
							Percent: num.MakePercentage(50, 3),
						},
					},
				},
			},
		},
	},
}
