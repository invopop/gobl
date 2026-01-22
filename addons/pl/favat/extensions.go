package favat

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Regime extension codes for local electronic formats.
const (
	ExtKeyTaxCategory           cbc.Key = "pl-favat-tax-category"
	ExtKeyEffectiveDate         cbc.Key = "pl-favat-effective-date"
	ExtKeyPaymentMeans          cbc.Key = "pl-favat-payment-means"           // for mapping to TFormaPlatnosci's codes (type of payment means - e.g. cash, bank transfer etc)
	ExtKeyInvoiceType           cbc.Key = "pl-favat-invoice-type"            // for mapping to TRodzajFaktury's codes (type of invoice - e.g. regular, in advance, correction etc)
	ExtKeyCashAccounting        cbc.Key = "pl-favat-cash-accounting"         // for mapping to P_16 field, indicating cash accounting
	ExtKeySelfBilling           cbc.Key = "pl-favat-self-billing"            // for mapping to P_17 field, indicating self-invoicing
	ExtKeyReverseCharge         cbc.Key = "pl-favat-reverse-charge"          // for mapping to P_18, indicating reverse charge
	ExtKeySplitPaymentMechanism cbc.Key = "pl-favat-split-payment-mechanism" // for mapping to P_18A, indicating split payment mechanism
	ExtKeyExemption             cbc.Key = "pl-favat-exemption"               // for mapping to P_19 and its subfields (P_19A, P_19B, P_19C), indicating exemption
	ExtKeyMarginScheme          cbc.Key = "pl-favat-margin-scheme"           // for mapping to P_PMarzy, indicating margin scheme
	ExtKeyJST                   cbc.Key = "pl-favat-jst"                     // for mapping to JST (Subordinate Local Government Unit)
	ExtKeyGroupVAT              cbc.Key = "pl-favat-group-vat"               // for mapping to GV (Group VAT member)
	ExtKeyThirdPartyRole        cbc.Key = "pl-favat-third-party-role"        // for third party roles (issuer, GV member, JST member)
)

var extensionKeys = []*cbc.Definition{
	{
		Key: ExtKeyTaxCategory,
		Name: i18n.String{
			i18n.EN: "Tax categories for KSeF",
			i18n.PL: "Kategorie podatkowe dla KSeF",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The ~pl-favat-tax-category~ extension specifies tax categories for Polish FA_VAT/KSeF invoices.
				Each category corresponds to a field in the FA_VAT XML schema (P_13_*, P_14_*, etc.).

				This extension is used at the **line tax level** (~lines[].taxes[]~) and is automatically
				normalized by GOBL based on the tax combo ~key~ and ~rate~ fields. The extension is required
				for all line items.

				Automatic mapping from GOBL tax combos:

				| GOBL Tax Key          | GOBL Rate        | FA_VAT Category | VAT % |
				| --------------------- | ---------------- | --------------- | ----- |
				| ~standard~            | ~general~        | 1               | 23%   |
				| ~standard~            | ~reduced~        | 2               | 8%    |
				| ~standard~            | ~super-reduced~  | 3               | 5%    |
				| ~zero~                | -                | 6.1             | 0%    |
				| ~intra-community~     | -                | 6.2             | 0%    |
				| ~export~              | -                | 6.3             | 0%    |
				| ~exempt~              | -                | 7               | -     |
				| ~outside-scope~       | -                | 8               | -     |
				| ~reverse-charge~      | -                | 9               | -     |

				Example:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					// ...
					"lines": [
						{
							// ...
							"taxes": [
								{
									"cat": "VAT",
									"key": "standard",
									"rate": "general",
									"percent": "23.0%",
									// Extension automatically set to "1"
									"ext": {
										"pl-favat-tax-category": "1"
									}
								}
							]
						}
					]
				}
				~~~

				The ones that are not automatically mapped can be set manually.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Base rate",
					i18n.PL: "Stawka podstawowa",
				},
				Desc: i18n.String{
					i18n.EN: "Sales of goods and services subject to VAT rate of 23%.",
					i18n.PL: "Sprzedaż towarów i usług objętych stawką podatku VAT 23%.",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "First reduced rate",
					i18n.PL: "Stawka obniżona pierwsza",
				},
				Desc: i18n.String{
					i18n.EN: "Sales of goods and services subject to VAT rate of 8%.",
					i18n.PL: "Sprzedaż towarów i usług objętych stawką podatku VAT 8%.",
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "Second reduced rate",
					i18n.PL: "Stawka obniżona druga",
				},
				Desc: i18n.String{
					i18n.EN: "Sales of goods and services subject to VAT rate of 5%.",
					i18n.PL: "Sprzedaż towarów i usług objętych stawką podatku VAT 5%.",
				},
			},
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "Taxi Rate",
					i18n.PL: "Ryczałt dla taksówek",
				},
				Desc: i18n.String{
					i18n.EN: "Special flat rate for taxi drivers.",
					i18n.PL: "Specjalna stawka ryczałtu dla taksówkarzy.",
				},
			},
			{
				Code: "5",
				Name: i18n.String{
					i18n.EN: "OSS (one stop shop)",
					i18n.PL: "Punkt kompleksowej obsługi (OSS)",
				},
				Desc: i18n.String{
					i18n.EN: "Special European Union procedure for the supply of certain goods and services",
					i18n.PL: "Specjalna procedura unijna dla sprzedaży niektórych towarów i usług",
				},
			},
			{
				Code: "6.1",
				Name: i18n.String{
					i18n.EN: "0% VAT",
					i18n.PL: "0% WDT",
				},
				Desc: i18n.String{
					i18n.EN: "Zero VAT excluding intra-community supply of goods and export",
					i18n.PL: "Zerowa stawka podatku VAT, z wyłączeniem wewnątrzwspólnotowej dostawy towarów i eksportu",
				},
			},
			{
				Code: "6.2",
				Name: i18n.String{
					i18n.EN: "0% VAT for intra-community supply of goods",
					i18n.PL: "0% VAT dla dostawy towarów wewnątrzwspólnotowych",
				},
			},
			{
				Code: "6.3",
				Name: i18n.String{
					i18n.EN: "0% VAT for export",
					i18n.PL: "0% VAT dla eksportu",
				},
				Desc: i18n.String{
					i18n.EN: "Zero VAT, export outside the European Union",
					i18n.PL: "Zerowa stawka podatku VAT, eksport poza Unii Europejskiej",
				},
			},
			{
				Code: "7",
				Name: i18n.String{
					i18n.EN: "Exempt",
					i18n.PL: "Zwolnienie",
				},
				Desc: i18n.String{
					i18n.EN: "Sales exempt from VAT",
					i18n.PL: "Sprzedaż zwolniona od podatku VAT",
				},
			},
			{
				Code: "8",
				Name: i18n.String{
					i18n.EN: "Foreign sales outside scope of Polish VAT",
					i18n.PL: "Sprzedaż zagraniczna poza zakresem polskiego VAT",
				},
				Desc: i18n.String{
					i18n.EN: "Total value of sales for the supply of goods and provision of services outside the country, excluding the amounts shown for values 5 and 9",
					i18n.PL: "Wartość sprzedaży towarów oraz świadczenia usług poza terytorium kraju, z wyłączeniem wartości pokazanych dla wartości 5 i 9",
				},
			},
			{
				Code: "9",
				Name: i18n.String{
					i18n.EN: "EU Reverse Charge",
					i18n.PL: "Odwrotne obciążenie UE",
				},
				Desc: i18n.String{
					i18n.EN: "Services provided within the European Union, where services are taxed in the customer's country",
					i18n.PL: "Świadczenia usług w Unii Europejskiej, gdzie usługi są opodatkowane w kraju nabywcy (VAT UE)",
				},
			},
			{
				// https://poradnikprzedsiebiorcy.pl/-odwrotne-obciazenie
				Code: "10",
				Name: i18n.String{
					i18n.EN: "Domestic Reverse charge",
					i18n.PL: "Odwrotne obciążenie krajowe",
				},
				Desc: i18n.String{
					i18n.EN: "Obligation to account for tax by the purchaser",
					i18n.PL: "Obowiązek rozliczenia podatku przez nabywcę",
				},
			},
			{
				Code: "11",
				Name: i18n.String{
					i18n.EN: "Margin scheme",
					i18n.PL: "Schemat marży",
				},
				// https://poradnikprzedsiebiorcy.pl/-faktura-vat-marza-kiedy-mozna-ja-wystawic
				Desc: i18n.String{
					i18n.EN: "Tax only on the margin, which the seller has calculated - applies to tourism services, sale of used goods, antiques, works of art",
					i18n.PL: "Podatek tylko od marży, którą naliczył sprzedawca - dotyczy usług turystycznych, sprzedaży towarów używanych, przedmiotów kolekcjonerskich i antyków, dzieł sztuki",
				},
			},
		},
	},
	{
		Key: ExtKeyInvoiceType,
		Name: i18n.String{
			i18n.EN: "Invoice type code for KSeF",
			i18n.PL: "Kod rodzaju faktury dla KSeF",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The ~pl-favat-invoice-type~ extension specifies the type of invoice for Polish FA_VAT/KSeF.
				This extension is used in the **invoice tax section** (~tax.ext~) and is automatically
				normalized by GOBL based on the invoice type and tags during the ~scenarios~ normalization step.

				Automatic mapping from GOBL invoice structure:

				| Code    | GOBL invoice type | Tags         | Description          |
				|---------|-------------------|--------------|----------------------|
				| VAT     | ~standard~        | -            | Regular invoice      |
				| ZAL     | ~standard~        | ~partial~    | Prepayment invoice   |
				| ROZ     | ~standard~        | ~settlement~ | Settlement invoice   |
				| UPR     | ~standard~        | ~simplified~ | Simplified invoice   |
				| KOR     | ~credit-note~     | -            | Credit note          |
				| KOR_ZAL | ~credit-note~     | ~partial~    | Prepayment credit    |
				| KOR_ROZ | ~credit-note~     | ~settlement~ | Settlement credit    |

				Example of a prepayment invoice:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					"$addons": ["pl-favat-v3"],
					"$tags": ["partial"],
					"type": "standard",
					// ...
					"tax": {
						// Extension automatically set to "ZAL"
						"ext": {
							"pl-favat-invoice-type": "ZAL"
						}
					}
				}
				~~~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "VAT",
				Name: i18n.String{
					i18n.EN: "Regular Invoice",
					i18n.PL: "Faktura Podstawowa",
				},
				Desc: i18n.String{
					i18n.EN: "Most commonly used invoice type.",
					i18n.PL: "Najczęściej używany typ faktury.",
				},
			},
			{
				Code: "ZAL",
				Name: i18n.String{
					i18n.EN: "Prepayment Invoice",
					i18n.PL: `Faktura Zaliczkowa`,
				},
				Desc: i18n.String{
					i18n.EN: "Invoice documenting receipt of payment or part of payment before performing the act, as well as an invoice issued in connection with article 106f paragraph 4 of the act (advance invoice)",
					i18n.PL: "Faktura dokumentująca otrzymanie zapłaty lub jej części przed dokonaniem czynności oraz faktura wystawiona w związku z art. 106f ust. 4 ustawy (faktura zaliczkowa)",
				},
			},
			{
				Code: "ROZ",
				Name: i18n.String{
					i18n.EN: "Settlement Invoice",
					i18n.PL: "Faktura Rozliczeniowa",
				},
				Desc: i18n.String{
					i18n.EN: "Invoice issued in connection with article 106f paragraph 3 of the act. Final invoice, which concludes the contract or transaction, includes all previous advance invoices and contains a full list of payments and transaction amount.",
					i18n.PL: "Faktura wystawiona w związku z art. 106f ust. 3 ustawy. Kończy umowę lub transakcję, uwzględnia wszystkie poprzednie faktury zaliczkowe, zawiera pełną listę płatności i kwotę transakcji.",
				},
			},
			{
				Code: "UPR",
				Name: i18n.String{
					i18n.EN: "Simplified Invoice",
					i18n.PL: "Faktura Uproszczona",
				},
				Desc: i18n.String{
					i18n.EN: "Invoice, as referred to in article 106e paragraph 5 point 3 of the act. Receipt up to 450 zł gross (100 euro) containing the buyer's NIP.",
					i18n.PL: "Faktura, o której mowa w art. 106e ust. 5 pkt 3 ustawy. Paragon fiskalny do kwoty 450 zł brutto (100 euro) zawierający NIP nabywcy.",
				},
			},
			{
				Code: "KOR",
				Name: i18n.String{
					i18n.EN: "Credit note",
					i18n.PL: "Faktura korygująca",
				},
				Desc: i18n.String{
					i18n.EN: "Corrects the original invoice.",
					i18n.PL: "Poprawia fakturę oryginalną.",
				},
			},
			{
				Code: "KOR_ZAL",
				Name: i18n.String{
					i18n.EN: "Prepayment credit note",
					i18n.PL: `Faktura korygująca fakturę zaliczkową`,
				},
				Desc: i18n.String{
					i18n.EN: "Invoice issued in connection with article 106f paragraph 4 of the act. Corrects the prepayment invoice (ZAL).",
					i18n.PL: "Faktura korygująca fakturę dokumentującą otrzymanie zapłaty lub jej części przed dokonaniem czynności oraz fakturę wystawioną w związku z art. 106f ust. 4 ustawy (faktura korygująca fakturę zaliczkową). Poprawia fakturę zaliczkową (ZAL).",
				},
			},
			{
				Code: "KOR_ROZ",
				Name: i18n.String{
					i18n.EN: "Settlement credit note",
					i18n.PL: "Faktura korygująca fakturę rozliczeniową",
				},
				Desc: i18n.String{
					i18n.EN: "Invoice issued in connection with article 106f paragraph 3 of the act. Corrects the settlement invoice (ROZ).",
					i18n.PL: "Faktura korygująca fakturę wystawioną w związku z art. 106f ust. 3 ustawy. Poprawia fakturę rozliczeniową (ROZ).",
				},
			},
		},
	},
	{
		Key: ExtKeyCashAccounting,
		Name: i18n.String{
			i18n.EN: "Cash accounting flag for KSeF",
			i18n.PL: "Flaga księgowania gotówki dla KSeF (P_16)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The ~pl-favat-cash-accounting~ extension indicates whether an invoice uses cash accounting
				for VAT purposes (kasowa metoda rozliczenia VAT). This extension is used in the **invoice
				tax section** (~tax.ext~).

				This extension is **not normalized automatically** and must be set manually by the user
				when the cash accounting method applies. According to Polish VAT law (Article 19a sec. 5
				item 1 or Article 21 sec. 1), certain small businesses may account for VAT on a cash
				basis rather than accrual basis.

				Values:
				- "1": Cash accounting applies
				- "2": Normal accounting (accrual basis) - default

				Example with cash accounting:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					"$addons": ["pl-favat-v3"],
					// ...
					"tax": {
						"ext": {
							"pl-favat-cash-accounting": "1"
						}
					}
				}
				~~~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Cash accounting",
					i18n.PL: "Księgowanie gotówkowe",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "No cash accounting",
					i18n.PL: "Bez księgowania gotówki",
				},
			},
		},
	},
	{
		Key: ExtKeySelfBilling,
		Name: i18n.String{
			i18n.EN: "Self-invoicing flag for KSeF",
			i18n.PL: "Flaga samofakturowania dla KSeF (P_17)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The ~pl-favat-self-billing~ extension indicates whether an invoice is self-billed
				(samofakturowanie), where the buyer issues the invoice on behalf of the supplier.
				This extension is used in the **invoice tax section** (~tax.ext~).

				This extension is **automatically normalized** by GOBL:
				- When the invoice has the ~self-billed~ tag, the value is set to "1"
				- Otherwise, the value defaults to "2"

				Self-billing is permitted under Article 106d sec. 1 of the Polish VAT Act when
				there is a prior agreement between the parties.

				Values:
				- "1": Self-billed invoice
				- "2": Regular invoice - default

				Example of a self-billed invoice:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					"$addons": ["pl-favat-v3"],
					"$tags": ["self-billed"],
					// ...
					"tax": {
						// Extension automatically set to "1"
						"ext": {
							"pl-favat-self-billing": "1"
						}
					}
				}
				~~~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Self-invoicing",
					i18n.PL: "Samofakturowanie",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Not self-invoicing",
					i18n.PL: "Bez samofakturowania",
				},
			},
		},
	},
	{
		Key: ExtKeyReverseCharge,
		Name: i18n.String{
			i18n.EN: "Reverse charge code for KSeF",
			i18n.PL: "Kod odwrotnego obciążenia dla KSeF (P_18)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The ~pl-favat-reverse-charge~ extension indicates whether the reverse charge mechanism
				applies to an invoice. This extension is used in the **invoice tax section** (~tax.ext~).

				This extension is **automatically normalized** by GOBL:
				- When the invoice has the ~reverse-charge~ tag, the value is set to "1"
				- Otherwise, the value defaults to "2"

				Under the reverse charge mechanism (odwrotne obciążenie), the buyer rather than the
				seller is responsible for accounting for the VAT. This applies to:
				- EU intra-community services (B2B cross-border services)
				- Domestic reverse charge for specific goods and services listed in Polish VAT law
				- Imports of services from outside the EU

				When using reverse charge, the line items should use the ~reverse-charge~ tax key,
				which automatically sets the tax category to "9" (EU Reverse Charge) or "10"
				(Domestic Reverse Charge) depending on the scenario.

				Values:
				- "1": Reverse charge applies
				- "2": Normal charge - default

				Example of an EU reverse charge invoice:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					"$addons": ["pl-favat-v3"],
					"$tags": ["reverse-charge"],
					// ...
					"lines": [
						{
							// ...
							"taxes": [
								{
									"cat": "VAT",
									"key": "reverse-charge"
								}
							]
						}
					],
					"tax": {
						// Extension automatically set to "1"
						"ext": {
							"pl-favat-reverse-charge": "1"
						}
					}
				}
				~~~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Reverse charge",
					i18n.PL: "Odwrotne obciążenie",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "No reverse charge",
					i18n.PL: "Bez odwrotnego obciążenia",
				},
			},
		},
	},
	{
		Key: ExtKeySplitPaymentMechanism,
		Name: i18n.String{
			i18n.EN: "Split payment mechanism flag for KSeF",
			i18n.PL: "Flaga mechanizmu płatności rozłożonej dla KSeF (P_18A)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The ~pl-favat-split-payment-mechanism~ extension indicates whether the mandatory split
				payment mechanism (mechanizm podzielonej płatności, MPP) applies to an invoice.
				This extension is used in the **invoice tax section** (~tax.ext~).

				This extension is **not normalized automatically** and must be set manually when the
				split payment mechanism applies.

				According to Polish VAT law, split payment is mandatory when ALL of the following
				conditions are met:
				- The invoice amount exceeds PLN 15,000 (or equivalent in foreign currency)
				- The invoice covers goods or services listed in Annex 15 to the VAT Act
				- The transaction is B2B (business to business)
				- The invoice bears the note "split payment mechanism"

				Under split payment, the buyer pays the VAT portion to a separate VAT account rather
				than to the seller directly.

				Values:
				- "1": Split payment mechanism applies
				- "2": Normal payment - default

				Example:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					"$addons": ["pl-favat-v3"],
					// ...
					"tax": {
						"ext": {
							"pl-favat-split-payment-mechanism": "1"
						}
					},
					"notes": [
						{
							"key": "general",
							"text": "Mechanizm podzielonej płatności"
						}
					]
				}
				~~~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Split payment mechanism",
					i18n.PL: "Mechanizm płatności rozłożonej",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "No split payment mechanism",
					i18n.PL: "Bez mechanizmu płatności rozłożonej",
				},
			},
		},
	},
	{
		Key: ExtKeyExemption,
		Name: i18n.String{
			i18n.EN: "Tax exemption code for KSeF",
			i18n.PL: "Kod oznaczający zwolnienie od podatku dla KSeF",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Extension used to indicate the type of reason for tax exemption code for KSeF. When the ~exempt~ tag
				is used in the invoice, having ~ext~ map's ~pl-favat-exemption~ property is required. Also, it is
				required to add descriptive text for the legal basis for exemption. To do this in GOBL, add a note
				to the invoice with the exemption reason, in the following format:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					// ...

					"ext": {
						"pl-favat-exemption": "A"
						// ...
					},
					"notes": [
						{
							"key": "legal",
							"code": "A",
							"src": "pl-favat-exemption",
							"text": "Art. 25a ust. 1 pkt 9 ustawy o VAT"
						}
					]
				}
				~~~

				In notes, code must match the code from the extension.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "A",
				Name: i18n.String{
					i18n.EN: "Law or act issued under the Polish law",
					i18n.PL: "Ustawa lub akt wydany na podstawie ustawy",
				},
			},
			{
				Code: "B",
				Name: i18n.String{
					i18n.EN: "Directive 2006/112/EC",
					i18n.PL: "Dyrektywa 2006/112/WE",
				},
			},
			{
				Code: "C",
				Name: i18n.String{
					i18n.EN: "Other legal basis",
					i18n.PL: "Inna podstawa prawna",
				},
			},
		},
	},
	{
		Key: ExtKeyMarginScheme,
		Name: i18n.String{
			i18n.EN: "Margin scheme code for KSeF",
			i18n.PL: "Kod oznaczający procedurę marży dla KSeF (P_PMarzy)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The ~pl-favat-margin-scheme~ extension specifies the type of margin scheme (procedura
				marży) applied to an invoice. This extension is used in the **invoice tax section**
				(~tax.ext~).

				This extension is **not normalized automatically** and must be set manually when a
				margin scheme applies.

				Under margin schemes, VAT is calculated only on the seller's margin (profit) rather
				than on the full sale price. This applies to specific business sectors in Poland:
				- Travel agencies (tour operators)
				- Sales of second-hand goods
				- Sales of works of art
				- Sales of collector's items and antiques

				When using a margin scheme, the tax category should typically be set to "11"
				(Margin scheme) in line items.

				Example for a travel agency:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					"$addons": ["pl-favat-v3"],
					// ...
					"tax": {
						"ext": {
							"pl-favat-margin-scheme": "2"
						}
					}
				}
				~~~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Travel agency",
					i18n.PL: "Biuro podróży",
				},
			},
			{
				Code: "3.1",
				Name: i18n.String{
					i18n.EN: "Second-hand goods",
					i18n.PL: "Towary używane",
				},
			},
			{
				Code: "3.2",
				Name: i18n.String{
					i18n.EN: "Works of art",
					i18n.PL: "Dzieła sztuki",
				},
			},
			{
				Code: "3.3",
				Name: i18n.String{
					i18n.EN: "Antiques and collectibles",
					i18n.PL: "Przedmioty kolekcjonerskie i antyki",
				},
			},
		},
	},
	{
		Key: ExtKeyPaymentMeans,
		Name: i18n.String{
			i18n.EN: "Payment method code for KSeF",
			i18n.PL: "Kod formy płatności dla KSeF (TFormaPlatnosci)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The ~pl-favat-payment-means~ extension specifies the payment method used in an invoice.
				This extension is used in **payment instructions** (~payment.instructions.ext~) or
				**payment advances** (~payment.advances[].ext~).

				This extension is **automatically normalized** by GOBL based on the payment means key.
				The following table shows the mapping:

				| Code | FA_VAT Name     | GOBL Payment Means Key    |
				| ---- | --------------- | ------------------------- |
				| 1    | Cash            | ~cash~                    |
				| 2    | Card            | ~card~                    |
				| 3    | Voucher         | ~other+voucher~           |
				| 4    | Cheque          | ~cheque~                  |
				| 5    | Credit/Loan     | ~other+credit~            |
				| 6    | Credit Transfer | ~credit-transfer~         |
				| 7    | Mobile          | ~online~                  |

				Example with bank transfer payment:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					"$addons": ["pl-favat-v3"],
					// ...
					"payment": {
						"instructions": {
							"key": "credit-transfer",
							// Extension automatically set to "6"
							"ext": {
								"pl-favat-payment-means": "6"
							},
							"credit_transfer": [
								{
									"iban": "PL61109010140000071219812874",
									"name": "Company Name"
								}
							]
						}
					}
				}
				~~~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Cash",
					i18n.PL: "Gotówka",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Card",
					i18n.PL: "Karta",
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "Coupon",
					i18n.PL: "Bon",
				},
			},
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "Cheque",
					i18n.PL: "Czek",
				},
			},
			{
				Code: "5",
				Name: i18n.String{
					i18n.EN: "Loan",
					i18n.PL: "Kredyt",
				},
			},
			{
				Code: "6",
				Name: i18n.String{
					i18n.EN: "Wire Transfer",
					i18n.PL: "Przelew",
				},
			},
			{
				Code: "7",
				Name: i18n.String{
					i18n.EN: "Mobile",
					i18n.PL: "Mobilna",
				},
			},
		},
	},
	{
		Key: ExtKeyJST,
		Name: i18n.String{
			i18n.EN: "Subordinate Local Government Unit flag",
			i18n.PL: "Flaga jednostki samorządu terytorialnego (P_22)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The ~pl-favat-jst~ extension indicates whether the customer is a Subordinate Local
				Government Unit (Jednostka Samorządu Terytorialnego - JST). This extension is used
				in the **customer party section** (~customer.ext~).

				This extension is **not normalized automatically** and must be set manually.

				When this extension is set to "1" (customer is JST), GOBL validates that the customer
				has an identity with role "8" (Local Government Unit - recipient) in the ~customer.identities~
				array. This identity must include a ~code~ field with the JST identifier.

				Values:
				- "1": Customer is a Subordinate Local Government Unit
				- "2": Customer is not a JST - default

				Example:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					"$addons": ["pl-favat-v3"],
					// ...
					"customer": {
						"name": "Gmina Warszawa",
						"tax_id": {
							"country": "PL",
							"code": "5252548806"
						},
						"ext": {
							"pl-favat-jst": "1"
						},
						"identities": [
							{
								"code": "146501",
								"ext": {
									"pl-favat-third-party-role": "8"
								}
							}
						]
					}
				}
				~~~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Customer is a Subordinate Local Government Unit",
					i18n.PL: "Nabywca jest jednostką samorządu terytorialnego",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Customer is not a Subordinate Local Government Unit",
					i18n.PL: "Nabywca nie jest jednostką samorządu terytorialnego",
				},
			},
		},
	},
	{
		Key: ExtKeyGroupVAT,
		Name: i18n.String{
			i18n.EN: "Group VAT member flag",
			i18n.PL: "Flaga członka grupy podatkowej (P_23)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The ~pl-favat-group-vat~ extension indicates whether the customer is a member of a VAT
				group (Grupa VAT - GV). This extension is used in the **customer party section**
				(~customer.ext~).

				This extension is **not normalized automatically** and must be set manually.

				VAT groups in Poland allow multiple legal entities to be treated as a single taxpayer
				for VAT purposes. When this extension is set to "1" (customer is a GV member), GOBL
				validates that the customer has an identity with role "10" (GV member - recipient) in
				the ~customer.identities~ array. This identity must include a ~code~ field with the
				VAT group member identifier.

				Values:
				- "1": Customer is a VAT Group member
				- "2": Customer is not a VAT Group member - default

				Example:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					"$addons": ["pl-favat-v3"],
					// ...
					"customer": {
						"name": "Member Company Ltd.",
						"tax_id": {
							"country": "PL",
							"code": "1234567890"
						},
						"ext": {
							"pl-favat-group-vat": "1"
						},
						"identities": [
							{
								"code": "GV-12345",
								"ext": {
									"pl-favat-third-party-role": "10"
								}
							}
						]
					}
				}
				~~~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Customer is a Group VAT member",
					i18n.PL: "Nabywca jest członkiem grupy podatkowej",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Customer is not a Group VAT member",
					i18n.PL: "Nabywca nie jest członkiem grupy podatkowej",
				},
			},
		},
	},
	{
		Key: ExtKeyThirdPartyRole,
		Name: i18n.String{
			i18n.EN: "Third party role",
			i18n.PL: "Rola trzeciej strony (RolaPodmiotu)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The ~pl-favat-third-party-role~ extension specifies the role of a third party or
				additional entity in an invoice transaction. This extension is used in **party
				identities** (~customer.identities[].ext~, ~supplier.identities[].ext~, or in
				additional parties).

				This extension is **not normalized automatically** and must be set manually based on
				the role of the entity in the transaction.

				Common use cases:
				- Role "8": Required when customer has ~pl-favat-jst~ = "1" (Local Government Unit)
				- Role "10": Required when customer has ~pl-favat-group-vat~ = "1" (VAT Group member)
				- Role "5": When invoice is issued by an entity on behalf of the taxpayer
				- Other roles: For factoring, recipients, payers, etc.

				Each identity with this extension should also include a ~code~ field with the
				identifier of the third party.

				Example with JST customer:

				~~~js
				{
					"customer": {
						// ...
						"ext": {
							"pl-favat-jst": "1"
						},
						"identities": [
							{
								"code": "146501",
								"ext": {
									"pl-favat-third-party-role": "8"
								}
							}
						]
					}
				}
				~~~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Factor",
					i18n.PL: "Faktor",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						When the invoice includes the factoring company details
					`),
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Recipient",
					i18n.PL: "Odbiorca",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						When the invoice includes internal units, branches, or divisions of the purchaser, which do not
						themselves qualify as purchasers under the Act.
					`),
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "Original entity",
					i18n.PL: "Podmiot pierwotny",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						When the invoice includes the details of an entity that was taken over or transformed and originally
						supplied the goods or services (excluding cases referred to in Article 106j sec. 2 item 3 of the Act,
						where these details are shown in section Podmiot1K).
					`),
				},
			},
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "Additional Purchaser",
					i18n.PL: "Dodatkowy nabywca",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						When the invoice includes details of further purchasers (other than the one indicated in Customer).
					`),
				},
			},
			{
				Code: "5",
				Name: i18n.String{
					i18n.EN: "Invoice Issuer",
					i18n.PL: "Wystawca faktury",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						When the invoice includes details of the
						entity issuing the invoice on behalf of the taxpayer.
					`),
				},
			},
			{
				Code: "6",
				Name: i18n.String{
					i18n.EN: "Payer",
					i18n.PL: "Zapłacający",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						When the invoice includes details of the
						entity paying the tax on behalf of the taxpayer.
					`),
				},
			},
			{
				Code: "7",
				Name: i18n.String{
					i18n.EN: "Local Government Unit (LGU) - issuer",
					i18n.PL: "Jednostka samorządu terytorialnego (JST) - wystawca",
				},
			},
			{
				Code: "8",
				Name: i18n.String{
					i18n.EN: "Local Government Unit (LGU) - recipient",
					i18n.PL: "Jednostka samorządu terytorialnego (JST) - odbiorca",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						When the invoice includes details of the
						entity receiving the goods or services from the Local Government Unit.
					`),
				},
			},
			{
				Code: "9",
				Name: i18n.String{
					i18n.EN: "GV member - issuer",
					i18n.PL: "Członek grupy podatkowej - wystawca",
				},
			},
			{
				Code: "10",
				Name: i18n.String{
					i18n.EN: "GV member - recipient",
					i18n.PL: "Członek grupy podatkowej - odbiorca",
				},
			},
			{
				Code: "11",
				Name: i18n.String{
					i18n.EN: "Employee",
					i18n.PL: "Pracownik",
				},
			},
		},
	},
	{
		Key: ExtKeyEffectiveDate,
		Name: i18n.String{
			i18n.EN: "Effective date code",
			i18n.PL: "Kod daty wejścia w życie (OkresFaKorygowanej)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The ~pl-favat-effective-date~ extension specifies when a correction invoice (credit note)
				becomes effective for VAT purposes. This extension is used in **credit note tax section**
				(~tax.ext~) when ~type~ is ~credit-note~.

				This extension is **not normalized automatically** and it is not mandatory. It can be set manually for credit
				notes based on when the correction should take effect.

				According to Polish VAT regulations, a correction invoice can be effective:
				- On the date of the original invoice (code "1")
				- On the date of the correction invoice itself (code "2")
				- On another specific date, or different dates for different line items (code "3")

				Values:
				- "1": Effective according to original invoice date
				- "2": Effective according to correction date
				- "3": Effective on another date (or multiple dates)

				Example of a credit note effective on correction date:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					"$addons": ["pl-favat-v3"],
					"type": "credit-note",
					// ...
					"tax": {
						"ext": {
							"pl-favat-invoice-type": "KOR",
							"pl-favat-effective-date": "2"
						}
					},
					"preceding": [
						{
							"code": "INVOICE-001",
							"issue_date": "2025-01-15"
						}
					]
				}
				~~~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Original",
					i18n.PL: "Pierwotna",
				},
				Desc: i18n.String{
					i18n.EN: "Effective according to date of the original invoice.",
					i18n.PL: "Faktura skutkująca w dacie ujęcia faktury pierwotnej.",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Correction",
					i18n.PL: "Korygująca",
				},
				Desc: i18n.String{
					i18n.EN: "Effective according to date of correction.",
					i18n.PL: "Faktura skutkująca w dacie ujęcia faktury korygującej.",
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "Other",
					i18n.PL: "Inna",
				},
				Desc: i18n.String{
					i18n.EN: "Correction has legal consequences in another date or the dates are different for different position on the invoice",
					i18n.PL: "Faktura skutkująca w innej dacie. W tym gdy dla różnych pozycji faktury korygującej data jest różna.",
				},
			},
		},
	},
}
