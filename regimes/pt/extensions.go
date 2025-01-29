package pt

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
)

// Special codes to be used inside rates.
const (
	ExtKeyRegion = "pt-region"
)

// Region codes
const (
	RegionMainland = "PT"
	RegionAzores   = "PT-AC"
	RegionMadeira  = "PT-MA"
)

var extensionKeys = []*cbc.Definition{
	{
		Key: ExtKeyRegion,
		Name: i18n.String{
			i18n.EN: "Region Code",
			i18n.PT: "Código da Região",
		},
		Values: regionDefs(),
	},
}

func regionDefs() []*cbc.Definition {
	regs := []*cbc.Definition{
		{
			Code: RegionMainland,
			Name: i18n.String{
				i18n.EN: "Mainland Portugal",
				i18n.PT: "Portugal Continental",
			},
		},
		{
			Code: RegionAzores,
			Name: i18n.String{
				i18n.EN: "Azores",
				i18n.PT: "Açores",
			},
		},
		{
			Code: RegionMadeira,
			Name: i18n.String{
				i18n.EN: "Madeira",
				i18n.PT: "Madeira",
			},
		},
	}
	for _, rd := range l10n.Countries().ISO() {
		if rd.Code == l10n.PT {
			continue
		}
		regs = append(regs, &cbc.Definition{
			Code: cbc.Code(rd.Code),
			Name: i18n.String{
				i18n.EN: rd.Name,
			},
		})
	}
	return regs
}
