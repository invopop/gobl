package mt

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// Rates run from 1 January 1999, when the VAT Act (Cap. 406) came into force. Malta's
// 1995 VAT was charged under the separate Value Added Tax Act, 1994, repealed in 1997,
// and is out of scope here. The 7% tier only appears in 2011.
//
// Since versions a percentage, not the supplies it applies to: the Eighth Schedule items
// named below are today's and are not versioned.
var taxCategories = []*tax.CategoryDef{
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.MT: "VAT",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.MT: "Taxxa fuq il-Valur Miżjud",
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
				},
				Description: i18n.String{
					i18n.EN: "Article 19(1) of the VAT Act. Applies to any taxable supply for which no reduced rate or exemption is specified.",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(180, 3), //18%
						Since:   cal.NewDate(2004, 1, 1),
					},
					{
						Percent: num.MakePercentage(150, 3), //15%
						Since:   cal.NewDate(1999, 1, 1),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateIntermediate,
				Name: i18n.String{
					i18n.EN: "12% Rate",
				},
				Description: i18n.String{
					i18n.EN: "Items 12 to 15 of the Eighth Schedule: custody and management of securities, management of credit and credit guarantees, short-term hiring of pleasure boats, and care of the human body by regulated health professionals.",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(120, 3), //12%
						Since:   cal.NewDate(2024, 1, 1),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "7% Rate",
				},
				Description: i18n.String{
					i18n.EN: "Items 1 and 11 of the Eighth Schedule: licensed tourist accommodation, and the use of sporting facilities.",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(70, 3), //7%
						Since:   cal.NewDate(2011, 1, 1),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "5% Rate",
				},
				Description: i18n.String{
					i18n.EN: "Items 2 to 10 of the Eighth Schedule: electricity, confectionery, medical accessories, printed matter, items for the exclusive use of the disabled, the importation of works of art, minor repairs, domestic care services, and admission to cultural events.",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(50, 3), //5%
						Since:   cal.NewDate(1999, 1, 1),
					},
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("VAT Act (Chapter 406) - Eighth Schedule, Rate of tax"),
				URL:   "https://legislation.mt/eli/cap/406/eng",
			},
			{
				Title: i18n.NewString("Legal Notice 231 of 2023 - Amendment of Eighth Schedule"),
				URL:   "https://legislation.mt/eli/ln/2023/231/eng",
			},
			{
				Title: i18n.NewString("MTCA - VAT Rates & Exemptions FAQs"),
				URL:   "https://mtca.gov.mt/docs/default-source/documents/business-tax/vat/faqs/vat-rates-exemptions---faqs.pdf",
			},
			{
				Title:       i18n.NewString("European Commission (DG TAXUD) - VAT rates applied in the Member States of the European Union, section VIII: evolution of the VAT rates"),
				URL:         "https://taxation-customs.ec.europa.eu/system/files/2020-10/vat_rates_en.pdf",
				ContentType: "application/pdf",
				At:          cal.NewDateTime(2026, 9, 4, 0, 0, 0),
			},
		},
	},
}
