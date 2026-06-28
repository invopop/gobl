package kr

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
			i18n.KO: "부가세",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.KO: "부가가치세",
		},
		Retained: false,
		// The standard VAT keys cover zero-rated supplies (e.g. exports) and
		// exemptions, so no separate rate definitions are needed for them.
		Keys: tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General Rate",
					i18n.KO: "기본세율",
				},
				// VAT has applied at a single rate since it was introduced in Korea on
				// 1 July 1977, with no subsequent rate changes.
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(100, 3),
						Since:   cal.NewDate(1977, 7, 1),
					},
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Value-Added Tax Act, Article 30 (tax rate)",
				},
				URL: "https://elaw.klri.re.kr/eng_service/lawView.do?hseq=53110&lang=ENG",
			},
			{
				Title: i18n.String{
					i18n.EN: "National Tax Service (NTS)",
					i18n.KO: "국세청",
				},
				URL: "https://www.nts.go.kr/english/",
			},
		},
	},
}
