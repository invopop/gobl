package sa

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Identification keys used for additional codes not covered by the standard fields
const (
	IdentityTypeTIN      cbc.Code = "TIN" // Tax Identification Number
	IdentityTypeCRN      cbc.Code = "CRN" // Commercial Registration Number
	IdentityTypeMom      cbc.Code = "MOM" // Ministry of Municipal, Rural Affairs and Housing Number
	IdentityTypeMLS      cbc.Code = "MLS" // Ministry of Human Resources and Social Development Number
	IdentityType700      cbc.Code = "700" // 700 Number
	IdentityTypeSAG      cbc.Code = "SAG" // Saudi Arabian General Authority Number
	IdentityTypeNational cbc.Code = "NAT" // National ID
	IdentityTypeGcc      cbc.Code = "GCC" // GCC ID
	IdentityTypeIqa      cbc.Code = "IQA" // Iqama Number (Resident ID)
	IdentityTypePassport cbc.Code = "PAS" // Passport Number
	IdentityTypeOTH      cbc.Code = "OTH" // Other ID
)

var (
	identitiesRegexp = regexp.MustCompile(`^[a-zA-Z0-9]*$`)

	// CustomerValidIdentities groups customer accepted identities by ZATCA
	CustomerValidIdentities = []cbc.Code{
		IdentityTypeTIN,
		IdentityTypeCRN,
		IdentityTypeMom,
		IdentityTypeMLS,
		IdentityType700,
		IdentityTypeSAG,
		IdentityTypeNational,
		IdentityTypeGcc,
		IdentityTypeIqa,
		IdentityTypePassport,
		IdentityTypeOTH,
	}

	supplierValidIdentities = []cbc.Code{
		IdentityTypeCRN,
		IdentityTypeMom,
		IdentityTypeMLS,
		IdentityType700,
		IdentityTypeSAG,
		IdentityTypeOTH,
	}
)

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(countryCode)),
			rules.Field("code",
				rules.Assert("01", "identity code must be valid",
					is.MatchesRegexp(identitiesRegexp),
				),
			),
		),
	)
}
