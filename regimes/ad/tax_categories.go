package ad

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// taxCategories defines the IGI tax categories for Andorra.
//
// All five rates and their descriptions are sourced directly from the
// Departament de Tributs i de Fronteres official page:
// https://www.e-tramits.ad/tramits/impostos/igi
//
// IGI was introduced on 1 January 2013 by Llei 11/2012, del 21 de juny.
// Rates have remained stable since introduction.
func taxCategories() []*tax.CategoryDef {
	return []*tax.CategoryDef{
		{
			Code: tax.CategoryVAT,
			Name: i18n.String{
				i18n.EN: "IGI",
				i18n.CA: "IGI",
				i18n.ES: "IGI",
			},
			Title: i18n.String{
				i18n.EN: "Indirect General Tax",
				i18n.CA: "Impost General Indirecte",
				i18n.ES: "Impuesto General Indirecto",
			},
			Sources: []*cbc.Source{
				{
					Title: i18n.NewString("Departament de Tributs i de Fronteres - IGI"),
					URL:   "https://www.e-tramits.ad/tramits/impostos/igi", // ("Quins són els tipus de gravamen de l'IGI?" section)
				},
			},
			Retained: false,
			Keys:     tax.GlobalVATKeys(),
			Rates: []*tax.RateDef{
				{
					// Tipus de gravamen general: 4,5%
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateGeneral,
					Name: i18n.String{
						i18n.EN: "General rate",
						i18n.CA: "Tipus de gravamen general",
						i18n.ES: "Tipo de gravamen general",
					},
					Description: i18n.String{
						i18n.EN: "Standard IGI rate applied to most goods and services. " +
							"At 4.5% it is the lowest standard indirect tax rate in Europe.",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2013, 1, 1),
							Percent: num.MakePercentage(45, 3), // 4.5%
						},
					},
				},
				{
					// Tipus de gravamen superreduït: 0%
					// Applies to: healthcare, education, social services,
					// residential rentals, CASS-reimbursable medications,
					// sports services by non-profits.
					Keys: []cbc.Key{tax.KeyZero},
					Rate: tax.RateZero,
					Name: i18n.String{
						i18n.EN: "Super-reduced rate",
						i18n.CA: "Tipus de gravamen superreduït",
						i18n.ES: "Tipo de gravamen superreducido",
					},
					Description: i18n.String{
						i18n.EN: "Zero-rated IGI applied to healthcare by public or para-public entities, " +
							"education at all levels, social assistance services, " +
							"residential property rentals, CASS-reimbursable medications, " +
							"and sports services provided by non-profit organisations.",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2013, 1, 1),
							Percent: num.MakePercentage(0, 2), // 0%
						},
					},
				},
				{
					// Tipus de gravamen reduït: 1%
					// Applies to: food and beverages for human/animal consumption
					// (excl. alcohol), water, books, newspapers and magazines.
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateReduced,
					Name: i18n.String{
						i18n.EN: "Reduced rate",
						i18n.CA: "Tipus de gravamen reduït",
						i18n.ES: "Tipo de gravamen reducido",
					},
					Description: i18n.String{
						i18n.EN: "Reduced IGI rate applied to food and beverages for human or animal " +
							"consumption (excluding alcohol), water, and books, newspapers and " +
							"magazines not consisting exclusively of advertising.",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2013, 1, 1),
							Percent: num.MakePercentage(1, 2), // 1%
						},
					},
				},
				{
					// Tipus de gravamen especial: 2,5%
					// Applies to: passenger transport, cultural and artistic
					// services not provided by public or non-profit entities,
					// antiques and collectibles.
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateSpecial,
					Name: i18n.String{
						i18n.EN: "Special rate",
						i18n.CA: "Tipus de gravamen especial",
						i18n.ES: "Tipo de gravamen especial",
					},
					Description: i18n.String{
						i18n.EN: "Special IGI rate applied to passenger transport, cultural and artistic " +
							"services not provided by public or non-profit entities, and antiques, " +
							"collectibles and works of art.",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2013, 1, 1),
							Percent: num.MakePercentage(25, 3), // 2.5%
						},
					},
				},
				{
					// Tipus de gravamen incrementat: 9,5%
					// Applies to: banking and financial services.
					// Note: GOBL has no dedicated "increased" rate constant;
					// RateOther is used as the closest equivalent.
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateOther,
					Name: i18n.String{
						i18n.EN: "Increased rate",
						i18n.CA: "Tipus de gravamen incrementat",
						i18n.ES: "Tipo de gravamen incrementado",
					},
					Description: i18n.String{
						i18n.EN: "Increased IGI rate applied to banking and financial services.",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2013, 1, 1),
							Percent: num.MakePercentage(95, 3), // 9.5%
						},
					},
				},
			},
		},
	}
}