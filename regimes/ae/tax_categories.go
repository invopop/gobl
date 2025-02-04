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
		Rates: []*tax.RateDef{
			{
				Key: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.AR: "معدل صفر",
				},
				Description: i18n.String{
					i18n.EN: "A VAT rate of 0% applicable to specific exports, designated areas, and essential services.",
					i18n.AR: "نسبة ضريبة قيمة مضافة 0٪ تطبق على الصادرات المحددة والمناطق المعينة والخدمات الأساسية.",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
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
			{
				Key: tax.RateExempt,
				Name: i18n.String{
					i18n.EN: "Exempt",
					i18n.AR: "معفى",
				},
				Exempt: true,
				Description: i18n.String{
					i18n.EN: "Certain goods and services, such as financial services and residential real estate, are exempt from VAT.",
					i18n.AR: "بعض السلع والخدمات، مثل الخدمات المالية والعقارات السكنية، معفاة من ضريبة القيمة المضافة.",
				},
			},
		},
	},
}
