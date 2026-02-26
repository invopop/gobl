package sa

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Seller identification scheme codes used for Saudi business registrations.
//
// Source: ZATCA XML Implementation Standard v1.2
// https://zatca.gov.sa/ar/E-Invoicing/SystemsDevelopers/Documents/20230519_ZATCA_Electronic_Invoice_XML_Implementation_Standard_%20vF.pdf
const (
	// IdentityTypeCRN is the Commercial Registration Number (10 digits)
	// issued by the Ministry of Commerce.
	IdentityTypeCRN cbc.Code = "CRN"
	// IdentityTypeMOM is a MOMRAH / Ministry of Municipalities and Housing license.
	IdentityTypeMOM cbc.Code = "MOM"
	// IdentityTypeMLS is an MHRSD / Ministry of Human Resources and Social Development license.
	IdentityTypeMLS cbc.Code = "MLS"
	// IdentityType700 is the Unified Number (10 digits starting with 7).
	IdentityType700 cbc.Code = "700"
	// IdentityTypeSAG is a MISA / Ministry of Investment license.
	IdentityTypeSAG cbc.Code = "SAG"
)

var (
	crnRegex    = regexp.MustCompile(`^\d{10}$`)
	num700Regex = regexp.MustCompile(`^7\d{9}$`)
)

var identityDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeCRN,
		Name: i18n.String{
			i18n.EN: "Commercial Registration Number",
			i18n.AR: "رقم السجل التجاري",
		},
	},
	{
		Code: IdentityTypeMOM,
		Name: i18n.String{
			i18n.EN: "MOMRAH License",
			i18n.AR: "ترخيص وزارة البلديات والإسكان",
		},
	},
	{
		Code: IdentityTypeMLS,
		Name: i18n.String{
			i18n.EN: "MHRSD License",
			i18n.AR: "ترخيص وزارة الموارد البشرية والتنمية الاجتماعية",
		},
	},
	{
		Code: IdentityType700,
		Name: i18n.String{
			i18n.EN: "700 Number",
			i18n.AR: "الرقم الموحد",
		},
	},
	{
		Code: IdentityTypeSAG,
		Name: i18n.String{
			i18n.EN: "MISA License",
			i18n.AR: "ترخيص وزارة الاستثمار",
		},
	},
}

// normalizeIdentity removes non-alphanumeric characters from identity codes.
func normalizeIdentity(id *org.Identity) {
	if id == nil || id.Code == "" {
		return
	}
	switch id.Type {
	case IdentityTypeCRN, IdentityType700:
		code := tax.IdentityCodeBadCharsRegexp.ReplaceAllString(id.Code.String(), "")
		id.Code = cbc.Code(code)
	}
}

// validateIdentity checks to ensure the organization identity is valid.
func validateIdentity(id *org.Identity) error {
	if id == nil {
		return nil
	}
	switch id.Type {
	case IdentityTypeCRN:
		return validation.Match(crnRegex).Error("must be a 10-digit number").Validate(id.Code)
	case IdentityType700:
		return validation.Match(num700Regex).Error("must be a 10-digit number starting with 7").Validate(id.Code)
	default:
		return nil
	}
}
