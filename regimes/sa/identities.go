package sa

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Party identification scheme codes for seller (BT-29-1) and buyer (BT-46-1).
const (
	// IdentityTypeCRN is the Commercial Registration Number (10 digits).
	IdentityTypeCRN cbc.Code = "CRN"
	// IdentityTypeMOM is a MOMRAH / Ministry of Municipalities and Housing license.
	IdentityTypeMOM cbc.Code = "MOM"
	// IdentityTypeMLS is an MHRSD / Ministry of Human Resources and Social Development license.
	IdentityTypeMLS cbc.Code = "MLS"
	// IdentityType700 is the Unified Number (10 digits starting with 7).
	IdentityType700 cbc.Code = "700"
	// IdentityTypeSAG is a MISA / Ministry of Investment license.
	IdentityTypeSAG cbc.Code = "SAG"
	// IdentityTypeNAT is the Saudi National ID (10 digits starting with 1).
	IdentityTypeNAT cbc.Code = "NAT"
	// IdentityTypeIQA is the Iqama residency permit (10 digits starting with 2).
	IdentityTypeIQA cbc.Code = "IQA"
	// IdentityTypePAS is a passport number.
	IdentityTypePAS cbc.Code = "PAS"
	// IdentityTypeGCC is a GCC member state national ID.
	IdentityTypeGCC cbc.Code = "GCC"
	// IdentityTypeTIN is a Tax Identification Number (buyer-only).
	IdentityTypeTIN cbc.Code = "TIN"
	// IdentityTypeOTH is any other form of identification.
	IdentityTypeOTH cbc.Code = "OTH"
)

var (
	crnRegex    = regexp.MustCompile(`^\d{10}$`)
	num700Regex = regexp.MustCompile(`^7\d{9}$`)
	natRegex    = regexp.MustCompile(`^1\d{9}$`)
	iqaRegex    = regexp.MustCompile(`^2\d{9}$`)
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
	{
		Code: IdentityTypeNAT,
		Name: i18n.String{
			i18n.EN: "National ID",
			i18n.AR: "الهوية الوطنية",
		},
	},
	{
		Code: IdentityTypeIQA,
		Name: i18n.String{
			i18n.EN: "Iqama",
			i18n.AR: "الإقامة",
		},
	},
	{
		Code: IdentityTypePAS,
		Name: i18n.String{
			i18n.EN: "Passport",
			i18n.AR: "جواز السفر",
		},
	},
	{
		Code: IdentityTypeGCC,
		Name: i18n.String{
			i18n.EN: "GCC ID",
			i18n.AR: "هوية مواطني دول مجلس التعاون",
		},
	},
	{
		Code: IdentityTypeTIN,
		Name: i18n.String{
			i18n.EN: "Tax Identification Number",
			i18n.AR: "رقم التعريف الضريبي",
		},
	},
	{
		Code: IdentityTypeOTH,
		Name: i18n.String{
			i18n.EN: "Other ID",
			i18n.AR: "هوية أخرى",
		},
	},
}

// normalizeIdentity removes non-alphanumeric characters from identity codes.
func normalizeIdentity(id *org.Identity) {
	if id == nil || id.Code == "" {
		return
	}
	switch id.Type {
	case IdentityTypeCRN, IdentityType700, IdentityTypeNAT, IdentityTypeIQA:
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
		return validation.ValidateStruct(id,
			validation.Field(&id.Code, validation.Match(crnRegex).Error("must be a 10-digit number")),
		)
	case IdentityType700:
		return validation.ValidateStruct(id,
			validation.Field(&id.Code, validation.Match(num700Regex).Error("must be a 10-digit number starting with 7")),
		)
	case IdentityTypeNAT:
		return validation.ValidateStruct(id,
			validation.Field(&id.Code, validation.Match(natRegex).Error("must be a 10-digit number starting with 1")),
		)
	case IdentityTypeIQA:
		return validation.ValidateStruct(id,
			validation.Field(&id.Code, validation.Match(iqaRegex).Error("must be a 10-digit number starting with 2")),
		)
	default:
		return nil
	}
}
