package pt

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Special codes to be used inside rates.
const (
	ExtKeyRegion = "pt-region"
)

var extensionKeys = []*cbc.KeyDefinition{
	{
		Key: ExtKeyRegion,
		Name: i18n.String{
			i18n.EN: "Region Code",
			i18n.PT: "Código da Região",
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "PT",
				Name: i18n.String{
					i18n.EN: "Mainland Portugal",
					i18n.PT: "Portugal Continental",
				},
			},
			{
				Value: "PT-AC",
				Name: i18n.String{
					i18n.EN: "Azores",
					i18n.PT: "Açores",
				},
			},
			{
				Value: "PT-MA",
				Name: i18n.String{
					i18n.EN: "Madeira",
					i18n.PT: "Madeira",
				},
			},
		},
	},
}
