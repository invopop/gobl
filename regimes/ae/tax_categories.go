// Package ae defines VAT tax categories specific to the United Arab Emirates.
package ae

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

// TaxKeyMarginScheme is the VAT key for a supply under the profit margin
// scheme, where VAT is due on the supplier's margin and the invoice shows no
// VAT amount. See the FTA's Profit Margin Scheme guide (VATGPM1):
// https://tax.gov.ae/Datafolder/Files/Pdf/2026/Guide/Profit%20Margin-Scheme-EN-02-01-2026-re.pdf
const TaxKeyMarginScheme cbc.Key = "margin-scheme"

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
		Keys:     vatKeys(),
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

// vatKeys returns the global VAT keys plus the profit margin scheme key.
func vatKeys() []*tax.KeyDef {
	keys := append([]*tax.KeyDef{}, tax.GlobalVATKeys()...)
	return append(keys, &tax.KeyDef{
		Key: TaxKeyMarginScheme,
		Name: i18n.String{
			i18n.EN: "Profit Margin Scheme",
			i18n.AR: "آلية هامش الربح",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Supply of goods under the profit margin scheme, such as second-hand
				goods, antiques or collectors' items. VAT is due on the margin, not
				on the price, and the invoice shows no VAT amount.
			`),
		},
		NoPercent: true,
	})
}
