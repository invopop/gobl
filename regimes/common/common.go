// Package common provides re-usable regime related structures and data.
package common

import (
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
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
