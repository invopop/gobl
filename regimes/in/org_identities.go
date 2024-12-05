package in

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
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
	identityRegexpPAN = regexp.MustCompile(`^[A-Z]{5}[0-9]{4}[A-Z]$`)
	identityRegexpHSN = regexp.MustCompile(`^(?:\d{4}|\d{6}|\d{8})$`)
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

func validateOrgIdentity(id *org.Identity) error {
	if id == nil {
		return nil
	}
	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.When(
				id.Type == IdentityTypePAN,
				validation.Match(identityRegexpPAN),
			),
			validation.When(
				id.Type == IdentityTypeHSN,
				validation.Match(identityRegexpHSN).Error("must be a 4, 6, or 8 digit number"),
			),
			validation.Skip,
		),
	)
}
