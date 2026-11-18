// Package ae provides tax scenarios specific to UAE VAT regulations.
package ae

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

// TagMarginScheme identifies an invoice for a supply under the profit margin
// scheme.
const TagMarginScheme cbc.Key = "margin-scheme"

var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*cbc.Definition{
		{
			Key: TagMarginScheme,
			Name: i18n.String{
				i18n.EN: "Profit Margin Scheme",
				i18n.AR: "آلية هامش الربح",
			},
			Desc: i18n.String{
				i18n.EN: here.Doc(`
					Supply of eligible goods, such as second-hand goods, antiques or
					collectors' items, under the profit margin scheme. Margin scheme
					lines use the "exempt" VAT key, as they do for Italy, Greece and
					Portugal. VAT is due on the margin and reported in the VAT return,
					so the invoice shows no VAT amount.
				`),
			},
		},
	},
}

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		// Reverse Charges
		{
			Tags:       []cbc.Key{tax.TagReverseCharge},
			Categories: []cbc.Code{tax.CategoryVAT},
			Note: &tax.Note{
				Category: tax.CategoryVAT,
				Key:      tax.KeyReverseCharge,
				Text:     "Reverse Charge",
			},
		},
		// Simplified Tax Invoice
		{
			Tags:       []cbc.Key{tax.TagSimplified},
			Categories: []cbc.Code{tax.CategoryVAT},
			Note: &tax.Note{
				Category: tax.CategoryVAT,
				Key:      tax.TagSimplified,
				Text:     "Simplified Tax Invoice",
			},
		},
		// Profit Margin Scheme
		{
			Tags:       []cbc.Key{TagMarginScheme},
			Categories: []cbc.Code{tax.CategoryVAT},
			Note: &tax.Note{
				Category: tax.CategoryVAT,
				Key:      TagMarginScheme,
				Text:     "VAT charged under the Profit Margin Scheme",
			},
		},
	},
}
