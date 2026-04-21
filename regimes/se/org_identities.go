package se

import (
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/luhn"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// IdentityTypeOrgNr defines the key for the Swedish Organization Number (Organisationsnummer).
	IdentityTypeOrgNr cbc.Code = "ON" // Officially SE-ON
	// IdentityTypePersonNr defines the key for the Swedish Person Number (Personnummer).
	IdentityTypePersonNr cbc.Code = "PN" // Unofficial
	// IdentityTypeCoordinationNr defines the key for the Swedish Coordination Number (Samordningsnummer).
	IdentityTypeCoordinationNr cbc.Code = "CN" // Unofficial

	// IdentityKeyFSkatt marks a Swedish supplier as approved for F-tax (F-skatt).
	// In UBL, this maps to a non-VAT cac:PartyTaxScheme entry whose cbc:CompanyID
	// carries the boilerplate text "Godkänd för F-skatt", as required by Peppol
	// rule SE-R-005.
	IdentityKeyFSkatt cbc.Key = "se-f-skatt"
)

// FSkattText is the literal Swedish boilerplate required by Peppol SE-R-005
// in the cac:PartyTaxScheme/cbc:CompanyID field.
const FSkattText = "Godkänd för F-skatt"

var (
	// ValidOrgIdentityTypes defines the keys for the Swedish organization identities.
	ValidOrgIdentityTypes = []cbc.Code{IdentityTypeOrgNr, IdentityTypePersonNr, IdentityTypeCoordinationNr}

	// Regular expressions for validating Swedish identity codes
	orgNrRegex        = regexp.MustCompile(`^\d{10}$`)              // 0123456789
	individualNrRegex = regexp.MustCompile(`^\d{6}[\-\+]{1}\d{4}$`) // 010101-0101 or 010101+0101
)

var identityTypeDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeOrgNr,
		Name: i18n.String{
			i18n.EN: "Organization Number",
			i18n.SE: "Organisationsnummer",
		},
		Desc: i18n.String{
			i18n.EN: "Swedish company registration number.",
			i18n.SE: "Svenskt företagsregistreringsnummer.",
		},
	},
	{
		Code: IdentityTypePersonNr,
		Name: i18n.String{
			i18n.EN: "Person Number",
			i18n.SE: "Personnummer",
		},
		Desc: i18n.String{
			i18n.EN: "Swedish personal registration number.",
			i18n.SE: "Svenskt personnummer.",
		},
	},
	{
		Code: IdentityTypeCoordinationNr,
		Name: i18n.String{
			i18n.EN: "Coordination Number",
			i18n.SE: "Samordningsnummer",
		},
		Desc: i18n.String{
			i18n.EN: "Swedish coordination number.",
			i18n.SE: "Svenskt samordningsnummer.",
		},
	},
	{
		Key: IdentityKeyFSkatt,
		Name: i18n.String{
			i18n.EN: "F-Tax Approval",
			i18n.SE: "Godkänd för F-skatt",
		},
		Desc: i18n.String{
			i18n.EN: "Swedish F-tax (F-skatt) approval. Setting this identity on a " +
				"supplier asserts that the business handles its own tax payments. " +
				"Required for Peppol BIS Billing 3.0 (SE-R-005); rendered as a " +
				"non-VAT party tax scheme entry with the boilerplate text " +
				"\"Godkänd för F-skatt\".",
			i18n.SE: "Svenskt godkännande för F-skatt. När denna identitet anges " +
				"på en leverantör betyder det att verksamheten hanterar sina egna " +
				"skattebetalningar. Krävs för Peppol BIS Billing 3.0 (SE-R-005).",
		},
	},
}

// normalizeOrgIdentity performs normalization specific to Swedish identity codes.
//
//   - For organization numbers, it returns a 10 digit number, removing any separators.
//   - For individual numbers, it returns a 10 digit number with the separator. If none are present, a hyphen will be added. If a plus sign (`+`) is present anywhere, it will be used as the separator.
//
// If too many or too few numbers are present, it does nothing.
func normalizeOrgIdentity(id *org.Identity) {
	if id == nil {
		return
	}
	// Key-based identities (e.g. F-skatt) carry a fixed boilerplate code.
	if id.Key == IdentityKeyFSkatt && id.Code == "" {
		id.Code = FSkattText
		return
	}
	switch id.Type {
	case IdentityTypeOrgNr:
		// Organization numbers should be numeric only, with no separators
		code := cbc.NormalizeNumericalCode(id.Code).String()
		// Only if we have 12 digits, i.e the check digits are present
		// can we safely remove them
		if len(code) == taxCodeLength {
			code = strings.TrimSuffix(code, taxCodeCheckDigit)
		}

		// If we don't have the expected number of digits, it's likely not valid and no safe operation
		// can be performed.
		if len(code) != taxCodeLengthWithoutCheckDigits {
			return
		}

		id.Code = cbc.Code(code)

	case IdentityTypePersonNr, IdentityTypeCoordinationNr:
		// Individual numbers should maintain separator (- or +)
		code := strings.TrimSpace(id.Code.String())

		// If there's no separator but we have the right number of digits,
		// insert a hyphen at the right position, since it's the most
		// statistically likely separator.
		if len(code) == taxCodeLengthWithoutCheckDigits && !strings.ContainsAny(code, "-+") {
			code = code[:6] + "-" + code[6:]
		} else {
			// Extract digits and keep the separator
			digitsOnly := ""
			for _, c := range code {
				if c >= '0' && c <= '9' {
					digitsOnly += string(c)
				}
			}

			// If we don't have the expected number of digits, it's likely not valid and no safe operation
			// can be performed.
			if len(digitsOnly) != taxCodeLengthWithoutCheckDigits {
				return
			}

			// Determine the separator. If a plus sign is present anywhere, preserve it.
			separator := "-"
			if strings.Contains(code, "+") {
				separator = "+"
			}
			code = digitsOnly[:6] + separator + digitsOnly[6:]
		}

		id.Code = cbc.Code(code)

	default:
		return
	}
}

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				org.IdentityTypeIn(IdentityTypeOrgNr),
				rules.Field("code",
					rules.Assert("01", "invalid organization number format",
						is.Func("valid org number", orgNrCodeValid),
					),
					rules.Assert("02", "invalid checksum",
						is.Func("luhn checksum", orgNrChecksumValid),
					),
				),
			),
			rules.When(
				org.IdentityTypeIn(IdentityTypePersonNr, IdentityTypeCoordinationNr),
				rules.Field("code",
					rules.Assert("03", "invalid person or coordination number format",
						is.Func("valid individual number", individualNrCodeValid),
					),
					rules.Assert("04", "invalid checksum",
						is.Func("luhn checksum", individualNrChecksumValid),
					),
				),
			),
		),
	)
}

func orgNrCodeValid(val any) bool {
	code, ok := val.(cbc.Code)
	return ok && code != "" && orgNrRegex.MatchString(cbc.NormalizeNumericalCode(code).String())
}

func orgNrChecksumValid(val any) bool {
	code, ok := val.(cbc.Code)
	if !ok || code == "" || !orgNrRegex.MatchString(cbc.NormalizeNumericalCode(code).String()) {
		return true // skip if format invalid; format assertion handles that
	}
	return luhn.Check(cbc.NormalizeNumericalCode(code))
}

func individualNrCodeValid(val any) bool {
	code, ok := val.(cbc.Code)
	return ok && code != "" && individualNrRegex.MatchString(code.String())
}

func individualNrChecksumValid(val any) bool {
	code, ok := val.(cbc.Code)
	if !ok || code == "" || !individualNrRegex.MatchString(code.String()) {
		return true // skip if format invalid; format assertion handles that
	}
	digitsOnly := cbc.NormalizeNumericalCode(code)
	return luhn.Check(digitsOnly)
}
