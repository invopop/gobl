package eu

import "github.com/invopop/gobl/tax"

// Standard tax categories that may be shared between countries.
const (
	TaxCategoryVAT tax.Code = "VAT"
)

// Standard VAT tax codes
const (
	TaxRateVATStandard     tax.Code = "STD"
	TaxRateVATExempt       tax.Code = "EXPT"
	TaxRateVATZero         tax.Code = "ZERO"
	TaxRateVATReduced      tax.Code = "RED"
	TaxRateVATSuperReduced tax.Code = "SRED"
)
