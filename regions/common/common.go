package common

import "github.com/invopop/gobl/tax"

// Standard tax categories that may be shared between countries.
const (
	TaxCategoryVAT tax.Code = "VAT"
)

// Most common VAT codes used througout the European Union.
const (
	TaxRateVATStandard     tax.Code = "STD"
	TaxRateVATReduced      tax.Code = "RED"
	TaxRateVATSuperReduced tax.Code = "SRD"
	TaxRateVATZero         tax.Code = "ZRO"
)
