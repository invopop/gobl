package bis

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
)

const (
	// IdentityKeyGreekMARK identifies a Greek MARK (Mydata Auto-Registration Key)
	// number issued by the Greek Independent Authority for Public Revenue.
	IdentityKeyGreekMARK cbc.Key = "gr-mark"

	// IdentityKeyFSkatt marks a Swedish supplier as approved for F-tax (F-skatt).
	// In UBL it surfaces as a non-VAT cac:PartyTaxScheme entry whose cbc:CompanyID
	// carries the literal text "Godkänd för F-skatt", as required by Peppol rule
	// SE-R-005. Lives here (not in regimes/se) because the assertion is a
	// Peppol-specific UBL artifact, not a general property of the SE regime.
	IdentityKeyFSkatt cbc.Key = "se-f-skatt"
)

// FSkattText is the literal Swedish boilerplate required by Peppol SE-R-005
// in the cac:PartyTaxScheme/cbc:CompanyID field.
const FSkattText = "Godkänd för F-skatt"

// FSkattTaxSchemeID is the non-VAT cac:TaxScheme/cbc:ID emitted alongside the
// F-skatt boilerplate. The schematron's only constraint is that the scheme
// be non-VAT (case-insensitive).
const FSkattTaxSchemeID cbc.Code = "TAX"

// identities lists the party identity types recognised by the Peppol addon.
var identities = []*cbc.Definition{
	{
		Key: IdentityKeyGreekMARK,
		Name: i18n.String{
			i18n.EN: "Greek MARK Number",
		},
		Desc: i18n.String{
			i18n.EN: "Mydata Auto-Registration Key assigned by the Greek Independent " +
				"Authority for Public Revenue (IAPR) to each registered invoice.",
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

// normalizeIdentity fills addon-specific identity defaults.
//
// For IdentityKeyFSkatt we populate Scope=tax, a non-VAT Type, and the
// Swedish boilerplate Code. These three fields together drive gobl.ubl's
// tax-scope identity path to emit the second cac:PartyTaxScheme block that
// Peppol SE-R-005 requires, without any converter-side changes.
func normalizeIdentity(id *org.Identity) {
	if id == nil {
		return
	}
	if id.Key == IdentityKeyFSkatt {
		if id.Scope == "" {
			id.Scope = org.IdentityScopeTax
		}
		if id.Type == "" {
			id.Type = FSkattTaxSchemeID
		}
		if id.Code == "" {
			id.Code = FSkattText
		}
	}
}
