package ae

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Identity types for the UAE identification numbers that can appear on
// invoices. Except for the TIN, they are the codes PINT AE uses for the types
// of commercial registration number.
const (
	// IdentityTypeTIN is the 10-digit Tax Identification Number the Federal
	// Tax Authority (FTA) issues to businesses. It is not the 15-digit VAT
	// number (TRN), which goes in the party's tax ID.
	IdentityTypeTIN cbc.Code = "TIN"
	// IdentityTypeTradeLicense is the number of the commercial or trade
	// licence issued by an emirate's licensing authority.
	IdentityTypeTradeLicense cbc.Code = "TL"
	// IdentityTypeEmiratesID is the number of the Emirates ID card issued to
	// UAE citizens and residents.
	IdentityTypeEmiratesID cbc.Code = "EID"
	// IdentityTypeCabinetDecision is the number of the Cabinet Decision that
	// establishes an entity.
	IdentityTypeCabinetDecision cbc.Code = "CD"
)

const (
	// tinPattern is the format of the TIN, as the Federal Tax Authority
	// registered it under ISO 6523 (ICD 0235): 10 digits, the first one 1,
	// the UAE's identifier under the GCC agreement. The last one is a Luhn
	// check digit.
	tinPattern = `^1\d{9}$`
	// emiratesIDPattern is the format of the Emirates ID number: 15 digits,
	// the first three 784, the international code of the UAE. The last one is
	// a verification digit whose algorithm is not published, so it is not
	// checked.
	emiratesIDPattern = `^784\d{12}$`
)

// identityDefinitions lists the identity types the AE regime defines.
var identityDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeTIN,
		Name: i18n.String{
			i18n.EN: "Tax Identification Number",
			i18n.AR: "رقم التعريف الضريبي",
		},
		Desc: i18n.String{
			i18n.EN: "10-digit number issued by the Federal Tax Authority, starting with 1 and ending in a Luhn check digit. Not the 15-digit VAT number (TRN), which starts with it and goes in the party's tax ID.",
		},
		Sources: []*cbc.Source{
			mofMandatoryFieldsSource,
			{
				Title: i18n.NewString("Peppol Participant Identifier Schemes - AE:TIN (ISO 6523 ICD 0235)"),
				URL:   "https://docs.peppol.eu/edelivery/codelists/",
			},
		},
	},
	{
		Code: IdentityTypeTradeLicense,
		Name: i18n.String{
			i18n.EN: "Trade License Number",
			i18n.AR: "رقم الرخصة التجارية",
		},
		Sources: []*cbc.Source{mofMandatoryFieldsSource},
	},
	{
		Code: IdentityTypeEmiratesID,
		Name: i18n.String{
			i18n.EN: "Emirates ID Number",
			i18n.AR: "رقم الهوية الإماراتية",
		},
		Desc: i18n.String{
			i18n.EN: "15-digit number of the Emirates ID card, starting with 784, unique to its holder and kept for life.",
		},
		Sources: []*cbc.Source{
			mofMandatoryFieldsSource,
			{
				Title: i18n.NewString("UAE Government Portal - Emirates ID"),
				URL:   "https://u.ae/en/information-and-services/visa-and-emirates-id/emirates-id",
			},
			{
				Title: i18n.NewString("Gulf News - ID number may not be amended, EIDA says"),
				URL:   "https://gulfnews.com/uae/government/id-number-may-not-be-amended-eida-says-1.1654705",
			},
		},
	},
	{
		Code: IdentityTypeCabinetDecision,
		Name: i18n.String{
			i18n.EN: "Cabinet Decision Number",
			i18n.AR: "رقم قرار مجلس الوزراء",
		},
		Sources: []*cbc.Source{mofMandatoryFieldsSource},
	},
}

// mofMandatoryFieldsSource is the Ministry of Finance's list of mandatory
// e-invoice fields. It defines the TIN and lists the types of commercial
// registration number issued in the UAE: trade licence, Emirates ID, passport
// and Cabinet Decision.
var mofMandatoryFieldsSource = &cbc.Source{
	Title: i18n.NewString("Ministry of Finance - UAE Electronic Invoice Mandatory Field Requirements"),
	URL:   "https://mof.gov.ae/wp-content/uploads/2026/10/UAE-Electronic-Invoice-Mandatory-Field-Requirements.pdf",
}

// orgIdentityRules checks the format of the UAE identities that have one.
func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				org.IdentityTypeIn(IdentityTypeTIN),
				rules.Field("code",
					rules.Assert("01", "UAE TIN must be 10 digits starting with 1",
						is.Matches(tinPattern),
					),
					rules.When(
						is.Matches(tinPattern),
						rules.AssertIfPresent("02", "UAE TIN must end in a valid check digit",
							is.Func("valid TIN check digit", tinCheckDigitIsValid),
						),
						rules.AssertIfPresent("04", "UAE TIN must have a business identifier that is not all zeros",
							is.Func("TIN business identifier set", tinBusinessIDIsSet),
						),
					),
				),
			),
			rules.When(
				org.IdentityTypeIn(IdentityTypeEmiratesID),
				rules.Field("code",
					rules.Assert("03", "Emirates ID number must be 15 digits starting with 784",
						is.Matches(emiratesIDPattern),
					),
				),
			),
		),
	)
}
