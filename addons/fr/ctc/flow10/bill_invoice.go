package flow10

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/dgfip"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// -- Normalisation --------------------------------------------------------

func normalizeInvoice(inv *bill.Invoice) {
	if inv == nil {
		return
	}
	normalizeParty(inv.Supplier)
	normalizeParty(inv.Customer)
	if invoiceIsB2CDoc(inv) {
		normalizeB2CCategoryOnInvoice(inv)
		return
	}
	normalizeBillingMode(inv)
}

// invoiceIsB2CDoc reports whether the invoice is a B2C transaction —
// Flow 10 distinguishes B2C from B2B by the absence of a Customer party.
func invoiceIsB2CDoc(inv *bill.Invoice) bool {
	return inv != nil && inv.Customer == nil
}

// normalizeBillingMode picks a sensible default for the billing-mode
// extension when the caller hasn't supplied one. M2 when the invoice
// is fully paid, M1 otherwise.
func normalizeBillingMode(inv *bill.Invoice) {
	if inv.Tax != nil && !inv.Tax.Ext.IsZero() && inv.Tax.Ext.Get(dgfip.ExtKeyBillingMode) != "" {
		return
	}
	mode := dgfip.BillingModeM1
	if inv.Totals != nil && inv.Totals.Paid() {
		mode = dgfip.BillingModeM2
	}
	if inv.Tax == nil {
		inv.Tax = &bill.Tax{}
	}
	inv.Tax.Ext = inv.Tax.Ext.Set(dgfip.ExtKeyBillingMode, mode)
}

// normalizeB2CCategoryOnInvoice defaults the B2C transaction category
// to TNT1 when the caller has not supplied one.
func normalizeB2CCategoryOnInvoice(inv *bill.Invoice) {
	if inv.Tax != nil && inv.Tax.Ext.Get(ExtKeyB2CCategory) != "" {
		return
	}
	if inv.Tax == nil {
		inv.Tax = &bill.Tax{}
	}
	inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyB2CCategory, B2CCategoryNotTaxable)
}

// -- Rule set -------------------------------------------------------------

// billInvoiceRules validates only the integrity of the addon's own
// extensions: the B2C category and UNTDID document type, when present,
// must be recognised Flow 10 values. The e-reporting business rules
// (G1.*/G2.*) are the converter's responsibility — see the package doc.
func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Field("tax",
			rules.Field("ext",
				rules.Assert("01", "invoice tax ext fr-ctc-flow10-b2c-category must be a recognised B2C transaction category",
					tax.ExtensionHasValidCode(ExtKeyB2CCategory),
				),
				rules.Assert("02", "invoice tax ext untdid-document-type must be a recognised Flow 10 document type",
					tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, allowedDocumentTypes...),
				),
			),
		),
	)
}
