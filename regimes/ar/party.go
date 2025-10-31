package ar

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
)

// normalizeParty applies Argentine-specific normalization to party data
func normalizeParty(p *org.Party) {
	if p == nil {
		return
	}

	// Normalize tax identity if present
	if p.TaxID != nil {
		normalizeTaxIdentity(p.TaxID)
	}

	// Normalize addresses
	normalizePartyAddresses(p)
}

// normalizePartyAddresses applies normalization to party addresses
func normalizePartyAddresses(p *org.Party) {
	if p == nil || len(p.Addresses) == 0 {
		return
	}

	for _, addr := range p.Addresses {
		if addr == nil {
			continue
		}

		// Normalize postal codes - remove spaces and special characters
		if addr.Code != "" {
			code := addr.Code.String()
			normalized := ""
			for _, char := range code {
				if char >= '0' && char <= '9' {
					normalized += string(char)
				}
			}
			addr.Code = cbc.Code(normalized)
		}

		// Ensure country is set for Argentine parties
		if addr.Country == "" && p.TaxID != nil && p.TaxID.Country.Code() == "AR" {
			addr.Country = "AR"
		}
	}
}

// validateParty validates party data according to Argentine requirements
// This function is reserved for future use and can be extended with
// additional Argentine-specific party validations when needed.
// Currently, tax ID validation is handled separately by validateTaxIdentity.
//
//nolint:unused,unparam // Reserved for future use
func validateParty(p *org.Party) error {
	if p == nil {
		return nil
	}

	// Tax ID validation is handled separately by validateTaxIdentity
	// Additional party-specific validations can be added here as needed

	return nil
}
