package gr

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pay"
)

// Regime Specific Payment Means Extension Keys
const (
	MeansKeyForeign cbc.Key = "foreign"
)

var paymentMeansKeys = []*cbc.KeyDefinition{
	{
		Key: pay.MeansKeyCreditTransfer,
		Name: i18n.String{
			i18n.EN: "Domestic Payments Account Number",
			i18n.EL: "Επαγ. Λογαριασμός Πληρωμών Ημεδαπής",
		},
		Map: cbc.CodeMap{
			KeyIAPRPaymentMethod: "1",
		},
	},
	{
		Key: pay.MeansKeyCreditTransfer.With(MeansKeyForeign),
		Name: i18n.String{
			i18n.EN: "Foreign Payments Account Number",
			i18n.EL: "Επαγ. Λογαριασμός Πληρωμών Αλλοδαπής",
		},
		Map: cbc.CodeMap{
			KeyIAPRPaymentMethod: "2",
		},
	},
	{
		Key: pay.MeansKeyCash,
		Name: i18n.String{
			i18n.EN: "Cash",
			i18n.EL: "Μετρητά",
		},
		Map: cbc.CodeMap{
			KeyIAPRPaymentMethod: "3",
		},
	},
	{
		Key: pay.MeansKeyCheque,
		Name: i18n.String{
			i18n.EN: "Check",
			i18n.EL: "Επιταγή",
		},
		Map: cbc.CodeMap{
			KeyIAPRPaymentMethod: "4",
		},
	},
	{
		Key: pay.MeansKeyPromissoryNote,
		Name: i18n.String{
			i18n.EN: "On credit",
			i18n.EL: "Επί Πιστώσει",
		},
		Map: cbc.CodeMap{
			KeyIAPRPaymentMethod: "5",
		},
	},
	{
		Key: pay.MeansKeyOnline,
		Name: i18n.String{
			i18n.EN: "Web Banking",
			i18n.EL: "Web Banking",
		},
		Map: cbc.CodeMap{
			KeyIAPRPaymentMethod: "6",
		},
	},
	{
		Key: pay.MeansKeyCard,
		Name: i18n.String{
			i18n.EN: "POS / e-POS",
			i18n.EL: "POS / e-POS",
		},
		Map: cbc.CodeMap{
			KeyIAPRPaymentMethod: "7",
		},
	},
}

var isValidPaymentMeanKey = cbc.InKeyDefs(paymentMeansKeys)
