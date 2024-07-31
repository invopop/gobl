// Package common provides re-usable regime related structures and data.
package common

import (
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
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
	// TaxCodeBadCharsRegexp is used to remove any characters that are not valid in a tax code.
	TaxCodeBadCharsRegexp = regexp.MustCompile(`[^A-Z0-9]+`)
)

// NormalizeTaxIdentity removes any whitespace or separation characters and ensures all letters are
// uppercase.
func NormalizeTaxIdentity(tID *tax.Identity, altCodes ...l10n.Code) error {
	code := strings.ToUpper(tID.Code.String())
	code = TaxCodeBadCharsRegexp.ReplaceAllString(code, "")
	code = strings.TrimPrefix(code, string(tID.Country))
	for _, alt := range altCodes {
		code = strings.TrimPrefix(code, string(alt))
	}
	tID.Code = cbc.Code(code)
	return nil
}
