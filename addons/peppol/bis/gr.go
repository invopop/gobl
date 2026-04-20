package bis

import (
	"regexp"
	"strconv"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

var (
	grTINRe      = regexp.MustCompile(`^\d{9}$`)
	grEndpointRe = regexp.MustCompile(`^\d{9}$`)
)

func billInvoiceRulesGR() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(supplierCountryIs(l10n.GR),
			// GR-R-001 (invoice ID segmentation) deferred to gobl.ubl — see deferred.go.
			//
			// GR-R-004: exactly one MARK identity with positive integer code.
			rules.Assert("GR-R-004-1", "Greek invoice must have exactly one MARK identity (GR-R-004-1)",
				is.Func("gr mark count", grMARKExactlyOne),
			),
			rules.Assert("GR-R-004-2", "Greek invoice MARK must be a positive integer (GR-R-004-2)",
				is.Func("gr mark positive", grMARKPositive),
			),
			// GR-R-008-2: at most one invoice URL attachment.
			rules.Assert("GR-R-008-2", "at most one Greek invoice URL attachment allowed (GR-R-008-2)",
				is.Func("gr url count", grInvoiceURLCardinality),
			),
		),
	)
}

func orgPartyRulesGR() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(supplierCountryIs(l10n.GR),
			rules.Field("supplier",
				// GR-R-002: supplier name required.
				rules.Field("name",
					rules.Assert("GR-R-002", "Greek supplier name is required (GR-R-002)", is.Present),
				),
				// GR-R-003: supplier VAT starts with EL and valid TIN.
				rules.Assert("GR-R-003", "Greek supplier VAT must start with EL and be a valid TIN (GR-R-003)",
					is.Func("gr vat format", grSupplierVATValid),
				),
				// GR-R-009: supplier inbox scheme 9933 with TIN as code.
				rules.Assert("GR-R-009", "Greek supplier inbox must use scheme 9933 with TIN as code (GR-R-009)",
					is.Func("gr supplier inbox", grSupplierInboxValid),
				),
			),
			rules.Field("customer",
				// GR-R-005: customer name required.
				rules.Field("name",
					rules.Assert("GR-R-005", "Greek customer name is required (GR-R-005)", is.Present),
				),
				// GR-R-006/R-010: Greek customer must have VAT and correct inbox.
				rules.Assert("GR-R-006", "Greek customer must have a VAT identifier when customer is Greek (GR-R-006)",
					is.Func("gr customer vat", grCustomerVATWhenGreek),
				),
				rules.Assert("GR-R-010", "Greek customer inbox must use scheme 9933 with TIN as code when customer is Greek (GR-R-010)",
					is.Func("gr customer inbox", grCustomerInboxWhenGreek),
				),
			),
		),
	)
}

// --- helpers ---

func grMARKIdentities(inv *bill.Invoice) []*org.Identity {
	if inv == nil || inv.Ordering == nil {
		return nil
	}
	var out []*org.Identity
	for _, id := range inv.Ordering.Identities {
		if id != nil && id.Key == IdentityKeyGreekMARK {
			out = append(out, id)
		}
	}
	return out
}

func grMARKExactlyOne(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	return len(grMARKIdentities(inv)) == 1
}

func grMARKPositive(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	marks := grMARKIdentities(inv)
	if len(marks) == 0 {
		return true
	}
	for _, m := range marks {
		n, err := strconv.Atoi(m.Code.String())
		if err != nil || n <= 0 {
			return false
		}
	}
	return true
}

func grInvoiceURLCardinality(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	count := 0
	for _, pre := range inv.Preceding {
		if pre != nil && pre.URL != "" {
			count++
		}
	}
	return count <= 1
}

func grSupplierVATValid(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil || p.TaxID == nil {
		return true
	}
	if p.TaxID.Country.Code() != l10n.GR {
		return true
	}
	return grTINRe.MatchString(p.TaxID.Code.String())
}

func grSupplierInboxValid(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	if len(p.Inboxes) == 0 {
		return true // presence is enforced by R010/R020 at base layer
	}
	tin := ""
	if p.TaxID != nil {
		tin = p.TaxID.Code.String()
	}
	for _, ib := range p.Inboxes {
		if ib == nil {
			continue
		}
		if ib.Scheme != "9933" {
			return false
		}
		if tin != "" && ib.Code.String() != tin {
			return false
		}
		if !grEndpointRe.MatchString(ib.Code.String()) {
			return false
		}
	}
	return true
}

func grCustomerVATWhenGreek(val any) bool {
	c, ok := val.(*org.Party)
	if !ok || c == nil {
		return true
	}
	if partyCountry(c) != l10n.GR {
		return true
	}
	return c.TaxID != nil && c.TaxID.Code != ""
}

func grCustomerInboxWhenGreek(val any) bool {
	c, ok := val.(*org.Party)
	if !ok || c == nil {
		return true
	}
	if partyCountry(c) != l10n.GR {
		return true
	}
	return grSupplierInboxValid(c)
}
