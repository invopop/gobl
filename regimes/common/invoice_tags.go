package common

import (
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

var invoiceTags = []*tax.KeyDefinition{
	// Simplified invoices are issued when the complete fiscal details of
	// a customer are not available.
	{
		Key: tax.TagSimplified,
		Name: i18n.String{
			i18n.EN: "Simplified Invoice",
			i18n.ES: "Factura Simplificada",
			i18n.IT: "Fattura Semplificata",
		},
		Desc: i18n.String{
			i18n.EN: "Used for B2C transactions when the client details are not available, check with local authorities for limits.",
			i18n.ES: "Usado para transacciones B2C cuando los detalles del cliente no están disponibles, consulte con las autoridades locales para los límites.",
			i18n.IT: "Utilizzato per le transazioni B2C quando i dettagli del cliente non sono disponibili, controllare con le autorità locali per i limiti.",
		},
	},

	// Reverse Charge mechanism is used when the supplier is not
	// required to charge VAT on the invoice and the customer is
	// responsible for paying the VAT to the tax authorities.
	{
		Key: tax.TagReverseCharge,
		Name: i18n.String{
			i18n.EN: "Reverse Charge",
			i18n.ES: "Inversión del Sujeto Pasivo",
			i18n.IT: "Inversione del soggetto passivo",
		},
	},

	// Self-billed invoices are issued by the customer instead of the
	// supplier. This is usually done when the customer is a large
	// company and the supplier is a small company.
	{
		Key: tax.TagSelfBilled,
		Name: i18n.String{
			i18n.EN: "Self-billed",
			i18n.ES: "Facturación por el destinatario",
			i18n.IT: "Autofattura",
		},
	},

	// Customer rates (mainly for digital goods inside EU)
	{
		Key: tax.TagCustomerRates,
		Name: i18n.String{
			i18n.EN: "Customer rates",
			i18n.ES: "Tarifas aplicables al destinatario",
		},
	},
}

// InvoiceTags returns a list of common invoice tag key
// definitions.
func InvoiceTags() []*tax.KeyDefinition {
	return invoiceTags
}

// InvoiceTagsWith appends the list of provided key definitions
// to the base list of tags and returns a new array.
func InvoiceTagsWith(tags []*tax.KeyDefinition) []*tax.KeyDefinition {
	return append(InvoiceTags(), tags...)
}
