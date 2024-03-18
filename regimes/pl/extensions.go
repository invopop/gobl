package pl

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Regime extension codes for local electronic formats.
const (
	ExtKeyKSeFVATZero    = "pl-ksef-vat-zero"
	ExtKeyKSeFVATSpecial = "pl-ksef-vat-special"
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
}
