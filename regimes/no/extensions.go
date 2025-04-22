package no

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Constants for extension keys
const (
	ExtKeyMarginScheme  cbc.Key = "no-margin-scheme"
	ExtKeyExemptionCode cbc.Key = "no-exemption-code"
	ExtKeyVOEC          cbc.Key = "no-voec"
)

// Tags for special cases
const (
	TagBooks        cbc.Key = "books"
	TagSecondHand   cbc.Key = "second-hand"
	TagArtworks     cbc.Key = "artworks"
	TagCollectibles cbc.Key = "collectibles"
	TagAntiques     cbc.Key = "antiques"
	TagECommerce    cbc.Key = "ecommerce"
)

// Define extensions
var extensionKeys = []*cbc.Definition{
	{
		Key: ExtKeyExemptionCode,
		Name: i18n.String{
			i18n.EN: "Norwegian VAT Exemption Code",
			i18n.NO: "Norsk kode for MVA-fritak",
		},
		Values: []*cbc.Definition{
			{
				Code: "E1",
				Name: i18n.String{
					i18n.EN: "Financial Services - § 3-6 MVAL",
					i18n.NO: "Finansielle tjenester - § 3-6 MVAL",
				},
			},
			{
				Code: "E2",
				Name: i18n.String{
					i18n.EN: "Insurance Services - § 3-10 MVAL",
					i18n.NO: "Forsikringstjenester - § 3-10 MVAL",
				},
			},
			{
				Code: "E3",
				Name: i18n.String{
					i18n.EN: "Books and Periodicals - § 6-4 MVAL",
					i18n.NO: "Bøker og tidsskrifter - § 6-4 MVAL",
				},
			},
			{
				Code: "M1",
				Name: i18n.String{
					i18n.EN: "Margin Scheme - Second-hand Goods - Chapter Va MVAL",
					i18n.NO: "Avansesystem - Brukte varer - Kapittel Va MVAL",
				},
			},
			{
				Code: "M2",
				Name: i18n.String{
					i18n.EN: "Margin Scheme - Works of Art - Chapter Va MVAL",
					i18n.NO: "Avansesystem - Kunstgjenstander - Kapittel Va MVAL",
				},
			},
			{
				Code: "V1",
				Name: i18n.String{
					i18n.EN: "VOEC Scheme - B2C E-commerce - § 3-30 MVAL",
					i18n.NO: "VOEC-ordning - B2C E-handel - § 3-30 MVAL",
				},
			},
		},
	},
	{
		Key: ExtKeyMarginScheme,
		Name: i18n.String{
			i18n.EN: "Norwegian Margin Scheme Type",
			i18n.NO: "Type norsk avansesystem",
		},
		Values: []*cbc.Definition{
			{
				Code: "second-hand",
				Name: i18n.String{
					i18n.EN: "Second-hand Goods",
					i18n.NO: "Brukte varer",
				},
			},
			{
				Code: "artworks",
				Name: i18n.String{
					i18n.EN: "Works of Art",
					i18n.NO: "Kunstgjenstander",
				},
			},
		},
	},
	{
		Key: ExtKeyVOEC,
		Name: i18n.String{
			i18n.EN: "VOEC Scheme Details",
			i18n.NO: "VOEC-ordning detaljer",
		},
		Values: []*cbc.Definition{
			{
				Code: "registered",
				Name: i18n.String{
					i18n.EN: "Registered in VOEC Scheme",
					i18n.NO: "Registrert i VOEC-ordning",
				},
			},
		},
	},
}
