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
	"github.com/invopop/gobl/tax"
)

const (
	// V2017 is the key for the EN16931-1:2017 specification.
	V2017 cbc.Key = "eu-en16931-v2017"
)

func init() {
	tax.RegisterAddonDef(newAddon())
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
			`),
		},
		Scenarios:  scenarios,
		Normalizer: normalize,
		Validator:  validate,
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
	case *org.Inbox:
		normalizeOrgInbox(obj)
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *pay.Instructions:
		return validatePayInstructions(obj)
	case *bill.Invoice:
		return validateBillInvoice(obj)
	case *bill.Line:
		return validateBillLine(obj)
	case *bill.PaymentDetails:
		return validateBillPayment(obj)
	case *tax.Combo:
		return validateTaxCombo(obj)
	case *bill.Discount:
		return validateBillDiscount(obj)
	case *bill.Charge:
		return validateBillCharge(obj)
	case *org.Item:
		return validateOrgItem(obj)
	case *org.Attachment:
		return validateOrgAttachment(obj)
	case *org.Party:
		return validateOrgParty(obj)
	case *org.Inbox:
		return validateOrgInbox(obj)
	}
	return nil
}
