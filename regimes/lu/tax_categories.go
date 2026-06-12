package lu

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// taxCategories defines Luxembourg's VAT rates with full rate history.
//
// Luxembourg applies a standard, intermediate ("parking"), reduced, and
// super-reduced rate.  In 2023 all four rates were temporarily reduced by
// one percentage point (cost-of-living measure; Law of 30 Nov 2022) and
// restored to their pre-2023 values on 1 January 2024.
var taxCategories = []*tax.CategoryDef{
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.FR: "TVA",
			i18n.LB: "TVA",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.FR: "Taxe sur la Valeur Ajoutée",
			i18n.LB: "Taxe sur la Valeur Ajoutée",
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.FR: "Taux normal",
					i18n.LB: "Normalsaz",
				},
				Description: i18n.String{
					i18n.EN: "General rate applicable to most goods and services.",
					i18n.FR: "Taux applicable à la plupart des biens et services.",
				},
				Values: []*tax.RateValueDef{
					// Restored after temporary 2023 reduction (Law of 19 Dec 2023).
					{Since: cal.NewDate(2024, 1, 1), Percent: num.MakePercentage(170, 3)},
					// Temporary 1 pp reduction throughout 2023 (Law of 30 Nov 2022).
					{Since: cal.NewDate(2023, 1, 1), Percent: num.MakePercentage(160, 3)},
					// Increased from 15% to 17% (Law of 19 Dec 2014, effective 2015).
					{Since: cal.NewDate(2015, 1, 1), Percent: num.MakePercentage(170, 3)},
					// Original rate when Luxembourg's TVA was introduced.
					{Since: cal.NewDate(1992, 1, 1), Percent: num.MakePercentage(150, 3)},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate Rate",
					i18n.FR: "Taux intermédiaire",
					i18n.LB: "Taux intermédiaire",
				},
				Description: i18n.String{
					i18n.EN: "Applies to wine, certain fuels, printed advertising material, and similar goods (EU parking rate).",
					i18n.FR: "Applicable aux vins, certains combustibles, imprimés publicitaires et biens similaires (taux parking UE).",
				},
				Values: []*tax.RateValueDef{
					{Since: cal.NewDate(2024, 1, 1), Percent: num.MakePercentage(140, 3)},
					{Since: cal.NewDate(2023, 1, 1), Percent: num.MakePercentage(130, 3)},
					{Since: cal.NewDate(2015, 1, 1), Percent: num.MakePercentage(140, 3)},
					{Since: cal.NewDate(1992, 1, 1), Percent: num.MakePercentage(120, 3)},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.FR: "Taux réduit",
					i18n.LB: "Reduzéierter Saz",
				},
				Description: i18n.String{
					i18n.EN: "Applies to food, non-alcoholic beverages, pharmaceuticals, books, newspapers, water, and passenger transport.",
					i18n.FR: "Applicable aux denrées alimentaires, boissons non alcoolisées, médicaments, livres, journaux, eau et transport de personnes.",
				},
				Values: []*tax.RateValueDef{
					{Since: cal.NewDate(2024, 1, 1), Percent: num.MakePercentage(80, 3)},
					{Since: cal.NewDate(2023, 1, 1), Percent: num.MakePercentage(70, 3)},
					{Since: cal.NewDate(2015, 1, 1), Percent: num.MakePercentage(80, 3)},
					{Since: cal.NewDate(1992, 1, 1), Percent: num.MakePercentage(60, 3)},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Super-reduced Rate",
					i18n.FR: "Taux super-réduit",
					i18n.LB: "Super-reduzéierter Saz",
				},
				Description: i18n.String{
					i18n.EN: "Applies to certain pharmaceutical products, footwear and clothing for children, social housing, natural gas, and electricity.",
					i18n.FR: "Applicable à certains médicaments, chaussures et vêtements pour enfants, logements sociaux, gaz naturel et électricité.",
				},
				// The super-reduced rate was not changed during the 2023 reduction.
				Values: []*tax.RateValueDef{
					{Since: cal.NewDate(2015, 1, 1), Percent: num.MakePercentage(30, 3)},
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "AED – VAT rates",
					i18n.FR: "AED – Taux de TVA",
				},
				URL: "https://www.aed.public.lu/en/tva/taux-tva.html",
			},
			{
				Title: i18n.String{
					i18n.EN: "Law of 30 November 2022 – Temporary rate reduction",
					i18n.FR: "Loi du 30 novembre 2022 – Réduction temporaire des taux",
				},
				URL: "https://www.legilux.public.lu/eli/etat/leg/loi/2022/11/30/a559/jo",
			},
			{
				Title: i18n.String{
					i18n.EN: "Law of 19 December 2014 – Rate increase to current levels",
					i18n.FR: "Loi du 19 décembre 2014 – Augmentation des taux",
				},
				URL: "https://www.legilux.public.lu/eli/etat/leg/loi/2014/12/19/n1/jo",
			},
			{
				Title: i18n.String{
					i18n.EN: "European Commission – VAT rates",
				},
				URL: "https://taxation-customs.ec.europa.eu/taxation/vat/eu-vat-rules/vat-rates_en",
			},
		},
	},
}
