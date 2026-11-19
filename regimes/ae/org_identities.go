package ae

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Identity types for UAE numbers on invoices.
const (
	IdentityTypeTIN             cbc.Code = "TIN" // the FTA's 10-digit TIN, not the 15-digit TRN of the tax ID
	IdentityTypeTradeLicense    cbc.Code = "TL"
	IdentityTypeEmiratesID      cbc.Code = "EID"
	IdentityTypeCabinetDecision cbc.Code = "CD"
)

// TINScheme is the ISO 6523 ICD the Federal Tax Authority registered for the
// TIN, which is also the scheme of UAE Peppol endpoints (0235:<TIN>).
const TINScheme cbc.Code = "0235"

const (
	// TINPattern is the TIN format: 10 digits starting with 1. Its Luhn check digit is checked apart.
	TINPattern = `^1\d{9}$`
	// emiratesIDPattern is the Emirates ID format; its check digit algorithm isn't published.
	emiratesIDPattern = `^784\d{12}$`
)

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
				Title: i18n.NewString("Peppol Code Lists v9.7 - Participant Identifier Schemes - AE:TIN (ISO 6523 ICD 0235)"),
				URL:   "https://docs.peppol.eu/edelivery/codelists/v9.7/Peppol%20Code%20Lists%20-%20Participant%20identifier%20schemes%20v9.7.json",
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

// mofMandatoryFieldsSource is the MoF list of mandatory e-invoice fields.
var mofMandatoryFieldsSource = &cbc.Source{
	Title: i18n.NewString("Ministry of Finance - UAE Electronic Invoice Mandatory Field Requirements"),
	URL:   "https://mof.gov.ae/wp-content/uploads/2026/10/UAE-Electronic-Invoice-Mandatory-Field-Requirements.pdf",
}

// orgIdentityRules validates identities on UAE documents.
func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				org.IdentityTypeIn(IdentityTypeTIN),
				rules.Field("code",
					rules.Assert("01", "UAE TIN must be 10 digits starting with 1",
						is.Matches(TINPattern),
					),
					rules.When(
						is.Matches(TINPattern),
						rules.AssertIfPresent("02", "UAE TIN must end in a valid check digit",
							is.Func("valid TIN check digit", tinCheckDigitIsValid),
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
