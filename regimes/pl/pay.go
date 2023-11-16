package pl

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
)

// Regime Specific Payment Means Extension Keys
const (
	MeansKeyCoupon     cbc.Key = "coupon"
	MeansKeyCheque     cbc.Key = "cheque"
	MeansKeyLoan       cbc.Key = "loan"
	MeansKeyDebtRelief cbc.Key = "credit-transfer"
	MeansKeyMobile     cbc.Key = "mobile"
)

var paymentMeansKeyDefinitions = []*tax.KeyDefinition{
	{
		Key: pay.MeansKeyCash,
		Name: i18n.String{
			i18n.EN: "Cash",
			i18n.PL: "Gotówka",
		},
		Map: cbc.CodeMap{
			KeyFA_VATFormaPlatnosci: "1",
		},
	},
	{
		Key: pay.MeansKeyCard,
		Name: i18n.String{
			i18n.EN: "Card",
			i18n.PL: "Karta",
		},
		Map: cbc.CodeMap{
			KeyFA_VATFormaPlatnosci: "2",
		},
	},
	{
		Key: pay.MeansKeyOther.With(MeansKeyCoupon),
		Name: i18n.String{
			i18n.EN: "Coupon",
			i18n.PL: "Bon",
		},
		Map: cbc.CodeMap{
			KeyFA_VATFormaPlatnosci: "3",
		},
	},
	{
		Key: pay.MeansKeyCheque,
		Name: i18n.String{
			i18n.EN: "Cheque",
			i18n.PL: "Czek",
		},
		Map: cbc.CodeMap{
			KeyFA_VATFormaPlatnosci: "4",
		},
	},
	{
		Key: pay.MeansKeyOnline.With(MeansKeyLoan),
		Name: i18n.String{
			i18n.EN: "Loan",
			i18n.PL: "Kredyt",
		},
		Map: cbc.CodeMap{
			KeyFA_VATFormaPlatnosci: "5",
		},
	},
	{
		Key: pay.MeansKeyCreditTransfer,
		Name: i18n.String{
			i18n.EN: "Wire Transfer",
			i18n.PL: "Przelew",
		},
		Map: cbc.CodeMap{
			KeyFA_VATFormaPlatnosci: "6",
		},
	},
	{
		Key: pay.MeansKeyOther.With(MeansKeyMobile),
		Name: i18n.String{
			i18n.EN: "Mobile",
			i18n.PL: "Mobilna",
		},
		Map: cbc.CodeMap{
			KeyFA_VATFormaPlatnosci: "7",
		},
	},
}
