package xrechnung

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

var (
	ExtKeyPostCode cbc.Key = "de-xrechnung-post-code"
	ExtKeyTaxRate  cbc.Key = "de-xrechnung-tax-rate"
)

var extensions = []*cbc.KeyDefinition{
	{
		Key: ExtKeyPostCode,
		Name: i18n.String{
			i18n.EN: "Post Code",
			i18n.DE: "Postleitzahl",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
			Post code of a supplier or customer to use instead of an address. Example: "55000".
		`),
			i18n.DE: here.Doc(`
			Postleitzahl eines Lieferanten oder Kunden, die anstelle einer Adresse verwendet wird. Beispiel: "55000".
		`),
		},
		Pattern: "^[0-9]{5}$",
	},
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
}
