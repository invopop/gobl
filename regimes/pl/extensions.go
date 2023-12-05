package pl

import (
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Regime extension codes for local electronic formats.
const (
	ExtKeyKSeFVATZero    = "pl-ksef-vat-zero"
	ExtKeyKSeFVATSpecial = "pl-ksef-vat-special"
)

var extensionKeys = []*tax.KeyDefinition{
	{
		Key: ExtKeyKSeFVATSpecial,
		Name: i18n.String{
			i18n.EN: "Special VAT Extensions for KSeF",
			i18n.PL: "Rozszerzenia specjalne dla KSeF",
		},
		Keys: []*tax.KeyDefinition{
			{
				Key: "taxi",
				Name: i18n.String{
					i18n.EN: "Taxi Rate",
					i18n.PL: "Ryczałt dla taksówek",
				},
			},
		},
	},
	{
		Key: ExtKeyKSeFVATZero,
		Name: i18n.String{
			i18n.EN: "Zero VAT Extensions for KSeF",
		},
		Keys: []*tax.KeyDefinition{
			{
				Key: "wdt",
				Name: i18n.String{
					i18n.EN: "WDT",
					i18n.PL: "WDT",
				},
			},
			{
				Key: "domestic",
				Name: i18n.String{
					i18n.EN: "Domestic",
					i18n.PL: "Krajowy",
				},
			},
			{
				Key: "export",
				Name: i18n.String{
					i18n.EN: "Export",
					i18n.PL: "Eksport",
				},
			},
		},
	},
}
