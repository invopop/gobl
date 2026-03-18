package in

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
)

const (
	// IdentityTypePAN represents the Indian Permanent Account Number (PAN). It is a unique identifier assigned
	// to individuals, companies, and other entities.
	IdentityTypePAN cbc.Code = "PAN"

	// IdentityTypeHSN represents the Harmonized System of Nomenclature (HSN) code. It is used to classify products
	// or services for taxation purposes.
	//
	// HSN codes for India can be found using the online service here:
	// https://services.gst.gov.in/services/searchhsnsac
	//
	// The SAC (Service Accounting Code) is a similar classification system for services, which has been replaced
	// by the HSN code.
	IdentityTypeHSN cbc.Code = "HSN"
)

var (
	identityPatternPAN = `^[A-Z]{5}[0-9]{4}[A-Z]$`
	identityPatternHSN = `^(?:\d{4}|\d{6}|\d{8})$`
)

func normalizeOrgIdentity(id *org.Identity) {
	if id == nil {
		return
	}
	switch id.Type {
	case IdentityTypePAN:
		id.Code = cbc.NormalizeAlphanumericalCode(id.Code)
	case IdentityTypeHSN:
		id.Code = cbc.NormalizeNumericalCode(id.Code)
	}
}

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			org.IdentityTypeIn(IdentityTypePAN),
			rules.Field("code",
				rules.Assert("01", "identity code must be a valid PAN format",
					rules.Matches(identityPatternPAN)),
			),
		),
		rules.When(
			org.IdentityTypeIn(IdentityTypeHSN),
			rules.Field("code",
				rules.Assert("02", "identity code must be a valid HSN format",
					rules.Matches(identityPatternHSN)),
			),
		),
	)
}
