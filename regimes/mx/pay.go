package mx

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
)

// Regime Specific Payment Means Extension Keys
const (
	MeansKeyWallet          cbc.Key = "wallet"
	MeansKeyGroceryVouchers cbc.Key = "grocery-vouchers"
	MeansKeyInKind          cbc.Key = "in-kind"
	MeansKeySubrogation     cbc.Key = "subrogation"
	MeansKeyConsignment     cbc.Key = "consignment"
	MeansKeyDebtRelief      cbc.Key = "debt-relief"
	MeansKeyNovation        cbc.Key = "novation"
	MeansKeyMerger          cbc.Key = "merger"
	MeansKeyRemission       cbc.Key = "remission"
	MeansKeyExpiration      cbc.Key = "expiration"
	MeansKeyExtingishment   cbc.Key = "extinguishment"
	MeansKeyDebit           cbc.Key = "debit"
	MeansKeyServices        cbc.Key = "services"
	MeansKeyAdvance         cbc.Key = "advance"
	MeansKeyIntermediary    cbc.Key = "intermediary"
)

var paymentMeansKeyDefinitions = []*tax.KeyDefinition{
	{
		Key: pay.MeansKeyCash,
		Name: i18n.String{
			i18n.EN: "Cash",
			i18n.ES: "Efectivo",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "01",
		},
	},
	{
		Key: pay.MeansKeyCheque,
		Name: i18n.String{
			i18n.EN: "Nominative cheque",
			i18n.ES: "Cheque nominativo", // nolint:misspell
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "02",
		},
	},
	{
		Key: pay.MeansKeyCreditTransfer,
		Name: i18n.String{
			i18n.EN: "Electronic funds transfer",
			i18n.ES: "Transferencia electrónica de fondos",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "03",
		},
	},
	{
		Key: pay.MeansKeyCard,
		Name: i18n.String{
			i18n.EN: "Credit card",
			i18n.ES: "Tarjeta de crédito",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "04",
		},
	},
	{
		Key: pay.MeansKeyOnline.With(MeansKeyWallet),
		Name: i18n.String{
			i18n.EN: "Electronic wallet",
			i18n.ES: "Monedero electrónico",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "05",
		},
	},
	{
		Key: pay.MeansKeyOnline,
		Name: i18n.String{
			i18n.EN: "Online or electronic payment",
			i18n.ES: "Dinero electrónico",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "06",
		},
	},
	{
		Key: pay.MeansKeyOther.With(MeansKeyGroceryVouchers),
		Name: i18n.String{
			i18n.EN: "Grocery vouchers",
			i18n.ES: "Vales de despensa",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "08",
		},
	},
	{
		Key: pay.MeansKeyOther.With(MeansKeyInKind),
		Name: i18n.String{
			i18n.EN: "Payment in kind",
			i18n.ES: "Dación en pago",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "12",
		},
	},
	{
		Key: pay.MeansKeyOther.With(MeansKeySubrogation),
		Name: i18n.String{
			i18n.EN: "Payment by subrogation",
			i18n.ES: "Pago por subrogación",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "13",
		},
	},
	{
		Key: pay.MeansKeyOther.With(MeansKeyConsignment),
		Name: i18n.String{
			i18n.EN: "Payment by consignment",
			i18n.ES: "Pago por consignación",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "14",
		},
	},
	{
		Key: pay.MeansKeyOther.With(MeansKeyDebtRelief),
		Name: i18n.String{
			i18n.EN: "Debt relief",
			i18n.ES: "Condonación",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "15",
		},
	},
	{
		Key: pay.MeansKeyNetting,
		Name: i18n.String{
			i18n.EN: "Netting",
			i18n.ES: "Compensación",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "17",
		},
	},
	{
		Key: pay.MeansKeyOther.With(MeansKeyNovation),
		Name: i18n.String{
			i18n.EN: "Novation",
			i18n.ES: "Novación",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "23",
		},
	},
	{
		Key: pay.MeansKeyOther.With(MeansKeyMerger),
		Name: i18n.String{
			i18n.EN: "Merger",
			i18n.ES: "Confusión",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "24",
		},
	},
	{
		Key: pay.MeansKeyOther.With(MeansKeyRemission),
		Name: i18n.String{
			i18n.EN: "Debt remission",
			i18n.ES: "Remisión de deuda",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "25",
		},
	},
	{
		Key: pay.MeansKeyOther.With(MeansKeyExpiration),
		Name: i18n.String{
			i18n.EN: "Expiration of payment obligation",
			i18n.ES: "Prescripción o caducidad",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "26",
		},
	},
	{
		Key: pay.MeansKeyOther.With(MeansKeyExtingishment),
		Name: i18n.String{
			i18n.EN: "Extinguishment of payment obligation",
			i18n.ES: "A satisfacción del acreedor",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "27",
		},
	},
	{
		Key: pay.MeansKeyCard.With(MeansKeyDebit),
		Name: i18n.String{
			i18n.EN: "Debit card",
			i18n.ES: "Tarjeta de débito",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "28",
		},
	},
	{
		Key: pay.MeansKeyCard.With(MeansKeyServices),
		Name: i18n.String{
			i18n.EN: "Services card",
			i18n.ES: "Tarjeta de servicios",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "29",
		},
	},
	{
		Key: pay.MeansKeyOther.With(MeansKeyAdvance),
		Name: i18n.String{
			i18n.EN: "Advance payment",
			i18n.ES: "Aplicación de anticipos",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "30",
		},
	},
	{
		Key: pay.MeansKeyOther.With(MeansKeyIntermediary),
		Name: i18n.String{
			i18n.EN: "Payment via intermediary",
			i18n.ES: "Intermediario pagos",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "31",
		},
	},
	{
		Key: pay.MeansKeyOther,
		Name: i18n.String{
			i18n.EN: "To be defined",
			i18n.ES: "Por definir",
		},
		Codes: cbc.CodeSet{
			KeySATFormaPago: "99",
		},
	},
}
