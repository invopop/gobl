package bis

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// seAllowedVATPercents lists the Swedish VAT rates allowed by SE-R-006. Compared
// numerically to the invoice percent so callers can store any equivalent
// representation (25%, 25.0%, 0.25, etc.) without tripping the rule.
var seAllowedVATPercents = []num.Percentage{
	num.MakePercentage(6, 2),
	num.MakePercentage(12, 2),
	num.MakePercentage(25, 2),
}

func billInvoiceRulesSE() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(supplierCountryIs(l10n.SE),
			rules.Assert("SE-01", "Swedish VAT rate must be 6, 12 or 25 (SE-R-006)",
				is.Func("se vat rate", seVATRateAllowed),
			),
			// SE-R-005 (F-skatt boilerplate in cac:PartyTaxScheme/cbc:CompanyID)
			// is not enforced here. The structured marker is IdentityKeyFSkatt
			// (see identities.go); the addon normalizer populates Scope=tax,
			// a non-VAT Type, and the boilerplate Code on that identity, which
			// drives gobl.ubl's existing tax-scope identity → cac:PartyTaxScheme
			// path to emit the block the schematron looks for.
		),
	)
}

func orgPartyRulesSE() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(supplierCountryIs(l10n.SE),
			rules.Field("supplier",
				rules.Assert("SE-02", "Swedish VAT must be 14 characters (SE-R-001)",
					is.Func("se vat length", seVATLength),
				),
				rules.Assert("SE-03", "Swedish VAT trailing 12 characters must be numeric (SE-R-002)",
					is.Func("se vat trailing digits", seVATTrailingDigits),
				),
				// SE-R-003 (numeric org number) is warning-level and intentionally
				// not enforced here.
				rules.Assert("SE-04", "Swedish organization number must be 10 characters (SE-R-004)",
					is.Func("se org length", seOrgLength),
				),
				rules.Assert("SE-05", "Swedish organization number last digit must be a valid Luhn checksum (SE-R-013)",
					is.Func("se org luhn", seOrgLuhn),
				),
			),
		),
	)
}

// payInstructionsRulesSE currently has no fatal-level rules. SE-R-012 (prefer
// code 30 for domestic credit transfers) and SE-R-007..R-011 (Plusgiro /
// Bankgiro format) are all warning-level and intentionally not enforced.
func payInstructionsRulesSE() *rules.Set {
	return rules.For(new(bill.Invoice))
}

// --- helpers ---

func seVATRateAllowed(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Totals == nil || inv.Totals.Taxes == nil {
		return true
	}
	for _, cat := range inv.Totals.Taxes.Categories {
		if cat == nil {
			continue
		}
		for _, rt := range cat.Rates {
			if rt == nil || rt.Percent == nil {
				continue
			}
			allowed := false
			for _, a := range seAllowedVATPercents {
				if rt.Percent.Equals(a) {
					allowed = true
					break
				}
			}
			if !allowed {
				return false
			}
		}
	}
	return true
}

func seVATLength(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil || p.TaxID == nil {
		return true
	}
	if p.TaxID.Country.Code() != l10n.SE {
		return true
	}
	code := p.TaxID.Code.String()
	// GOBL typically stores the bare numeric; the 14-character full form includes SE prefix + "01".
	// Accept either the bare 10-digit number or the full 14-char version.
	if code == "" {
		return true
	}
	if len(code) == 10 || len(code) == 14 {
		return true
	}
	return false
}

func seVATTrailingDigits(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil || p.TaxID == nil {
		return true
	}
	if p.TaxID.Country.Code() != l10n.SE {
		return true
	}
	code := p.TaxID.Code.String()
	// For the 14-char form, trailing 12 must be digits.
	if len(code) == 14 {
		return onlyDigits(code[2:])
	}
	// For bare 10-digit form, all must be digits (implies trailing digits are fine).
	if len(code) == 10 {
		return onlyDigits(code)
	}
	return true
}

func seOrgLength(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	for _, id := range p.Identities {
		if id == nil || id.Scope != org.IdentityScopeLegal {
			continue
		}
		if len(id.Code.String()) != 10 {
			return false
		}
	}
	return true
}

func seOrgLuhn(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	for _, id := range p.Identities {
		if id == nil || id.Scope != org.IdentityScopeLegal {
			continue
		}
		code := id.Code.String()
		if !onlyDigits(code) || len(code) != 10 {
			continue // handled by SE-R-003/R-004
		}
		if !luhnValid(code) {
			return false
		}
	}
	return true
}
