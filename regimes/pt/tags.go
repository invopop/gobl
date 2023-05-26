package pt

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

// Invoice type tags
const (
	TagInvoiceReceipt cbc.Key = "invoice-receipt"
)

var invoiceTags = []*tax.KeyDefinition{
	{
		Key: TagInvoiceReceipt,
		Name: i18n.String{
			i18n.EN: "Invoice-receipt",
			i18n.PT: "Fatura-recibo",
		},
	},
	{
		Key: common.TagSimplified,
		Name: i18n.String{
			i18n.EN: "Simplified invoice",
			i18n.PT: "Fatura simplificada",
		},
	},
}
