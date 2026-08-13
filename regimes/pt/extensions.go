package pt

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/pkg/here"
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
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				SAF-T's ~TaxCountryRegion~ (País ou região do imposto) specifies the region of taxation
				(Portugal mainland, Açores, Madeira or any ISO country) in a Portuguese invoice. Each
				region has their own tax rates which can be determined automatically.

				To set the specific a region different to Portugal mainland, the ~pt-region~ extension of
				each line's VAT tax should be set to one of the following values:

				| Code  | Description                                         |
				| ----- | --------------------------------------------------- |
				| PT    | Mainland Portugal (default)                         |
				| PT-AC | Açores                                              |
				| PT-MA | Madeira                                             |
				| ...   | Any ISO country code (e.g. ES, FR, DE, etc.)        |

				For example:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					// ...
					"lines": [
						{
							// ...
							"item": {
								"name": "Some service",
								"price": "25.00"
							},
							"tax": [
								{
										"cat": "VAT",
										"rate": "exempt",
										"ext": {
											"pt-region": "PT-AC",
											// ...
										}
								}
							]
						}
					]
				}
				~~~
			`),
		},
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
