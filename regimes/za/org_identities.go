package za

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// IdentityTypeCIPC represents a company registration number issued by the
// Companies and Intellectual Property Commission (CIPC). Section 32(4) of
// the Companies Act 71 of 2008 requires registered companies to display
// this number - in addition to their VAT number, where applicable - on
// invoices and other business correspondence, on penalty of a fine or
// imprisonment.
//
// Reference: https://marxgore.co.za/wp-content/uploads/2020/01/Section-32-Use-of-company-name-and-registration-number.pdf
const IdentityTypeCIPC cbc.Code = "CIPC"

// cipcRegexp matches CIPC company registration numbers in the format
// YYYY/NNNNNN/XX, e.g. 2020/123456/07 (the suffix identifies the type of
// company, 07 for a private company). The number is sequential plus a type
// suffix; like the VAT number, CIPC does not publish a check digit
// algorithm for it.
var cipcRegexp = regexp.MustCompile(`^\d{4}/\d{6}/\d{2}$`)

var identityTypeDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeCIPC,
		Name: i18n.String{
			i18n.EN: "CIPC Registration Number",
		},
		Desc: i18n.String{
			i18n.EN: "South African company registration number issued by the Companies and Intellectual Property Commission (CIPC), format YYYY/NNNNNN/XX.",
		},
	},
}

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				is.Func("is CIPC", isCIPCIdentity),
				rules.Field("code",
					rules.Assert("01", "invalid CIPC registration number format",
						is.MatchesRegexp(cipcRegexp),
					),
				),
			),
		),
	)
}

func isCIPCIdentity(val any) bool {
	id, _ := val.(*org.Identity)
	return id != nil && id.Type == IdentityTypeCIPC
}
