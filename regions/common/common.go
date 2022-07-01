package common

import (
	"github.com/invopop/gobl/org"
)

// Standard tax categories that may be shared between countries.
const (
	TaxCategoryVAT org.Code = "VAT"
)

// Most commonly used codes. Local regions may add their own rate codes.
const (
	TaxRateZero         org.Key = "zero"
	TaxRateStandard     org.Key = "standard"
	TaxRateReduced      org.Key = "reduced"
	TaxRateSuperReduced org.Key = "super-reduced"
)

// Standard scheme definitions
const (
	SchemeReverseCharge org.Key = "reverse-charge"
	SchemeCustomerRates org.Key = "customer-rates"
)

// Common inbox keys
const (
	InboxKeyPEPPOL org.Key = "peppol-id"
)
