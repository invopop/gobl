package eu

import "github.com/invopop/gobl/tax"

// Standard tax categories that may be shared between countries.
const (
	TaxCategoryVAT tax.Code = "VAT"
)

// Standard VAT tax codes
const (
	TaxRateVATZero         tax.Code = "zro"
	TaxRateVATStandard     tax.Code = "std"
	TaxRateVATReduced      tax.Code = "red"
	TaxRateVATSuperReduced tax.Code = "srd"
)
