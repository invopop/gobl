package favat

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Regime extension codes for local electronic formats.
const (
	ExtKeyVATZero       cbc.Key = "pl-favat-vat-zero"
	ExtKeyVATSpecial    cbc.Key = "pl-favat-vat-special"
	ExtKeyEffectiveDate cbc.Key = "pl-favat-effective-date"
	ExtKeyPaymentMeans  cbc.Key = "pl-favat-payment-means" // for mapping to TFormaPlatnosci's codes
	ExtKeyInvoiceType   cbc.Key = "pl-favat-invoice-type"  // for mapping to TRodzajFaktury's codes
)

var extensionKeys = []*cbc.Definition{
	{
		Key: ExtKeyVATSpecial,
		Name: i18n.String{
			i18n.EN: "Special VAT Extensions for KSeF",
			i18n.PL: "Rozszerzenia specjalne dla KSeF",
		},
		Values: []*cbc.Definition{
			{
				Code: "taxi",
				Name: i18n.String{
					i18n.EN: "Taxi Rate",
					i18n.PL: "Ryczałt dla taksówek",
				},
				Desc: i18n.String{
					i18n.EN: "Special flat rate for taxi drivers.",
					i18n.PL: "Specjalna stawka ryczałtu dla taksówkarzy.",
				},
			},
		},
	},
	{
		Key: ExtKeyVATZero,
		Name: i18n.String{
			i18n.EN: "Zero VAT Extensions for KSeF",
		},
		Values: []*cbc.Definition{
			{
				Code: "wdt",
				Name: i18n.String{
					i18n.EN: "WDT",
					i18n.PL: "WDT",
				},
				Desc: i18n.String{
					i18n.EN: "Intra-community supply of goods",
					i18n.PL: "Wewnątrzwspólnotowa dostawa towarów",
				},
			},
			{
				Code: "domestic",
				Name: i18n.String{
					i18n.EN: "Domestic",
					i18n.PL: "Krajowy",
				},
				Desc: i18n.String{
					i18n.EN: "Zero VAT, excluding WDT and export",
					i18n.PL: "Zerowa stawka podatku z wyłączeniem WDT i eksportu",
				},
			},
			{
				Code: "export",
				Name: i18n.String{
					i18n.EN: "Export",
					i18n.PL: "Eksport",
				},
				Desc: i18n.String{
					i18n.EN: "Export outside the EU",
					i18n.PL: "Eksport poza Unią Europejską",
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
}
