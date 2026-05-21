package flow10

import (
	"slices"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/dgfip"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// allowedVATPercents is the whitelist of VAT percentages authorised on
// a Flow 10 invoice / payment (G1.24).
var allowedVATPercents = mustParsePercentages(
	"0%", "0.9%", "1.05%", "1.75%", "2.1%", "5.5%", "7%", "8.5%",
	"9.2%", "9.6%", "10%", "13%", "19.6%", "20%", "20.6%",
)

func mustParsePercentages(values ...string) []num.Percentage {
	out := make([]num.Percentage, len(values))
	for i, v := range values {
		p, err := num.PercentageFromString(v)
		if err != nil {
			panic(err)
		}
		out[i] = p
	}
	return out
}

func percentageInList(p num.Percentage, list []num.Percentage) bool {
	for _, a := range list {
		if p.Compare(a) == 0 {
			return true
		}
	}
	return false
}

// vatKeyToUNTDIDCategory maps each supported GOBL VAT rate key to its
// UNTDID 5305 category code.
var vatKeyToUNTDIDCategory = map[cbc.Key]cbc.Code{
	tax.KeyStandard:       "S",
	tax.KeyZero:           "Z",
	tax.KeyExempt:         "E",
	tax.KeyReverseCharge:  "AE",
	tax.KeyIntraCommunity: "K",
	tax.KeyExport:         "G",
	tax.KeyOutsideScope:   "O",
}

// advancePaymentDocumentTypes are the UNTDID 1001 codes representing
// advance-payment invoices (forbidden combined with B4/S4/M4 modes).
var advancePaymentDocumentTypes = []cbc.Code{
	"386", // Advance payment invoice
	"500", // Self-billed advance payment
	"503", // Self-billed credit for claim
}

// finalAfterAdvanceBillingModes are the billing-mode codes that mark
// an invoice as a "final invoice after down payment".
var finalAfterAdvanceBillingModes = []cbc.Code{
	dgfip.BillingModeB4, dgfip.BillingModeS4, dgfip.BillingModeM4,
}

// -- Normalisation --------------------------------------------------------

func normalizeInvoice(inv *bill.Invoice) {
	if inv == nil {
		return
	}
	normalizeInvoiceTaxCategories(inv)
	normalizeParty(inv.Supplier)
	normalizeParty(inv.Customer)
	if invoiceIsB2C(inv) {
		normalizeB2CCategoryOnInvoice(inv)
		return
	}
	normalizeBillingMode(inv)
}

// invoiceIsB2C reports whether the invoice is a business-to-consumer
// transaction. Flow 10 distinguishes B2C from B2B by the presence of
// a Customer party.
func invoiceIsB2C(inv *bill.Invoice) bool {
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

// normalizeInvoiceTaxCategories sets the UNTDID 5305 category extension
// on each VAT combo based on its rate key.
func normalizeInvoiceTaxCategories(inv *bill.Invoice) {
	for _, line := range inv.Lines {
		if line == nil {
			continue
		}
		for _, combo := range line.Taxes {
			if combo == nil || combo.Category != tax.CategoryVAT {
				continue
			}
			if combo.Ext.Get(untdid.ExtKeyTaxCategory) != "" {
				continue
			}
			if code, ok := vatKeyToUNTDIDCategory[combo.Key]; ok {
				combo.Ext = combo.Ext.Set(untdid.ExtKeyTaxCategory, code)
			}
		}
	}
}

// -- Rules ----------------------------------------------------------------

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Assert("01", "invoice must be in EUR or provide an exchange rate to EUR",
			currency.CanConvertTo(currency.EUR),
		),
		// B2C rules: category, supplier SIREN, VAT percent whitelist.
		rules.When(
			is.Func("B2C invoice", invoiceIsB2CAny),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("02", "B2C transaction category extension (fr-ctc-flow10-b2c-category) is required on B2C invoices (G1.68)",
						is.Func("has B2C category", extensionsHaveB2CCategory),
					),
				),
			),
			rules.Field("supplier",
				rules.Assert("03", "supplier is required on B2C invoice",
					is.Present,
				),
				rules.Assert("04", "supplier must have a SIREN identity (ISO/IEC 6523 scheme 0002) on a B2C invoice",
					is.Func("party has SIREN", partyHasSIREN),
				),
			),
			rules.Assert("05", "every VAT line percent must be one of the Flow 10 permitted values (G1.24): 0, 0.9, 1.05, 1.75, 2.1, 5.5, 7, 8.5, 9.2, 9.6, 10, 13, 19.6, 20, 20.6",
				is.Func("allowed Flow 10 VAT percents", invoiceVATPercentsAllowed),
			),
		),
		rules.Field("supplier",
			rules.Field("addresses",
				rules.Each(
					rules.Field("country",
						rules.Assert("06", "supplier address must include country",
							is.Present,
						),
					),
				),
			),
		),
		rules.Field("customer",
			rules.Field("addresses",
				rules.Each(
					rules.Field("country",
						rules.Assert("07", "customer address must include country",
							is.Present,
						),
					),
				),
			),
		),
		// B2B reporting: supplier + customer must each carry a legal
		// identity declaring an allowed ICD 6523 scheme, plus matching
		// TaxID when scheme is SIREN / EU-VAT.
		rules.When(
			is.Func("cross-border B2B invoice", invoiceIsCrossBorderB2BAny),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("08", "invoice document type must be one of the Flow 10 permitted UNTDID 1001 codes",
						is.Func("allowed Flow 10 document type", invoiceDocumentTypeAllowed),
					),
					rules.Assert("09", "billing mode extension (dgfip-billing-mode) is required (G1.02)",
						is.Func("has billing mode", extensionsHaveBillingMode),
					),
				),
			),
			rules.When(
				is.Func("billing mode is final-after-advance (B4/S4/M4)", invoiceIsFinalAfterAdvance),
				rules.Field("tax",
					rules.Field("ext",
						rules.Assert("10", "final-after-advance billing mode (B4/S4/M4) cannot be combined with an advance-payment document type (386/500/503) (G1.60)",
							is.Func("not advance-payment doc type", invoiceNotAdvancePaymentDocType),
						),
					),
				),
			),
			rules.Field("supplier",
				rules.Assert("11", "supplier is required for Flow 10 B2B invoice (G2.19)",
					is.Present,
				),
				rules.Assert("12", "supplier must declare a legal identity with an allowed ICD 6523 scheme (G2.19): 0002, 0223, 0227, 0228 or 0229",
					is.Func("party has allowed legal scheme", partyHasAllowedLegalScheme),
				),
				rules.Assert("13", "supplier TaxID is required when legal identity scheme is SIREN (0002) or EU VAT (0223) (G2.33)",
					is.Func("party has TaxID when required", partyHasTaxIDWhenRequired),
				),
			),
			rules.When(
				is.Func("invoice has exempt (E) VAT category", invoiceHasExemptCombo),
				rules.Assert("14", "supplier VAT ID or ordering.seller (tax representative) VAT ID is required when the invoice VAT breakdown contains an exempt (E) category",
					is.Func("supplier or tax rep has VAT ID", invoiceHasSellerVATIDForExempt),
				),
				rules.Assert("15", "invoice with an exempt (E) VAT category must include an exemption reason in tax.notes (key=exempt, non-empty text)",
					is.Func("has exempt tax note", invoiceHasExemptTaxNote),
				),
			),
			rules.Field("customer",
				rules.Assert("16", "customer is required for Flow 10 B2B invoice (G2.19)",
					is.Present,
				),
				rules.Assert("17", "customer must declare a legal identity with an allowed ICD 6523 scheme (G2.19): 0002, 0223, 0227, 0228 or 0229",
					is.Func("party has allowed legal scheme", partyHasAllowedLegalScheme),
				),
				rules.Assert("18", "customer TaxID is required when legal identity scheme is SIREN (0002) or EU VAT (0223) (G2.33)",
					is.Func("party has TaxID when required", partyHasTaxIDWhenRequired),
				),
			),
		),
	)
}

// -- Predicates ------------------------------------------------------------

func invoiceIsB2CAny(v any) bool {
	inv, ok := v.(*bill.Invoice)
	return ok && invoiceIsB2C(inv)
}

func invoiceIsCrossBorderB2BAny(v any) bool {
	inv, ok := v.(*bill.Invoice)
	return ok && !invoiceIsB2C(inv)
}

func invoiceVATPercentsAllowed(v any) bool {
	inv, ok := v.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	for _, line := range inv.Lines {
		if line == nil {
			continue
		}
		for _, combo := range line.Taxes {
			if combo == nil || combo.Category != tax.CategoryVAT || combo.Percent == nil {
				continue
			}
			if !percentageInList(*combo.Percent, allowedVATPercents) {
				return false
			}
		}
	}
	return true
}

func extensionsHaveB2CCategory(v any) bool {
	return extValue(v).Get(ExtKeyB2CCategory) != ""
}

func extensionsHaveBillingMode(v any) bool {
	return extValue(v).Get(dgfip.ExtKeyBillingMode) != ""
}

func invoiceIsFinalAfterAdvance(v any) bool {
	inv, ok := v.(*bill.Invoice)
	if !ok || inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
		return false
	}
	return slices.Contains(finalAfterAdvanceBillingModes, inv.Tax.Ext.Get(dgfip.ExtKeyBillingMode))
}

func invoiceNotAdvancePaymentDocType(v any) bool {
	return !slices.Contains(advancePaymentDocumentTypes, extValue(v).Get(untdid.ExtKeyDocumentType))
}

func invoiceDocumentTypeAllowed(v any) bool {
	return slices.Contains(allowedDocumentTypes, extValue(v).Get(untdid.ExtKeyDocumentType))
}

func invoiceHasSellerVATIDForExempt(v any) bool {
	inv, ok := v.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	if partyHasVATCode(inv.Supplier) {
		return true
	}
	if inv.Ordering != nil && partyHasVATCode(inv.Ordering.Seller) {
		return true
	}
	return false
}

func invoiceHasExemptCombo(v any) bool {
	inv, ok := v.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	for _, line := range inv.Lines {
		if line == nil {
			continue
		}
		for _, combo := range line.Taxes {
			if combo == nil || combo.Category != tax.CategoryVAT {
				continue
			}
			if combo.Ext.Get(untdid.ExtKeyTaxCategory) == "E" {
				return true
			}
		}
	}
	return false
}

func invoiceHasExemptTaxNote(v any) bool {
	inv, ok := v.(*bill.Invoice)
	if !ok || inv == nil || inv.Tax == nil {
		return false
	}
	for _, n := range inv.Tax.Notes {
		if n != nil && n.Key == tax.KeyExempt && n.Text != "" {
			return true
		}
	}
	return false
}

// quiet linter
var _ = org.Party{}
