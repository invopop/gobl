package common

import (
	"regexp"
	"strings"

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

var (
	taxCodeBadCharsRegexp = regexp.MustCompile(`[^A-Z0-9]+`)
)

// NormalizeTaxIdentity removes any whitespace or separation characters and ensures all letters are
// uppercase.
func NormalizeTaxIdentity(tID *org.TaxIdentity) error {
	code := strings.ToUpper(tID.Code)
	code = taxCodeBadCharsRegexp.ReplaceAllString(code, "")
	tID.Code = code
	return nil
}
