package favat

import (
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func normalizeTaxCombo(tc *tax.Combo) {
	switch tc.Key {
	case tax.KeyStandard:
		switch tc.Rate {
		case tax.RateGeneral:
			tc.Ext = tc.Ext.Set(ExtKeyTaxCategory, "1") // Base rate
		case tax.RateReduced:
			tc.Ext = tc.Ext.Set(ExtKeyTaxCategory, "2") // First reduced rate
		case tax.RateSuperReduced:
			tc.Ext = tc.Ext.Set(ExtKeyTaxCategory, "3") // Second reduced rate
		}
	case tax.KeyZero:
		tc.Ext = tc.Ext.Set(ExtKeyTaxCategory, "6.1") // 0% VAT
	case tax.KeyIntraCommunity:
		tc.Ext = tc.Ext.Set(ExtKeyTaxCategory, "6.2") // 0% VAT for intra-community supply of goods
	case tax.KeyExport:
		tc.Ext = tc.Ext.Set(ExtKeyTaxCategory, "6.3") // 0% VAT for export
	case tax.KeyExempt:
		tc.Ext = tc.Ext.Set(ExtKeyTaxCategory, "7") // Exempt
	case tax.KeyOutsideScope:
		tc.Ext = tc.Ext.Set(ExtKeyTaxCategory, "8") // Foreign sales
	case tax.KeyReverseCharge:
		tc.Ext = tc.Ext.Set(ExtKeyTaxCategory, "9") // EU Reverse Charge
	}
}

func validateTaxCombo(tc *tax.Combo) error {
	return validation.ValidateStruct(tc,
		validation.Field(&tc.Ext,
			tax.ExtensionsRequire(ExtKeyTaxCategory),
		),
	)
}
