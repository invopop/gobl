package mydata

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Regime Specific Payment Means Extension Keys
const (
	MeansKeyForeign cbc.Key = "foreign"
)

// PaymentMeansExtensions returns the mapping of payment means to their
// extension values used by myDATA.
func PaymentMeansExtensions() tax.Extensions {
	return paymentMeansMap
}

var paymentMeansMap = tax.Extensions{
	pay.MeansKeyCreditTransfer:                       "1",
	pay.MeansKeyCreditTransfer.With(MeansKeyForeign): "2",
	pay.MeansKeyCash:                                 "3",
	pay.MeansKeyCheque:                               "4",
	pay.MeansKeyPromissoryNote:                       "5",
	pay.MeansKeyOnline:                               "6",
	pay.MeansKeyCard:                                 "7",
}

func normalizePayInstructions(i *pay.Instructions) {
	if i == nil {
		return
	}
	extVal := paymentMeansMap[i.Key]
	if extVal != "" {
		if i.Ext == nil {
			i.Ext = make(tax.Extensions)
		}
		i.Ext[ExtKeyPaymentMeans] = extVal
	}
}

func normalizePayAdvance(a *pay.Advance) {
	if a == nil {
		return
	}
	extVal := paymentMeansMap[a.Key]
	if extVal != "" {
		if a.Ext == nil {
			a.Ext = make(tax.Extensions)
		}
		a.Ext[ExtKeyPaymentMeans] = extVal
	}
}

func validatePayInstructions(value any) error {
	i, ok := value.(*pay.Instructions)
	if !ok || i == nil {
		return nil
	}
	return validation.ValidateStruct(i,
		validation.Field(&i.Key,
			validation.Required,
			validation.Skip,
		),
		validation.Field(&i.Ext,
			tax.ExtensionsRequire(ExtKeyPaymentMeans),
			validation.Skip,
		),
	)
}

func validatePayAdvance(value any) error {
	a, ok := value.(*pay.Advance)
	if !ok || a == nil {
		return nil
	}
	return validation.ValidateStruct(a,
		validation.Field(&a.Key,
			validation.Required,
			validation.Skip,
		),
		validation.Field(&a.Ext,
			tax.ExtensionsRequire(ExtKeyPaymentMeans),
			validation.Skip,
		),
	)
}
