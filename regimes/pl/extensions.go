package pl

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Regime extension codes for local electronic formats.
const (
	ExtKeyKSeFVATZero    = "pl-ksef-vat-zero"
	ExtKeyKSeFVATSpecial = "pl-ksef-vat-special"
	ExtKeyKSEFCorrection = "pl-ksef-correction"
)

var extensionKeys = []*cbc.KeyDefinition{
	{
		Key: ExtKeyKSeFVATSpecial,
		Name: i18n.String{
			i18n.EN: "Special VAT Extensions for KSeF",
			i18n.PL: "Rozszerzenia specjalne dla KSeF",
		},
		Keys: []*cbc.KeyDefinition{
			{
				Key: "taxi",
				Name: i18n.String{
					i18n.EN: "Taxi Rate",
					i18n.PL: "Ryczałt dla taksówek",
				},
				Desc: i18n.String{
					i18n.EN: "Special flat rate for taxi drivers.",
					i18n.PL: "Specjalna stawka ryczałtu dla taksówkarzy.",
				},
			},
		},
	},
	{
		Key: ExtKeyKSeFVATZero,
		Name: i18n.String{
			i18n.EN: "Zero VAT Extensions for KSeF",
		},
		Keys: []*cbc.KeyDefinition{
			{
				Key: "wdt",
				Name: i18n.String{
					i18n.EN: "WDT",
					i18n.PL: "WDT",
				},
				// TODO: description required
			},
			{
				Key: "domestic",
				Name: i18n.String{
					i18n.EN: "Domestic",
					i18n.PL: "Krajowy",
				},
				// TODO: description required
			},
			{
				Key: "export",
				Name: i18n.String{
					i18n.EN: "Export",
					i18n.PL: "Eksport",
				},
				// TODO: description required
			},
		},
	},
	{
		Key: ExtKeyKSEFCorrection,
		Name: i18n.String{
			i18n.EN: "Zero VAT Extensions for KSeF",
		},
		Codes: []*cbc.CodeDefinition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Original",
					i18n.PL: "Pierwotna",
				},
				Desc: i18n.String{
					i18n.EN: "Correction has legal consequences in the date of the original invoice",
					i18n.PL: "Faktura skutkująca w dacie ujęcia faktury pierwotnej",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Correction",
					i18n.PL: "Korygująca",
				},
				Desc: i18n.String{
					i18n.EN: "Correction has legal consequences in the date of the correction invoice",
					i18n.PL: "Faktura skutkująca w dacie ujęcia faktury korygującej",
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "Other",
					i18n.PL: "Inna",
				},
				Desc: i18n.String{
					i18n.EN: "Correction has legal consequences in another date or the dates are different for different position on the invoice",
					i18n.PL: "Faktura skutkująca w innej dacie. W tym gdy dla różnych pozycji faktury korygującej data jest różna.",
				},
			},
		},
	},
}
