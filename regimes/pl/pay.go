package pl

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pay"
)

// Regime Specific Payment Means Extension Keys
const (
	MeansKeyCoupon     cbc.Key = "coupon"
	MeansKeyCheque     cbc.Key = "cheque"
	MeansKeyLoan       cbc.Key = "loan"
	MeansKeyDebtRelief cbc.Key = "credit-transfer"
	MeansKeyMobile     cbc.Key = "mobile"
)

var paymentMeansKeyDefinitions = []*cbc.Definition{
	{
		Key: pay.MeansKeyCash,
		Name: i18n.String{
			i18n.EN: "Cash",
			i18n.PL: "Gotówka",
		},
		Map: cbc.CodeMap{
			KeyFAVATPaymentType: "1",
		},
	},
	{
		Key: pay.MeansKeyCard,
		Name: i18n.String{
			i18n.EN: "Card",
			i18n.PL: "Karta",
		},
		Map: cbc.CodeMap{
			KeyFAVATPaymentType: "2",
		},
	},
	{
		Key: pay.MeansKeyOther.With(MeansKeyCoupon),
		Name: i18n.String{
			i18n.EN: "Coupon",
			i18n.PL: "Bon",
		},
		Map: cbc.CodeMap{
			KeyFAVATPaymentType: "3",
		},
	},
	{
		Key: pay.MeansKeyCheque,
		Name: i18n.String{
			i18n.EN: "Cheque",
			i18n.PL: "Czek",
		},
		Map: cbc.CodeMap{
			KeyFAVATPaymentType: "4",
		},
	},
	{
		Key: pay.MeansKeyOnline.With(MeansKeyLoan),
		Name: i18n.String{
			i18n.EN: "Loan",
			i18n.PL: "Kredyt",
		},
		Map: cbc.CodeMap{
			KeyFAVATPaymentType: "5",
		},
	},
	{
		Key: pay.MeansKeyCreditTransfer,
		Name: i18n.String{
			i18n.EN: "Wire Transfer",
			i18n.PL: "Przelew",
		},
		Map: cbc.CodeMap{
			KeyFAVATPaymentType: "6",
		},
	},
	{
		Key: pay.MeansKeyOther.With(MeansKeyMobile),
		Name: i18n.String{
			i18n.EN: "Mobile",
			i18n.PL: "Mobilna",
		},
		Map: cbc.CodeMap{
			KeyFAVATPaymentType: "7",
		},
	},
}
