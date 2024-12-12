package osa

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Special Values to be used inside rates
const (
	ExtKeyExemptionValue = "hu-exemption-Value"
)

var extensions = []*cbc.KeyDefinition{
	{
		Key: ExtKeyExemptionValue,
		Name: i18n.String{
			i18n.EN: "Tax exemption reason Value",
			i18n.HU: "Adómentesség okának kódja",
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "AAM",
				Name: i18n.String{
					i18n.EN: "Personal tax exemption: Chapter XIII of the VAT Act",
					i18n.HU: "Személyi adómentesség: ÁFA törvény XIII. fejezete",
				},
			},
			{
				Value: "TAM",
				Name: i18n.String{
					i18n.EN: "Public interest: Section 85 and 86 of the VAT Act",
					i18n.HU: "Közérdekűség: ÁFA törvény 85. és 86. §",
				},
			},
			{
				Value: "KBAET",
				Name: i18n.String{
					i18n.EN: "Intra-Community supply (no new means of transport): Section 89 of the VAT Act",
					i18n.HU: "Közösségi beszerzés (nem új közlekedési eszköz): ÁFA törvény 89. §",
				},
			},
			{
				Value: "KBAUK",
				Name: i18n.String{
					i18n.EN: "Intra-Community supply (new means of transport): Section 89 of the VAT Act",
					i18n.HU: "Közösségi beszerzés (új közlekedési eszköz): ÁFA törvény 89. §",
				},
			},
			{
				Value: "EAM",
				Name: i18n.String{
					i18n.EN: "Export to non-EU countries: Sections 98 to 109 of the VAT Act",
					i18n.HU: "Export az EU-n kívüli országokba: ÁFA törvény 98-109. §",
				},
			},
			{
				Value: "NAM",
				Name: i18n.String{
					i18n.EN: "Other international transaction: Sections 110 to 118 of the VAT Act",
					i18n.HU: "Egyéb nemzetközi ügylet: ÁFA törvény 110-118. §",
				},
			},
			{
				Value: "ATK",
				Name: i18n.String{
					i18n.EN: "Outside Scope of VAT act: Sections 2 and 3 of the VAT Act",
					i18n.HU: "ÁFA törvény hatálya alól mentesített: ÁFA törvény 2. és 3. §",
				},
			},
			{
				Value: "EUFAD37",
				Name: i18n.String{
					i18n.EN: "Reverse charge in another member state: Section 37 of the VAT Act",
					i18n.HU: "Fordított adózás más tagállamban: ÁFA törvény 37. §",
				},
			},
			{
				Value: "EUFADE",
				Name: i18n.String{
					i18n.EN: "Reverse charge in another member state: Not subject to section 37 of the VAT Act",
					i18n.HU: "Fordított adózás más tagállamban: Nem tartozik az ÁFA törvény 37. § hatálya alá",
				},
			},
			{
				Value: "EUE",
				Name: i18n.String{
					i18n.EN: "Non-reverse charge in another member state",
					i18n.HU: "Nem fordított adózás más tagállamban",
				},
			},
			{
				Value: "HO",
				Name: i18n.String{
					i18n.EN: "Transaction in a 3rd country",
					i18n.HU: "Ügylet harmadik országban",
				},
			},
			{
				Value: "UNKNOWN",
				Name: i18n.String{
					i18n.EN: "It can be used for modifying or cancelling invoices or if unknown",
					i18n.HU: "Számla módosítására, törlésére vagy ismeretlen esetén használható",
				},
			},
		},
	},
}
