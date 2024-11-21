package nfse

import (
	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

const (
	// ISSLiabilityDefault is the default value for the ISS liability extension
	ISSLiabilityDefault = "1" // Liable
)

func validateTaxCombo(tc *tax.Combo) error {
	return validation.ValidateStruct(tc,
		validation.Field(&tc.Ext,
			validation.When(tc.Category == br.TaxCategoryISS,
				tax.ExtensionsRequires(ExtKeyISSLiability),
			),
		),
	)
}

func normalizeTaxCombo(tc *tax.Combo) {
	if tc == nil || tc.Category != br.TaxCategoryISS {
		return
	}

	if !tc.Ext.Has(ExtKeyISSLiability) {
		if tc.Ext == nil {
			tc.Ext = make(tax.Extensions)
		}
		tc.Ext[ExtKeyISSLiability] = ISSLiabilityDefault
	}
}
