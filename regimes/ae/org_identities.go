package ae

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Identity keys for the UAE identification numbers that can appear on
// invoices.
const (
	// IdentityKeyTIN is the 10-digit Tax Identification Number the Federal Tax
	// Authority (FTA) issues to businesses. It is not the 15-digit VAT number
	// (TRN), which goes in the party's tax ID.
	IdentityKeyTIN cbc.Key = "ae-tin"
	// IdentityKeyTradeLicense is the number of the commercial or trade licence
	// issued by an emirate's licensing authority.
	IdentityKeyTradeLicense cbc.Key = "ae-trade-license"
	// IdentityKeyEmiratesID is the number of the Emirates ID card issued to
	// UAE citizens and residents.
	IdentityKeyEmiratesID cbc.Key = "ae-emirates-id"
	// IdentityKeyCabinetDecision is the number of the Cabinet Decision that
	// establishes an entity.
	IdentityKeyCabinetDecision cbc.Key = "ae-cabinet-decision"
)

// emiratesIDPattern is the format of the Emirates ID number: 15 digits, the
// first three 784, the international code of the UAE. The last one is a
// verification digit whose algorithm is not published, so it is not checked.
const emiratesIDPattern = `^784\d{12}$`

// identityDefinitions lists the identity keys the AE regime defines.
var identityDefinitions = []*cbc.Definition{
	{
		Key: IdentityKeyTIN,
		Name: i18n.String{
			i18n.EN: "Tax Identification Number",
			i18n.AR: "رقم التعريف الضريبي",
		},
		Desc: i18n.String{
			i18n.EN: "10-digit number issued by the Federal Tax Authority: a 1, an 8-digit business identifier and a Luhn check digit. Not the 15-digit VAT number (TRN), which starts with it and goes in the party's tax ID.",
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
		Key: IdentityKeyTradeLicense,
		Name: i18n.String{
			i18n.EN: "Trade License Number",
			i18n.AR: "رقم الرخصة التجارية",
		},
		Sources: []*cbc.Source{mofMandatoryFieldsSource},
	},
	{
		Key: IdentityKeyEmiratesID,
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
		Key: IdentityKeyCabinetDecision,
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
				is.Func("is TIN", identityHasKey(IdentityKeyTIN)),
				rules.Field("code",
					rules.Assert("01", "UAE TIN must be 10 digits starting with 1, with a valid check digit",
						is.Func("valid TIN", isValidTIN),
					),
				),
			),
			rules.When(
				is.Func("is Emirates ID", identityHasKey(IdentityKeyEmiratesID)),
				rules.Field("code",
					rules.Assert("02", "Emirates ID number must be 15 digits starting with 784",
						is.Matches(emiratesIDPattern),
					),
				),
			),
		),
	)
}

// isValidTIN reports whether an identity code is a valid TIN.
func isValidTIN(val any) bool {
	code, ok := val.(cbc.Code)
	return ok && validTIN(code)
}

// identityHasKey returns a predicate that reports whether an identity has the
// given key.
func identityHasKey(key cbc.Key) func(any) bool {
	return func(val any) bool {
		id, ok := val.(*org.Identity)
		return ok && id != nil && id.Key == key
	}
}
