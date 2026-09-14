package cr

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// TaxCategoryISC is the code for the Impuesto Selectivo de Consumo.
const TaxCategoryISC cbc.Code = "ISC"

// Rates and official names per Ley 6826, arts. 10-11, and "Anexos y
// Estructuras v4.4", nota 8.1.
var taxCategories = []*tax.CategoryDef{
	//
	// IVA
	//
	{
		Code:     tax.CategoryVAT,
		Retained: false,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.ES: "IVA",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.ES: "Impuesto al Valor Agregado",
		},
		Keys: tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General Rate",
					i18n.ES: "Tarifa General",
				},
				Description: i18n.String{
					i18n.EN: "Officially \"tarifa general 13%\": all operations subject to the tax not covered by a reduced rate (art. 10, Ley 6826).",
					i18n.ES: "Oficialmente \"tarifa general 13%\": todas las operaciones sujetas al impuesto no cubiertas por una tarifa reducida (art. 10, Ley 6826).",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2019, 7, 1),
						Percent: num.MakePercentage(130, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate Rate",
					i18n.ES: "Tarifa Intermedia",
				},
				Description: i18n.String{
					i18n.EN: "Officially \"tarifa reducida 4%\": air tickets with origin or destination in Costa Rica, and private health services (art. 11.1, Ley 6826).",
					i18n.ES: "Oficialmente \"tarifa reducida 4%\": boletos aéreos con origen o destino en Costa Rica y servicios de salud privados (art. 11.1, Ley 6826).",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2019, 7, 1),
						Percent: num.MakePercentage(40, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.ES: "Tarifa Reducida",
				},
				Description: i18n.String{
					i18n.EN: "Officially \"tarifa reducida 2%\": medicines and their production inputs, non-exempt private education, personal insurance premiums, and purchases by state higher-education institutions (art. 11.2, Ley 6826).",
					i18n.ES: "Oficialmente \"tarifa reducida 2%\": medicamentos y sus insumos de producción, educación privada no exenta, primas de seguros personales y compras de las instituciones estatales de educación superior (art. 11.2, Ley 6826).",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2019, 7, 1),
						Percent: num.MakePercentage(20, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Super-Reduced Rate",
					i18n.ES: "Tarifa Superreducida",
				},
				Description: i18n.String{
					i18n.EN: "Officially \"tarifa reducida 1%\": the basic tax basket (canasta básica) and its production chain, named grains for animal feed, veterinary products, agricultural and non-sport fishing inputs, and reef-safe sunscreens (art. 11.3, Ley 6826).",
					i18n.ES: "Oficialmente \"tarifa reducida 1%\": la canasta básica tributaria y su cadena de producción, granos para alimentación animal, productos veterinarios, insumos agropecuarios y de pesca no deportiva, y protectores solares que no dañan los arrecifes (art. 11.3, Ley 6826).",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2019, 7, 1),
						Percent: num.MakePercentage(10, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSpecial,
				Name: i18n.String{
					i18n.EN: "Special Rate",
					i18n.ES: "Tarifa Especial",
				},
				Description: i18n.String{
					i18n.EN: "Officially \"tarifa reducida 0,5%\": registered and certified organic agricultural and agro-industrial products, and production inputs of organized organic producer groups (art. 11.4, Ley 6826, added by Ley 10256).",
					i18n.ES: "Oficialmente \"tarifa reducida 0,5%\": productos agropecuarios y agroindustriales orgánicos registrados y certificados, e insumos de producción de grupos de personas productoras orgánicas organizadas (art. 11.4, Ley 6826, adicionado por la Ley 10256).",
				},
				Values: []*tax.RateValueDef{
					{
						// Ley 10256, published in La Gaceta 185, Alcance 207.
						Since:   cal.NewDate(2022, 9, 29),
						Percent: num.MakePercentage(5, 3),
					},
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Ley 6826, Ley del Impuesto al Valor Agregado, arts. 10-11"),
				URL:   "https://sinalevi.go.cr/ResultadosNormativa/Informacion?param1=32526&param2=150649&param3=1",
			},
			{
				Title: i18n.NewString("Anexos y Estructuras de Comprobantes Electrónicos v4.4, nota 8.1"),
				URL:   "https://www.hacienda.go.cr/docs/ANEXOS_Y_ESTRUCTURAS_V4.4.pdf",
			},
		},
	},
	//
	// ISC - rates are set per product by decree, so invoice lines supply the
	// percentage directly.
	//
	{
		Code:     TaxCategoryISC,
		Retained: false,
		Name: i18n.String{
			i18n.EN: "ISC",
			i18n.ES: "ISC",
		},
		Title: i18n.String{
			i18n.EN: "Selective Consumption Tax",
			i18n.ES: "Impuesto Selectivo de Consumo",
		},
		Description: &i18n.String{
			i18n.EN: "Single-stage ad valorem excise on selected goods established by Title II of Ley 4961, with rates set per product by decree.",
			i18n.ES: "Impuesto monofásico ad valorem sobre mercancías seleccionadas establecido por el Título II de la Ley 4961, con tarifas fijadas por decreto para cada producto.",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Ley 4961, Ley de Reforma Tributaria, Título II (Consolidación de Impuestos Selectivos de Consumo)"),
				URL:   "https://sinalevi.go.cr/ResultadosNormativa/Informacion?param1=18507&param2=151936&param3=1",
			},
		},
	},
}
