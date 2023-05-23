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
	TagSanMarinoPaper  cbc.Key = "san-marino-paper"

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
	TagSanMarino                  cbc.Key = "san-marino"
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
		Key: TagSanMarinoPaper,
		Name: i18n.String{
			i18n.EN: "Purchases from San Marino with VAT and paper invoice",
			i18n.IT: "Acquisti da San Marino con IVA e fattura cartacea",
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
		Key: TagNotTaxable.With(TagSanMarino),
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
			i18n.EN: "Self-employment services falling within the exercise of habitual art or profession",
			i18n.IT: "Prestazioni di lavoro autonomo rientranti nell'esercizio di arte o professione abituale;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "A",
		},
	},
	{
		Key: TagAuthorIPUsage,
		Name: i18n.String{
			i18n.EN: "Economic use of intellectual works, industrial patents, and processes, formulas or information related to experiences gained in the industrial, commercial or scientific field, by the author or inventor",
			i18n.IT: "Utilizzazione economica, da parte dell'autore o dell'inventore, di opere dell'ingegno, di brevetti industriali e di processi, formule o informazioni relativi ad esperienze acquisite in campo industriale, commerciale o scientifico;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "B",
		},
	},
	{
		Key: TagPartnershipAgreements,
		Name: i18n.String{
			i18n.EN: "Profits deriving from contracts of association in participation and from contracts of co-interest, when the contribution consists exclusively of the provision of labor",
			i18n.IT: "Utili derivanti da contratti di associazione in partecipazione e da contratti di cointeressenza, quando l'apporto è costituito esclusivamente dalla prestazione di lavoro;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "C",
		},
	},
	{
		Key: TagFounderLimitedCompany,
		Name: i18n.String{
			i18n.EN: "Profits due to the promoting partners and founding partners of capital companies",
			i18n.IT: "Utili spettanti ai soci promotori ed ai soci fondatori delle società di capitali;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "D",
		},
	},
	{
		Key: TagCertificationDishonoredBills,
		Name: i18n.String{
			i18n.EN: "Bills of exchange protests levied by municipal secretaries",
			i18n.IT: "Levata di protesti cambiari da parte dei segretari comunali;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "E",
		},
	},
	{
		Key: TagHonoraryJudicialOfficers,
		Name: i18n.String{
			i18n.EN: "Allowances paid to honorary justices of the peace and honorary deputy prosecutors",
			i18n.IT: "Indennità corrisposte ai giudici onorari di pace e ai vice procuratori onorari;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "F",
		},
	},
	{
		Key: TagCessationSports,
		Name: i18n.String{
			i18n.EN: "Allowances paid for the cessation of professional sports activities",
			i18n.IT: "Indennità corrisposte per la cessazione di attività sportiva professionale;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "G",
		},
	},
	{
		Key: TagCessationAgency,
		Name: i18n.String{
			i18n.EN: "Allowances paid for the termination of agency relationships of individuals and partnerships, excluding amounts accrued up to December 31, 2003, already allocated for competence and taxed as business income",
			i18n.IT: "Indennità corrisposte per la cessazione dei rapporti di agenzia delle persone fisiche e delle società di persone con esclusione delle somme maturate entro il 31 dicembre 2003, già imputate per competenza e tassate come reddito d'impresa;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "H",
		},
	},
	{
		Key: TagCessationNotary,
		Name: i18n.String{
			i18n.EN: "Allowances paid for the cessation of notarial functions",
			i18n.IT: "Indennità corrisposte per la cessazione da funzioni notarili;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "I",
		},
	},
	{
		Key: TagTruffleGathering,
		Name: i18n.String{
			i18n.EN: "Fees paid to occasional truffle collectors not identified for value-added tax purposes, in relation to the sale of truffles",
			i18n.IT: "Compensi corrisposti ai raccoglitori occasionali di tartufi non identificati ai fini dell'imposta sul valore aggiunto, in relazione alla cessione di tartufi;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "J",
		},
	},
	{
		Key: TagCivilService,
		Name: i18n.String{
			i18n.EN: "Universal civil service checks referred to in Article 16 of Legislative Decree no. 40 of March 6, 2017",
			i18n.IT: "Assegni di servizio civile universale di cui all'art.16 del d.lgs. n. 40 del 6 marzo 2017;", //nolint:misspell
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "K",
		},
	},
	{
		Key: TagEntitledIPUsage,
		Name: i18n.String{
			i18n.EN: "Income deriving from the economic use of intellectual works, industrial patents, and processes, formulas, and information related to experiences gained in the industrial, commercial or scientific field, which are received by those entitled free of charge (e.g. heirs and legatees of the author and inventor)",
			i18n.IT: "Redditi derivanti dall'utilizzazione economica di opere dell'ingegno, di brevetti industriali e di processi, formule e informazioni relativi a esperienze acquisite in campo industriale, commerciale o scientifico, che sono percepiti dagli aventi causa a titolo gratuito (ad es. eredi e legatari dell'autore e inventore);",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "L",
		},
	},
	{
		Key: TagPurchasedIPUsage,
		Name: i18n.String{
			i18n.EN: "Income deriving from the economic use of intellectual works, industrial patents, and processes, formulas, and information related to experiences gained in the industrial, commercial or scientific field, which are received by subjects who have purchased the rights to their use for valuable consideration",
			i18n.IT: "Redditi derivanti dall'utilizzazione economica di opere dell'ingegno, di brevetti industriali e di processi, formule e informazioni relativi a esperienze acquisite in campo industriale, commerciale o scientifico, che sono percepiti da soggetti che abbiano acquistato a titolo oneroso i diritti alla loro utilizzazione;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "L1",
		},
	},
	{
		Key: TagOccasionalSelfEmployment,
		Name: i18n.String{
			i18n.EN: "Self-employment services not carried out habitually",
			i18n.IT: "Prestazioni di lavoro autonomo non esercitate abitualmente;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "M",
		},
	},
	{
		Key: TagAssumptionObligations,
		Name: i18n.String{
			i18n.EN: "Income deriving from the assumption of obligations to do, not to do, or to allow",
			i18n.IT: "Redditi derivanti dall'assunzione di obblighi di fare, di non fare o permettere;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "M1",
		},
	},
	{
		Key: TagENPAPISelfEmployment,
		Name: i18n.String{
			i18n.EN: "Self-employment services not carried out habitually for which there is an obligation to register with the Separate ENPAPI Management",
			i18n.IT: "Prestazioni di lavoro autonomo non esercitate abitualmente per le quali sussiste l'obbligo di iscrizione alla gestione separata enpapi;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "M2",
		},
	},
	{
		Key: TagAmateurSports,
		Name: i18n.String{
			i18n.EN: "Travel allowances, flat-rate reimbursement of expenses, prizes, and fees paid: - in the direct exercise of amateur sports activities; - in relation to coordinated and continuous collaboration relationships of an administrative-managerial nature, not professional, provided in favor of amateur sports companies and associations, and choirs, bands, and amateur theater groups by the director and technical collaborators",
			i18n.IT: "Indennità di trasferta, rimborso forfetario di spese, premi e compensi erogati: - nell'esercizio diretto di attività sportive dilettantistiche; - in relazione a rapporti di collaborazione coordinata e continuativa di carattere amministrativo-gestionale di natura non professionale resi a favore di società e associazioni sportive dilettantistiche e di cori, bande e filodrammatiche da parte del direttore e dei collaboratori tecnici;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "N",
		},
	},
	{
		Key: TagNonENPAPISelfEmployment,
		Name: i18n.String{
			i18n.EN: "Self-employment services not carried out habitually, for which there is no obligation to register with the separate management (Circ. INPS n. 104/2001)",
			i18n.IT: "Prestazioni di lavoro autonomo non esercitate abitualmente, per le quali non sussiste l'obbligo di iscrizione alla gestione separata (circ. inps n. 104/2001);",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "O",
		},
	},
	{
		Key: TagNonENPAPIObligations,
		Name: i18n.String{
			i18n.EN: "Income deriving from the assumption of obligations to do, not to do, or to allow, for which there is no obligation to register with the separate management (Circ. INPS n. 104/2001)",
			i18n.IT: "Redditi derivanti dall'assunzione di obblighi di fare, di non fare o permettere, per le quali non sussiste l'obbligo di iscrizione alla gestione separata (circ. inps n. 104/2001);",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "O1",
		},
	},
	{
		Key: TagSwissEquipmentsUse,
		Name: i18n.String{
			i18n.EN: "Fees paid to non-resident subjects without a permanent establishment for the use or concession of use of industrial, commercial or scientific equipment located in the State's territory or to Swiss companies or permanent establishments of Swiss companies meeting the requirements of Article 15, paragraph 2 of the Agreement between the European Community and the Swiss Confederation of October 26, 2004 (published in G.U.C.E. of December 29, 2004, no. L385/30)",
			i18n.IT: "Compensi corrisposti a soggetti non residenti privi di stabile organizzazione per l'uso o la concessione in uso di attrezzature industriali, commerciali o scientifiche che si trovano nel territorio dello stato ovvero a società svizzere o stabili organizzazioni di società svizzere che possiedono i requisiti di cui all'art. 15, comma 2 dell'accordo tra la comunità europea e la confederazione svizzera del 26 ottobre 2004 (pubblicato in g.u.c.e. del 29 dicembre 2004 n. l385/30);",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "P",
		},
	},
	{
		Key: TagSingleMandateAgent,
		Name: i18n.String{
			i18n.EN: "Commissions paid to a single-mandate agent or commercial representative",
			i18n.IT: "Provvigioni corrisposte ad agente o rappresentante di commercio monomandatario;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "Q",
		},
	},
	{
		Key: TagMultiMandateAgent,
		Name: i18n.String{
			i18n.EN: "Commissions paid to a multi-mandate agent or commercial representative",
			i18n.IT: "Provvigioni corrisposte ad agente o rappresentante di commercio plurimandatario;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "R",
		},
	},
	{
		Key: TagCommissionAgent,
		Name: i18n.String{
			i18n.EN: "Commissions paid to a commission agent",
			i18n.IT: "Provvigioni corrisposte a commissionario;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "S",
		},
	},
	{
		Key: TagComissionBroker,
		Name: i18n.String{
			i18n.EN: "Commissions paid to a broker",
			i18n.IT: "Provvigioni corrisposte a mediatore;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "T",
		},
	},
	{
		Key: TagBusinessReferrer,
		Name: i18n.String{
			i18n.EN: "Commissions paid to a business finder",
			i18n.IT: "Provvigioni corrisposte a procacciatore di affari;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "U",
		},
	},
	{
		Key: TagHomeSales,
		Name: i18n.String{
			i18n.EN: "Commissions paid to a home sales agent; commissions paid to an agent for door-to-door and street sales of daily newspapers and periodicals (Law of February 25, 1987, no. 67)",
			i18n.IT: "Provvigioni corrisposte a incaricato per le vendite a domicilio; provvigioni corrisposte a incaricato per la vendita porta a porta e per la vendita ambulante di giornali quotidiani e periodici (l. 25 febbraio 1987, n. 67);",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "V",
		},
	},
	{
		Key: TagOccasionalCommercial,
		Name: i18n.String{
			i18n.EN: "Income deriving from non-habitual commercial activities (e.g. commissions paid for occasional services to agents or commercial representatives, brokers, business finders)",
			i18n.IT: "Redditi derivanti da attività commerciali non esercitate abitualmente (ad esempio, provvigioni corrisposte per prestazioni occasionali ad agente o rappresentante di commercio, mediatore, procacciatore d'affari);",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "V1",
		},
	},
	{
		Key: TagHomeSalesNonHabitual,
		Name: i18n.String{
			i18n.EN: "Income from non-habitual services provided by direct home sales agents",
			i18n.IT: "Redditi derivanti dalle prestazioni non esercitate abitualmente rese dagli incaricati alla vendita diretta a domicilio;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "V2",
		},
	},
	{
		Key: TagContractWork2021,
		Name: i18n.String{
			i18n.EN: "Considerations paid in 2021 for services related to subcontracting contracts to which the provisions contained in Article 25-ter of Presidential Decree no. 600 of September 29, 1973, have been applied",
			i18n.IT: "Corrispettivi erogati nel 2021 per prestazioni relative a contratti d'appalto cui si sono resi applicabili le disposizioni contenute nell'art. 25-ter del d.p.r. n. 600 del 29 settembre 1973;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "W",
		},
	},
	{
		Key: TagEUFees2004,
		Name: i18n.String{
			i18n.EN: "Fees paid in 2004 by resident companies or entities or by permanent establishments of foreign companies referred to in Article 26-quater, paragraph 1, letters a) and b) of Presidential Decree 600 of September 29, 1973, to companies or permanent establishments of companies located in another EU Member State meeting the requirements of the aforementioned Article 26-quater of Presidential Decree 600 of September 29, 1973, for which a refund of the withholding tax was made in 2006 pursuant to Article 4 of Legislative Decree no. 143 of May 30, 2005",
			i18n.IT: "Canoni corrisposti nel 2004 da società o enti residenti ovvero da stabili organizzazioni di società estere di cui all'art. 26-quater, comma 1, lett. a) e b) del d.p.r. 600 del 29 settembre 1973, a società o stabili organizzazioni di società, situate in altro stato membro dell'unione europea in presenza dei requisiti di cui al citato art. 26-quater, del d.p.r. 600 del 29 settembre 1973, per i quali è stato effettuato, nell'anno 2006, il rimborso della ritenuta ai sensi dell'art. 4 del d.lgs. 30 maggio 2005 n. 143;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "X",
		},
	},
	{
		Key: TagEUFees2005H1,
		Name: i18n.String{
			i18n.EN: "Fees paid from January 1, 2005, to July 26, 2005, by resident companies or entities or by permanent establishments of foreign companies referred to in Article 26-quater, paragraph 1, letters a) and b) of Presidential Decree no. 600 of September 29, 1973, to companies or permanent establishments of companies located in another EU Member State meeting the requirements of the aforementioned Article 26-quater of Presidential Decree 600 of September 29, 1973, for which a refund of the withholding tax was made in 2006 pursuant to Article 4 of Legislative Decree no. 143 of May 30, 2005",
			i18n.IT: "Canoni corrisposti dal 1° gennaio 2005 al 26 luglio 2005 da società o enti residenti ovvero da stabili organizzazioni di società estere di cui all'art. 26-quater, comma 1, lett. a) e b) del d.p.r. n. 600 del 29 settembre 1973, a società o stabili organizzazioni di società, situate in altro stato membro dell'unione europea in presenza dei requisiti di cui al citato art. 26-quater, del d.p.r. n. 600 del 29 settembre 1973, per i quali è stato effettuato, nell'anno 2006, il rimborso della ritenuta ai sensi dell'art. 4 del d.lgs. 30 maggio 2005 n. 143;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "Y",
		},
	},
	{
		Key: TagOtherTitle,
		Name: i18n.String{
			i18n.EN: "Different title from the previous ones",
			i18n.IT: "Titolo diverso dai precedenti;",
		},
		Meta: cbc.Meta{
			KeyFatturaPACausalePagamento: "ZO",
		},
	},
}
