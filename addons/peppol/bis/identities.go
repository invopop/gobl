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
	// SE-R-005. The addon's normalizer fills the boilerplate code when the key
	// is set without one. Lives here (not in regimes/se) because the assertion
	// is a Peppol-specific UBL artifact, not a general property of the SE regime.
	IdentityKeyFSkatt cbc.Key = "se-f-skatt"
)

// FSkattText is the literal Swedish boilerplate required by Peppol SE-R-005
// in the cac:PartyTaxScheme/cbc:CompanyID field.
const FSkattText = "Godkänd för F-skatt"

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

// normalizeIdentity fills addon-specific identity defaults — currently just
// the F-skatt boilerplate code when the key is set without one.
func normalizeIdentity(id *org.Identity) {
	if id == nil {
		return
	}
	if id.Key == IdentityKeyFSkatt && id.Code == "" {
		id.Code = FSkattText
	}
}
