package es

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regions/common"
	"github.com/invopop/gobl/tax"
)

// Local tax category definitions which are not considered standard.
const (
	TaxCategoryIRPF org.Code = "IRPF"
	TaxCategoryIGIC org.Code = "IGIC"
	TaxCategoryIPSI org.Code = "IPSI"
)

// Specific tax rate codes.
const (
	// IRPF non-standard Rates (usually for self-employed)
	TaxRatePro                org.Key = "pro"                 // Professional Services
	TaxRateProStart           org.Key = "pro-start"           // Professionals, first 2 years
	TaxRateModules            org.Key = "modules"             // Module system
	TaxRateAgriculture        org.Key = "agriculture"         // Agricultural
	TaxRateAgricultureSpecial org.Key = "agriculture-special" // Agricultural special
	TaxRateCapital            org.Key = "capital"             // Rental or Interest

	// Special tax rate surcharge extension
	TaxRateEquivalence org.Key = "eqs"
)

// Scheme key definitions
const (
	SchemeSimplified      org.Key = "simplified"
	SchemeCustomerIssued  org.Key = "customer-issued"
	SchemeTravelAgency    org.Key = "travel-agency"
	SchemeSecondHandGoods org.Key = "second-hand-goods"
	SchemeArt             org.Key = "art"
	SchemeAntiques        org.Key = "antiques"
	SchemeCashBasis       org.Key = "cash-basis"
)

// Inbox key and role definitions
const (
	InboxKeyFACE org.Key = "face"

	// Main roles defined in FACE
	InboxRoleFiscal    org.Key = "fiscal"    // Fiscal / 01
	InboxRoleRecipient org.Key = "recipient" // Receptor / 02
	InboxRolePayer     org.Key = "payer"     // Pagador / 03
	InboxRoleCustomer  org.Key = "customer"  // Comprador / 04

)

// Locality code definitions for Spain
const (
	LocalityVI l10n.Code = "VI" // (01) Álava
	LocalityAB l10n.Code = "AB" // (02) Albacete
	LocalityA  l10n.Code = "A"  // (03) Alicante
	LocalityAL l10n.Code = "AL" // (04) Almería
	LocalityAV l10n.Code = "AV" // (05) Ávila
	LocalityBA l10n.Code = "BA" // (06) Badajoz
	LocalityPM l10n.Code = "PM" // (07) Baleares
	LocalityIB l10n.Code = "IB" // (07) Baleares
	LocalityB  l10n.Code = "B"  // (08) Barcelona
	LocalityBU l10n.Code = "BU" // (09) Burgos
	LocalityCC l10n.Code = "CC" // (10) Cáceres
	LocalityCA l10n.Code = "CA" // (11) Cádiz
	LocalityCS l10n.Code = "CS" // (12) Castellon
	LocalityCR l10n.Code = "CR" // (13) Ciudad Real
	LocalityCO l10n.Code = "CO" // (14) Cordoba
	LocalityC  l10n.Code = "C"  // (15) La Coruña
	LocalityCU l10n.Code = "CU" // (16) Cuenca
	LocalityGE l10n.Code = "GE" // (17) Gerona
	LocalityGI l10n.Code = "GI" // (17) Girona
	LocalityGR l10n.Code = "GR" // (18) Granada
	LocalityGU l10n.Code = "GU" // (19) Guadalajara
	LocalitySS l10n.Code = "SS" // (20) Guipúzcoa
	LocalityH  l10n.Code = "H"  // (21) Huelva
	LocalityHU l10n.Code = "HU" // (22) Huesca
	LocalityJ  l10n.Code = "J"  // (23) Jaén
	LocalityLE l10n.Code = "LE" // (24) León
	LocalityL  l10n.Code = "L"  // (25) Lérida / Lleida
	LocalityLO l10n.Code = "LO" // (26) La Rioja
	LocalityLU l10n.Code = "LU" // (27) Lugo
	LocalityM  l10n.Code = "M"  // (28) Madrid
	LocalityMA l10n.Code = "MA" // (29) Málaga
	LocalityMU l10n.Code = "MU" // (30) Murcia
	LocalityNA l10n.Code = "NA" // (31) Navarra
	LocalityOR l10n.Code = "OR" // (32) Orense
	LocalityOU l10n.Code = "OU" // (32) Orense
	LocalityO  l10n.Code = "O"  // (33) Asturias
	LocalityP  l10n.Code = "P"  // (34) Palencia
	LocalityGC l10n.Code = "GC" // (35) Las Palmas
	LocalityPO l10n.Code = "PO" // (36) Pontevedra
	LocalitySA l10n.Code = "SA" // (37) Salamanca
	LocalityTF l10n.Code = "TF" // (38) Santa Cruz de Tenerife
	LocalityS  l10n.Code = "S"  // (39) Cantabria
	LocalitySG l10n.Code = "SG" // (40) Segovia
	LocalitySE l10n.Code = "SE" // (41) Sevilla
	LocalitySO l10n.Code = "SO" // (42) Soria
	LocalityT  l10n.Code = "T"  // (43) Tarragona
	LocalityTE l10n.Code = "TE" // (44) Teruel
	LocalityTO l10n.Code = "TO" // (45) Toledo
	LocalityV  l10n.Code = "V"  // (46) Valencia
	LocalityVA l10n.Code = "VA" // (47) Valladolid
	LocalityBI l10n.Code = "BI" // (48) Vizcaya
	LocalityZA l10n.Code = "ZA" // (49) Zamora
	LocalityZ  l10n.Code = "Z"  // (50) Zaragoza
	LocalityCE l10n.Code = "CE" // (51) Ceuta
	LocalityML l10n.Code = "ML" // (52) Melilla
)

// Custom keys used typically in meta information.
const (
	KeyPost org.Key = "post"
)

// Region provides the Spanish region definition
func Region() *tax.Region {
	return &tax.Region{
		Country:  l10n.ES,
		Currency: "EUR",
		Name: i18n.String{
			i18n.EN: "Spain",
			i18n.ES: "España",
		},
		ValidateDocument: Validate,
		Localities: tax.Localities{
			{
				Code: LocalityVI,
				Name: i18n.String{i18n.ES: "Ávila"},
				Meta: org.Meta{KeyPost: "01"},
			},
			{
				Code: LocalityAB,
				Name: i18n.String{i18n.ES: "Albacete"},
				Meta: org.Meta{KeyPost: "02"},
			},
			{
				Code: LocalityA,
				Name: i18n.String{i18n.ES: "Alicante"},
				Meta: org.Meta{KeyPost: "03"},
			},
			{
				Code: LocalityAL,
				Name: i18n.String{i18n.ES: "Almería"},
				Meta: org.Meta{KeyPost: "04"},
			},
			{
				Code: LocalityAV,
				Name: i18n.String{i18n.ES: "Ávila"},
				Meta: org.Meta{KeyPost: "05"},
			},
			{
				Code: LocalityBA,
				Name: i18n.String{i18n.ES: "Badajoz"},
				Meta: org.Meta{KeyPost: "06"},
			},
			{
				Code: LocalityPM,
				Name: i18n.String{i18n.ES: "Baleares"},
				Meta: org.Meta{KeyPost: "07"},
			},
			{
				Code: LocalityIB,
				Name: i18n.String{i18n.ES: "Baleares"},
				Meta: org.Meta{KeyPost: "07"},
			},
			{
				Code: LocalityB,
				Name: i18n.String{i18n.ES: "Barcelona"},
				Meta: org.Meta{KeyPost: "08"},
			},
			{
				Code: LocalityBU,
				Name: i18n.String{i18n.ES: "Burgos"},
				Meta: org.Meta{KeyPost: "09"},
			},
			{
				Code: LocalityCC,
				Name: i18n.String{i18n.ES: "Cáceres"},
				Meta: org.Meta{KeyPost: "10"},
			},
			{
				Code: LocalityCA,
				Name: i18n.String{i18n.ES: "Cadiz"},
				Meta: org.Meta{KeyPost: "11"},
			},
			{
				Code: LocalityCS,
				Name: i18n.String{i18n.ES: "Castellón"},
				Meta: org.Meta{KeyPost: "12"},
			},
			{
				Code: LocalityCR,
				Name: i18n.String{i18n.ES: "Ciudad Real"},
				Meta: org.Meta{KeyPost: "13"},
			},
			{
				Code: LocalityCO,
				Name: i18n.String{i18n.ES: "Cordoba"},
				Meta: org.Meta{KeyPost: "14"},
			},
			{
				Code: LocalityC,
				Name: i18n.String{i18n.ES: "La Coruña"},
				Meta: org.Meta{KeyPost: "15"},
			},
			{
				Code: LocalityCU,
				Name: i18n.String{i18n.ES: "Cuenca"},
				Meta: org.Meta{KeyPost: "16"},
			},
			{
				Code: LocalityGE,
				Name: i18n.String{i18n.ES: "Gerona"},
				Meta: org.Meta{KeyPost: "17"},
			},
			{
				Code: LocalityGI,
				Name: i18n.String{i18n.ES: "Girona"},
				Meta: org.Meta{KeyPost: "17"},
			},
			{
				Code: LocalityGR,
				Name: i18n.String{i18n.ES: "Granada"},
				Meta: org.Meta{KeyPost: "18"},
			},
			{
				Code: LocalityGU,
				Name: i18n.String{i18n.ES: "Guadalajara"},
				Meta: org.Meta{KeyPost: "19"},
			},
			{
				Code: LocalitySS,
				Name: i18n.String{i18n.ES: "Guipúzcoa"},
				Meta: org.Meta{KeyPost: "20"},
			},
			{
				Code: LocalityH,
				Name: i18n.String{i18n.ES: "Huelva"},
				Meta: org.Meta{KeyPost: "21"},
			},
			{
				Code: LocalityHU,
				Name: i18n.String{i18n.ES: "Huesca"},
				Meta: org.Meta{KeyPost: "22"},
			},
			{
				Code: LocalityJ,
				Name: i18n.String{i18n.ES: "Jaén"},
				Meta: org.Meta{KeyPost: "23"},
			},
			{
				Code: LocalityLE,
				Name: i18n.String{i18n.ES: "León"},
				Meta: org.Meta{KeyPost: "24"},
			},
			{
				Code: LocalityL,
				Name: i18n.String{i18n.ES: "Lérida / Lleida"},
				Meta: org.Meta{KeyPost: "25"},
			},
			{
				Code: LocalityLO,
				Name: i18n.String{i18n.ES: "La Rioja"},
				Meta: org.Meta{KeyPost: "26"},
			},
			{
				Code: LocalityLU,
				Name: i18n.String{i18n.ES: "Lugo"},
				Meta: org.Meta{KeyPost: "27"},
			},
			{
				Code: LocalityM,
				Name: i18n.String{i18n.ES: "Madrid"},
				Meta: org.Meta{KeyPost: "28"},
			},
			{
				Code: LocalityMA,
				Name: i18n.String{i18n.ES: "Málaga"},
				Meta: org.Meta{KeyPost: "29"},
			},
			{
				Code: LocalityMU,
				Name: i18n.String{i18n.ES: "Murcia"},
				Meta: org.Meta{KeyPost: "30"},
			},
			{
				Code: LocalityNA,
				Name: i18n.String{i18n.ES: "Navarra"},
				Meta: org.Meta{KeyPost: "31"},
			},
			{
				Code: LocalityOR,
				Name: i18n.String{i18n.ES: "Orense"},
				Meta: org.Meta{KeyPost: "32"},
			},
			{
				Code: LocalityOU,
				Name: i18n.String{i18n.ES: "Orense"},
				Meta: org.Meta{KeyPost: "32"},
			},
			{
				Code: LocalityO,
				Name: i18n.String{i18n.ES: "Asturias"},
				Meta: org.Meta{KeyPost: "33"},
			},
			{
				Code: LocalityP,
				Name: i18n.String{i18n.ES: "Palencia"},
				Meta: org.Meta{KeyPost: "34"},
			},
			{
				Code: LocalityGC,
				Name: i18n.String{i18n.ES: "Las Palmas"},
				Meta: org.Meta{KeyPost: "35"},
			},
			{
				Code: LocalityPO,
				Name: i18n.String{i18n.ES: "Pontevedra"},
				Meta: org.Meta{KeyPost: "36"},
			},
			{
				Code: LocalitySA,
				Name: i18n.String{i18n.ES: "Salamanca"},
				Meta: org.Meta{KeyPost: "37"},
			},
			{
				Code: LocalityTF,
				Name: i18n.String{i18n.ES: "Santa Cruz de Tenerife"},
				Meta: org.Meta{KeyPost: "38"},
			},
			{
				Code: LocalityS,
				Name: i18n.String{i18n.ES: "Cantabria"},
				Meta: org.Meta{KeyPost: "39"},
			},
			{
				Code: LocalitySG,
				Name: i18n.String{i18n.ES: "Segovia"},
				Meta: org.Meta{KeyPost: "40"},
			},
			{
				Code: LocalitySE,
				Name: i18n.String{i18n.ES: "Sevilla"},
				Meta: org.Meta{KeyPost: "41"},
			},
			{
				Code: LocalitySO,
				Name: i18n.String{i18n.ES: "Soria"},
				Meta: org.Meta{KeyPost: "42"},
			},
			{
				Code: LocalityT,
				Name: i18n.String{i18n.ES: "Tarragona"},
				Meta: org.Meta{KeyPost: "43"},
			},
			{
				Code: LocalityTE,
				Name: i18n.String{i18n.ES: "Teruel"},
				Meta: org.Meta{KeyPost: "44"},
			},
			{
				Code: LocalityTO,
				Name: i18n.String{i18n.ES: "Toledo"},
				Meta: org.Meta{KeyPost: "45"},
			},
			{
				Code: LocalityV,
				Name: i18n.String{i18n.ES: "Valencia"},
				Meta: org.Meta{KeyPost: "46"},
			},
			{
				Code: LocalityVA,
				Name: i18n.String{i18n.ES: "Valladolid"},
				Meta: org.Meta{KeyPost: "47"},
			},
			{
				Code: LocalityBI,
				Name: i18n.String{i18n.ES: "Vizcaya"},
				Meta: org.Meta{KeyPost: "48"},
			},
			{
				Code: LocalityZA,
				Name: i18n.String{i18n.ES: "Zamora"},
				Meta: org.Meta{KeyPost: "49"},
			},
			{
				Code: LocalityZ,
				Name: i18n.String{i18n.ES: "Zaragoza"},
				Meta: org.Meta{KeyPost: "50"},
			},
			{
				Code: LocalityCE,
				Name: i18n.String{i18n.ES: "Ceuta"},
				Meta: org.Meta{KeyPost: "51"},
			},
			{
				Code: LocalityML,
				Name: i18n.String{i18n.ES: "Melilla"},
				Meta: org.Meta{KeyPost: "52"},
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
				Categories: []org.Code{
					common.TaxCategoryVAT,
				},
				Note: &org.Note{
					Key:  org.NoteKeyLegal,
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
				Note: &org.Note{
					Key:  org.NoteKeyLegal,
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
				Note: &org.Note{
					Key:  org.NoteKeyLegal,
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
				Note: &org.Note{
					Key:  org.NoteKeyLegal,
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
				Note: &org.Note{
					Key:  org.NoteKeyLegal,
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
				Note: &org.Note{
					Key:  org.NoteKeyLegal,
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
				Note: &org.Note{
					Key:  org.NoteKeyLegal,
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
				Note: &org.Note{
					Key:  org.NoteKeyLegal,
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
				Code: common.TaxCategoryVAT,
				Name: i18n.String{
					i18n.EN: "VAT",
					i18n.ES: "IVA",
				},
				Desc: i18n.String{
					i18n.EN: "Value Added Tax",
					i18n.ES: "Impuesto sobre el Valor Añadido",
				},
				Retained: false,
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
	}
	return nil
}
