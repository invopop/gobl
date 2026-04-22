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

// nlAllowedSchemes are the ISO 6523 scheme IDs accepted for Dutch legal
// entities (NL-R-003, NL-R-005): KVK ("0106") and OIN ("0190").
var nlAllowedSchemes = []cbc.Code{"0106", "0190"}

// nlAllowedPaymentMeans are the UNTDID 4461 codes permitted for Dutch
// suppliers by NL-R-008.
var nlAllowedPaymentMeans = []cbc.Code{"30", "48", "49", "57", "58", "59"}

func billInvoiceRulesNL() *rules.Set {
	return rules.For(new(bill.Invoice),
		// NL-R-001 applies only when both parties are Dutch.
		rules.When(bothCountriesAre(l10n.NL),
			rules.Assert("NL-01", "Dutch credit note must reference a preceding invoice (NL-R-001)",
				is.Func("nl credit note preceding", nlCreditNoteHasPreceding),
			),
		),
		// NL-R-007 / NL-R-009: supplier-scoped.
		rules.When(supplierCountryIs(l10n.NL),
			rules.Assert("NL-02", "Dutch supplier must provide payment instructions (NL-R-007)",
				is.Func("has payment instructions", hasPaymentInstructions),
			),
			rules.Assert("NL-03", "Dutch invoice with line order references must have an ordering code (NL-R-009)",
				is.Func("nl order ref requires ordering code", nlLineOrderRefRequiresOrderingCode),
			),
		),
	)
}

// orgPartyRulesNL covers the party-shape checks NL-R-002..R-005. The
// schematron for all four requires both supplier and customer to be Dutch;
// the field-scoped assertions enforce the supplier / customer shape from
// there.
func orgPartyRulesNL() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(bothCountriesAre(l10n.NL),
			rules.Field("supplier",
				rules.Field("addresses",
					rules.Assert("NL-04", "Dutch supplier address must have street, city and postcode (NL-R-002)",
						is.Func("nl supplier addr", firstAddressStreetLocalityCode),
					),
				),
				rules.Assert("NL-05", "Dutch supplier legal entity must use scheme 0106 (KVK) or 0190 (OIN) (NL-R-003)",
					is.Func("nl supplier legal scheme", nlPartyLegalScheme),
				),
			),
			rules.Field("customer",
				rules.Field("addresses",
					rules.AssertIfPresent("NL-06", "Dutch customer address must have street, city and postcode (NL-R-004)",
						is.Func("nl customer addr", firstAddressStreetLocalityCode),
					),
				),
				rules.Assert("NL-07", "Dutch customer legal entity must use scheme 0106 (KVK) or 0190 (OIN) (NL-R-005)",
					is.Func("nl customer legal scheme", nlPartyLegalScheme),
				),
			),
		),
	)
}

// payInstructionsRulesNL carries NL-R-008 (customer-scoped per schematron).
func payInstructionsRulesNL() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(customerCountryIs(l10n.NL),
			rules.Field("payment",
				rules.Field("instructions",
					rules.Assert("NL-08", "Dutch payment means code must be in the allowed subset (NL-R-008)",
						is.Func("nl payment means", nlPaymentMeansAllowed),
					),
				),
			),
		),
	)
}

// --- helpers ---

func nlCreditNoteHasPreceding(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Type != bill.InvoiceTypeCreditNote {
		return true
	}
	for _, pre := range inv.Preceding {
		if pre != nil && pre.Code != "" {
			return true
		}
	}
	return false
}

func nlLineOrderRefRequiresOrderingCode(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	hasLineRef := false
	for _, line := range inv.Lines {
		if line != nil && line.Order != "" {
			hasLineRef = true
			break
		}
	}
	if !hasLineRef {
		return true
	}
	return inv.Ordering != nil && inv.Ordering.Code != ""
}

func firstAddressStreetLocalityCode(val any) bool {
	addrs, ok := val.([]*org.Address)
	if !ok || len(addrs) == 0 {
		return true
	}
	a := addrs[0]
	return a != nil && a.Street != "" && a.Locality != "" && a.Code != ""
}

func nlPartyLegalScheme(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	if partyCountry(p) != l10n.NL {
		return true
	}
	for _, id := range p.Identities {
		if id == nil || id.Scope != org.IdentityScopeLegal {
			continue
		}
		scheme := id.Ext.Get(iso.ExtKeySchemeID)
		if scheme.In(nlAllowedSchemes...) {
			return true
		}
	}
	return false
}

func nlPaymentMeansAllowed(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return true
	}
	code := instr.Ext.Get(untdid.ExtKeyPaymentMeans)
	if code == "" {
		return true
	}
	return code.In(nlAllowedPaymentMeans...)
}
