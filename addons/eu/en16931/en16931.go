// Package en16931 defines an addon that will apply rules from the EN 16931 specification to
// GOBL documents.
package en16931

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// V2017 is the key for the EN16931-1:2017 specification.
	V2017 cbc.Key = "eu-en16931-v2017"
)

func init() {
	tax.RegisterAddonDef(newAddon())
	rules.RegisterWithGuard(
		V2017.String(),
		rules.GOBL.Add("EU-EN16931-V2017"),
		is.HasContext(tax.AddonIn(V2017)),
		billInvoiceRules(),
		billDiscountRules(),
		billLineDiscountRules(),
		billChargeRules(),
		billLineChargeRules(),
		payInstructionsRules(),
		payTermsRules(),
		orgItemRules(),
		orgAttachmentRules(),
		orgPartyRules(),
		orgInboxRules(),
		orgAddressRules(),
		taxComboRules(),
	)
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V2017,
		Name: i18n.String{
			i18n.EN: "EN 16931-1:2017",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Support for the European Norm (EN) 16931-1:2017 standard for electronic invoicing.

				This addon ensures the basic rules and mappings are applied to the GOBL document
				ensure that it is compliant and easily convertible to other formats.

				We strongly recommend checking the output and specifically the extension codes
				used to ensure that any assumptions do not need be adjusted.

				## Tax Code Extension Mappings

				The following tables show how GOBL tax keys/categories are mapped to UNTDID 5305 tax category codes:

				### VAT

				| GOBL Tax Key | UNTDID 5305 Code | Description |
				|--------------|------------------|-------------|
				| standard | S | Standard rate |
				| zero | Z | Zero rated goods |
				| exempt | E | Exempt from tax |
				| reverse-charge | AE | VAT Reverse Charge |
				| intra-community | K | Intra-community supply |
				| export | G | Export outside the EU |
				| outside-scope | O | Not subject to VAT |

				### Other

				For Spanish special territories, **IGIC** (Canary Islands) maps to code **L** and **IPSI** (Ceuta and Melilla) maps to code **M**.
				Any other tax category defaults to UNTDID 5305 code **O** (Outside Scope).
			`),
		},
		Scenarios:  scenarios,
		Normalizer: normalize,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeBillInvoice(obj)
	case *pay.Instructions:
		normalizePayInstructions(obj)
	case *tax.Combo:
		normalizeTaxCombo(obj)
	case *bill.Discount:
		normalizeBillDiscount(obj)
	case *bill.LineDiscount:
		normalizeBillLineDiscount(obj)
	case *bill.Charge:
		normalizeBillCharge(obj)
	case *bill.LineCharge:
		normalizeBillLineCharge(obj)
	case *org.Note:
		normalizeOrgNote(obj)
	case *org.Item:
		normalizeOrgItem(obj)
	case *org.Identity:
		normalizeOrgIdentity(obj)
	case *org.Inbox:
		normalizeOrgInbox(obj)
	}
}
