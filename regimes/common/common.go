package common

import (
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

// Standard tax categories that may be shared between countries.
const (
	TaxCategoryVAT cbc.Code = "VAT"
)

// Most commonly used codes. Local regions may add their own rate codes.
const (
	TaxRateZero         cbc.Key = "zero"
	TaxRateStandard     cbc.Key = "standard"
	TaxRateReduced      cbc.Key = "reduced"
	TaxRateSuperReduced cbc.Key = "super-reduced"
)

// Standard scheme definitions
const (
	SchemeReverseCharge cbc.Key = "reverse-charge"
	SchemeCustomerRates cbc.Key = "customer-rates"
)

// Common inbox keys
const (
	InboxKeyPEPPOL cbc.Key = "peppol-id"
)

var (
	taxCodeBadCharsRegexp = regexp.MustCompile(`[^A-Z0-9]+`)
)

// NormalizeTaxIdentity removes any whitespace or separation characters and ensures all letters are
// uppercase.
func NormalizeTaxIdentity(tID *tax.Identity) error {
	code := strings.ToUpper(tID.Code)
	code = taxCodeBadCharsRegexp.ReplaceAllString(code, "")
	tID.Code = code
	return nil
}
