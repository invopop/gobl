package bis

import (
	"regexp"
	"strconv"
	"strings"

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
			// GR-R-001-1 enforces the 6-segment shape. The segment contents
			// (GR-R-001-2..-7: supplier TIN, YYYYMMDD, sequence, doc type,
			// free-form) are not checked here — gobl.ubl should build the
			// Peppol-visible ID from the structured Supplier.TaxID + IssueDate
			// + sequence fields rather than parsing them back out.
			rules.Assert("GR-01", "Greek invoice ID must have 6 underscore-delimited segments when joining series and code (GR-R-001-1)",
				is.Func("gr id segments", grIDSixSegments),
			),
			rules.Assert("GR-02", "Greek invoice must have exactly one MARK identity (GR-R-004-1)",
				is.Func("gr mark count", grMARKExactlyOne),
			),
			rules.Assert("GR-03", "Greek invoice MARK must be a positive integer (GR-R-004-2)",
				is.Func("gr mark positive", grMARKPositive),
			),
			rules.Assert("GR-04", "at most one Greek invoice URL attachment allowed (GR-R-008-2)",
				is.Func("gr url count", grInvoiceURLCardinality),
			),
		),
	)
}

// grFullInvoiceID joins Series and Code with an underscore so the count check
// works whether the caller stores all six segments in Code or splits them
// across Series + Code.
func grFullInvoiceID(inv *bill.Invoice) string {
	if inv == nil {
		return ""
	}
	if inv.Series == "" {
		return inv.Code.String()
	}
	return inv.Series.String() + "_" + inv.Code.String()
}

// grIDSixSegments enforces GR-R-001-1: the joined identifier must split into
// exactly 6 non-empty-or-empty segments on `_`. We deliberately do not check
// the contents of each segment here — that is gobl.ubl's job at UBL emit.
func grIDSixSegments(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	id := grFullInvoiceID(inv)
	if id == "" {
		return true
	}
	return len(strings.Split(id, "_")) == 6
}

func orgPartyRulesGR() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(supplierCountryIs(l10n.GR),
			rules.Field("supplier",
				rules.Field("name",
					rules.Assert("GR-05", "Greek supplier name is required (GR-R-002)", is.Present),
				),
				rules.Assert("GR-06", "Greek supplier VAT must start with EL and be a valid TIN (GR-R-003)",
					is.Func("gr vat format", grSupplierVATValid),
				),
				rules.Assert("GR-07", "Greek supplier inbox must use scheme 9933 with TIN as code (GR-R-009)",
					is.Func("gr supplier inbox", grSupplierInboxValid),
				),
			),
			rules.Field("customer",
				rules.Field("name",
					rules.Assert("GR-08", "Greek customer name is required (GR-R-005)", is.Present),
				),
				rules.Assert("GR-09", "Greek customer must have a VAT identifier when customer is Greek (GR-R-006)",
					is.Func("gr customer vat", grCustomerVATWhenGreek),
				),
				rules.Assert("GR-10", "Greek customer inbox must use scheme 9933 with TIN as code when customer is Greek (GR-R-010)",
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
