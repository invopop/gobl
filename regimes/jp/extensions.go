package jp

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Extension keys for the reduced consumption tax rate (軽減税率).
//
// Under the Qualified Invoice System, invoices with mixed rates must clearly identify which line items are subject to
// the reduced 8% rate.
const (
	// ExtKeyReducedRate is the extension key used on line items to indicate they qualify for the 8% reduced consumption
	// tax rate instead of the standard 10% rate.
	ExtKeyReducedRate = "jp-reduced-rate"
)

// Reduced rate item codes identifying the category of goods eligible for the 8% rate, as defined by the NTA
// reduced-rate guidance.
const (
	// ExtCodeFoodBeverage identifies food and non-alcoholic beverages purchased for takeaway, delivery, or retail.
	// Dining in at restaurants and alcoholic beverages are excluded and taxed at the standard 10% rate.
	ExtCodeFoodBeverage = "food-beverage"

	// ExtCodeNewspaper identifies newspaper subscriptions issued at least twice a week under a subscription contract.
	ExtCodeNewspaper = "newspaper"
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyReducedRate,
		Name: i18n.String{
			i18n.EN: "Reduced Rate Item",
			i18n.JA: "軽減税率対象品目",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Extensions for items subject to the reduced consumption tax rate (8%).

				Items subject to the reduced tax rate include:
				1. Food and drink excluding alcoholic beverages and dining out
				2. Newspapers issued at least twice a week by subscription

				Use these extensions to identify line items that qualify for the 8% rate
				instead of the standard 10% rate.
			`),
			i18n.JA: here.Doc(`
				軽減税率（8%）の対象となる品目の拡張。

				軽減税率の対象品目：
				1. アルコール飲料および外食を除く飲食料品
				2. 定期購読契約により週2回以上発行される新聞

				これらの拡張を使用して、標準税率10%ではなく8%の対象となる
				明細項目を識別します。
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: ExtCodeFoodBeverage,
				Name: i18n.String{
					i18n.EN: "Food and Non-Alcoholic Beverages (Takeaway)",
					i18n.JA: "飲食料品（テイクアウト）",
				},
				Desc: i18n.String{
					i18n.EN: "Food and drink excluding alcoholic beverages and dining out. Applies to takeaway, delivery, and grocery purchases.",
					i18n.JA: "アルコール飲料および外食を除く飲食料品。テイクアウト、配達、食料品の購入に適用されます。",
				},
			},
			{
				Code: ExtCodeNewspaper,
				Name: i18n.String{
					i18n.EN: "Newspaper Subscription",
					i18n.JA: "新聞購読",
				},
				Desc: i18n.String{
					i18n.EN: "Newspapers issued at least twice a week by subscription contract.",
					i18n.JA: "定期購読契約により週2回以上発行される新聞。",
				},
			},
		},
	},
}
