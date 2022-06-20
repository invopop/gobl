package common

import (
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// Standard tax categories that may be shared between countries.
const (
	TaxCategoryVAT tax.Code = "VAT"
)

// Most commonly used codes. Local regions may add their own rate codes.
const (
	TaxRateZero         tax.Key = "zero"
	TaxRateStandard     tax.Key = "standard"
	TaxRateReduced      tax.Key = "reduced"
	TaxRateSuperReduced tax.Key = "super-reduced"
)

// Standard scheme definitions
const (
	SchemeReverseCharge tax.Key = "reverse-charge"
	SchemeCustomerRates tax.Key = "customer-rates"
)

// Common inbox keys
const (
	InboxKeyPEPPOL org.InboxKey = "peppol-id"
)
