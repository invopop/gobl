package flow10

import (
	"fmt"
	"slices"
	"strings"

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

// advancePaymentDocumentTypes are the UNTDID 1001 codes representing
// advance-payment invoices (forbidden combined with B4/S4/M4 modes).
var advancePaymentDocumentTypes = []cbc.Code{
	"386", // Advance payment invoice
	"500", // Self-billed advance payment
	"503", // Self-billed credit for claim
}

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

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Assert("01", "invoice must be in EUR or provide an exchange rate to EUR",
			currency.CanConvertTo(currency.EUR),
		),
		rules.Field("supplier",
			rules.Field("addresses",
				rules.Each(
					rules.Field("country",
						rules.Assert("02", "invoice supplier address country is required",
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
						rules.Assert("03", "invoice customer address country is required",
							is.Present,
						),
					),
				),
			),
		),
		// B2C invoices — no Customer party.
		rules.When(
			invoiceIsB2C(),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("04", "invoice tax ext fr-ctc-flow10-b2c-category is required on B2C invoices (G1.68)",
						tax.ExtensionsRequire(ExtKeyB2CCategory),
					),
				),
			),
			rules.Field("supplier",
				rules.Assert("05", "invoice supplier is required on B2C invoices",
					is.Present,
				),
				rules.Assert("06", "invoice supplier must have a SIREN identity (ISO/IEC 6523 scheme 0002) on a B2C invoice",
					is.Func("party has SIREN", partyHasSIREN),
				),
			),
			rules.Assert("07", "invoice VAT line percent must be one of the Flow 10 permitted values 0%, 0.9%, 1.05%, 1.75%, 2.1%, 5.5%, 7%, 8.5%, 9.2%, 9.6%, 10%, 13%, 19.6%, 20%, 20.6% (G1.24)",
				is.Func("allowed Flow 10 VAT percents", invoiceVATPercentsAllowed),
			),
		),
		// Cross-border B2B invoices — Customer present.
		rules.When(
			invoiceIsB2B(),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("08", "invoice tax ext untdid-document-type must be one of the Flow 10 permitted UNTDID 1001 codes",
						tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, allowedDocumentTypes...),
					),
					rules.Assert("09", "invoice tax ext dgfip-billing-mode is required (G1.02)",
						tax.ExtensionsRequire(dgfip.ExtKeyBillingMode),
					),
				),
			),
			rules.When(
				invoiceTaxExtIn(dgfip.ExtKeyBillingMode, dgfip.BillingModeB4, dgfip.BillingModeS4, dgfip.BillingModeM4),
				rules.Field("tax",
					rules.Field("ext",
						rules.Assert("10", "invoice tax ext untdid-document-type must not be an advance-payment code (386, 500, 503) when billing mode is final-after-advance (B4, S4, M4) (G1.60)",
							tax.ExtensionsExcludeCodes(untdid.ExtKeyDocumentType, advancePaymentDocumentTypes...),
						),
					),
				),
			),
			rules.Field("supplier",
				rules.Assert("11", "invoice supplier is required for Flow 10 B2B invoices (G2.19)",
					is.Present,
				),
				rules.Assert("12", "invoice supplier must declare a legal identity with an allowed ICD 6523 scheme (0002, 0223, 0227, 0228 or 0229) (G2.19)",
					is.Func("party has allowed legal scheme", partyHasAllowedLegalScheme),
				),
				rules.Assert("13", "invoice supplier tax_id is required when legal identity scheme is SIREN (0002) or EU VAT (0223) (G2.33)",
					is.Func("party has TaxID when required", partyHasTaxIDWhenRequired),
				),
			),
			rules.When(
				is.Func("invoice has exempt (E) VAT category", invoiceHasExemptCombo),
				rules.Assert("14", "invoice supplier tax_id or ordering.seller tax_id is required when the VAT breakdown contains an exempt (E) category",
					is.Func("supplier or tax rep has VAT ID", invoiceHasSellerVATIDForExempt),
				),
				rules.Assert("15", "invoice tax.notes must include an exempt-reason entry (key=exempt with non-empty text) when the VAT breakdown contains an exempt (E) category",
					is.Func("has exempt tax note", invoiceHasExemptTaxNote),
				),
			),
			rules.Field("customer",
				rules.Assert("16", "invoice customer is required for Flow 10 B2B invoices (G2.19)",
					is.Present,
				),
				rules.Assert("17", "invoice customer must declare a legal identity with an allowed ICD 6523 scheme (0002, 0223, 0227, 0228 or 0229) (G2.19)",
					is.Func("party has allowed legal scheme", partyHasAllowedLegalScheme),
				),
				rules.Assert("18", "invoice customer tax_id is required when legal identity scheme is SIREN (0002) or EU VAT (0223) (G2.33)",
					is.Func("party has TaxID when required", partyHasTaxIDWhenRequired),
				),
			),
		),
	)
}

// -- Rule-level guards ----------------------------------------------------

// invoiceIsB2C returns a Test that passes when the invoice has no
// customer party (B2C transaction).
func invoiceIsB2C() rules.Test {
	return is.Func("invoice is B2C (no customer)", func(v any) bool {
		inv, ok := v.(*bill.Invoice)
		return ok && invoiceIsB2CDoc(inv)
	})
}

// invoiceIsB2B returns a Test that passes when the invoice has a
// customer party. Within Flow 10's scope these are cross-border B2B
// invoices (domestic B2B clearance is Flow 2's territory).
func invoiceIsB2B() rules.Test {
	return is.Func("invoice is B2B (has customer)", func(v any) bool {
		inv, ok := v.(*bill.Invoice)
		return ok && !invoiceIsB2CDoc(inv)
	})
}

// invoiceTaxExtIn returns a Test that passes when bill.Invoice.Tax.Ext[key]
// matches one of the provided codes. Used to gate per-document-type or
// per-billing-mode branches.
func invoiceTaxExtIn(key cbc.Key, codes ...cbc.Code) rules.Test {
	return is.Func(
		fmt.Sprintf("invoice tax ext %s in [%s]", key, joinCodes(codes)),
		func(v any) bool {
			inv, ok := v.(*bill.Invoice)
			if !ok || inv == nil || inv.Tax == nil {
				return false
			}
			return slices.Contains(codes, inv.Tax.Ext.Get(key))
		},
	)
}

// -- Imperative predicates ------------------------------------------------

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

// joinCodes formats a code slice for use in guard-name strings.
func joinCodes(codes []cbc.Code) string {
	ss := make([]string, len(codes))
	for i, c := range codes {
		ss[i] = string(c)
	}
	return strings.Join(ss, ", ")
}

// quiet linter
var _ = org.Party{}
