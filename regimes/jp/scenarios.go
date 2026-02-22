package jp

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// invoiceTags returns the set of invoice-level tags supported by the JP regime. Each tag controls validation behaviour
// and triggers a corresponding scenario.
func invoiceTags() *tax.TagSet {
	return &tax.TagSet{
		Schema: bill.ShortSchemaInvoice,
		List: []*cbc.Definition{
			{
				Key: TagExport,
				Name: i18n.String{
					i18n.EN: "Export",
					i18n.JA: "輸出",
				},
				Desc: i18n.String{
					i18n.EN: "Zero-rated export invoice. All consumption tax lines must use the zero rate.",
					i18n.JA: "輸出免税請求書。すべての消費税行はゼロ税率を使用する必要があります。",
				},
			},
			{
				Key: TagSimplified,
				Name: i18n.String{
					i18n.EN: "Simplified Invoice",
					i18n.JA: "簡易適格請求書",
				},
				Desc: i18n.String{
					i18n.EN: "Simplified qualified invoice for retail, restaurant, taxi and similar businesses. Buyer name is not required.",
					i18n.JA: "小売、飲食店、タクシーなどの事業者向けの簡易適格請求書。購入者名は不要です。",
				},
			},
			{
				Key: TagSelfBilling,
				Name: i18n.String{
					i18n.EN: "Self-billed Invoice",
					i18n.JA: "仕入明細書",
				},
				Desc: i18n.String{
					i18n.EN: "Buyer-issued invoice that substitutes for a qualified invoice once confirmed by the supplier.",
					i18n.JA: "供給者の確認後、適格請求書の代わりとなる購入者発行の請求書。",
				},
			},
		},
	}
}

// invoiceScenarios returns scenarios that inject bilingual legal notes into invoices during Calculate() based on the
// tags set on the invoice.
func invoiceScenarios() *tax.ScenarioSet {
	return &tax.ScenarioSet{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			{
				Tags: []cbc.Key{TagExport},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  TagExport,
					Text: "Export - Zero rated consumption tax / 輸出取引（免税）",
				},
			},
			{
				Tags: []cbc.Key{TagSimplified},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  TagSimplified,
					Text: "Simplified Qualified Invoice / 簡易適格請求書",
				},
			},
			{
				Tags: []cbc.Key{TagSelfBilling},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  TagSelfBilling,
					Text: "Self-billed invoice / 仕入明細書",
				},
			},
		},
	}
}
