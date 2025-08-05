// Package ae defines VAT tax categories specific to the United Arab Emirates.
package ae

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.AR: "ضريبة القيمة المضافة",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.AR: "ضريبة القيمة المضافة",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Federal Tax Authority - UAE VAT Regulations",
					i18n.AR: "الهيئة الاتحادية للضرائب",
				},
				URL: "https://www.tax.gov.ae",
			},
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General Rate",
					i18n.AR: "معدل قياسي",
				},
				Description: i18n.String{
					i18n.EN: "Applies to most goods and services unless specified otherwise.",
					i18n.AR: "ينطبق على معظم السلع والخدمات ما لم ينص على خلاف ذلك.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2018, 1, 1),
						Percent: num.MakePercentage(5, 2),
					},
				},
			},
		},
	},
}
