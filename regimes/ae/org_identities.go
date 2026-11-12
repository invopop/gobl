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

const (
	// tinPattern is the format of the TIN: the first 10 digits of the 15-digit
	// TRN, as defined by the Ministry of Finance's e-invoicing mandatory fields.
	tinPattern = `^\d{10}$`
	// emiratesIDPattern is the format of the Emirates ID number: 15 digits, as
	// stated on the UAE government portal.
	emiratesIDPattern = `^\d{15}$`
)

// identityDefinitions lists the identity keys the AE regime defines.
var identityDefinitions = []*cbc.Definition{
	{
		Key: IdentityKeyTIN,
		Name: i18n.String{
			i18n.EN: "Tax Identification Number",
			i18n.AR: "رقم التعريف الضريبي",
		},
		Desc: i18n.String{
			i18n.EN: "10-digit number issued by the Federal Tax Authority. Not the 15-digit VAT number (TRN), which goes in the party's tax ID.",
		},
		Sources: []*cbc.Source{mofMandatoryFieldsSource},
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
			i18n.EN: "15-digit number of the Emirates ID card, unique to its holder and kept for life.",
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
					rules.Assert("01", "UAE TIN must be 10 digits",
						is.Matches(tinPattern),
					),
				),
			),
			rules.When(
				is.Func("is Emirates ID", identityHasKey(IdentityKeyEmiratesID)),
				rules.Field("code",
					rules.Assert("02", "Emirates ID number must be 15 digits",
						is.Matches(emiratesIDPattern),
					),
				),
			),
		),
	)
}

// identityHasKey returns a predicate that reports whether an identity has the
// given key.
func identityHasKey(key cbc.Key) func(any) bool {
	return func(val any) bool {
		id, ok := val.(*org.Identity)
		return ok && id != nil && id.Key == key
	}
}
