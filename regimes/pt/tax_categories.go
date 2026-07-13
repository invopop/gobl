package pt

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

// Local tax category definitions which are not considered standard.
const (
	TaxCategoryIRS cbc.Code = "IRS" // Imposto sobre o Rendimento das Pessoas Singulares
	TaxCategoryIRC cbc.Code = "IRC" // Imposto sobre o Rendimento das Pessoas Coletivas
)

var taxCategories = []*tax.CategoryDef{
	//
	// VAT
	//
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.PT: "IVA",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.PT: "Imposto sobre o Valor Acrescentado",
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General Rate",
					i18n.PT: "Tipo Geral",
				},
				Values: []*tax.RateValueDef{
					{
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							ExtKeyRegion: RegionAzores,
						}),
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(160, 3),
					},
					{
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							ExtKeyRegion: RegionMadeira,
						}),
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(220, 3),
					},
					{
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(230, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate Rate",
					i18n.PT: "Taxa Intermédia", //nolint:misspell
				},
				Values: []*tax.RateValueDef{
					{
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							ExtKeyRegion: RegionAzores,
						}),
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(90, 3),
					},
					{
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							ExtKeyRegion: RegionMadeira,
						}),
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(120, 3),
					},
					{
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(130, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.PT: "Taxa Reduzida",
				},
				Values: []*tax.RateValueDef{
					{
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							ExtKeyRegion: RegionAzores,
						}),
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(40, 3),
					},
					{
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							ExtKeyRegion: RegionMadeira,
						}),
						Since:   cal.NewDate(2024, 10, 1),
						Percent: num.MakePercentage(40, 3),
					},
					{
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							ExtKeyRegion: RegionMadeira,
						}),
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(50, 3),
					},
					{
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(60, 3),
					},
				},
			},
			{
				// Other is a special case for rates that are not defined.
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateOther,
				Name: i18n.String{
					i18n.EN: "Other",
					i18n.PT: "Outro",
				},
				Values: []*tax.RateValueDef{},
			},
		},
	},

	//
	// IRS
	//
	{
		Code:     TaxCategoryIRS,
		Retained: true,
		Name: i18n.String{
			i18n.EN: "IRS",
			i18n.PT: "IRS",
		},
		Title: i18n.String{
			i18n.EN: "Personal income tax",
			i18n.PT: "Imposto sobre o Rendimento das Pessoas Singulares",
		},
		Description: &i18n.String{
			i18n.EN: here.Doc(`
				Personal income tax withheld at source from payments made to individuals
				and self-employed workers. The Portuguese payer retains the tax on each
				payment and remits it to the AT on behalf of the recipient.
			`),
			i18n.PT: here.Doc(`
				Imposto sobre o rendimento retido na fonte sobre pagamentos efetuados a
				pessoas singulares e trabalhadores independentes. O pagador português retém
				o imposto em cada pagamento e entrega-o à AT por conta do titular do
				rendimento.
			`),
		},
	},

	//
	// IRC
	//
	{
		Code:     TaxCategoryIRC,
		Retained: true,
		Name: i18n.String{
			i18n.EN: "IRC",
			i18n.PT: "IRC",
		},
		Title: i18n.String{
			i18n.EN: "Corporate income tax",
			i18n.PT: "Imposto sobre o Rendimento das Pessoas Coletivas",
		},
		Description: &i18n.String{
			i18n.EN: here.Doc(`
				Corporate income tax withheld at source from payments made to legal
				persons. The Portuguese payer retains the tax on each payment and remits it
				to the AT on behalf of the recipient.
			`),
			i18n.PT: here.Doc(`
				Imposto sobre o rendimento retido na fonte sobre pagamentos efetuados a
				pessoas coletivas. O pagador português retém o imposto em cada pagamento e
				entrega-o à AT por conta do titular do rendimento.
			`),
		},
	},
}
