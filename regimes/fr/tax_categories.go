package fr

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.Category{
	//
	// VAT
	//
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.FR: "TVA",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.FR: "Taxe sur la Valeur Ajoutée",
		},
		Sources: []*tax.Source{
			{
				Title: i18n.String{
					i18n.EN: "What are the current VAT rates in France and in the European Union?",
					i18n.FR: "Quels sont les taux de TVA en vigueur en France et dans l'Union européenne?",
				},
				URL: "https://www.economie.gouv.fr/cedef/taux-tva-france-et-union-europeenne",
			},
		},
		Retained: false,
		Rates: []*tax.Rate{
			{
				Key: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard rate",
					i18n.FR: "Taux normal",
				},
				Description: i18n.String{
					i18n.EN: "For the majority of sales of goods and services: it applies to all products or services for which no other rate is expressly provided.",
					i18n.FR: "Pour la majorité des ventes de biens et des prestations de services : il s'applique à tous les produits ou services pour lesquels aucun autre taux n'est expressément prévu.",
				},
				Values: []*tax.RateValue{
					{
						Since:   cal.NewDate(2011, 1, 4),
						Percent: num.MakePercentage(20, 2),
					},
				},
			},
			{
				Key: tax.RateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate rate",
					i18n.FR: "Taux intermédiaire",
				},
				Description: i18n.String{
					i18n.EN: "Applicable in particular to unprocessed agricultural products, firewood, housing improvement works which do not benefit from the 5.5% rate, to certain accommodation and camping services, to fairs and exhibitions, fairground games and rides, to entrance fees to museums, zoos, monuments, to passenger transport, to the processing of waste, restoration.",
					i18n.FR: "Notamment applicable aux produits agricoles non transformés, au bois de chauffage, aux travaux d'amélioration du logement qui ne bénéficient pas du taux de 5,5%, à certaines prestations de logement et de camping, aux foires et salons, jeux et manèges forains, aux droits d'entrée des musées, zoo, monuments, aux transports de voyageurs, au traitement des déchets, à la restauration.",
				},
				Values: []*tax.RateValue{
					{
						Since:   cal.NewDate(2014, 1, 1),
						Percent: num.MakePercentage(10, 2),
					},
				},
			},
			{
				Key: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced rate",
					i18n.FR: "Taux réduit",
				},
				Description: i18n.String{
					i18n.EN: "Concerns most food products, feminine hygiene protection products, equipment and services for the disabled, books on any medium, gas and electricity subscriptions, supply of heat from renewable energies, supply of meals in school canteens, ticketing for live shows and cinemas, certain imports and deliveries of works of art, improvement works the energy quality of housing, social or emergency housing, home ownership.",
					i18n.FR: "Concerne l'essentiel des produits alimentaires, les produits de protection hygiénique féminine, équipements et services pour handicapés, livres sur tout support, abonnements gaz et électricité, fourniture de chaleur issue d’énergies renouvelables, fourniture de repas dans les cantines scolaires, billeterie de spectacle vivant et de cinéma, certaines importations et livraisons d'œuvres d'art, travaux d’amélioration de la qualité énergétique des logements, logements sociaux ou d'urgence, accession à la propriété.",
				},
				Values: []*tax.RateValue{
					{
						Since:   cal.NewDate(2014, 1, 1),
						Percent: num.MakePercentage(55, 3),
					},
				},
			},
			{
				Key: tax.RateSpecial,
				Name: i18n.String{
					i18n.EN: "Special rate",
					i18n.FR: "Taux particulier",
				},
				Description: i18n.String{
					i18n.EN: "Reserved for medicines reimbursable by social security, sales of live animals for slaughter and charcuterie to non-taxable persons, the television license fee, certain shows and press publications registered with the Joint Commission for Publications and Press Agencies.",
					i18n.FR: "Réservé aux médicaments remboursables par la sécurité sociale, aux ventes d’animaux vivants de boucherie et de charcuterie à des non assujettis, à la redevance télévision, à certains spectacles et aux publications de presse inscrites à la Commission paritaire des publications et agences de presse.",
				},
				Values: []*tax.RateValue{
					{
						Since:   cal.NewDate(2014, 1, 1),
						Percent: num.MakePercentage(21, 3),
					},
				},
			},
		},
	},
}
