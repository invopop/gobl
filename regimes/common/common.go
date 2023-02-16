package common

import (
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

// Standard tax categories that may be shared between countries.
const (
	TaxCategoryST  cbc.Code = "ST"  // Sales Tax
	TaxCategoryVAT cbc.Code = "VAT" // Value Added Tax
	TaxCategoryGST cbc.Code = "GST" // Goods and Services Tax
)

// Most commonly used codes. Local regions may add their own rate codes.
const (
	TaxRateZero         cbc.Key = "zero"
	TaxRateStandard     cbc.Key = "standard"
	TaxRateIntermediate cbc.Key = "intermediate"
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

// Common Identity Type Codes that are not country specific.
const (
	IdentityTypeDUNS cbc.Code = "DUNS" // Dun & Bradstreet - Data Universal Numbering System
)

var (
	taxCodeBadCharsRegexp = regexp.MustCompile(`[^A-Z0-9]+`)
)

// NormalizeTaxIdentity removes any whitespace or separation characters and ensures all letters are
// uppercase.
func NormalizeTaxIdentity(tID *tax.Identity) error {
	code := strings.ToUpper(tID.Code.String())
	code = taxCodeBadCharsRegexp.ReplaceAllString(code, "")
	code = strings.TrimPrefix(code, string(tID.Country))
	tID.Code = cbc.Code(code)
	return nil
}
