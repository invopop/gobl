package jp

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

const (
	// RateExempt identifies the exempt rate key for non-taxable transactions
	// such as financial services, medical care, and education.
	RateExempt cbc.Key = "exempt"
)

// taxCategories defines Japan's consumption tax as VAT (the same category code as other
// non-EU regimes; the local name is JCT / 消費税). Rate "general" aligns with
// GOBL's VAT combo normalization (standard → general).
// Sources: NTA (National Tax Agency) and Japan Customs official guidance.
var taxCategories = []*tax.CategoryDef{
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "JCT",
			i18n.JA: "消費税",
		},
		Title: i18n.String{
			i18n.EN: "Japanese Consumption Tax",
			i18n.JA: "消費税",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "National Tax Agency - Consumption Tax",
					i18n.JA: "国税庁 - 消費税",
				},
				URL: "https://www.nta.go.jp/english/taxes/consumption_tax/",
			},
			{
				Title: i18n.String{
					i18n.EN: "Qualified Invoice System",
					i18n.JA: "適格請求書保存方式",
				},
				URL: "https://www.nta.go.jp/taxes/shiraberu/zeimokubetsu/shohi/keigenzeiritsu1.htm",
			},
			{
				Title: i18n.String{
					i18n.EN: "NTA — Invoice system (reduced tax rate)",
					i18n.JA: "国税庁 — インボイス制度（軽減税率）",
				},
				URL: "https://www.nta.go.jp/taxes/shiraberu/zeimokubetsu/shohi/keigenzeiritsu/invoice.htm",
			},
			{
				Title: i18n.String{
					i18n.EN: "Japan Customs — Consumption tax exemption on exports",
					i18n.JA: "税関 — 輸出品に対する消費税の免除",
				},
				URL: "https://www.customs.go.jp/english/c-answer_e/extsukan/5003_e.htm",
			},
			{
				Title: i18n.String{
					i18n.EN: "NTA — Qualified invoice issuer registration",
					i18n.JA: "国税庁 — 適格請求書発行事業者の登録",
				},
				URL: "https://www.invoice-kohyo.nta.go.jp",
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
					i18n.JA: "標準税率",
				},
				Description: i18n.String{
					i18n.EN: here.Doc(`
							The standard consumption tax rate applies to most goods and services
							in Japan, including electronics, clothing, professional services,
							restaurant meals, and alcoholic beverages.

							Total: 10% (National: 7.8% + Local: 2.2%)
						`),
					i18n.JA: here.Doc(`
							標準税率は、電化製品、衣類、専門サービス、外食、
							アルコール飲料を含む日本のほとんどの商品およびサービスに適用されます。

							合計：10%（国税：7.8% + 地方税：2.2%）
						`),
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2019, 10, 1),
						Percent: num.MakePercentage(10, 2),
					},
					{
						Since:   cal.NewDate(2014, 4, 1),
						Percent: num.MakePercentage(8, 2),
					},
					{
						Since:   cal.NewDate(1997, 4, 1),
						Percent: num.MakePercentage(5, 2),
					},
					{
						Since:   cal.NewDate(1989, 4, 1),
						Percent: num.MakePercentage(3, 2),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.JA: "軽減税率",
				},
				Description: i18n.String{
					i18n.EN: here.Doc(`
							The reduced consumption tax rate applies to specific goods and services:

							- Food products (excluding restaurant meals and alcoholic beverages)
							- Non-alcoholic beverages
							- Newspaper subscriptions (published at least twice weekly)

							Total: 8% (National: 6.24% + Local: 1.76%)

							Note: Dining in at restaurants is subject to the standard 10% rate,
							while takeout and delivery of food items may qualify for the reduced 8% rate.
						`),
					i18n.JA: here.Doc(`
							軽減税率は、以下の特定の商品およびサービスに適用されます：

							- 食料品（外食およびアルコール飲料を除く）
							- 非アルコール飲料
							- 新聞購読（週2回以上発行）

							合計：8%（国税：6.24% + 地方税：1.76%）

							注：レストランでの店内飲食は標準税率10%の対象ですが、
							テイクアウトやデリバリーは軽減税率8%の対象となる場合があります。
						`),
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2019, 10, 1),
						Percent: num.MakePercentage(8, 2),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyZero},
				Rate: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate (Export)",
					i18n.JA: "輸出等（免税）",
				},
				Description: i18n.String{
					i18n.EN: here.Doc(`
							Exports and certain international transactions are subject to a zero
							consumption tax rate. This includes:

							- Export of goods
							- International transportation services
							- Services provided to non-residents

							Businesses can still claim input tax credits on purchases related
							to zero-rated supplies.
						`),
					i18n.JA: here.Doc(`
							輸出および特定の国際取引は、消費税0%の対象となります。これには以下が含まれます：

							- 貨物の輸出
							- 国際輸送サービス
							- 非居住者に提供されるサービス

							事業者は、免税取引に関連する購入について、
							仕入税額控除を請求することができます。
						`),
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(1989, 4, 1),
						Percent: num.MakePercentage(0, 2),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyExempt},
				Rate: RateExempt,
				Name: i18n.String{
					i18n.EN: "Exempt",
					i18n.JA: "非課税",
				},
				Description: i18n.String{
					i18n.EN: here.Doc(`
							Certain transactions are exempt from consumption tax, including:

							- Financial services (banking, insurance, securities)
							- Medical services under health insurance
							- Educational services
							- Land sale and leasing
							- Postal services

							Businesses making only exempt supplies cannot claim input tax credits.
						`),
					i18n.JA: here.Doc(`
							特定の取引は消費税の非課税対象となります。これには以下が含まれます：

							- 金融サービス（銀行、保険、証券）
							- 健康保険下の医療サービス
							- 教育サービス
							- 土地の売買および賃貸
							- 郵便サービス

							非課税取引のみを行う事業者は、仕入税額控除を受けることができません。
						`),
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(1989, 4, 1),
						Percent: num.MakePercentage(0, 2),
					},
				},
			},
		},
	},
}
