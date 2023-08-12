package it

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

// Local tax category definitions which are not considered standard.
// There is a 6th retained tax type, RT06 "Other contributions", which is
// currently not supported.
const (
	// https://www.agenziaentrate.gov.it/portale/imposta-sul-reddito-delle-persone-fisiche-irpef-/aliquote-e-calcolo-dell-irpef
	TaxCategoryIRPEF    cbc.Code = "IRPEF"
	TaxCategoryIRES     cbc.Code = "IRES"
	TaxCategoryINPS     cbc.Code = "INPS"
	TaxCategoryENASARCO cbc.Code = "ENASARCO"
	TaxCategoryENPAM    cbc.Code = "ENPAM"
)

var categories = []*tax.Category{
	{
		Code:     common.TaxCategoryVAT,
		Retained: false,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.IT: "IVA",
		},
		Desc: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.IT: "Imposta sul Valore Aggiunto",
		},
		Rates: []*tax.Rate{
			{
				Key: common.TaxRateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.IT: "Aliquota Zero",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Key: common.TaxRateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Minimum Rate",
					i18n.IT: "Aliquota Minima",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(40, 3),
					},
				},
			},
			{
				Key: common.TaxRateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.IT: "Aliquota Ridotta",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(50, 3),
					},
				},
			},
			{
				Key: common.TaxRateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate Rate",
					i18n.IT: "Aliquota Intermedia",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(100, 3),
					},
				},
			},
			{
				Key: common.TaxRateStandard,
				Name: i18n.String{
					i18n.EN: "Ordinary Rate",
					i18n.IT: "Aliquota Ordinaria",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(220, 3),
					},
				},
			},
			{
				Key:    common.TaxRateExempt,
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "None",
					i18n.IT: "Natura",
				},
				Codes: []*tax.CodeDefinition{
					{
						Code: "N1",
						Name: i18n.String{
							i18n.EN: "Excluded pursuant to Art. 15, DPR 633/72",
							i18n.IT: "Escluse ex. art. 15 del D.P.R. 633/1972",
						},
					},
					{
						Code: "N2.1",
						Name: i18n.String{
							i18n.EN: "Not subject pursuant to Art. 7, DPR 633/72",
							i18n.IT: "Non soggette ex. art. 7 del D.P.R. 633/72",
						},
					},
					{
						Code: "N2.2",
						Name: i18n.String{
							i18n.EN: "Not subject - other",
							i18n.IT: "Non soggette - altri casi",
						},
					},
					{
						Code: "N3.1",
						Name: i18n.String{
							i18n.EN: "Not taxable - exports",
							i18n.IT: "Non imponibili - esportazioni",
						},
					},
					{
						Code: "N3.2",
						Name: i18n.String{
							i18n.EN: "Not taxable - intra-community supplies",
							i18n.IT: "Non imponibili - cessioni intracomunitarie",
						},
					},
					{
						Code: "N3.3",
						Name: i18n.String{
							i18n.EN: "Not taxable - transfers to San Marino",
							i18n.IT: "Non imponibili - cessioni verso San Marino",
						},
					},
					{
						Code: "N3.4",
						Name: i18n.String{
							i18n.EN: "Not taxable - export supplies of goods and services",
							i18n.IT: "Non Imponibili - operazioni assimilate alle cessioni all'esportazione",
						},
					},
					{
						Code: "N3.5",
						Name: i18n.String{
							i18n.EN: "Not taxable - declaration of intent",
							i18n.IT: "Non imponibili - dichiarazioni d'intento",
						},
					},
					{
						Code: "N3.6",
						Name: i18n.String{
							i18n.EN: "Not taxable - other",
							i18n.IT: "Non imponibili - altre operazioni che non concorrono alla formazione del plafond",
						},
					},
					{
						Code: "N4",
						Name: i18n.String{
							i18n.EN: "Exempt",
							i18n.IT: "Esenti",
						},
					},
					{
						Code: "N5",
						Name: i18n.String{
							i18n.EN: "Margin regime / VAT not exposed",
							i18n.IT: "Regime del margine/IVA non esposta in fattura",
						},
					},
					{
						Code: "N6.1",
						Name: i18n.String{
							i18n.EN: "Reverse charge - Transfer of scrap and of other recyclable materials",
							i18n.IT: "Inversione contabile - cessione di rottami e altri materiali di recupero",
						},
					},
					{
						Code: "N6.2",
						Name: i18n.String{
							i18n.EN: "Reverse charge - Transfer of gold and pure silver pursuant to law 7/2000 as well as used jewelery to OPO",
							i18n.IT: "Inversione contabile - cessione di oro e argento ai sensi della legge 7/2000 nonché di oreficeria usata ad OPO",
						},
					},
					{
						Code: "N6.3",
						Name: i18n.String{
							i18n.EN: "Reverse charge - Construction subcontracting",
							i18n.IT: "Inversione contabile - subappalto nel settore edile",
						},
					},
					{
						Code: "N6.4",
						Name: i18n.String{
							i18n.EN: "Reverse charge - Transfer of buildings",
							i18n.IT: "Inversione contabile - cessione di fabbricati",
						},
					},
					{
						Code: "N6.5",
						Name: i18n.String{
							i18n.EN: "Reverse charge - Transfer of mobile phones",
							i18n.IT: "Inversione contabile - cessione di telefoni cellulari",
						},
					},
					{
						Code: "N6.6",
						Name: i18n.String{
							i18n.EN: "Reverse charge - Transfer of electronic products",
							i18n.IT: "Inversione contabile - cessione di prodotti elettronici",
						},
					},
					{
						Code: "N6.7",
						Name: i18n.String{
							i18n.EN: "Reverse charge - provisions in the construction and related sectors",
							i18n.IT: "Inversione contabile - prestazioni comparto edile e settori connessi",
						},
					},
					{
						Code: "N6.8",
						Name: i18n.String{
							i18n.EN: "Reverse charge - transactions in the energy sector",
							i18n.IT: "Inversione contabile - operazioni settore energetico",
						},
					},
					{
						Code: "N6.9",
						Name: i18n.String{
							i18n.EN: "Reverse charge - other cases",
							i18n.IT: "Inversione contabile - altri casi",
						},
					},
					{
						Code: "N7",
						Name: i18n.String{
							i18n.EN: "VAT paid in other EU countries (telecommunications, tele-broadcasting and electronic services provision pursuant to Art. 7 -octies letter a, b, art. 74-sexies Italian Presidential Decree 633/72)",
							i18n.IT: "IVA assolta in altro stato UE (prestazione di servizi di telecomunicazioni, tele-radiodiffusione ed elettronici ex art. 7-octies lett. a, b, art. 74-sexies DPR 633/72)",
						},
					},
				},
			},
		},
	},
	{
		// IT: https://www.agenziaentrate.gov.it/portale/imposta-sul-reddito-delle-persone-fisiche-irpef-/aliquote-e-calcolo-dell-irpef
		// EN: https://www.agenziaentrate.gov.it/portale/web/english/information-for-specific-categories-of-workers
		Code:     TaxCategoryIRPEF,
		Retained: true,
		Name: i18n.String{
			i18n.EN: "IRPEF",
			i18n.IT: "IRPEF",
		},
		Desc: i18n.String{
			i18n.EN: "Personal Income Tax",
			i18n.IT: "Imposta sul Reddito delle Persone Fisiche",
		},
		Map: cbc.CodeSet{
			KeyFatturaPATipoRitenuta: "RT01",
		},
		RateCodes: retainedTaxRateCodes,
	},
	{
		Code:     TaxCategoryIRES,
		Retained: true,
		Name: i18n.String{
			i18n.EN: "IRES",
			i18n.IT: "IRES",
		},
		Desc: i18n.String{
			i18n.EN: "Corporate Income Tax",
			i18n.IT: "Imposta sul Reddito delle Società",
		},
		Map: cbc.CodeSet{
			KeyFatturaPATipoRitenuta: "RT02",
		},
		RateCodes: retainedTaxRateCodes,
	},
	{
		Code:     TaxCategoryINPS,
		Retained: true,
		Name: i18n.String{
			i18n.EN: "INPS Contribution",
			i18n.IT: "Contributo INPS", // nolint:misspell
		},
		Desc: i18n.String{
			i18n.EN: "Contribution to the National Social Security Institute",
			i18n.IT: "Contributo Istituto Nazionale della Previdenza Sociale", // nolint:misspell
		},
		RateCodes: retainedTaxRateCodes,
		Map: cbc.CodeSet{
			KeyFatturaPATipoRitenuta: "RT03",
		},
	},
	{
		Code:     TaxCategoryENASARCO,
		Retained: true,
		Name: i18n.String{
			i18n.EN: "ENASARCO Contribution",
			i18n.IT: "Contributo ENASARCO", // nolint:misspell
		},
		Desc: i18n.String{
			i18n.EN: "Contribution to the National Welfare Board for Sales Agents and Representatives",
			i18n.IT: "Contributo Ente Nazionale Assistenza Agenti e Rappresentanti di Commercio", // nolint:misspell
		},
		RateCodes: retainedTaxRateCodes,
		Map: cbc.CodeSet{
			KeyFatturaPATipoRitenuta: "RT04",
		},
	},
	{
		Code:     TaxCategoryENPAM,
		Retained: true,
		Name: i18n.String{
			i18n.EN: "ENPAM Contribution",
			i18n.IT: "Contributo ENPAM", // nolint:misspell
		},
		Desc: i18n.String{
			i18n.EN: "Contribution to the National Pension and Welfare Board for Doctors",
			i18n.IT: "Contributo - Ente Nazionale Previdenza e Assistenza Medici", // nolint:misspell
		},
		RateCodes: retainedTaxRateCodes,
		Map: cbc.CodeSet{
			KeyFatturaPATipoRitenuta: "RT05",
		},
	},
}

// Retained tax tag keys determined from the "CausalePagamento" field from FatturaPA.
// Source: https://www.agenziaentrate.gov.it/portale/documents/20143/4115385/CU_istr_2022.pdf
// Section VII, Part 2
var retainedTaxRateCodes = []*tax.CodeDefinition{
	{
		Code: "A",
		Name: i18n.String{
			i18n.EN: "Self-employment services falling within the exercise of habitual art or profession",
			i18n.IT: "Prestazioni di lavoro autonomo rientranti nell'esercizio di arte o professione abituale;",
		},
	},
	{
		Code: "B",
		Name: i18n.String{
			i18n.EN: "Economic use of intellectual works, industrial patents, and processes, formulas or information related to experiences gained in the industrial, commercial or scientific field, by the author or inventor",
			i18n.IT: "Utilizzazione economica, da parte dell'autore o dell'inventore, di opere dell'ingegno, di brevetti industriali e di processi, formule o informazioni relativi ad esperienze acquisite in campo industriale, commerciale o scientifico;",
		},
	},
	{
		Code: "C",
		Name: i18n.String{
			i18n.EN: "Profits deriving from contracts of association in participation and from contracts of co-interest, when the contribution consists exclusively of the provision of labor",
			i18n.IT: "Utili derivanti da contratti di associazione in partecipazione e da contratti di cointeressenza, quando l'apporto è costituito esclusivamente dalla prestazione di lavoro;",
		},
	},
	{
		Code: "D",
		Name: i18n.String{
			i18n.EN: "Profits due to the promoting partners and founding partners of capital companies",
			i18n.IT: "Utili spettanti ai soci promotori ed ai soci fondatori delle società di capitali;",
		},
	},
	{
		Code: "E",
		Name: i18n.String{
			i18n.EN: "Bills of exchange protests levied by municipal secretaries",
			i18n.IT: "Levata di protesti cambiari da parte dei segretari comunali;",
		},
	},
	{
		Code: "F",
		Name: i18n.String{
			i18n.EN: "Allowances paid to honorary justices of the peace and honorary deputy prosecutors",
			i18n.IT: "Indennità corrisposte ai giudici onorari di pace e ai vice procuratori onorari;",
		},
	},
	{
		Code: "G",
		Name: i18n.String{
			i18n.EN: "Allowances paid for the cessation of professional sports activities",
			i18n.IT: "Indennità corrisposte per la cessazione di attività sportiva professionale",
		},
	},
	{
		Code: "H",
		Name: i18n.String{
			i18n.EN: "Allowances paid for the termination of agency relationships of individuals and partnerships, excluding amounts accrued up to December 31, 2003, already allocated for competence and taxed as business income",
			i18n.IT: "Indennità corrisposte per la cessazione dei rapporti di agenzia delle persone fisiche e delle società di persone con esclusione delle somme maturate entro il 31 dicembre 2003, già imputate per competenza e tassate come reddito d'impresa",
		},
	},
	{
		Code: "I",
		Name: i18n.String{
			i18n.EN: "Allowances paid for the cessation of notarial functions",
			i18n.IT: "Indennità corrisposte per la cessazione da funzioni notarili",
		},
	},
	{
		Code: "J",
		Name: i18n.String{
			i18n.EN: "Fees paid to occasional truffle collectors not identified for value-added tax purposes, in relation to the sale of truffles",
			i18n.IT: "Compensi corrisposti ai raccoglitori occasionali di tartufi non identificati ai fini dell'imposta sul valore aggiunto, in relazione alla cessione di tartufi",
		},
	},
	{
		Code: "K",
		Name: i18n.String{
			i18n.EN: "Universal civil service checks referred to in Article 16 of Legislative Decree no. 40 of March 6, 2017",
			i18n.IT: "Assegni di servizio civile universale di cui all'art.16 del d.lgs. n. 40 del 6 marzo 2017", //nolint:misspell
		},
	},
	{
		Code: "L",
		Name: i18n.String{
			i18n.EN: "Income deriving from the economic use of intellectual works, industrial patents, and processes, formulas, and information related to experiences gained in the industrial, commercial or scientific field, which are received by those entitled free of charge (e.g. heirs and legatees of the author and inventor)",
			i18n.IT: "Redditi derivanti dall'utilizzazione economica di opere dell'ingegno, di brevetti industriali e di processi, formule e informazioni relativi a esperienze acquisite in campo industriale, commerciale o scientifico, che sono percepiti dagli aventi causa a titolo gratuito (ad es. eredi e legatari dell'autore e inventore)",
		},
	},
	{
		Code: "L1",
		Name: i18n.String{
			i18n.EN: "Income deriving from the economic use of intellectual works, industrial patents, and processes, formulas, and information related to experiences gained in the industrial, commercial or scientific field, which are received by subjects who have purchased the rights to their use for valuable consideration",
			i18n.IT: "Redditi derivanti dall'utilizzazione economica di opere dell'ingegno, di brevetti industriali e di processi, formule e informazioni relativi a esperienze acquisite in campo industriale, commerciale o scientifico, che sono percepiti da soggetti che abbiano acquistato a titolo oneroso i diritti alla loro utilizzazione",
		},
	},
	{
		Code: "M",
		Name: i18n.String{
			i18n.EN: "Self-employment services not carried out habitually",
			i18n.IT: "Prestazioni di lavoro autonomo non esercitate abitualmente",
		},
	},
	{
		Code: "M1",
		Name: i18n.String{
			i18n.EN: "Income deriving from the assumption of obligations to do, not to do, or to allow",
			i18n.IT: "Redditi derivanti dall'assunzione di obblighi di fare, di non fare o permettere",
		},
	},
	{
		Code: "M2",
		Name: i18n.String{
			i18n.EN: "Self-employment services not carried out habitually for which there is an obligation to register with the Separate ENPAPI Management",
			i18n.IT: "Prestazioni di lavoro autonomo non esercitate abitualmente per le quali sussiste l'obbligo di iscrizione alla gestione separata enpapi",
		},
	},
	{
		Code: "N",
		Name: i18n.String{
			i18n.EN: "Travel allowances, flat-rate reimbursement of expenses, prizes, and fees paid: - in the direct exercise of amateur sports activities; - in relation to coordinated and continuous collaboration relationships of an administrative-managerial nature, not professional, provided in favor of amateur sports companies and associations, and choirs, bands, and amateur theater groups by the director and technical collaborators",
			i18n.IT: "Indennità di trasferta, rimborso forfetario di spese, premi e compensi erogati: - nell'esercizio diretto di attività sportive dilettantistiche; - in relazione a rapporti di collaborazione coordinata e continuativa di carattere amministrativo-gestionale di natura non professionale resi a favore di società e associazioni sportive dilettantistiche e di cori, bande e filodrammatiche da parte del direttore e dei collaboratori tecnici;",
		},
	},
	{
		Code: "O",
		Name: i18n.String{
			i18n.EN: "Self-employment services not carried out habitually, for which there is no obligation to register with the separate management (Circ. INPS n. 104/2001)",
			i18n.IT: "Prestazioni di lavoro autonomo non esercitate abitualmente, per le quali non sussiste l'obbligo di iscrizione alla gestione separata (circ. inps n. 104/2001)",
		},
	},
	{
		Code: "O1",
		Name: i18n.String{
			i18n.EN: "Income deriving from the assumption of obligations to do, not to do, or to allow, for which there is no obligation to register with the separate management (Circ. INPS n. 104/2001)",
			i18n.IT: "Redditi derivanti dall'assunzione di obblighi di fare, di non fare o permettere, per le quali non sussiste l'obbligo di iscrizione alla gestione separata (circ. inps n. 104/2001)",
		},
	},
	{
		Code: "P",
		Name: i18n.String{
			i18n.EN: "Fees paid to non-resident subjects without a permanent establishment for the use or concession of use of industrial, commercial or scientific equipment located in the State's territory or to Swiss companies or permanent establishments of Swiss companies meeting the requirements of Article 15, paragraph 2 of the Agreement between the European Community and the Swiss Confederation of October 26, 2004 (published in G.U.C.E. of December 29, 2004, no. L385/30)",
			i18n.IT: "Compensi corrisposti a soggetti non residenti privi di stabile organizzazione per l'uso o la concessione in uso di attrezzature industriali, commerciali o scientifiche che si trovano nel territorio dello stato ovvero a società svizzere o stabili organizzazioni di società svizzere che possiedono i requisiti di cui all'art. 15, comma 2 dell'accordo tra la comunità europea e la confederazione svizzera del 26 ottobre 2004 (pubblicato in g.u.c.e. del 29 dicembre 2004 n. l385/30)",
		},
	},
	{
		Code: "Q",
		Name: i18n.String{
			i18n.EN: "Commissions paid to a single-mandate agent or commercial representative",
			i18n.IT: "Provvigioni corrisposte ad agente o rappresentante di commercio monomandatario",
		},
	},
	{
		Code: "R",
		Name: i18n.String{
			i18n.EN: "Commissions paid to a multi-mandate agent or commercial representative",
			i18n.IT: "Provvigioni corrisposte ad agente o rappresentante di commercio plurimandatario",
		},
	},
	{
		Code: "S",
		Name: i18n.String{
			i18n.EN: "Commissions paid to a commission agent",
			i18n.IT: "Provvigioni corrisposte a commissionario",
		},
	},
	{
		Code: "T",
		Name: i18n.String{
			i18n.EN: "Commissions paid to a broker",
			i18n.IT: "Provvigioni corrisposte a mediatore",
		},
	},
	{
		Code: "U",
		Name: i18n.String{
			i18n.EN: "Commissions paid to a business finder",
			i18n.IT: "Provvigioni corrisposte a procacciatore di affari",
		},
	},
	{
		Code: "V",
		Name: i18n.String{
			i18n.EN: "Commissions paid to a home sales agent; commissions paid to an agent for door-to-door and street sales of daily newspapers and periodicals (Law of February 25, 1987, no. 67)",
			i18n.IT: "Provvigioni corrisposte a incaricato per le vendite a domicilio; provvigioni corrisposte a incaricato per la vendita porta a porta e per la vendita ambulante di giornali quotidiani e periodici (l. 25 febbraio 1987, n. 67);",
		},
	},
	{
		Code: "V1",
		Name: i18n.String{
			i18n.EN: "Income deriving from non-habitual commercial activities (e.g. commissions paid for occasional services to agents or commercial representatives, brokers, business finders)",
			i18n.IT: "Redditi derivanti da attività commerciali non esercitate abitualmente (ad esempio, provvigioni corrisposte per prestazioni occasionali ad agente o rappresentante di commercio, mediatore, procacciatore d'affari);",
		},
	},
	{
		Code: "V2",
		Name: i18n.String{
			i18n.EN: "Income from non-habitual services provided by direct home sales agents",
			i18n.IT: "Redditi derivanti dalle prestazioni non esercitate abitualmente rese dagli incaricati alla vendita diretta a domicilio;",
		},
	},
	{
		Code: "W",
		Name: i18n.String{
			i18n.EN: "Considerations paid in 2021 for services related to subcontracting contracts to which the provisions contained in Article 25-ter of Presidential Decree no. 600 of September 29, 1973, have been applied",
			i18n.IT: "Corrispettivi erogati nel 2021 per prestazioni relative a contratti d'appalto cui si sono resi applicabili le disposizioni contenute nell'art. 25-ter del d.p.r. n. 600 del 29 settembre 1973;",
		},
	},
	{
		Code: "X",
		Name: i18n.String{
			i18n.EN: "Fees paid in 2004 by resident companies or entities or by permanent establishments of foreign companies referred to in Article 26-quater, paragraph 1, letters a) and b) of Presidential Decree 600 of September 29, 1973, to companies or permanent establishments of companies located in another EU Member State meeting the requirements of the aforementioned Article 26-quater of Presidential Decree 600 of September 29, 1973, for which a refund of the withholding tax was made in 2006 pursuant to Article 4 of Legislative Decree no. 143 of May 30, 2005",
			i18n.IT: "Canoni corrisposti nel 2004 da società o enti residenti ovvero da stabili organizzazioni di società estere di cui all'art. 26-quater, comma 1, lett. a) e b) del d.p.r. 600 del 29 settembre 1973, a società o stabili organizzazioni di società, situate in altro stato membro dell'unione europea in presenza dei requisiti di cui al citato art. 26-quater, del d.p.r. 600 del 29 settembre 1973, per i quali è stato effettuato, nell'anno 2006, il rimborso della ritenuta ai sensi dell'art. 4 del d.lgs. 30 maggio 2005 n. 143;",
		},
	},
	{
		Code: "Y",
		Name: i18n.String{
			i18n.EN: "Fees paid from January 1, 2005, to July 26, 2005, by resident companies or entities or by permanent establishments of foreign companies referred to in Article 26-quater, paragraph 1, letters a) and b) of Presidential Decree no. 600 of September 29, 1973, to companies or permanent establishments of companies located in another EU Member State meeting the requirements of the aforementioned Article 26-quater of Presidential Decree 600 of September 29, 1973, for which a refund of the withholding tax was made in 2006 pursuant to Article 4 of Legislative Decree no. 143 of May 30, 2005",
			i18n.IT: "Canoni corrisposti dal 1° gennaio 2005 al 26 luglio 2005 da società o enti residenti ovvero da stabili organizzazioni di società estere di cui all'art. 26-quater, comma 1, lett. a) e b) del d.p.r. n. 600 del 29 settembre 1973, a società o stabili organizzazioni di società, situate in altro stato membro dell'unione europea in presenza dei requisiti di cui al citato art. 26-quater, del d.p.r. n. 600 del 29 settembre 1973, per i quali è stato effettuato, nell'anno 2006, il rimborso della ritenuta ai sensi dell'art. 4 del d.lgs. 30 maggio 2005 n. 143;",
		},
	},
	{
		Code: "ZO",
		Name: i18n.String{
			i18n.EN: "Different title from the previous ones",
			i18n.IT: "Titolo diverso dai precedenti;",
		},
	},
}
