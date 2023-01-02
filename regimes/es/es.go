package es

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegime(New())
}

// Local tax category definitions which are not considered standard.
const (
	TaxCategoryIRPF cbc.Code = "IRPF"
	TaxCategoryIGIC cbc.Code = "IGIC"
	TaxCategoryIPSI cbc.Code = "IPSI"
)

// Specific tax rate codes.
const (
	// IRPF non-standard Rates (usually for self-employed)
	TaxRatePro                cbc.Key = "pro"                 // Professional Services
	TaxRateProStart           cbc.Key = "pro-start"           // Professionals, first 2 years
	TaxRateModules            cbc.Key = "modules"             // Module system
	TaxRateAgriculture        cbc.Key = "agriculture"         // Agricultural
	TaxRateAgricultureSpecial cbc.Key = "agriculture-special" // Agricultural special
	TaxRateCapital            cbc.Key = "capital"             // Rental or Interest

	// Special tax rate surcharge extension
	TaxRateEquivalence cbc.Key = "eqs"
)

// Scheme key definitions
const (
	SchemeSimplified      cbc.Key = "simplified"
	SchemeCustomerIssued  cbc.Key = "customer-issued"
	SchemeTravelAgency    cbc.Key = "travel-agency"
	SchemeSecondHandGoods cbc.Key = "second-hand-goods"
	SchemeArt             cbc.Key = "art"
	SchemeAntiques        cbc.Key = "antiques"
	SchemeCashBasis       cbc.Key = "cash-basis"
)

// Official stamps or codes validated by government agencies
const (
	// TicketBAI (Basque Country) codes used for stamps.
	StampProviderTBAICode cbc.Key = "tbai-code"
	StampProviderTBAIQR   cbc.Key = "tbai-qr"
)

// Inbox key and role definitions
const (
	InboxKeyFACE cbc.Key = "face"

	// Main roles defined in FACE
	InboxRoleFiscal    cbc.Key = "fiscal"    // Fiscal / 01
	InboxRoleRecipient cbc.Key = "recipient" // Receptor / 02
	InboxRolePayer     cbc.Key = "payer"     // Pagador / 03
	InboxRoleCustomer  cbc.Key = "customer"  // Comprador / 04

)

// Zone code definitions for Spain
const (
	ZoneVI l10n.Code = "VI" // (01) Álava
	ZoneAB l10n.Code = "AB" // (02) Albacete
	ZoneA  l10n.Code = "A"  // (03) Alicante
	ZoneAL l10n.Code = "AL" // (04) Almería
	ZoneAV l10n.Code = "AV" // (05) Ávila
	ZoneBA l10n.Code = "BA" // (06) Badajoz
	ZonePM l10n.Code = "PM" // (07) Baleares
	ZoneIB l10n.Code = "IB" // (07) Baleares
	ZoneB  l10n.Code = "B"  // (08) Barcelona
	ZoneBU l10n.Code = "BU" // (09) Burgos
	ZoneCC l10n.Code = "CC" // (10) Cáceres
	ZoneCA l10n.Code = "CA" // (11) Cádiz
	ZoneCS l10n.Code = "CS" // (12) Castellon
	ZoneCR l10n.Code = "CR" // (13) Ciudad Real
	ZoneCO l10n.Code = "CO" // (14) Cordoba
	ZoneC  l10n.Code = "C"  // (15) La Coruña
	ZoneCU l10n.Code = "CU" // (16) Cuenca
	ZoneGE l10n.Code = "GE" // (17) Gerona
	ZoneGI l10n.Code = "GI" // (17) Girona
	ZoneGR l10n.Code = "GR" // (18) Granada
	ZoneGU l10n.Code = "GU" // (19) Guadalajara
	ZoneSS l10n.Code = "SS" // (20) Guipúzcoa
	ZoneH  l10n.Code = "H"  // (21) Huelva
	ZoneHU l10n.Code = "HU" // (22) Huesca
	ZoneJ  l10n.Code = "J"  // (23) Jaén
	ZoneLE l10n.Code = "LE" // (24) León
	ZoneL  l10n.Code = "L"  // (25) Lérida / Lleida
	ZoneLO l10n.Code = "LO" // (26) La Rioja
	ZoneLU l10n.Code = "LU" // (27) Lugo
	ZoneM  l10n.Code = "M"  // (28) Madrid
	ZoneMA l10n.Code = "MA" // (29) Málaga
	ZoneMU l10n.Code = "MU" // (30) Murcia
	ZoneNA l10n.Code = "NA" // (31) Navarra
	ZoneOR l10n.Code = "OR" // (32) Orense
	ZoneOU l10n.Code = "OU" // (32) Orense
	ZoneO  l10n.Code = "O"  // (33) Asturias
	ZoneP  l10n.Code = "P"  // (34) Palencia
	ZoneGC l10n.Code = "GC" // (35) Las Palmas
	ZonePO l10n.Code = "PO" // (36) Pontevedra
	ZoneSA l10n.Code = "SA" // (37) Salamanca
	ZoneTF l10n.Code = "TF" // (38) Santa Cruz de Tenerife
	ZoneS  l10n.Code = "S"  // (39) Cantabria
	ZoneSG l10n.Code = "SG" // (40) Segovia
	ZoneSE l10n.Code = "SE" // (41) Sevilla
	ZoneSO l10n.Code = "SO" // (42) Soria
	ZoneT  l10n.Code = "T"  // (43) Tarragona
	ZoneTE l10n.Code = "TE" // (44) Teruel
	ZoneTO l10n.Code = "TO" // (45) Toledo
	ZoneV  l10n.Code = "V"  // (46) Valencia
	ZoneVA l10n.Code = "VA" // (47) Valladolid
	ZoneBI l10n.Code = "BI" // (48) Vizcaya
	ZoneZA l10n.Code = "ZA" // (49) Zamora
	ZoneZ  l10n.Code = "Z"  // (50) Zaragoza
	ZoneCE l10n.Code = "CE" // (51) Ceuta
	ZoneML l10n.Code = "ML" // (52) Melilla
)

// Custom keys used typically in meta information.
const (
	KeyAddressCode cbc.Key = "post"
)

// New provides the Spanish tax regime definition
func New() *tax.Regime {
	return &tax.Regime{
		Country:  l10n.ES,
		Currency: "EUR",
		Name: i18n.String{
			i18n.EN: "Spain",
			i18n.ES: "España",
		},
		Validator:  Validate,
		Calculator: Calculate,
		Zones: []tax.Zone{
			{
				Code: ZoneVI,
				Name: i18n.String{i18n.ES: "Ávila"},
				Meta: cbc.Meta{KeyAddressCode: "01"},
			},
			{
				Code: ZoneAB,
				Name: i18n.String{i18n.ES: "Albacete"},
				Meta: cbc.Meta{KeyAddressCode: "02"},
			},
			{
				Code: ZoneA,
				Name: i18n.String{i18n.ES: "Alicante"},
				Meta: cbc.Meta{KeyAddressCode: "03"},
			},
			{
				Code: ZoneAL,
				Name: i18n.String{i18n.ES: "Almería"},
				Meta: cbc.Meta{KeyAddressCode: "04"},
			},
			{
				Code: ZoneAV,
				Name: i18n.String{i18n.ES: "Ávila"},
				Meta: cbc.Meta{KeyAddressCode: "05"},
			},
			{
				Code: ZoneBA,
				Name: i18n.String{i18n.ES: "Badajoz"},
				Meta: cbc.Meta{KeyAddressCode: "06"},
			},
			{
				Code: ZonePM,
				Name: i18n.String{i18n.ES: "Baleares"},
				Meta: cbc.Meta{KeyAddressCode: "07"},
			},
			{
				Code: ZoneIB,
				Name: i18n.String{i18n.ES: "Baleares"},
				Meta: cbc.Meta{KeyAddressCode: "07"},
			},
			{
				Code: ZoneB,
				Name: i18n.String{i18n.ES: "Barcelona"},
				Meta: cbc.Meta{KeyAddressCode: "08"},
			},
			{
				Code: ZoneBU,
				Name: i18n.String{i18n.ES: "Burgos"},
				Meta: cbc.Meta{KeyAddressCode: "09"},
			},
			{
				Code: ZoneCC,
				Name: i18n.String{i18n.ES: "Cáceres"},
				Meta: cbc.Meta{KeyAddressCode: "10"},
			},
			{
				Code: ZoneCA,
				Name: i18n.String{i18n.ES: "Cadiz"},
				Meta: cbc.Meta{KeyAddressCode: "11"},
			},
			{
				Code: ZoneCS,
				Name: i18n.String{i18n.ES: "Castellón"},
				Meta: cbc.Meta{KeyAddressCode: "12"},
			},
			{
				Code: ZoneCR,
				Name: i18n.String{i18n.ES: "Ciudad Real"},
				Meta: cbc.Meta{KeyAddressCode: "13"},
			},
			{
				Code: ZoneCO,
				Name: i18n.String{i18n.ES: "Cordoba"},
				Meta: cbc.Meta{KeyAddressCode: "14"},
			},
			{
				Code: ZoneC,
				Name: i18n.String{i18n.ES: "La Coruña"},
				Meta: cbc.Meta{KeyAddressCode: "15"},
			},
			{
				Code: ZoneCU,
				Name: i18n.String{i18n.ES: "Cuenca"},
				Meta: cbc.Meta{KeyAddressCode: "16"},
			},
			{
				Code: ZoneGE,
				Name: i18n.String{i18n.ES: "Gerona"},
				Meta: cbc.Meta{KeyAddressCode: "17"},
			},
			{
				Code: ZoneGI,
				Name: i18n.String{i18n.ES: "Girona"},
				Meta: cbc.Meta{KeyAddressCode: "17"},
			},
			{
				Code: ZoneGR,
				Name: i18n.String{i18n.ES: "Granada"},
				Meta: cbc.Meta{KeyAddressCode: "18"},
			},
			{
				Code: ZoneGU,
				Name: i18n.String{i18n.ES: "Guadalajara"},
				Meta: cbc.Meta{KeyAddressCode: "19"},
			},
			{
				Code: ZoneSS,
				Name: i18n.String{i18n.ES: "Guipúzcoa"},
				Meta: cbc.Meta{KeyAddressCode: "20"},
			},
			{
				Code: ZoneH,
				Name: i18n.String{i18n.ES: "Huelva"},
				Meta: cbc.Meta{KeyAddressCode: "21"},
			},
			{
				Code: ZoneHU,
				Name: i18n.String{i18n.ES: "Huesca"},
				Meta: cbc.Meta{KeyAddressCode: "22"},
			},
			{
				Code: ZoneJ,
				Name: i18n.String{i18n.ES: "Jaén"},
				Meta: cbc.Meta{KeyAddressCode: "23"},
			},
			{
				Code: ZoneLE,
				Name: i18n.String{i18n.ES: "León"},
				Meta: cbc.Meta{KeyAddressCode: "24"},
			},
			{
				Code: ZoneL,
				Name: i18n.String{i18n.ES: "Lérida / Lleida"},
				Meta: cbc.Meta{KeyAddressCode: "25"},
			},
			{
				Code: ZoneLO,
				Name: i18n.String{i18n.ES: "La Rioja"},
				Meta: cbc.Meta{KeyAddressCode: "26"},
			},
			{
				Code: ZoneLU,
				Name: i18n.String{i18n.ES: "Lugo"},
				Meta: cbc.Meta{KeyAddressCode: "27"},
			},
			{
				Code: ZoneM,
				Name: i18n.String{i18n.ES: "Madrid"},
				Meta: cbc.Meta{KeyAddressCode: "28"},
			},
			{
				Code: ZoneMA,
				Name: i18n.String{i18n.ES: "Málaga"},
				Meta: cbc.Meta{KeyAddressCode: "29"},
			},
			{
				Code: ZoneMU,
				Name: i18n.String{i18n.ES: "Murcia"},
				Meta: cbc.Meta{KeyAddressCode: "30"},
			},
			{
				Code: ZoneNA,
				Name: i18n.String{i18n.ES: "Navarra"},
				Meta: cbc.Meta{KeyAddressCode: "31"},
			},
			{
				Code: ZoneOR,
				Name: i18n.String{i18n.ES: "Orense"},
				Meta: cbc.Meta{KeyAddressCode: "32"},
			},
			{
				Code: ZoneOU,
				Name: i18n.String{i18n.ES: "Orense"},
				Meta: cbc.Meta{KeyAddressCode: "32"},
			},
			{
				Code: ZoneO,
				Name: i18n.String{i18n.ES: "Asturias"},
				Meta: cbc.Meta{KeyAddressCode: "33"},
			},
			{
				Code: ZoneP,
				Name: i18n.String{i18n.ES: "Palencia"},
				Meta: cbc.Meta{KeyAddressCode: "34"},
			},
			{
				Code: ZoneGC,
				Name: i18n.String{i18n.ES: "Las Palmas"},
				Meta: cbc.Meta{KeyAddressCode: "35"},
			},
			{
				Code: ZonePO,
				Name: i18n.String{i18n.ES: "Pontevedra"},
				Meta: cbc.Meta{KeyAddressCode: "36"},
			},
			{
				Code: ZoneSA,
				Name: i18n.String{i18n.ES: "Salamanca"},
				Meta: cbc.Meta{KeyAddressCode: "37"},
			},
			{
				Code: ZoneTF,
				Name: i18n.String{i18n.ES: "Santa Cruz de Tenerife"},
				Meta: cbc.Meta{KeyAddressCode: "38"},
			},
			{
				Code: ZoneS,
				Name: i18n.String{i18n.ES: "Cantabria"},
				Meta: cbc.Meta{KeyAddressCode: "39"},
			},
			{
				Code: ZoneSG,
				Name: i18n.String{i18n.ES: "Segovia"},
				Meta: cbc.Meta{KeyAddressCode: "40"},
			},
			{
				Code: ZoneSE,
				Name: i18n.String{i18n.ES: "Sevilla"},
				Meta: cbc.Meta{KeyAddressCode: "41"},
			},
			{
				Code: ZoneSO,
				Name: i18n.String{i18n.ES: "Soria"},
				Meta: cbc.Meta{KeyAddressCode: "42"},
			},
			{
				Code: ZoneT,
				Name: i18n.String{i18n.ES: "Tarragona"},
				Meta: cbc.Meta{KeyAddressCode: "43"},
			},
			{
				Code: ZoneTE,
				Name: i18n.String{i18n.ES: "Teruel"},
				Meta: cbc.Meta{KeyAddressCode: "44"},
			},
			{
				Code: ZoneTO,
				Name: i18n.String{i18n.ES: "Toledo"},
				Meta: cbc.Meta{KeyAddressCode: "45"},
			},
			{
				Code: ZoneV,
				Name: i18n.String{i18n.ES: "Valencia"},
				Meta: cbc.Meta{KeyAddressCode: "46"},
			},
			{
				Code: ZoneVA,
				Name: i18n.String{i18n.ES: "Valladolid"},
				Meta: cbc.Meta{KeyAddressCode: "47"},
			},
			{
				Code: ZoneBI,
				Name: i18n.String{i18n.ES: "Vizcaya"},
				Meta: cbc.Meta{KeyAddressCode: "48"},
			},
			{
				Code: ZoneZA,
				Name: i18n.String{i18n.ES: "Zamora"},
				Meta: cbc.Meta{KeyAddressCode: "49"},
			},
			{
				Code: ZoneZ,
				Name: i18n.String{i18n.ES: "Zaragoza"},
				Meta: cbc.Meta{KeyAddressCode: "50"},
			},
			{
				Code: ZoneCE,
				Name: i18n.String{i18n.ES: "Ceuta"},
				Meta: cbc.Meta{KeyAddressCode: "51"},
			},
			{
				Code: ZoneML,
				Name: i18n.String{i18n.ES: "Melilla"},
				Meta: cbc.Meta{KeyAddressCode: "52"},
			},
		},
		Schemes: []*tax.Scheme{
			// Reverse Charge Scheme
			{
				Key: common.SchemeReverseCharge,
				Name: i18n.String{
					i18n.EN: "Reverse Charge",
					i18n.ES: "Inversión del sujeto pasivo",
				},
				Categories: []cbc.Code{
					common.TaxCategoryVAT,
				},
				Note: &cbc.Note{
					Key:  cbc.NoteKeyLegal,
					Src:  string(common.SchemeReverseCharge),
					Text: "Reverse Charge / Inversión del sujeto pasivo.",
				},
			},
			// Customer Rates Scheme (digital goods)
			{
				Key: common.SchemeCustomerRates,
				Name: i18n.String{
					i18n.EN: "Customer Country Rates",
					i18n.ES: "Tasas del País del Cliente",
				},
				Description: i18n.String{
					i18n.EN: "Use the customers country to determine tax rates.",
				},
			},
			// Simplified Regime
			{
				Key: SchemeSimplified,
				Name: i18n.String{
					i18n.EN: "Simplified tax scheme",
					i18n.ES: "Contribuyente en régimen simplificado",
				},
				Note: &cbc.Note{
					Key:  cbc.NoteKeyLegal,
					Src:  string(SchemeSimplified),
					Text: "Factura expedida por contibuyente en régimen simplificado.",
				},
			},
			// Customer issued invoices
			{
				Key: SchemeCustomerIssued,
				Name: i18n.String{
					i18n.EN: "Customer issued invoice",
					i18n.ES: "Facturación por el destinatario",
				},
				Note: &cbc.Note{
					Key:  cbc.NoteKeyLegal,
					Src:  string(SchemeCustomerIssued),
					Text: "Facturación por el destinatario.",
				},
			},
			// Travel agency
			{
				Key: SchemeTravelAgency,
				Name: i18n.String{
					i18n.EN: "Special scheme for travel agencies",
					i18n.ES: "Régimen especial de las agencias de viajes",
				},
				Note: &cbc.Note{
					Key:  cbc.NoteKeyLegal,
					Src:  string(SchemeTravelAgency),
					Text: "Régimen especial de las agencias de viajes.",
				},
			},
			// Secondhand stuff
			{
				Key: SchemeSecondHandGoods,
				Name: i18n.String{
					i18n.EN: "Special scheme for second-hand goods",
					i18n.ES: "Régimen especial de los bienes usados",
				},
				Note: &cbc.Note{
					Key:  cbc.NoteKeyLegal,
					Src:  string(SchemeSecondHandGoods),
					Text: "Régimen especial de los bienes usados.",
				},
			},
			// Art
			{
				Key: SchemeArt,
				Name: i18n.String{
					i18n.EN: "Special scheme of works of art",
					i18n.ES: "Régimen especial de los objetos de arte",
				},
				Note: &cbc.Note{
					Key:  cbc.NoteKeyLegal,
					Src:  string(SchemeArt),
					Text: "Régimen especial de los objetos de arte.",
				},
			},
			// Antiques
			{
				Key: SchemeAntiques,
				Name: i18n.String{
					i18n.EN: "Special scheme of antiques and collectables",
					i18n.ES: "Régimen especial de las antigüedades y objetos de colección",
				},
				Note: &cbc.Note{
					Key:  cbc.NoteKeyLegal,
					Src:  string(SchemeAntiques),
					Text: "Régimen especial de las antigüedades y objetos de colección.",
				},
			},
			// Special Regime of "Cash Criteria"
			{
				Key: SchemeCashBasis,
				Name: i18n.String{
					i18n.EN: "Special scheme on cash basis",
					i18n.ES: "Régimen especial del criterio de caja",
				},
				Note: &cbc.Note{
					Key:  cbc.NoteKeyLegal,
					Src:  string(SchemeCashBasis),
					Text: "Régimen especial del criterio de caja.",
				},
			},
		},
		Categories: []*tax.Category{
			//
			// VAT
			//
			{
				Code:     common.TaxCategoryVAT,
				Retained: false,
				Name: i18n.String{
					i18n.EN: "VAT",
					i18n.ES: "IVA",
				},
				Desc: i18n.String{
					i18n.EN: "Value Added Tax",
					i18n.ES: "Impuesto sobre el Valor Añadido",
				},
				Rates: []*tax.Rate{
					{
						Key: common.TaxRateZero,
						Name: i18n.String{
							i18n.EN: "Zero Rate",
							i18n.ES: "Tipo Zero",
						},
						Values: []*tax.RateValue{
							{
								Percent: num.MakePercentage(0, 3),
							},
						},
					},
					{
						Key: common.TaxRateStandard,
						Name: i18n.String{
							i18n.EN: "Standard Rate",
							i18n.ES: "Tipo General",
						},
						Values: []*tax.RateValue{
							{
								Since:   cal.NewDate(2012, 9, 1),
								Percent: num.MakePercentage(210, 3),
							},
							{
								Since:   cal.NewDate(2010, 7, 1),
								Percent: num.MakePercentage(180, 3),
							},
							{
								Since:   cal.NewDate(1995, 1, 1),
								Percent: num.MakePercentage(160, 3),
							},
							{
								Since:   cal.NewDate(1993, 1, 1),
								Percent: num.MakePercentage(150, 3),
							},
						},
					},
					{
						Key: common.TaxRateStandard.With(TaxRateEquivalence),
						Name: i18n.String{
							i18n.EN: "Standard Rate + Equivalence Surcharge",
							i18n.ES: "Tipo General + Recargo de Equivalencia",
						},
						Values: []*tax.RateValue{
							{
								Since:     cal.NewDate(2012, 9, 1),
								Percent:   num.MakePercentage(210, 3),
								Surcharge: num.NewPercentage(52, 3),
							},
							{
								Since:     cal.NewDate(2010, 7, 1),
								Percent:   num.MakePercentage(180, 3),
								Surcharge: num.NewPercentage(40, 3),
							},
						},
					},
					{
						Key: common.TaxRateReduced,
						Name: i18n.String{
							i18n.EN: "Reduced Rate",
							i18n.ES: "Tipo Reducido",
						},
						Values: []*tax.RateValue{
							{
								Since:   cal.NewDate(2012, 9, 1),
								Percent: num.MakePercentage(100, 3),
							},
							{
								Since:   cal.NewDate(2010, 7, 1),
								Percent: num.MakePercentage(80, 3),
							},
							{
								Since:   cal.NewDate(1995, 1, 1),
								Percent: num.MakePercentage(70, 3),
							},
							{
								Since:   cal.NewDate(1993, 1, 1),
								Percent: num.MakePercentage(60, 3),
							},
						},
					},
					{
						Key: common.TaxRateReduced.With(TaxRateEquivalence),
						Name: i18n.String{
							i18n.EN: "Reduced Rate + Equivalence Surcharge",
							i18n.ES: "Tipo Reducido + Recargo de Equivalencia",
						},
						Values: []*tax.RateValue{
							{
								Since:     cal.NewDate(2012, 9, 1),
								Percent:   num.MakePercentage(100, 3),
								Surcharge: num.NewPercentage(14, 3),
							},
							{
								Since:     cal.NewDate(2010, 7, 1),
								Percent:   num.MakePercentage(80, 3),
								Surcharge: num.NewPercentage(10, 3),
							},
						},
					},
					{
						Key: common.TaxRateSuperReduced,
						Name: i18n.String{
							i18n.EN: "Super-Reduced Rate",
							i18n.ES: "Tipo Superreducido",
						},
						Values: []*tax.RateValue{
							{
								Since:   cal.NewDate(1995, 1, 1),
								Percent: num.MakePercentage(40, 3),
							},
							{
								Since:   cal.NewDate(1993, 1, 1),
								Percent: num.MakePercentage(30, 3),
							},
						},
					},
					{
						Key: common.TaxRateSuperReduced.With(TaxRateEquivalence),
						Name: i18n.String{
							i18n.EN: "Super-Reduced Rate + Equivalence Surcharge",
							i18n.ES: "Tipo Superreducido + Recargo de Equivalencia",
						},
						Values: []*tax.RateValue{
							{
								Since:     cal.NewDate(1995, 1, 1),
								Percent:   num.MakePercentage(40, 3),
								Surcharge: num.NewPercentage(5, 3),
							},
						},
					},
				},
			},

			//
			// IGIC
			//
			{
				Code:     TaxCategoryIGIC,
				Retained: false,
				Name: i18n.String{
					i18n.EN: "IGIC",
					i18n.ES: "IGIC",
				},
				Desc: i18n.String{
					i18n.EN: "Canary Island General Indirect Tax",
					i18n.ES: "Impuesto General Indirecto Canario",
				},
				// This is a subset of the possible rates.
				Rates: []*tax.Rate{
					{
						Key: common.TaxRateZero,
						Name: i18n.String{
							i18n.EN: "Zero Rate",
							i18n.ES: "Tipo Zero",
						},
						Values: []*tax.RateValue{
							{
								Percent: num.MakePercentage(0, 3),
							},
						},
					},
					{
						Key: common.TaxRateStandard,
						Name: i18n.String{
							i18n.EN: "Standard Rate",
							i18n.ES: "Tipo General",
						},
						Values: []*tax.RateValue{
							{
								Percent: num.MakePercentage(70, 3),
							},
						},
					},
					{
						Key: common.TaxRateReduced,
						Name: i18n.String{
							i18n.EN: "Reduced Rate",
							i18n.ES: "Tipo Reducido",
						},
						Values: []*tax.RateValue{
							{
								Percent: num.MakePercentage(30, 3),
							},
						},
					},
				},
			},

			//
			// IPSI
			//
			{
				Code:     TaxCategoryIPSI,
				Retained: false,
				Name: i18n.String{
					i18n.EN: "IPSI",
					i18n.ES: "IPSI",
				},
				Desc: i18n.String{
					i18n.EN: "Production, Services, and Import Tax",
					i18n.ES: "Impuesto sobre la Producción, los Servicios y la Importación",
				},
				// IPSI rates are complex and don't align well regular rates. Users are
				// recommended to include whatever percentage applies to their situation
				// directly in the invoice.
				Rates: []*tax.Rate{},
			},

			//
			// IRPF
			//
			{
				Code:     TaxCategoryIRPF,
				Retained: true,
				Name: i18n.String{
					i18n.EN: "IRPF",
					i18n.ES: "IRPF",
				},
				Desc: i18n.String{
					i18n.EN: "Personal income tax.",
					i18n.ES: "Impuesto sobre la renta de las personas físicas.",
				},
				Rates: []*tax.Rate{
					{
						Key: TaxRatePro,
						Name: i18n.String{
							i18n.EN: "Professional Rate",
							i18n.ES: "Professionales",
						},
						Values: []*tax.RateValue{
							{
								Since:   cal.NewDate(2015, 7, 12),
								Percent: num.MakePercentage(150, 3),
							},
							{
								Since:   cal.NewDate(2015, 1, 1),
								Percent: num.MakePercentage(190, 3),
							},
							{
								Since:   cal.NewDate(2012, 9, 1),
								Percent: num.MakePercentage(210, 3),
							},
							{
								Since:   cal.NewDate(2007, 1, 1),
								Percent: num.MakePercentage(150, 3),
							},
						},
					},
					{
						Key: TaxRateProStart,
						Name: i18n.String{
							i18n.EN: "Professional Starting Rate",
							i18n.ES: "Professionales Inicio",
						},
						Values: []*tax.RateValue{
							{
								Since:   cal.NewDate(2007, 1, 1),
								Percent: num.MakePercentage(70, 3),
							},
						},
					},
					{
						Key: TaxRateCapital,
						Name: i18n.String{
							i18n.EN: "Rental or Interest Capital",
							i18n.ES: "Alquileres o Intereses de Capital",
						},
						Values: []*tax.RateValue{
							{
								Since:   cal.NewDate(2007, 1, 1),
								Percent: num.MakePercentage(190, 3),
							},
						},
					},
					{
						Key: TaxRateModules,
						Name: i18n.String{
							i18n.EN: "Modules Rate",
							i18n.ES: "Tipo Modulos",
						},
						Values: []*tax.RateValue{
							{
								Since:   cal.NewDate(2007, 1, 1),
								Percent: num.MakePercentage(10, 3),
							},
						},
					},
				},
			},
		},
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Calculate will perform any regime specific calculations.
func Calculate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return normalizeTaxIdentity(obj)
	}
	return nil
}
