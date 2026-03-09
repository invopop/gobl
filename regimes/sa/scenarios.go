package sa

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

var invoiceTags = []*tax.TagSet{
	{
		Schema: bill.ShortSchemaInvoice,
		List: []*cbc.Definition{
			{
				Key: tax.TagExport,
				Name: i18n.String{
					i18n.EN: "Export",
					i18n.AR: "تصدير",
				},
			},
		},
	},
}

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		// Reverse Charges
		{
			Tags: []cbc.Key{tax.TagReverseCharge},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  tax.TagReverseCharge,
				Text: "Reverse Charge / آلية الاحتساب العكسي",
			},
		},
		// Simplified Tax Invoice
		{
			Tags: []cbc.Key{tax.TagSimplified},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  tax.TagSimplified,
				Text: "Simplified Tax Invoice / فاتورة ضريبية مبسطة",
			},
		},
		// Export (zero-rated per Chapter 6 of the VAT Implementing Regulations)
		{
			Tags: []cbc.Key{tax.TagExport},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  tax.TagExport,
				Text: "Export of goods or services, zero-rated per VAT Implementing Regulations / تصدير سلع أو خدمات، خاضع لنسبة صفر وفقاً للائحة التنفيذية لضريبة القيمة المضافة",
			},
		},
	},
}
