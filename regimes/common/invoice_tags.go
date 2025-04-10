package common

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*cbc.Definition{
		// Simplified invoices are issued when the complete fiscal details of
		// a customer are not available.
		{
			Key: tax.TagSimplified,
			Name: i18n.String{
				i18n.EN: "Simplified Invoice",
				i18n.ES: "Factura Simplificada",
				i18n.IT: "Fattura Semplificata",
				i18n.DE: "Vereinfachte Rechnung",
				i18n.NO: "Forenklet faktura",
			},
			Desc: i18n.String{
				i18n.EN: "Used for B2C transactions when the client details are not available, check with local authorities for limits.",
				i18n.ES: "Usado para transacciones B2C cuando los detalles del cliente no están disponibles, consulte con las autoridades locales para los límites.",
				i18n.IT: "Utilizzato per le transazioni B2C quando i dettagli del cliente non sono disponibili, controllare con le autorità locali per i limiti.",
				i18n.DE: "Wird für B2C-Transaktionen verwendet, wenn die Kundendaten nicht verfügbar sind. Bitte wenden Sie sich an die örtlichen Behörden, um die Grenzwerte zu ermitteln.",
				i18n.NO: "Brukt for B2C-transaksjoner når kundeopplysningene ikke er tilgjengelige, sjekk med lokale myndigheter for grenser.",
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
				i18n.DE: "Umkehr der Steuerschuld",
				i18n.NO: "Omvendt avgiftsplikt",
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
				i18n.DE: "Rechnung durch den Leistungsempfänger",
				i18n.NO: "Selv-fakturert",
			},
		},

		// Customer rates (mainly for digital goods inside EU)
		{
			Key: tax.TagCustomerRates,
			Name: i18n.String{
				i18n.EN: "Customer rates",
				i18n.ES: "Tarifas aplicables al destinatario",
				i18n.IT: "Aliquote applicabili al destinatario",
				i18n.DE: "Kundensätze",
				i18n.NO: "Kundensatser",
			},
		},

		// Partial invoice document, implying that this is only a first part
		// and a final invoice for the remaining amount will be made later.
		// A few regimes use this tag to classify invoices, notably Italy.
		{
			Key: tax.TagPartial,
			Name: i18n.String{
				i18n.EN: "Partial",
				i18n.ES: "Parcial",
				i18n.IT: "Parziale",
				i18n.DE: "Teilweise",
				i18n.NO: "Delvis",
			},
		},
	},
}

// InvoiceTags returns a base tag set for invoices.
func InvoiceTags() *tax.TagSet {
	return invoiceTags
}
