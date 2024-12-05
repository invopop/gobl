package untdid

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

const (
	// ExtKeyPaymentMeans is used to identify the UNTDID 4461 payment means code.
	ExtKeyPaymentMeans cbc.Key = "untdid-payment-means"
)

var extPaymentMeans = &cbc.Definition{
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
	Values: []*cbc.Definition{
		{
			Code: "1",
			Name: i18n.NewString("Instrument not defined"),
		},
		{
			Code: "2",
			Name: i18n.NewString("Automated clearing house credit"),
		},
		{
			Code: "3",
			Name: i18n.NewString("Automated clearing house debit"),
		},
		{
			Code: "4",
			Name: i18n.NewString("ACH demand debit reversal"),
		},
		{
			Code: "5",
			Name: i18n.NewString("ACH demand credit reversal"),
		},
		{
			Code: "6",
			Name: i18n.NewString("ACH demand credit"),
		},
		{
			Code: "7",
			Name: i18n.NewString("ACH demand debit"),
		},
		{
			Code: "8",
			Name: i18n.NewString("Hold"),
		},
		{
			Code: "9",
			Name: i18n.NewString("National or regional clearing"),
		},
		{
			Code: "10",
			Name: i18n.NewString("In cash"),
		},
		{
			Code: "11",
			Name: i18n.NewString("ACH savings credit reversal"),
		},
		{
			Code: "12",
			Name: i18n.NewString("ACH savings debit reversal"),
		},
		{
			Code: "13",
			Name: i18n.NewString("ACH savings credit"),
		},
		{
			Code: "14",
			Name: i18n.NewString("ACH savings debit"),
		},
		{
			Code: "15",
			Name: i18n.NewString("Bookentry credit"),
		},
		{
			Code: "16",
			Name: i18n.NewString("Bookentry debit"),
		},
		{
			Code: "17",
			Name: i18n.NewString("ACH demand cash concentration/disbursement (CCD) credit"),
		},
		{
			Code: "18",
			Name: i18n.NewString("ACH demand cash concentration/disbursement (CCD) debit"),
		},
		{
			Code: "19",
			Name: i18n.NewString("ACH demand corporate trade payment (CTP) credit"),
		},
		{
			Code: "20",
			Name: i18n.NewString("Cheque"),
		},
		{
			Code: "21",
			Name: i18n.NewString("Banker's draft"),
		},
		{
			Code: "22",
			Name: i18n.NewString("Certified banker's draft"),
		},
		{
			Code: "23",
			Name: i18n.NewString("Bank cheque (issued by a banking or similar establishment)"),
		},
		{
			Code: "24",
			Name: i18n.NewString("Bill of exchange awaiting acceptance"),
		},
		{
			Code: "25",
			Name: i18n.NewString("Certified cheque"),
		},
		{
			Code: "26",
			Name: i18n.NewString("Local cheque"),
		},
		{
			Code: "27",
			Name: i18n.NewString("ACH demand corporate trade payment (CTP) debit"),
		},
		{
			Code: "28",
			Name: i18n.NewString("ACH demand corporate trade exchange (CTX) credit"),
		},
		{
			Code: "29",
			Name: i18n.NewString("ACH demand corporate trade exchange (CTX) debit"),
		},
		{
			Code: "30",
			Name: i18n.NewString("Credit transfer"),
		},
		{
			Code: "31",
			Name: i18n.NewString("Debit transfer"),
		},
		{
			Code: "32",
			Name: i18n.NewString("ACH demand cash concentration/disbursement plus (CCD+)"),
		},
		{
			Code: "33",
			Name: i18n.NewString("ACH demand cash concentration/disbursement plus (CCD+)"),
		},
		{
			Code: "34",
			Name: i18n.NewString("ACH prearranged payment and deposit (PPD)"),
		},
		{
			Code: "35",
			Name: i18n.NewString("ACH savings cash concentration/disbursement (CCD) credit"),
		},
		{
			Code: "36",
			Name: i18n.NewString("ACH savings cash concentration/disbursement (CCD) debit"),
		},
		{
			Code: "37",
			Name: i18n.NewString("ACH savings corporate trade payment (CTP) credit"),
		},
		{
			Code: "38",
			Name: i18n.NewString("ACH savings corporate trade payment (CTP) debit"),
		},
		{
			Code: "39",
			Name: i18n.NewString("ACH savings corporate trade exchange (CTX) credit"),
		},
		{
			Code: "40",
			Name: i18n.NewString("ACH savings corporate trade exchange (CTX) debit"),
		},
		{
			Code: "41",
			Name: i18n.NewString("ACH savings cash concentration/disbursement plus (CCD+)"),
		},
		{
			Code: "42",
			Name: i18n.NewString("Payment to bank account"),
		},
		{
			Code: "43",
			Name: i18n.NewString("ACH savings cash concentration/disbursement plus (CCD+)"),
		},
		{
			Code: "44",
			Name: i18n.NewString("Accepted bill of exchange"),
		},
		{
			Code: "45",
			Name: i18n.NewString("Referenced home-banking credit transfer"),
		},
		{
			Code: "46",
			Name: i18n.NewString("Interbank debit transfer"),
		},
		{
			Code: "47",
			Name: i18n.NewString("Home-banking debit transfer"),
		},
		{
			Code: "48",
			Name: i18n.NewString("Bank card"),
		},
		{
			Code: "49",
			Name: i18n.NewString("Direct debit"),
		},
		{
			Code: "50",
			Name: i18n.NewString("Payment by postgiro"),
		},
		{
			Code: "51",
			Name: i18n.NewString("FR, norme 6 97-Telereglement CFONB (French Organisation for"),
		},
		{
			Code: "52",
			Name: i18n.NewString("Urgent commercial payment"),
		},
		{
			Code: "53",
			Name: i18n.NewString("Urgent Treasury Payment"),
		},
		{
			Code: "54",
			Name: i18n.NewString("Credit card"),
		},
		{
			Code: "55",
			Name: i18n.NewString("Debit card"),
		},
		{
			Code: "56",
			Name: i18n.NewString("Bankgiro"),
		},
		{
			Code: "57",
			Name: i18n.NewString("Standing agreement"),
		},
		{
			Code: "58",
			Name: i18n.NewString("SEPA credit transfer"),
		},
		{
			Code: "59",
			Name: i18n.NewString("SEPA direct debit"),
		},
		{
			Code: "60",
			Name: i18n.NewString("Promissory note"),
		},
		{
			Code: "61",
			Name: i18n.NewString("Promissory note signed by the debtor"),
		},
		{
			Code: "62",
			Name: i18n.NewString("Promissory note signed by the debtor and endorsed by a bank"),
		},
		{
			Code: "63",
			Name: i18n.NewString("Promissory note signed by the debtor and endorsed by a"),
		},
		{
			Code: "64",
			Name: i18n.NewString("Promissory note signed by a bank"),
		},
		{
			Code: "65",
			Name: i18n.NewString("Promissory note signed by a bank and endorsed by another"),
		},
		{
			Code: "66",
			Name: i18n.NewString("Promissory note signed by a third party"),
		},
		{
			Code: "67",
			Name: i18n.NewString("Promissory note signed by a third party and endorsed by a"),
		},
		{
			Code: "68",
			Name: i18n.NewString("Online payment service"),
		},
		{
			Code: "69",
			Name: i18n.NewString("Transfer Advice"),
		},
		{
			Code: "70",
			Name: i18n.NewString("Bill drawn by the creditor on the debtor"),
		},
		{
			Code: "74",
			Name: i18n.NewString("Bill drawn by the creditor on a bank"),
		},
		{
			Code: "75",
			Name: i18n.NewString("Bill drawn by the creditor, endorsed by another bank"),
		},
		{
			Code: "76",
			Name: i18n.NewString("Bill drawn by the creditor on a bank and endorsed by a"),
		},
		{
			Code: "77",
			Name: i18n.NewString("Bill drawn by the creditor on a third party"),
		},
		{
			Code: "78",
			Name: i18n.NewString("Bill drawn by creditor on third party, accepted and"),
		},
		{
			Code: "91",
			Name: i18n.NewString("Not transferable banker's draft"),
		},
		{
			Code: "92",
			Name: i18n.NewString("Not transferable local cheque"),
		},
		{
			Code: "93",
			Name: i18n.NewString("Reference giro"),
		},
		{
			Code: "94",
			Name: i18n.NewString("Urgent giro"),
		},
		{
			Code: "95",
			Name: i18n.NewString("Free format giro"),
		},
		{
			Code: "96",
			Name: i18n.NewString("Requested method for payment was not used"),
		},
		{
			Code: "97",
			Name: i18n.NewString("Clearing between partners"),
		},
		{
			Code: "98",
			Name: i18n.NewString("JP, Electronically Recorded Monetary Claims"),
		},
		{
			Code: "ZZZ",
			Name: i18n.NewString("Mutually defined"),
		},
	},
}
