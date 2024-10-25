package untdid

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

const (
	// ExtKeyDocumentType is used to identify the UNTDID 1001 document type code.
	ExtKeyDocumentType cbc.Key = "untdid-document-type"
	// ExtKeyPaymentMeans is used to identify the UNTDID 4461 payment means code.
	ExtKeyPaymentMeans cbc.Key = "untdid-payment-means"
	// ExtKeyTaxCategory is used to identify the UNTDID 5305 duty/tax/fee category code.
	ExtKeyTaxCategory cbc.Key = "untdid-tax-category"
)

var extensions = []*cbc.KeyDefinition{
	{
		Key: ExtKeyDocumentType,
		Name: i18n.String{
			i18n.EN: "UNTDID 1001 Document Type",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				UNTDID 1001 code used to describe the type of document. Ths list is based
				on the [EN16931 code list](https://ec.europa.eu/digital-building-blocks/sites/display/DIGITAL/Registry+of+supporting+artefacts+to+implement+EN16931#RegistryofsupportingartefactstoimplementEN16931-Codelists)
				values table which focusses on invoices and payments.

				Other tax regimes and addons may use their own subset of codes.
			`),
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "71",
				Name:  i18n.NewString("Request for payment"),
			},
			{
				Value: "80",
				Name:  i18n.NewString("Debit note related to goods or services"),
			},
			{
				Value: "81",
				Name:  i18n.NewString("Credit note related to goods or services"),
			},
			{
				Value: "82",
				Name:  i18n.NewString("Metered services invoice"),
			},
			{
				Value: "83",
				Name:  i18n.NewString("Credit note related to financial adjustments"),
			},
			{
				Value: "84",
				Name:  i18n.NewString("Debit note related to financial adjustments"),
			},
			{
				Value: "102",
				Name:  i18n.NewString("Tax notification"),
			},
			{
				Value: "130",
				Name:  i18n.NewString("Invoicing data sheet"),
			},
			{
				Value: "202",
				Name:  i18n.NewString("Direct payment valuation"),
			},
			{
				Value: "203",
				Name:  i18n.NewString("Provisional payment valuation"),
			},
			{
				Value: "204",
				Name:  i18n.NewString("Payment valuation"),
			},
			{
				Value: "211",
				Name:  i18n.NewString("Interim application for payment"),
			},
			{
				Value: "218",
				Name:  i18n.NewString("Final payment request based on completion of work"),
			},
			{
				Value: "219",
				Name:  i18n.NewString("Payment request for completed units"),
			},
			{
				Value: "261",
				Name:  i18n.NewString("Self billed credit note"),
			},
			{
				Value: "262",
				Name:  i18n.NewString("Consolidated credit note - goods and services"),
			},
			{
				Value: "295",
				Name:  i18n.NewString("Price variation invoice"),
			},
			{
				Value: "296",
				Name:  i18n.NewString("Credit note for price variation"),
			},
			{
				Value: "308",
				Name:  i18n.NewString("Delcredere credit note"),
			},
			{
				Value: "325",
				Name:  i18n.NewString("Proforma invoice"),
			},
			{
				Value: "326",
				Name:  i18n.NewString("Partial invoice"),
			},
			{
				Value: "380",
				Name:  i18n.NewString("Standard Invoice"),
			},
			{
				Value: "381",
				Name:  i18n.NewString("Credit note"),
			},
			{
				Value: "382",
				Name:  i18n.NewString("Commission note"),
			},
			{
				Value: "383",
				Name:  i18n.NewString("Debit note"),
			},
			{
				Value: "384",
				Name:  i18n.NewString("Corrected invoice"),
			},
			{
				Value: "385",
				Name:  i18n.NewString("Consolidated invoice"),
			},
			{
				Value: "386",
				Name:  i18n.NewString("Prepayment invoice"),
			},
			{
				Value: "387",
				Name:  i18n.NewString("Hire invoice"),
			},
			{
				Value: "388",
				Name:  i18n.NewString("Tax invoice"),
			},
			{
				Value: "389",
				Name:  i18n.NewString("Self-billed invoice"),
			},
			{
				Value: "390",
				Name:  i18n.NewString("Delcredere invoice"),
			},
			{
				Value: "393",
				Name:  i18n.NewString("Factored invoice"),
			},
			{
				Value: "394",
				Name:  i18n.NewString("Lease invoice"),
			},
			{
				Value: "395",
				Name:  i18n.NewString("Consignment invoice"),
			},
			{
				Value: "396",
				Name:  i18n.NewString("Factored credit note"),
			},
			{
				Value: "420",
				Name:  i18n.NewString("Optical Character Reading (OCR) payment credit note"),
			},
			{
				Value: "456",
				Name:  i18n.NewString("Debit advice"),
			},
			{
				Value: "457",
				Name:  i18n.NewString("Reversal of debit"),
			},
			{
				Value: "458",
				Name:  i18n.NewString("Reversal of credit"),
			},
			{
				Value: "527",
				Name:  i18n.NewString("Self billed debit note"),
			},
			{
				Value: "532",
				Name:  i18n.NewString("Forwarder's credit note"),
			},
			{
				Value: "553",
				Name:  i18n.NewString("Forwarder's invoice discrepancy report"),
			},
			{
				Value: "575",
				Name:  i18n.NewString("Insurer's invoice"),
			},
			{
				Value: "623",
				Name:  i18n.NewString("Forwarder's invoice"),
			},
			{
				Value: "633",
				Name:  i18n.NewString("Port charges documents"),
			},
			{
				Value: "751",
				Name:  i18n.NewString("Invoice information for accounting purposes"),
			},
			{
				Value: "780",
				Name:  i18n.NewString("Freight invoice"),
			},
			{
				Value: "817",
				Name:  i18n.NewString("Claim notification"),
			},
			{
				Value: "870",
				Name:  i18n.NewString("Consular invoice"),
			},
			{
				Value: "875",
				Name:  i18n.NewString("Partial construction invoice"),
			},
			{
				Value: "876",
				Name:  i18n.NewString("Partial final construction invoice"),
			},
			{
				Value: "877",
				Name:  i18n.NewString("Final construction invoice"),
			},
			{
				Value: "935",
				Name:  i18n.NewString("Customs invoice"),
			},
		},
	},
	{
		Key: ExtKeyPaymentMeans,
		Name: i18n.String{
			i18n.EN: "UNTDID 4461 Payment Means",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				UNTDID 4461 code used to describe the means of payment. This list is based on the
				[EN16931 code list](https://ec.europa.eu/digital-building-blocks/sites/display/DIGITAL/Registry+of+supporting+artefacts+to+implement+EN16931#RegistryofsupportingartefactstoimplementEN16931-Codelists)
				values table which focusses on invoices and payments.
			`),
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "1",
				Name:  i18n.NewString("Instrument not defined"),
			},
			{
				Value: "2",
				Name:  i18n.NewString("Automated clearing house credit"),
			},
			{
				Value: "3",
				Name:  i18n.NewString("Automated clearing house debit"),
			},
			{
				Value: "4",
				Name:  i18n.NewString("ACH demand debit reversal"),
			},
			{
				Value: "5",
				Name:  i18n.NewString("ACH demand credit reversal"),
			},
			{
				Value: "6",
				Name:  i18n.NewString("ACH demand credit"),
			},
			{
				Value: "7",
				Name:  i18n.NewString("ACH demand debit"),
			},
			{
				Value: "8",
				Name:  i18n.NewString("Hold"),
			},
			{
				Value: "9",
				Name:  i18n.NewString("National or regional clearing"),
			},
			{
				Value: "10",
				Name:  i18n.NewString("In cash"),
			},
			{
				Value: "11",
				Name:  i18n.NewString("ACH savings credit reversal"),
			},
			{
				Value: "12",
				Name:  i18n.NewString("ACH savings debit reversal"),
			},
			{
				Value: "13",
				Name:  i18n.NewString("ACH savings credit"),
			},
			{
				Value: "14",
				Name:  i18n.NewString("ACH savings debit"),
			},
			{
				Value: "15",
				Name:  i18n.NewString("Bookentry credit"),
			},
			{
				Value: "16",
				Name:  i18n.NewString("Bookentry debit"),
			},
			{
				Value: "17",
				Name:  i18n.NewString("ACH demand cash concentration/disbursement (CCD) credit"),
			},
			{
				Value: "18",
				Name:  i18n.NewString("ACH demand cash concentration/disbursement (CCD) debit"),
			},
			{
				Value: "19",
				Name:  i18n.NewString("ACH demand corporate trade payment (CTP) credit"),
			},
			{
				Value: "20",
				Name:  i18n.NewString("Cheque"),
			},
			{
				Value: "21",
				Name:  i18n.NewString("Banker's draft"),
			},
			{
				Value: "22",
				Name:  i18n.NewString("Certified banker's draft"),
			},
			{
				Value: "23",
				Name:  i18n.NewString("Bank cheque (issued by a banking or similar establishment)"),
			},
			{
				Value: "24",
				Name:  i18n.NewString("Bill of exchange awaiting acceptance"),
			},
			{
				Value: "25",
				Name:  i18n.NewString("Certified cheque"),
			},
			{
				Value: "26",
				Name:  i18n.NewString("Local cheque"),
			},
			{
				Value: "27",
				Name:  i18n.NewString("ACH demand corporate trade payment (CTP) debit"),
			},
			{
				Value: "28",
				Name:  i18n.NewString("ACH demand corporate trade exchange (CTX) credit"),
			},
			{
				Value: "29",
				Name:  i18n.NewString("ACH demand corporate trade exchange (CTX) debit"),
			},
			{
				Value: "30",
				Name:  i18n.NewString("Credit transfer"),
			},
			{
				Value: "31",
				Name:  i18n.NewString("Debit transfer"),
			},
			{
				Value: "32",
				Name:  i18n.NewString("ACH demand cash concentration/disbursement plus (CCD+)"),
			},
			{
				Value: "33",
				Name:  i18n.NewString("ACH demand cash concentration/disbursement plus (CCD+)"),
			},
			{
				Value: "34",
				Name:  i18n.NewString("ACH prearranged payment and deposit (PPD)"),
			},
			{
				Value: "35",
				Name:  i18n.NewString("ACH savings cash concentration/disbursement (CCD) credit"),
			},
			{
				Value: "36",
				Name:  i18n.NewString("ACH savings cash concentration/disbursement (CCD) debit"),
			},
			{
				Value: "37",
				Name:  i18n.NewString("ACH savings corporate trade payment (CTP) credit"),
			},
			{
				Value: "38",
				Name:  i18n.NewString("ACH savings corporate trade payment (CTP) debit"),
			},
			{
				Value: "39",
				Name:  i18n.NewString("ACH savings corporate trade exchange (CTX) credit"),
			},
			{
				Value: "40",
				Name:  i18n.NewString("ACH savings corporate trade exchange (CTX) debit"),
			},
			{
				Value: "41",
				Name:  i18n.NewString("ACH savings cash concentration/disbursement plus (CCD+)"),
			},
			{
				Value: "42",
				Name:  i18n.NewString("Payment to bank account"),
			},
			{
				Value: "43",
				Name:  i18n.NewString("ACH savings cash concentration/disbursement plus (CCD+)"),
			},
			{
				Value: "44",
				Name:  i18n.NewString("Accepted bill of exchange"),
			},
			{
				Value: "45",
				Name:  i18n.NewString("Referenced home-banking credit transfer"),
			},
			{
				Value: "46",
				Name:  i18n.NewString("Interbank debit transfer"),
			},
			{
				Value: "47",
				Name:  i18n.NewString("Home-banking debit transfer"),
			},
			{
				Value: "48",
				Name:  i18n.NewString("Bank card"),
			},
			{
				Value: "49",
				Name:  i18n.NewString("Direct debit"),
			},
			{
				Value: "50",
				Name:  i18n.NewString("Payment by postgiro"),
			},
			{
				Value: "51",
				Name:  i18n.NewString("FR, norme 6 97-Telereglement CFONB (French Organisation for"),
			},
			{
				Value: "52",
				Name:  i18n.NewString("Urgent commercial payment"),
			},
			{
				Value: "53",
				Name:  i18n.NewString("Urgent Treasury Payment"),
			},
			{
				Value: "54",
				Name:  i18n.NewString("Credit card"),
			},
			{
				Value: "55",
				Name:  i18n.NewString("Debit card"),
			},
			{
				Value: "56",
				Name:  i18n.NewString("Bankgiro"),
			},
			{
				Value: "57",
				Name:  i18n.NewString("Standing agreement"),
			},
			{
				Value: "58",
				Name:  i18n.NewString("SEPA credit transfer"),
			},
			{
				Value: "59",
				Name:  i18n.NewString("SEPA direct debit"),
			},
			{
				Value: "60",
				Name:  i18n.NewString("Promissory note"),
			},
			{
				Value: "61",
				Name:  i18n.NewString("Promissory note signed by the debtor"),
			},
			{
				Value: "62",
				Name:  i18n.NewString("Promissory note signed by the debtor and endorsed by a bank"),
			},
			{
				Value: "63",
				Name:  i18n.NewString("Promissory note signed by the debtor and endorsed by a"),
			},
			{
				Value: "64",
				Name:  i18n.NewString("Promissory note signed by a bank"),
			},
			{
				Value: "65",
				Name:  i18n.NewString("Promissory note signed by a bank and endorsed by another"),
			},
			{
				Value: "66",
				Name:  i18n.NewString("Promissory note signed by a third party"),
			},
			{
				Value: "67",
				Name:  i18n.NewString("Promissory note signed by a third party and endorsed by a"),
			},
			{
				Value: "68",
				Name:  i18n.NewString("Online payment service"),
			},
			{
				Value: "69",
				Name:  i18n.NewString("Transfer Advice"),
			},
			{
				Value: "70",
				Name:  i18n.NewString("Bill drawn by the creditor on the debtor"),
			},
			{
				Value: "74",
				Name:  i18n.NewString("Bill drawn by the creditor on a bank"),
			},
			{
				Value: "75",
				Name:  i18n.NewString("Bill drawn by the creditor, endorsed by another bank"),
			},
			{
				Value: "76",
				Name:  i18n.NewString("Bill drawn by the creditor on a bank and endorsed by a"),
			},
			{
				Value: "77",
				Name:  i18n.NewString("Bill drawn by the creditor on a third party"),
			},
			{
				Value: "78",
				Name:  i18n.NewString("Bill drawn by creditor on third party, accepted and"),
			},
			{
				Value: "91",
				Name:  i18n.NewString("Not transferable banker's draft"),
			},
			{
				Value: "92",
				Name:  i18n.NewString("Not transferable local cheque"),
			},
			{
				Value: "93",
				Name:  i18n.NewString("Reference giro"),
			},
			{
				Value: "94",
				Name:  i18n.NewString("Urgent giro"),
			},
			{
				Value: "95",
				Name:  i18n.NewString("Free format giro"),
			},
			{
				Value: "96",
				Name:  i18n.NewString("Requested method for payment was not used"),
			},
			{
				Value: "97",
				Name:  i18n.NewString("Clearing between partners"),
			},
			{
				Value: "98",
				Name:  i18n.NewString("JP, Electronically Recorded Monetary Claims"),
			},
			{
				Value: "ZZZ",
				Name:  i18n.NewString("Mutually defined"),
			},
		},
	},
	{
		Key: ExtKeyTaxCategory,
		Name: i18n.String{
			i18n.EN: "UNTDID 3505 Tax Category",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				UNTDID 5305 code used to describe the applicable duty/tax/fee category. There are
				multiple versions and subsets of this table so regimes and addons may need to filter
				options for a specific subset of values.

				Data from https://unece.org/fileadmin/DAM/trade/untdid/d16b/tred/tred5305.htm.
			`),
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "A",
				Name: i18n.String{
					i18n.EN: "Mixed tax rate",
				},
			},
			{
				Value: "AA",
				Name: i18n.String{
					i18n.EN: "Lower rate",
				},
			},
			{
				Value: "AB",
				Name: i18n.String{
					i18n.EN: "Exempt for resale",
				},
			},
			{
				Value: "AC",
				Name: i18n.String{
					i18n.EN: "Exempt for resale",
				},
			},
			{
				Value: "AD",
				Name: i18n.String{
					i18n.EN: "Value Added Tax (VAT) due from a previous invoice",
				},
			},
			{
				Value: "AE",
				Name: i18n.String{
					i18n.EN: "VAT Reverse Charge",
				},
			},
			{
				Value: "B",
				Name: i18n.String{
					i18n.EN: "Transferred (VAT)",
				},
			},
			{
				Value: "C",
				Name: i18n.String{
					i18n.EN: "Duty paid by supplier",
				},
			},
			{
				Value: "D",
				Name: i18n.String{
					i18n.EN: "Value Added Tax (VAT) margin scheme - travel agents",
				},
			},
			{
				Value: "E",
				Name: i18n.String{
					i18n.EN: "Exempt from tax",
				},
			},
			{
				Value: "F",
				Name: i18n.String{
					i18n.EN: "Value Added Tax (VAT) margin scheme - second-hand goods",
				},
			},
			{
				Value: "G",
				Name: i18n.String{
					i18n.EN: "Free export item, tax not charged",
				},
			},
			{
				Value: "H",
				Name: i18n.String{
					i18n.EN: "Higher rate",
				},
			},
			{
				Value: "I",
				Name: i18n.String{
					i18n.EN: "Value Added Tax (VAT) margin scheme - works of art",
				},
			},
			{
				Value: "J",
				Name: i18n.String{
					i18n.EN: "Value Added Tax (VAT) margin scheme - collector's items and antiques",
				},
			},
			{
				Value: "K",
				Name: i18n.String{
					i18n.EN: "VAT exempt for EEA intra-community supply of goods and services",
				},
			},
			{
				Value: "L",
				Name: i18n.String{
					i18n.EN: "Canary Islands general indirect tax",
				},
			},
			{
				Value: "M",
				Name: i18n.String{
					i18n.EN: "Tax for production, services and importation in Ceuta and Melilla",
				},
			},
			{
				Value: "O",
				Name: i18n.String{
					i18n.EN: "Services outside scope of tax",
				},
			},
			{
				Value: "S",
				Name: i18n.String{
					i18n.EN: "Standard Rate",
				},
			},
			{
				Value: "Z",
				Name: i18n.String{
					i18n.EN: "Zero rated goods",
				},
			},
		},
	},
}
