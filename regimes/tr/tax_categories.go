package tr

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	//
	// VAT
	//
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.TR: "KDV",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.TR: "Katma Değer Vergisi",
		},
		Description: &i18n.String{
			i18n.EN: here.Doc(`
			Known in Turkish as "Katma Değer Vergisi" (KDV), governed by Law No. 3065.
			Applies to the supply of goods and services in Türkiye and imports, with three
			rates: standard (20%), reduced (10%) for basic foodstuffs, textiles and books,
			and super-reduced (1%) for certain agricultural products and essential services.
			Historical rates are recorded back to November 1999.
			`),
			i18n.TR: here.Doc(`
			3065 sayılı KDV Kanunu kapsamında mal ve hizmet teslimlerinde uygulanan vergi.
			Genel oran %20, indirimli oran %10 (temel gıda, tekstil, kitap), özel indirimli
			oran %1 (bazı tarım ürünleri ve zorunlu hizmetler).
			Geçmiş oranlar Kasım 1999'a kadar kayıtlıdır.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Revenue Administration",
					i18n.TR: "Gelir İdaresi Başkanlığı",
				},
				URL: "https://www.gib.gov.tr",
			},
			{
				Title: i18n.String{
					i18n.EN: "Council of Ministers Decree No. 2001/2344 - Standard rate raised from 17% to 18%",
					i18n.TR: "Bakanlar Kurulu Kararı 2001/2344 - Genel oran %17'den %18'e",
				},
				URL: "https://www.alomaliye.com/2001/05/10/bakanlar-kurulu-karari-20012344/",
			},
			{
				Title: i18n.String{
					i18n.EN: "EY Tax Alert - Türkiye increases VAT rates (Decree No. 7346, July 2023)",
				},
				URL: "https://www.ey.com/en_gl/technical/tax-alerts/turkiye-increases-vat-rates-on-goods-and-services",
			},
			{
				Title: i18n.String{
					i18n.EN: "VAT Rates in Turkey - Global VAT Compliance",
				},
				URL: "https://www.globalvatcompliance.com/vat-rates-in-turkey/",
			},
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.TR: "Genel Oran",
				},
				Values: []*tax.RateValueDef{
					{
						// Raised from 18% to 20% by Presidential Decree No. 7346,
						// published in the Official Gazette on 7 July 2023,
						// effective 10 July 2023.
						Since:   cal.NewDate(2023, 7, 10),
						Percent: num.MakePercentage(20, 2),
					},
					{
						// Raised from 17% to 18% by Council of Ministers Decree No. 2001/2344,
						// published in Official Gazette No. 24398 on 10 May 2001,
						// effective 15 May 2001.
						Since:   cal.NewDate(2001, 5, 15),
						Percent: num.MakePercentage(18, 2),
					},
					{
						// 17% rate set by Council of Ministers Decree No. 99/13648,
						// effective 28 November 1999.
						Since:   cal.NewDate(1999, 11, 28),
						Percent: num.MakePercentage(17, 2),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.TR: "İndirimli Oran",
				},
				Description: i18n.String{
					i18n.EN: "Applies to goods and services in List No. II of the KDV Law, including basic foodstuffs, textiles, and books.",
					i18n.TR: "KDV Kanunu'nun II sayılı listesindeki mal ve hizmetlere uygulanır.",
				},
				Values: []*tax.RateValueDef{
					{
						// Raised from 8% to 10% by Presidential Decree No. 7346,
						// published in the Official Gazette on 7 July 2023,
						// effective 10 July 2023.
						Since:   cal.NewDate(2023, 7, 10),
						Percent: num.MakePercentage(10, 2),
					},
					{
						// 8% reduced rate, unchanged across Decrees 99/13648 and 2001/2344.
						Since:   cal.NewDate(1999, 11, 28),
						Percent: num.MakePercentage(8, 2),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Super-Reduced Rate",
					i18n.TR: "Özel İndirimli Oran",
				},
				Description: i18n.String{
					i18n.EN: "Applies to goods and services in List No. I of the KDV Law, including certain agricultural products, bread, and funeral services.",
					i18n.TR: "KDV Kanunu'nun I sayılı listesindeki mal ve hizmetlere uygulanır.",
				},
				Values: []*tax.RateValueDef{
					{
						// 1% super-reduced rate, unchanged across all decrees since 1999.
						Since:   cal.NewDate(1999, 11, 28),
						Percent: num.MakePercentage(1, 2),
					},
				},
			},
		},
	},
}
