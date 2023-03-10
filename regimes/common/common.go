// Package common provides re-usable regime related structures and data.
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

// Standard tax tags
const (
	TagSimplified    cbc.Key = "simplified"
	TagReverseCharge cbc.Key = "reverse-charge"
	TagCustomerRates cbc.Key = "customer-rates"
	TagSelfBilled    cbc.Key = "self-billed"
	TagPartial       cbc.Key = "partial"
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

// ComputeLuhnCheckDigit expects as argument a number string excluding the check
// digit. The returned integer should be checked against the check digit by the
// caller.
// Luhn Algorithm definition: https://en.wikipedia.org/wiki/Luhn_algorithm
func ComputeLuhnCheckDigit(number string) int {
	sum := 0
	pos := 0

	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')

		if pos%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		pos++
	}

	return (10 - (sum % 10)) % 10
}
