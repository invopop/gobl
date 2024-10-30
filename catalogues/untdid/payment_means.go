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

var extPaymentMeans = &cbc.KeyDefinition{
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
}
