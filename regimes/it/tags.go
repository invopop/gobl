package it

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

// Document tag keys
const (
	// Tags for document type
	TagFreelance       cbc.Key = "freelance"
	TagCeilingExceeded cbc.Key = "ceiling-exceeded"
	TagSanMerinoPaper  cbc.Key = "san-merino-paper"

	// Tags for Fiscal Regime
	TagMinimumTaxPayers cbc.Key = "minimum-tax-payers"
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

// Retained tax tag keys determined from the "CausalePagamento" field from FatturaPA.
// Source: https://www.agenziaentrate.gov.it/portale/documents/20143/4115385/CU_istr_2022.pdf
// Section VII, Part 2
const (
	TagSelfEmployedHabitual         cbc.Key = "self-employed-habitual"         // A
	TagAuthorIPUsage                cbc.Key = "author-ip-usage"                // B
	TagPartnershipAgreements        cbc.Key = "partnership-agreements"         // C
	TagFounderLimitedCompany        cbc.Key = "founder-limited-company"        // D
	TagCertificationDishonoredBills cbc.Key = "certification-dishonored-bills" // E
	TagHonoraryJudicialOfficers     cbc.Key = "honorary-judicial-officers"     // F
	TagCessationSports              cbc.Key = "cessation-sports"               // G
	TagCessationAgency              cbc.Key = "cessation-agency"               // H
	TagCessationNotary              cbc.Key = "cessation-notary"               // I
	TagTruffleGathering             cbc.Key = "truffle-gathering"              // J
	TagCivilService                 cbc.Key = "civil-service"                  // K
	TagEntitledIPUsage              cbc.Key = "entitled-ip-usage"              // L
	TagPurchasedIPUsage             cbc.Key = "purchased-ip-usage"             // L1
	TagOccasionalSelfEmployment     cbc.Key = "occasional-self-employment"     // M
	TagAssumptionObligations        cbc.Key = "assumption-obligations"         // M1
	TagENPAPISelfEmployment         cbc.Key = "enpapi-self-employment"         // M2
	TagAmateurSports                cbc.Key = "amateur-sports"                 // N
	TagNonENPAPISelfEmployment      cbc.Key = "non-enpapi-self-employment"     // O
	TagNonENPAPIObligations         cbc.Key = "non-enpapi-obligations"         // O1
	TagSwissEquipmentsUse           cbc.Key = "swiss-equipments-use"           // P
	TagSingleMandateAgent           cbc.Key = "single-mandate-agent"           // Q
	TagMultiMandateAgent            cbc.Key = "multi-mandate-agent"            // R
	TagCommissionAgent              cbc.Key = "commission-agent"               // S
	TagComissionBroker              cbc.Key = "commission-broker"              // T
	TagBusinessReferrer             cbc.Key = "business-referrer"              // U
	TagHomeSales                    cbc.Key = "home-sales"                     // V
	TagOccasionalCommercial         cbc.Key = "occasional-commercial"          // V1
	TagHomeSalesNonHabitual         cbc.Key = "home-sales-non-habitual"        // V2
	TagContractWork2021             cbc.Key = "contract-work-2021"             // W
	TagEUFees2004                   cbc.Key = "eu-fees-2004"                   // X
	TagEUFees2005H1                 cbc.Key = "eu-fees-2005-h1"                // Y
	TagOtherTitle                   cbc.Key = "other-title"                    // ZO
)

// This is only a partial list of all the potential tags that
// could be available for use in Italy. Given the complexity
// involved, we've focussed here on the most useful.
var invoiceTags = []*tax.Tag{
	// *** Document Type Tags ***
	{
		Key: TagFreelance,
		Name: i18n.String{
			i18n.EN: "Freelancer",
			i18n.IT: "Parcella",
		},
	},
	{
		Key: common.TagReverseCharge,
		Name: i18n.String{
			i18n.EN: "Reverse Charge",
			i18n.IT: "Inversione del soggetto passivo",
		},
	},
	{
		Key: common.TagSelfBilled,
		Name: i18n.String{
			i18n.EN: "Self-billed",
			i18n.IT: "Autofattura",
		},
	},
	{
		Key: TagCeilingExceeded,
		Name: i18n.String{
			i18n.EN: "Ceiling exceeded",
			i18n.IT: "Splafonamento",
		},
	},
	{
		Key: TagSanMerinoPaper,
		Name: i18n.String{
			i18n.EN: "Purchases from San Merino with VAT and paper invoice",
			i18n.IT: "Acquisti da San Merino con IVA e fattura cartacea",
		},
	},

	// **** Fiscal Regime Tags ****
	{
		Key: TagMinimumTaxPayers,
		Name: i18n.String{
			i18n.EN: "Minimum Taxpayers",
			i18n.IT: "Contribuenti minimi",
		},
	},
}

var vatTaxTags = []*tax.Tag{
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

var retainedTaxTags = []*tax.Tag{
	{
		Key: TagSelfEmployedHabitual,
		Name: i18n.String{
			i18n.EN: "Self-employed work falling within the habitual practice of an art or profession",
			i18n.IT: "prestazioni di lavoro autonomo rientranti nell'esercizio di arte o professione abituale",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "A",
		},
	},
	{
		Key: TagAuthorIPUsage,
		Name: i18n.String{
			i18n.EN: "Use of intellectual property by the author",
			i18n.IT: "uso di beni immateriali ad opera dell'autore",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "B",
		},
	},
	{
		Key: TagPartnershipAgreements,
		Name: i18n.String{
			i18n.EN: "Partnership agreements in the exercise of an art or profession",
			i18n.IT: "accordi di collaborazione coordinata e continuativa nello svolgimento di attività artistiche o professionali",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "C",
		},
	},
	{
		Key: TagFounderLimitedCompany,
		Name: i18n.String{
			i18n.EN: "Payments made to the founder of a limited company",
			i18n.IT: "versamenti effettuati al socio fondatore di società a responsabilità limitata",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "D",
		},
	},
	{
		Key: TagCertificationDishonoredBills,
		Name: i18n.String{
			i18n.EN: "Certification of dishonored bills and protests",
			i18n.IT: "certificazione dei protesti e delle cambiali disonorate",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "E",
		},
	},
	{
		Key: TagHonoraryJudicialOfficers,
		Name: i18n.String{
			i18n.EN: "Payments made to honorary judges and prosecutors",
			i18n.IT: "versamenti effettuati ai magistrati e ai pubblici ministeri onorari",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "F",
		},
	},
	{
		Key: TagCessationSports,
		Name: i18n.String{
			i18n.EN: "Payments made to sports clubs upon cessation of activity",
			i18n.IT: "versamenti effettuati alle associazioni sportive dilettantistiche in caso di cessazione dell'attività",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "G",
		},
	},
	{
		Key: TagCessationAgency,
		Name: i18n.String{
			i18n.EN: "Payments made to agents upon cessation of activity",
			i18n.IT: "versamenti effettuati agli agenti in caso di cessazione dell'attività",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "H",
		},
	},
	{
		Key: TagCessationNotary,
		Name: i18n.String{
			i18n.EN: "Payments made to notaries upon cessation of activity",
			i18n.IT: "versamenti effettuati ai notai in caso di cessazione dell'attività",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "I",
		},
	},
	{
		Key: TagTruffleGathering,
		Name: i18n.String{
			i18n.EN: "Truffle gathering and sale",
			i18n.IT: "raccolta e vendita di tartufi",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "J",
		},
	},
	{
		Key: TagCivilService,
		Name: i18n.String{
			i18n.EN: "Civil service",
			i18n.IT: "prestazione di servizio civile",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "K",
		},
	},
	{
		Key: TagEntitledIPUsage,
		Name: i18n.String{
			i18n.EN: "Use of intellectual property by entitled parties free of charge",
			i18n.IT: "uso di beni immateriali ad opera di soggetti aventi diritto a titolo gratuito",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "L",
		},
	},
	{
		Key: TagPurchasedIPUsage,
		Name: i18n.String{
			i18n.EN: "Use of intellectual property by parties who paid for it",
			i18n.IT: "uso di beni immateriali ad opera di soggetti aventi diritto a titolo oneroso",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "L1",
		},
	},
	{
		Key: TagOccasionalSelfEmployment,
		Name: i18n.String{
			i18n.EN: "Occasional self-employment",
			i18n.IT: "lavoro autonomo occasionale",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "M",
		},
	},
	{
		Key: TagAssumptionObligations,
		Name: i18n.String{
			i18n.EN: "Assumption of obligations",
			i18n.IT: "assunzione di obbligazioni",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "M1",
		},
	},
	{
		Key: TagENPAPISelfEmployment,
		Name: i18n.String{
			i18n.EN: "ENPAPI self-employment",
			i18n.IT: "lavoro autonomo ENPAPI",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "M2",
		},
	},
	{
		Key: TagAmateurSports,
		Name: i18n.String{
			i18n.EN: "Amateur sports activities",
			i18n.IT: "attività sportive dilettantistiche",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "N",
		},
	},
	{
		Key: TagNonENPAPISelfEmployment,
		Name: i18n.String{
			i18n.EN: "Non-ENPAPI self-employment",
			i18n.IT: "lavoro autonomo non ENPAPI",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "O",
		},
	},
	{
		Key: TagNonENPAPIObligations,
		Name: i18n.String{
			i18n.EN: "Non-ENPAPI obligations",
			i18n.IT: "obbligazioni non ENPAPI",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "O1",
		},
	},
	{
		Key: TagSwissEquipmentsUse,
		Name: i18n.String{
			i18n.EN: "Use of Swiss manufactured equipment",
			i18n.IT: "uso di apparecchiature di fabbricazione svizzera",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "P",
		},
	},
	{
		Key: TagSingleMandateAgent,
		Name: i18n.String{
			i18n.EN: "Single-mandate agent",
			i18n.IT: "agente con mandato singolo",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "Q",
		},
	},
	{
		Key: TagMultiMandateAgent,
		Name: i18n.String{
			i18n.EN: "Multi-mandate agent",
			i18n.IT: "agente con mandato multiplo",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "R",
		},
	},
	{
		Key: TagCommissionAgent,
		Name: i18n.String{
			i18n.EN: "Commission agent",
			i18n.IT: "agente di commissione",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "S",
		},
	},
	{
		Key: TagComissionBroker,
		Name: i18n.String{
			i18n.EN: "Commission broker",
			i18n.IT: "mediatore di commissione",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "T",
		},
	},
	{
		Key: TagBusinessReferrer,
		Name: i18n.String{
			i18n.EN: "Business referrer",
			i18n.IT: "promotore di affari",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "U",
		},
	},
	{
		Key: TagHomeSales,
		Name: i18n.String{
			i18n.EN: "Home sales",
			i18n.IT: "vendita di case",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "V",
		},
	},
	{
		Key: TagOccasionalCommercial,
		Name: i18n.String{
			i18n.EN: "Occasional commercial activity",
			i18n.IT: "attività commerciale occasionale",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "V1",
		},
	},
	{
		Key: TagHomeSalesNonHabitual,
		Name: i18n.String{
			i18n.EN: "Non-habitual home sales",
			i18n.IT: "vendita di case non abituale",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "V2",
		},
	},
	{
		Key: TagContractWork2021,
		Name: i18n.String{
			i18n.EN: "Contract work in 2021",
			i18n.IT: "lavoro a contratto nel 2021",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "W",
		},
	},
	{
		Key: TagEUFees2004,
		Name: i18n.String{
			i18n.EN: "EU fees in 2004",
			i18n.IT: "tariffe UE nel 2004",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "X",
		},
	},
	{
		Key: TagEUFees2005H1,
		Name: i18n.String{
			i18n.EN: "EU fees in the first half of 2005",
			i18n.IT: "tariffe UE nella prima metà del 2005",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "Y",
		},
	},
	{
		Key: TagOtherTitle,
		Name: i18n.String{
			i18n.EN: "Other titles",
			i18n.IT: "altri titoli",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "ZO",
		},
	},
}
