package xrechnung

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// ExtKeyTaxRate is the key for the tax rate extension in XRechnung
const (
	ExtKeyTaxRate cbc.Key = "de-xrechnung-tax-rate"
	ExtKeyDocType cbc.Key = "de-xrechnung-doc-type"
)

var extensions = []*cbc.KeyDefinition{
	{
		Key: ExtKeyTaxRate,
		Name: i18n.String{
			i18n.EN: "Tax Rate",
			i18n.DE: "Steuersatz",
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "S",
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.DE: "Standardsteuersatz",
				},
			},
			{
				Value: "Z",
				Name: i18n.String{
					i18n.EN: "Zero rated goods",
					i18n.DE: "Güter mit Nullbewertung",
				},
			},
			{
				Value: "E",
				Name: i18n.String{
					i18n.EN: "Exempt from tax",
					i18n.DE: "von der Steuer befreit",
				},
			},
			{
				Value: "AE",
				Name: i18n.String{
					i18n.EN: "VAT Reverse Charge",
					i18n.DE: "Mehrwertsteuer Umkehrung der Steuerschuldnerschaft",
				},
			},
			{
				Value: "K",
				Name: i18n.String{
					i18n.EN: "VAT exempt for EEA intra-community supply of goods and services",
					i18n.DE: "Mehrwertsteuerbefreiung für innergemeinschaftliche Lieferungen von Gegenständen und Dienstleistungen im EWR",
				},
			},
			{
				Value: "G",
				Name: i18n.String{
					i18n.EN: "Free export item, tax not charged",
					i18n.DE: "Kostenlose Ausfuhrsendung, ohne Steuer",
				},
			},
			{
				Value: "O",
				Name: i18n.String{
					i18n.EN: "Services outside scope of tax",
					i18n.DE: "Dienstleistungen, die nicht unter die Steuer fallen",
				},
			},
			{
				Value: "L",
				Name: i18n.String{
					i18n.EN: "Canary Islands general indirect tax",
					i18n.DE: "Allgemeine indirekte Steuer der Kanarischen Inseln",
				},
			},
			{
				Value: "M",
				Name: i18n.String{
					i18n.EN: "Tax for production, services and importation in Ceuta and Melilla",
					i18n.DE: "Steuer auf Produktion, Dienstleistungen und Importe in Ceuta und Melilla",
				},
			},
		},
	},
	{
		Key: ExtKeyDocType,
		Name: i18n.String{
			i18n.EN: "Document Type",
			i18n.DE: "Dokumentenart",
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: string(invoiceTypeSelfBilled),
				Name: i18n.String{
					i18n.EN: "Self-Billed Invoice",
					i18n.DE: "Gutschrift",
				},
			},
			{
				Value: string(invoiceTypePartial),
				Name: i18n.String{
					i18n.EN: "Partial Invoice",
					i18n.DE: "Teilrechnung",
				},
			},
			{
				Value: string(invoiceTypePartialConstruction),
				Name: i18n.String{
					i18n.EN: "Partial Construction Invoice",
					i18n.DE: "Teilrechnung für Bauleistungen",
				},
			},
			{
				Value: string(invoiceTypePartialFinalConstruction),
				Name: i18n.String{
					i18n.EN: "Partial Final Construction Invoice",
					i18n.DE: "Schlussrechnung für Bauleistungen mit Teilrechnungen",
				},
			},
			{
				Value: string(invoiceTypeFinalConstruction),
				Name: i18n.String{
					i18n.EN: "Final Construction Invoice",
					i18n.DE: "Schlussrechnung für Bauleistungen",
				},
			},
		},
	},
}
