package flow10

import (
	"slices"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// A VAT rate key other than "standard" or "zero" is treated as an
// exemption for Flow 10 purposes: it must be paired with a matching
// exemption reason (tax.Note with the same Key and non-empty Text).
// Translating the key to the final UNTDID category code is the
// converter's job, not this addon's.

// finalAfterAdvanceBillingModes are the billing-mode codes that mark an
// invoice as a "final invoice after down payment" (G1.60): B4, S4, M4.
// Under these modes the invoice may not be an advance-payment document
// type (386/500/503).
var finalAfterAdvanceBillingModes = []cbc.Code{
	BillingModeB4, BillingModeS4, BillingModeM4,
}

// advancePaymentDocumentTypes are the UNTDID 1001 codes representing
// advance-payment invoices and their credit memo (G1.60 forbids them
// combined with B4/S4/M4 billing modes).
var advancePaymentDocumentTypes = []cbc.Code{
	"386", // Advance payment invoice
	"500", // Self-billed advance payment
	"503", // Down-payment credit memo
}

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		// B2C rules: category, supplier SIREN, VAT rate whitelist.
		rules.When(
			is.Func("B2C invoice", invoiceIsB2CAny),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("16", "B2C transaction category extension (fr-ctc-b2c-category) is required on B2C invoices (G1.68)",
						is.Func("has B2C category", extensionsHaveB2CCategory),
					),
				),
			),
			rules.Field("supplier",
				rules.Assert("17", "supplier is required on B2C invoice",
					is.Present,
				),
				rules.Assert("18", "supplier must have a SIREN identity (ISO/IEC 6523 scheme 0002) on a B2C invoice",
					is.Func("party has SIREN", partyHasSIREN),
				),
			),
			rules.Assert("19", "every VAT line rate must be one of the Flow 10 permitted percentages (G1.24): 0, 0.9, 1.05, 1.75, 2.1, 5.5, 7, 8.5, 9.2, 9.6, 10, 13, 19.6, 20, 20.6",
				is.Func("allowed Flow 10 VAT rates", invoiceVATRatesAllowed),
			),
		),
		// Flow 10 reports to the French authority in EUR: if the invoice is
		// issued in a different currency, an exchange rate must be provided.
		rules.Assert("10", "invoice must be in EUR or provide an exchange rate to EUR",
			currency.CanConvertTo(currency.EUR),
		),
		// When a party carries any postal address, the country on that
		// address must be populated. The address itself remains optional.
		rules.Field("supplier",
			rules.Field("addresses",
				rules.Each(
					rules.Field("country",
						rules.Assert("13", "supplier address must include country",
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
						rules.Assert("14", "customer address must include country",
							is.Present,
						),
					),
				),
			),
		),
		// B2B: both supplier and customer must be present, each with a legal
		// identity declaring an allowed ICD 6523 scheme (G2.19). If that scheme
		// is SIREN (0002) or EU VAT (0223), a matching TaxID must also be set
		// (G2.33).
		rules.When(
			is.Func("B2B invoice", invoiceIsB2BAny),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("09", "invoice document type must be one of the Flow 10 permitted UNTDID 1001 codes (380, 389, 393, 501, 386, 500, 384, 471, 472, 473, 381, 261, 396, 502, 503)",
						is.Func("allowed Flow 10 document type", invoiceDocumentTypeAllowed),
					),
					rules.Assert("11", "billing mode extension (fr-ctc-billing-mode) is required (G1.02)",
						is.Func("has billing mode", extensionsHaveBillingMode),
					),
				),
			),
			rules.When(
				is.Func("billing mode is final-after-advance (B4/S4/M4)", invoiceIsFinalAfterAdvance),
				rules.Field("tax",
					rules.Field("ext",
						rules.Assert("12", "final-after-advance billing mode (B4/S4/M4) cannot be combined with an advance-payment document type (386/500/503) (G1.60)",
							is.Func("not advance-payment doc type", invoiceNotAdvancePaymentDocType),
						),
					),
				),
			),
			rules.Field("supplier",
				rules.Assert("01", "supplier is required for Flow 10 B2B invoice (G2.19)",
					is.Present,
				),
				rules.Assert("02", "supplier must declare a legal identity with an allowed ICD 6523 scheme (G2.19): 0002, 0223, 0227, 0228 or 0229",
					is.Func("party has allowed legal scheme", partyHasAllowedLegalScheme),
				),
				rules.Assert("03", "supplier TaxID is required when legal identity scheme is SIREN (0002) or EU VAT (0223) (G2.33)",
					is.Func("party has TaxID when required", partyHasTaxIDWhenRequired),
				),
			),
			rules.When(
				is.Func("invoice has exempt (E) VAT category", invoiceHasExemptCombo),
				rules.Assert("07", "supplier VAT ID or ordering.seller (tax representative) VAT ID is required when the invoice VAT breakdown contains an exempt (E) category",
					is.Func("supplier or tax rep has VAT ID", invoiceHasSellerVATIDForExempt),
				),
				rules.Assert("15", "invoice with an exempt (E) VAT category must include an exemption reason in tax.notes (key=exempt, non-empty text)",
					is.Func("has exempt tax note", invoiceHasExemptTaxNote),
				),
			),
			rules.Field("customer",
				rules.Assert("04", "customer is required for Flow 10 B2B invoice (G2.19)",
					is.Present,
				),
				rules.Assert("05", "customer must declare a legal identity with an allowed ICD 6523 scheme (G2.19): 0002, 0223, 0227, 0228 or 0229",
					is.Func("party has allowed legal scheme", partyHasAllowedLegalScheme),
				),
				rules.Assert("06", "customer TaxID is required when legal identity scheme is SIREN (0002) or EU VAT (0223) (G2.33)",
					is.Func("party has TaxID when required", partyHasTaxIDWhenRequired),
				),
			),
		),
	)
}

func invoiceIsB2BAny(v any) bool {
	inv, ok := v.(*bill.Invoice)
	return ok && !invoiceIsB2C(inv)
}

func invoiceIsB2CAny(v any) bool {
	inv, ok := v.(*bill.Invoice)
	return ok && invoiceIsB2C(inv)
}

// allowedVATRates is the whitelist of VAT percentages authorised on a
// Flow 10 invoice (G1.24). Comparison is numeric — "20", "20.0" and
// "20.00" are all equivalent, handled by num.Percentage.Compare.
var allowedVATRates = mustParsePercentages(
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

func partyHasSIREN(v any) bool {
	party, ok := v.(*org.Party)
	if !ok || party == nil {
		return false
	}
	for _, id := range party.Identities {
		if id == nil || id.Ext.IsZero() {
			continue
		}
		if id.Ext.Get(iso.ExtKeySchemeID).String() == schemeIDSIREN {
			return true
		}
	}
	return false
}

func invoiceVATRatesAllowed(v any) bool {
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
			if !percentageInList(*combo.Percent, allowedVATRates) {
				return false
			}
		}
	}
	return true
}

func percentageInList(p num.Percentage, list []num.Percentage) bool {
	for _, a := range list {
		if p.Compare(a) == 0 {
			return true
		}
	}
	return false
}

func extensionsHaveB2CCategory(v any) bool {
	return extensionsValue(v).Get(ExtKeyB2CCategory) != ""
}

// extensionsValue extracts a tax.Extensions from either a value- or
// pointer-typed argument. The rules engine currently passes fields of
// struct type by pointer, so both forms must be handled.
func extensionsValue(v any) tax.Extensions {
	switch ext := v.(type) {
	case tax.Extensions:
		return ext
	case *tax.Extensions:
		if ext == nil {
			return tax.Extensions{}
		}
		return *ext
	default:
		return tax.Extensions{}
	}
}

func partyHasAllowedLegalScheme(v any) bool {
	party, ok := v.(*org.Party)
	if !ok || party == nil {
		return false
	}
	return slices.Contains(allowedPartySchemeIDs, partyLegalSchemeID(party))
}

func extensionsHaveBillingMode(v any) bool {
	return extensionsValue(v).Get(ExtKeyBillingMode) != ""
}

func invoiceIsFinalAfterAdvance(v any) bool {
	inv, ok := v.(*bill.Invoice)
	if !ok || inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
		return false
	}
	return slices.Contains(finalAfterAdvanceBillingModes, inv.Tax.Ext.Get(ExtKeyBillingMode))
}

func invoiceNotAdvancePaymentDocType(v any) bool {
	return !slices.Contains(advancePaymentDocumentTypes, extensionsValue(v).Get(untdid.ExtKeyDocumentType))
}

// invoiceDocumentTypeAllowed reads the untdid-document-type extension set
// by the Flow 10 scenarios and confirms it is one of the permitted codes.
func invoiceDocumentTypeAllowed(v any) bool {
	return slices.Contains(allowedDocumentTypes, extensionsValue(v).Get(untdid.ExtKeyDocumentType))
}

// invoiceHasSellerVATIDForExempt returns true if either the supplier or
// the ordering.seller (treated as the supplier's tax representative)
// carries a non-empty TaxID code. Per the Flow 10 spec, invoices with an
// exempt VAT breakdown must carry at least one of these two VAT IDs.
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

func partyHasVATCode(p *org.Party) bool {
	return p != nil && p.TaxID != nil && p.TaxID.Code != ""
}

// invoiceHasExemptVATCategory reports whether any line on the invoice
// carries a VAT combo tagged with UNTDID 5305 category code "E" (exempt).
// We inspect line-level combos rather than the aggregated totals because
// the untdid-tax-category extension is carried on the combo itself, and
// the totals breakdown is only populated after Calculate has run.
// invoiceHasExemptCombo reports whether the invoice has any VAT combo
// whose UNTDID 5305 tax-category extension is "E" (exempt). Reading the
// extension rather than the combo Key lets upstream converters or
// manual entries declare exemption directly via the UNTDID code.
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

// invoiceHasExemptTaxNote checks for at least one tax.Note with
// Key=exempt and non-empty Text.
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

func partyHasTaxIDWhenRequired(v any) bool {
	party, ok := v.(*org.Party)
	if !ok || party == nil {
		return true
	}
	scheme := partyLegalSchemeID(party)
	if !slices.Contains(schemeIDsRequiringVAT, scheme) {
		return true
	}
	return party.TaxID != nil && party.TaxID.Code != ""
}

// vatKeyToUNTDIDCategory maps each supported GOBL VAT rate key to its
// UNTDID 5305 category code. The Canary Islands (IGIC / "L") and
// Ceuta/Melilla (IPSI / "M") categories are intentionally absent since
// they are not applicable to Flow 10.
var vatKeyToUNTDIDCategory = map[cbc.Key]cbc.Code{
	tax.KeyStandard:       "S",
	tax.KeyZero:           "Z",
	tax.KeyExempt:         "E",
	tax.KeyReverseCharge:  "AE",
	tax.KeyIntraCommunity: "K",
	tax.KeyExport:         "G",
	tax.KeyOutsideScope:   "O",
}

func invoiceIsB2C(inv *bill.Invoice) bool {
	return inv != nil && inv.Tags.HasTags(TagB2C)
}

func normalizeInvoice(inv *bill.Invoice) {
	if inv == nil {
		return
	}
	normalizeInvoiceTaxCategories(inv)
	if invoiceIsB2C(inv) {
		normalizeB2CCategoryOnInvoice(inv)
		return
	}
	normalizeParty(inv.Supplier)
	normalizeParty(inv.Customer)
	normalizeInvoiceBillingMode(inv)
}

// normalizeB2CCategoryOnInvoice defaults the B2C transaction category to
// TNT1 (not subject to French VAT) when the caller has not supplied one.
// TNT1 is the safest default: it covers B2C sales that would otherwise
// require explicit per-case classification (intra-EU distance sales,
// out-of-scope, etc.), and a user wanting a narrower code must set it
// explicitly.
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
// on each VAT combo based on its rate key. Combos whose key we do not
// map (IGIC / IPSI, or unknown) are left untouched.
func normalizeInvoiceTaxCategories(inv *bill.Invoice) {
	for _, line := range inv.Lines {
		if line == nil {
			continue
		}
		for _, combo := range line.Taxes {
			if combo == nil || combo.Category != tax.CategoryVAT {
				continue
			}
			if code, ok := vatKeyToUNTDIDCategory[combo.Key]; ok {
				combo.Ext = combo.Ext.Set(untdid.ExtKeyTaxCategory, code)
			}
		}
	}
}

// normalizeInvoiceBillingMode picks a sensible default for the Flow 10
// billing-mode extension when the user has not supplied one. We default
// to the Mixed (M) prefix since it is the safest without line-level
// analysis: M2 when the invoice is already paid in full, M1 otherwise.
// The user can override by setting the extension explicitly.
func normalizeInvoiceBillingMode(inv *bill.Invoice) {
	if inv.Tax != nil && !inv.Tax.Ext.IsZero() && inv.Tax.Ext.Get(ExtKeyBillingMode) != "" {
		return
	}
	mode := BillingModeM1
	if inv.Totals != nil && inv.Totals.Paid() {
		mode = BillingModeM2
	}
	if inv.Tax == nil {
		inv.Tax = &bill.Tax{}
	}
	inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyBillingMode, mode)
}
