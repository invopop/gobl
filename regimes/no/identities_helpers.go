package no

import (
	"strings"

	"github.com/invopop/gobl/cbc"
)

// cleanNorwayOrgNr extracts the organization number digits from common formats.
// Returns (digits, true) only when exactly 9 digits are present after normalization.
func cleanNorwayOrgNr(code cbc.Code) (string, bool) {
	s := strings.TrimSpace(strings.ToUpper(code.String()))
	digits := cbc.NormalizeNumericalCode(cbc.Code(s)).String()
	return digits, len(digits) == 9
}

// cleanNorwayTaxCode extracts the organization number digits from common Norwegian VAT formats.
// Accepts inputs like:
//   - "974760673"
//   - "974760673MVA"
//   - "NO974760673MVA"
//
// with optional spaces and mixed case.
// Returns (digits, true) only when exactly 9 digits are present after normalization.
func cleanNorwayTaxCode(code cbc.Code) (string, bool) {
	s := strings.TrimSpace(strings.ToUpper(code.String()))

	// Allow optional country prefix "NO" (common in e-invoicing payloads)
	s = strings.TrimPrefix(s, "NO")

	// Allow optional "MVA" suffix
	s = strings.TrimSuffix(s, "MVA")

	return cleanNorwayOrgNr(cbc.Code(s))
}
