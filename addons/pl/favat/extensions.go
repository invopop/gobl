package favat

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Regime extension codes for local electronic formats.
const (
	ExtKeyTaxCategory   cbc.Key = "pl-favat-tax-category"
	ExtKeyEffectiveDate cbc.Key = "pl-favat-effective-date"
	ExtKeyPaymentMeans  cbc.Key = "pl-favat-payment-means"  // for mapping to TFormaPlatnosci's codes (type of payment means - e.g. cash, bank transfer etc)
	ExtKeyInvoiceType   cbc.Key = "pl-favat-invoice-type"   // for mapping to TRodzajFaktury's codes (type of invoice - e.g. regular, in advance, correction etc)
	ExtKeySelfBilling   cbc.Key = "pl-favat-self-billing"   // for mapping to P_17 field, indicating self-invoicing
	ExtKeyReverseCharge cbc.Key = "pl-favat-reverse-charge" // for mapping to P_18, indicating reverse charge
	ExtKeyMarginScheme  cbc.Key = "pl-favat-margin-scheme"  // for mapping to P_PMarzy, indicating margin scheme
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
				Specifies tax categories including various special cases, reduced rates, and exemptions. Each of these categories corresponds to a field in the XML, with name beginning with P_13_, P_14_ etc.

				For example, on an invoice for intra-community supply of goods, we fill field P_13_6_1.
			`),
			i18n.PL: here.Doc(`
				Określa kategorie podatkowe, w tym różne przypadki szczególne, stawki obniżone i zwolnienia. Każda z tych kategorii odpowiada polu w XML, z nazwą zaczynającą się od P_13_, P_14_ itd.

				Na przykład, na fakturze dla dostawy towarów wewnątrzwspólnotowych (WDT) wypełniamy pole P_13_6_1.
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
				Code: "6_1",
				Name: i18n.String{
					i18n.EN: "0% WDT",
					i18n.PL: "0% WDT",
				},
				Desc: i18n.String{
					i18n.EN: "Zero VAT, intra-community supply of goods",
					i18n.PL: "Zerowa stawka podatku VAT, wewnątrzwspólnotowa dostawa towarów",
				},
			},
			{
				Code: "6_2",
				Name: i18n.String{
					i18n.EN: "0% Domestic",
					i18n.PL: "0% Krajowy",
				},
				Desc: i18n.String{
					i18n.EN: "Zero VAT, excluding intra-community supply of goods and export",
					i18n.PL: "Zerowa stawka podatku VAT, z wyłączeniem wewnątrzwspólnotowej dostawy towarów i eksportu",
				},
			},
			{
				Code: "6_3",
				Name: i18n.String{
					i18n.EN: "0% Export",
					i18n.PL: "0% Eksport",
				},
				Desc: i18n.String{
					i18n.EN: "Zero VAT, export outside the EU",
					i18n.PL: "Zerowa stawka podatku VAT, eksport poza Unię Europejską",
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
					i18n.EN: "Export",
					i18n.PL: "Eksport",
				},
				Desc: i18n.String{
					i18n.EN: "Sales of goods and services outside the country, excluding EU VAT and OSS",
					i18n.PL: "Sprzedaż towarów oraz świadczenia usług poza terytorium kraju, z wyłączeniem VAT UE i OSS",
				},
			},
			{
				Code: "9",
				Name: i18n.String{
					i18n.EN: "EU VAT",
					i18n.PL: "VAT UE",
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
					i18n.EN: "Reverse charge",
					i18n.PL: "Odwrotne obciążenie",
				},
				Desc: i18n.String{
					i18n.EN: "Obligation to account for tax by the purchaser",
					i18n.PL: "Obowiązek rozliczenia podatku przez nabywcę",
				},
			},
			{
				Code: "11",
				Name: i18n.String{
					i18n.EN: "EU VAT",
					i18n.PL: "Faktura VAT marża",
				},
				// https://poradnikprzedsiebiorcy.pl/-faktura-vat-marza-kiedy-mozna-ja-wystawic
				Desc: i18n.String{
					i18n.EN: "Tax only on the margin, which the seller has calculated - applies to tourism services, sale of used goods, antiques, works of art",
					i18n.PL: "Podatek tylko od marży, którą naliczył sprzedawca - dotyczy usług turystycznych, sprzedaży towarów używanych, antyków, dzieł sztuki",
				},
			},
		},
	},
	{
		Key: ExtKeyEffectiveDate,
		Name: i18n.String{
			i18n.EN: "Effective date code.",
			i18n.PL: "Kod daty wejścia w życie.",
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
	{
		Key: ExtKeyInvoiceType,
		Name: i18n.String{
			i18n.EN: "Invoice type code for KSeF",
			i18n.PL: "Kod rodzaju faktury dla KSeF",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Code for invoice type for KSeF. If not provided, GOBL will determine the appropriate code based on the invoice tags and type:

				| Code    | GOBL invoice type | Tags         |
				|---------|-------------------|--------------|
				| VAT			| ~standard~        | -            |
				| ZAL			| ~standard~        | ~partial~    |
				| ROZ			| ~standard~        | ~settlement~ |
				| UPR			| ~standard~        | ~simplified~ |
				| KOR			| ~credit_note~     | -            |
				| KOR_ZAL	| ~credit_note~     | ~partial~    |
				| KOR_ROZ	| ~credit_note~     | ~settlement~ |
			`),
			i18n.PL: here.Doc(`
				Kod rodzaju faktury dla KSeF. Jeśli nie jest podany, GOBL wyznaczy odpowiedni kod na podstawie tagów i typu faktury:

				| Kod 		| Typ faktury w GOBL | Tagi         |
				|---------|--------------------|--------------|
				| VAT			| ~standard~         | -            |
				| ZAL			| ~standard~         | ~partial~    |
				| ROZ			| ~standard~         | ~settlement~ |
				| UPR			| ~standard~         | ~simplified~ |
				| KOR			| ~credit_note~      | -            |
				| KOR_ZAL	| ~credit_note~      | ~partial~    |
				| KOR_ROZ	| ~credit_note~      | ~settlement~ |
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
		Key: ExtKeyPaymentMeans,
		Name: i18n.String{
			i18n.EN: "Payment method code for KSeF",
			i18n.PL: "Kod formy płatności dla KSeF",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Code for payment method KSeF. If not provided, GOBL will determine the appropriate code:

				Code | Meaning | GOBL code | Code from GOBL extension
				1    | Cash    | cash      | -
				2    | Card    | card      | -
				3    | Coupon  | other     | coupon
				4    | Cheque  | cheque    | -
				5    | Loan    | other     | loan
				6    | Transfer| transfer  | -
				7    | Mobile  | other     | mobile
			`),
			i18n.PL: here.Doc(`
				Kod formy płatności dla KSeF. Jeśli nie jest podany, GOBL wyznaczy odpowiedni kod:

				Kod  | Znaczenie | Kod w GOBL | Kod z rozszerzenia GOBL
				1    | Gotówka   | cash       | -
				2    | Karta     | card       | -
				3    | Bon       | other      | coupon
				4    | Czek      | cheque     | -
				5    | Kredyt    | other      | loan
				6    | Przelew   | transfer   | -
				7    | Mobilna   | other      | mobile
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
		Key: ExtKeySelfBilling,
		Name: i18n.String{
			i18n.EN: "Self-invoicing code for KSeF",
			i18n.PL: "Kod samofakturowania dla KSeF",
		},
		Desc: i18n.String{
			i18n.EN: "Code for indicating self-invoicing.",
			i18n.PL: "Kod wskazujący na samofakturowanie.",
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
			i18n.PL: "Kod odwrotnego obciążenia dla KSeF",
		},
		Desc: i18n.String{
			i18n.EN: "Code for indicating reverse charge.",
			i18n.PL: "Kod wskazujący na odwrotne obciążenie.",
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
		Key: ExtKeyMarginScheme,
		Name: i18n.String{
			i18n.EN: "Margin scheme code for KSeF",
			i18n.PL: "Kod oznaczający procedurę marży dla KSeF",
		},
		Desc: i18n.String{
			i18n.EN: "Code for indicating margin scheme.",
			i18n.PL: "Kod wskazujący na procedurę marży.",
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
					i18n.EN: "Used goods",
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
}
