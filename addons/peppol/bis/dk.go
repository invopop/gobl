package bis

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// allowedDKPaymentMeans lists the UNTDID 4461 payment means codes permitted
// for Danish suppliers under DK-R-005.
var allowedDKPaymentMeans = []cbc.Code{"1", "10", "31", "42", "48", "49", "50", "58", "59", "93", "97"}

func billInvoiceRulesDK() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(supplierCountryIs(l10n.DK),
			// DK-R-016: credit notes must have a non-negative payable amount.
			rules.Assert("DK-R-016", "Danish credit note cannot have a negative total (DK-R-016)",
				is.Func("dk credit note total", dkCreditNoteNonNegative),
			),
		),
	)
}

func orgPartyRulesDK() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(supplierCountryIs(l10n.DK),
			// DK-R-002 / DK-R-014: Danish supplier must provide CVR (scheme 0184) as legal entity.
			rules.Field("supplier",
				rules.Assert("DK-R-002", "Danish supplier must provide a CVR identity with scheme 0184 (DK-R-002, DK-R-014)",
					is.Func("dk supplier cvr", partyHasCVRIdentity),
				),
				// DK-R-013: every supplier identity must carry a scheme ID.
				rules.Field("identities",
					rules.Each(
						rules.Assert("DK-R-013", "Danish supplier identities must specify an ISO 6523 scheme ID (DK-R-013)",
							is.Func("identity has scheme", identityHasSchemeID),
						),
					),
				),
			),
		),
	)
}

// orgItemRulesDK is intentionally a no-op. DK-R-003 requires UNSPSC
// classifications to use version 19.05.01 or 26.08.01, but GOBL has no
// structured slot for a classification-scheme version; gobl.ubl owns
// cac:CommodityClassification/@listVersionID and enforces the allowed set
// there.
func orgItemRulesDK() *rules.Set {
	return rules.For(new(bill.Invoice))
}

func payInstructionsRulesDK() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(supplierCountryIs(l10n.DK),
			rules.Field("payment",
				rules.Field("instructions",
					// DK-R-005: restrict allowed UNTDID payment means codes.
					rules.Assert("DK-R-005", "Danish payment means code must be in the allowed DK subset (DK-R-005)",
						is.Func("dk payment means", dkPaymentMeansAllowed),
					),
					// DK-R-006: payment means 31/42 require CreditTransfer account details.
					rules.Assert("DK-R-006", "Danish payment means 31 or 42 requires bank account and registration (DK-R-006)",
						is.Func("dk 31/42", dkCreditTransferComplete),
					),
					// DK-R-007: payment means 49 requires direct debit mandate and account.
					rules.Assert("DK-R-007", "Danish payment means 49 requires mandate reference and account (DK-R-007)",
						is.Func("dk 49", dkDirectDebit49Complete),
					),
				),
			),
		),
	)
}

// --- helpers ---

func dkCreditNoteNonNegative(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Type != bill.InvoiceTypeCreditNote || inv.Totals == nil {
		return true
	}
	return !inv.Totals.Payable.IsNegative()
}

func partyHasCVRIdentity(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	// Accept TaxID country DK as CVR signal, or any identity with scheme 0184.
	for _, id := range p.Identities {
		if id == nil {
			continue
		}
		if id.Ext.Get(iso.ExtKeySchemeID) == schemeDKCVR {
			return true
		}
	}
	// Also accept a DK tax ID with a code (legal CVR through TaxID).
	if p.TaxID != nil && p.TaxID.Country.Code() == l10n.DK && p.TaxID.Code != "" {
		return true
	}
	return false
}

func identityHasSchemeID(val any) bool {
	id, ok := val.(*org.Identity)
	if !ok || id == nil {
		return true
	}
	return id.Ext.Get(iso.ExtKeySchemeID) != ""
}

func customerIsDK(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return false
	}
	return partyCountry(p) == l10n.DK
}

func dkPaymentMeansAllowed(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return true
	}
	code := instr.Ext.Get(untdid.ExtKeyPaymentMeans)
	if code == "" {
		return true
	}
	return code.In(allowedDKPaymentMeans...)
}

func dkCreditTransferComplete(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return true
	}
	code := instr.Ext.Get(untdid.ExtKeyPaymentMeans)
	if code != "31" && code != "42" {
		return true
	}
	if len(instr.CreditTransfer) == 0 {
		return false
	}
	for _, ct := range instr.CreditTransfer {
		// GOBL's canonical SEPA field is IBAN; Number is the non-IBAN fallback.
		// Either satisfies DK-R-006's "bank account" requirement.
		if ct == nil || (ct.IBAN == "" && ct.Number == "") {
			return false
		}
	}
	return true
}

func dkDirectDebit49Complete(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return true
	}
	code := instr.Ext.Get(untdid.ExtKeyPaymentMeans)
	if code != "49" {
		return true
	}
	if instr.DirectDebit == nil {
		return false
	}
	return instr.DirectDebit.Ref != "" && instr.DirectDebit.Account != ""
}
