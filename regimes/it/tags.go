package it

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Document tag keys
const (
	TagFreelance cbc.Key = "freelance"
)

// Category tag keys determined from the "Natura" field from FatturaPA.
const (
	TagExcluded                   cbc.Key = "excluded"
	TagNotSubject                 cbc.Key = "not-subject"
	TagNotTaxable                 cbc.Key = "not-taxable"
	TagExempt                     cbc.Key = "exempt"
	TagMarginRegime               cbc.Key = "margin-regime"
	TagReverseCharge              cbc.Key = "reverse-charge"
	TagVATEU                      cbc.Key = "vat-eu"
	TagOther                      cbc.Key = "other"
	TagArticle7                   cbc.Key = "article-7"
	TagExport                     cbc.Key = "export"
	TagIntraCommunity             cbc.Key = "intra-community"
	TagSanMerino                  cbc.Key = "san-merino"
	TagExportSupplies             cbc.Key = "export-supplies"
	TagDeclarationOfIntent        cbc.Key = "declaration-of-intent"
	TagScrap                      cbc.Key = "scrap"
	TagPreciousMetals             cbc.Key = "precious-metals"
	TagConstructionSubcontracting cbc.Key = "construction-subcontracting"
	TagBuildings                  cbc.Key = "buildings"
	TagMobile                     cbc.Key = "mobile"
	TagElectronics                cbc.Key = "electronics"
	TagConstruction               cbc.Key = "construction"
	TagEnergy                     cbc.Key = "energy"
)

var invoiceTags = []*tax.TagDef{
	{
		Key: TagFreelance,
		Name: i18n.String{
			i18n.EN: "Freelancer",
			i18n.IT: "Parcella",
		},
	},
}

var vatZeroTaxTags = []*tax.TagDef{
	{
		Key: TagExcluded,
		Name: i18n.String{
			i18n.EN: "Excluded pursuant to Art. 15, DPR 633/72",
			i18n.IT: "Escluse ex. art. 15 del D.P.R. 633/1972",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N1",
		},
	},
	{
		Key: TagNotSubject,
		Name: i18n.String{
			i18n.EN: "Not subject (this code is no longer permitted to use on invoices emitted from 1 January 2021)",
			i18n.IT: "Non soggette (questo codice non è più utilizzabile a partire dal 1° gennaio 2021)",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N2",
		},
	},
	{
		Key: TagNotSubject.With(TagArticle7),
		Name: i18n.String{
			i18n.EN: "Not subject pursuant to Art. 7, DPR 633/72",
			i18n.IT: "Non soggette ex. art. 7 del D.P.R. 633/72",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N2.1",
		},
	},
	{
		Key: TagNotSubject.With(TagOther),
		Name: i18n.String{
			i18n.EN: "Not subject - other",
			i18n.IT: "Non soggette - altri casi",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N2.2",
		},
	},
	{
		Key: TagNotTaxable,
		Name: i18n.String{
			i18n.EN: "Not taxable (this code is no longer permitted to use on invoices emitted from 1 January 2021)",
			i18n.IT: "Non imponibili (questo codice non è più utilizzabile a partire dal 1° gennaio 2021)",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N3",
		},
	},
	{
		Key: TagNotTaxable.With(TagExport),
		Name: i18n.String{
			i18n.EN: "Not taxable - exports",
			i18n.IT: "Non imponibili - esportazioni",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N3.1",
		},
	},
	{
		Key: TagNotTaxable.With(TagIntraCommunity),
		Name: i18n.String{
			i18n.EN: "Not taxable - intra-community supplies",
			i18n.IT: "Non imponibili - cessioni intracomunitarie",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N3.2",
		},
	},
	{
		Key: TagNotTaxable.With(TagSanMerino),
		Name: i18n.String{
			i18n.EN: "Not taxable - transfers to San Marino",
			i18n.IT: "Non imponibili - cessioni verso San Marino",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N3.3",
		},
	},
	{
		Key: TagNotTaxable.With(TagExportSupplies),
		Name: i18n.String{
			i18n.EN: "Not taxable - export supplies of goods and services",
			i18n.IT: "Non Imponibili - operazioni assimilate alle cessioni all'esportazione",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N3.4",
		},
	},
	{
		Key: TagNotTaxable.With(TagDeclarationOfIntent),
		Name: i18n.String{
			i18n.EN: "Not taxable - declaration of intent",
			i18n.IT: "Non imponibili - dichiarazioni d'intento",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N3.5",
		},
	},
	{
		Key: TagNotTaxable.With(TagOther),
		Name: i18n.String{
			i18n.EN: "Not taxable - other",
			i18n.IT: "Non imponibili - altre operazioni che non concorrono alla formazione del plafond",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N3.6",
		},
	},
	{
		Key: TagExempt,
		Name: i18n.String{
			i18n.EN: "Exempt",
			i18n.IT: "Esenti",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N4",
		},
	},
	{
		Key: TagMarginRegime,
		Name: i18n.String{
			i18n.EN: "Margin regime / VAT not exposed",
			i18n.IT: "Regime del margine/IVA non esposta in fattura",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N5",
		},
	},
	{
		Key: TagReverseCharge,
		Name: i18n.String{
			i18n.EN: "Reverse charge (for transactions in reverse charge or for self invoicing for purchase of extra UE services or for import of goods only in the cases provided for) — (this code is no longer permitted to use on invoices emitted from 1 January 2021)",
			i18n.IT: "Inversione contabile (per operazioni in regime di inversione contabile o per autofattura per acquisti di servizi extra UE o per importazioni di beni solo nei casi previsti) — (questo codice non è più utilizzabile a partire dal 1° gennaio 2021)",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N6",
		},
	},
	{
		Key: TagReverseCharge.With(TagScrap),
		Name: i18n.String{
			i18n.EN: "Reverse charge - Transfer of scrap and of other recyclable materials",
			i18n.IT: "Inversione contabile - cessione di rottami e altri materiali di recupero",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N6.1",
		},
	},
	{
		Key: TagReverseCharge.With(TagPreciousMetals),
		Name: i18n.String{
			i18n.EN: "Reverse charge - Transfer of gold and pure silver pursuant to law 7/2000 as well as used jewelery to OPO",
			i18n.IT: "Inversione contabile - cessione di oro e argento ai sensi della legge 7/2000 nonché di oreficeria usata ad OPO",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N6.2",
		},
	},
	{
		Key: TagReverseCharge.With(TagConstructionSubcontracting),
		Name: i18n.String{
			i18n.EN: "Reverse charge - Construction subcontracting",
			i18n.IT: "Inversione contabile - subappalto nel settore edile",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N6.3",
		},
	},
	{
		Key: TagReverseCharge.With(TagBuildings),
		Name: i18n.String{
			i18n.EN: "Reverse charge - Transfer of buildings",
			i18n.IT: "Inversione contabile - cessione di fabbricati",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N6.4",
		},
	},
	{
		Key: TagReverseCharge.With(TagMobile),
		Name: i18n.String{
			i18n.EN: "Reverse charge - Transfer of mobile phones",
			i18n.IT: "Inversione contabile - cessione di telefoni cellulari",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N6.5",
		},
	},
	{
		Key: TagReverseCharge.With(TagElectronics),
		Name: i18n.String{
			i18n.EN: "Reverse charge - Transfer of electronic products",
			i18n.IT: "Inversione contabile - cessione di prodotti elettronici",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N6.6",
		},
	},
	{
		Key: TagReverseCharge.With(TagConstruction),
		Name: i18n.String{
			i18n.EN: "Reverse charge - provisions in the construction and related sectors",
			i18n.IT: "Inversione contabile - prestazioni comparto edile e settori connessi",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N6.7",
		},
	},
	{
		Key: TagReverseCharge.With(TagEnergy),
		Name: i18n.String{
			i18n.EN: "Reverse charge - transactions in the energy sector",
			i18n.IT: "Inversione contabile - operazioni settore energetico",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N6.8",
		},
	},
	{
		Key: TagReverseCharge.With(TagOther),
		Name: i18n.String{
			i18n.EN: "Reverse charge - other cases",
			i18n.IT: "Inversione contabile - altri casi",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N6.9",
		},
	},
	{
		Key: TagVATEU,
		Name: i18n.String{
			i18n.EN: "VAT paid in other EU countries (telecommunications, tele-broadcasting and electronic services provision pursuant to Art. 7 -octies letter a, b, art. 74-sexies Italian Presidential Decree 633/72)",
			i18n.IT: "IVA assolta in altro stato UE (prestazione di servizi di telecomunicazioni, tele-radiodiffusione ed elettronici ex art. 7-octies lett. a, b, art. 74-sexies DPR 633/72)",
		},
		Meta: cbc.Meta{
			KeyFatturaPANatura: "N7",
		},
	},
}
