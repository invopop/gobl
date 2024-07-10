package it

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/validation"
)

// Regime Specific Payment Means Extension Keys
const (
	MeansKeyTreasury       cbc.Key = "treasury"
	MeansKeyPaymentSlip    cbc.Key = "payment-slip"
	MeansKeyRID            cbc.Key = "rid"
	MeansKeyRIDUtility     cbc.Key = "rid-utility"
	MeansKeyRIDFast        cbc.Key = "rid-fast"
	MeansKeyRIBA           cbc.Key = "riba"
	MeansKeyTaxReceipt     cbc.Key = "tax-receipt"
	MeansKeySpecialAccount cbc.Key = "special-account"
	MeansKeyPostOffice     cbc.Key = "post-office"
	MeansKeySEPA           cbc.Key = "sepa"
	MeansKeySEPACore       cbc.Key = "sepa-core"
	MeansKeySEPAB2B        cbc.Key = "sepa-b2b"
	MeansKeyPagoPA         cbc.Key = "pagopa"
)

var paymentMeansKeyDefinitions = []*cbc.KeyDefinition{
	{
		Key: pay.MeansKeyCash,
		Name: i18n.String{
			i18n.EN: "Cash",
			i18n.IT: "Contanti", // nolint:misspell
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP01",
		},
	},
	{
		Key: pay.MeansKeyCheque,
		Name: i18n.String{
			i18n.EN: "Cheque",
			i18n.IT: "Assegno",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP02",
		},
	},
	{
		Key: pay.MeansKeyBankDraft,
		Name: i18n.String{
			i18n.EN: "Banker's Draft",
			i18n.IT: "Assegno circolare",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP03",
		},
	},
	{
		Key: pay.MeansKeyCash.With(MeansKeyTreasury),
		Name: i18n.String{
			i18n.EN: "Cash at Treasury",
			i18n.IT: "Contanti presso Tesoreria", // nolint:misspell
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP04",
		},
	},
	{
		Key: pay.MeansKeyCreditTransfer,
		Name: i18n.String{
			i18n.EN: "Bank Transfer",
			i18n.IT: "Bonifico",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP05",
		},
	},
	{
		Key: pay.MeansKeyPromissoryNote,
		Name: i18n.String{
			i18n.EN: "Promissory Note",
			i18n.IT: "Vaglia cambiario",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP06",
		},
	},
	{
		Key: pay.MeansKeyOther.With(MeansKeyPaymentSlip),
		Name: i18n.String{
			i18n.EN: "Bank payment slip",
			i18n.IT: "Bollettino bancario",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP07",
		},
	},
	{
		Key: pay.MeansKeyCard,
		Name: i18n.String{
			i18n.EN: "Payment card",
			i18n.IT: "Carta di pagamento",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP08",
		},
	},
	{
		Key: pay.MeansKeyDirectDebit.With(MeansKeyRID),
		Name: i18n.String{
			i18n.EN: "Direct Debit (RID)",
			i18n.IT: "RID",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP09",
		},
	},
	{
		Key: pay.MeansKeyDirectDebit.With(MeansKeyRIDUtility),
		Name: i18n.String{
			i18n.EN: "Utilities Direct Debit (RID utenze)",
			i18n.IT: "RID utenze",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP10",
		},
	},
	{
		Key: pay.MeansKeyDirectDebit.With(MeansKeyRIDFast),
		Name: i18n.String{
			i18n.EN: "Fast Direct Debit (RID veloce)",
			i18n.IT: "RID veloce",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP11",
		},
	},
	{
		Key: pay.MeansKeyDirectDebit.With(MeansKeyRIBA),
		Name: i18n.String{
			i18n.EN: "Direct Debit (RIBA)",
			i18n.IT: "RIBA",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP12",
		},
	},
	{
		Key: pay.MeansKeyDebitTransfer,
		Name: i18n.String{
			i18n.EN: "Debit Transfer (MAV)",
			i18n.IT: "MAV",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP13",
		},
	},
	{
		Key: pay.MeansKeyOther.With(MeansKeyTaxReceipt),
		Name: i18n.String{
			i18n.EN: "Tax Receipt",
			i18n.IT: "Quietanza erario",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP14",
		},
	},
	{
		Key: pay.MeansKeyOther.With(MeansKeySpecialAccount),
		Name: i18n.String{
			i18n.EN: "Transfer on special account",
			i18n.IT: "Giroconto su conti di contabilità speciale",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP15",
		},
	},
	{
		Key: pay.MeansKeyDirectDebit,
		Name: i18n.String{
			i18n.EN: "Direct Debit",
			i18n.IT: "Domiciliazione Bancaria",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP16",
		},
	},
	{
		Key: pay.MeansKeyDirectDebit.With(MeansKeyPostOffice),
		Name: i18n.String{
			i18n.EN: "Direct Debit Post Office",
			i18n.IT: "Domiciliazione Postale",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP17",
		},
	},
	{
		Key: pay.MeansKeyCheque.With(MeansKeyPostOffice),
		Name: i18n.String{
			i18n.EN: "Post Office Cheque",
			i18n.IT: "Bollettino di c/c postale",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP18",
		},
	},
	{
		Key: pay.MeansKeyDirectDebit.With(MeansKeySEPA),
		Name: i18n.String{
			i18n.EN: "SEPA Direct Debit",
			i18n.IT: "SEPA Direct Debit",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP19",
		},
	},
	{
		Key: pay.MeansKeyDirectDebit.With(MeansKeySEPACore),
		Name: i18n.String{
			i18n.EN: "SEPA Core Direct Debit",
			i18n.IT: "SEPA Direct Debit Core",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP20",
		},
	},
	{
		Key: pay.MeansKeyDirectDebit.With(MeansKeySEPAB2B),
		Name: i18n.String{
			i18n.EN: "SEPA B2B Direct Debit",
			i18n.IT: "SEPA Direct Debit B2B",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP21",
		},
	},
	{
		Key: pay.MeansKeyNetting,
		Name: i18n.String{
			i18n.EN: "Deductible Netting",
			i18n.IT: "Trattenuta su somme già riscosse",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP22",
		},
	},
	{
		Key: pay.MeansKeyOnline.With(MeansKeyPagoPA),
		Name: i18n.String{
			i18n.EN: "PagoPA",
			i18n.IT: "PagoPA",
		},
		Map: cbc.CodeMap{
			KeyFatturaPAModalitaPagamento: "MP23",
		},
	},
	{
		Key: pay.MeansKeyOnline,
		Name: i18n.String{
			i18n.EN: "Online",
			i18n.IT: "Online",
		},
		Map: cbc.CodeMap{
			// Using "card" code
			KeyFatturaPAModalitaPagamento: "MP08",
		},
	},
	{
		Key: pay.MeansKeyOther,
		Name: i18n.String{
			i18n.EN: "Other",
			i18n.IT: "Altro",
		},
		Map: cbc.CodeMap{
			// Anything else assume is Cash
			KeyFatturaPAModalitaPagamento: "MP01",
		},
	},
}

var isValidPaymentMeanKey = cbc.InKeyDefs(paymentMeansKeyDefinitions)

func validatePayAdvance(a *pay.Advance) error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Key,
			validation.Required,
			isValidPaymentMeanKey,
		),
	)
}

func validatePayInstructions(i *pay.Instructions) error {
	return validation.ValidateStruct(i,
		validation.Field(&i.Key,
			validation.Required,
			isValidPaymentMeanKey,
		),
	)
}
