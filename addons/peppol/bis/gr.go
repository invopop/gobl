package bis

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// grDocumentTypes lists the Greek Peppol document type codes allowed in the
// fourth segment of the invoice ID under GR-R-001-5.
var grDocumentTypes = []string{
	"1.1", "1.2", "1.3", "1.4", "1.5", "1.6", "2.1", "2.2", "2.3", "2.4",
	"5.1", "5.2", "6.1", "6.2", "7.1",
	"11.1", "11.2", "11.3", "11.4", "11.5", "13.1", "13.2", "13.3", "13.4", "13.31",
}

var (
	grTINRe = regexp.MustCompile(`^\d{9}$`)
	// Greek invoice ID: 6 _-delimited segments. Splitter used at parse time; here
	// we match the overall shape.
	grInvoiceIDRe = regexp.MustCompile(`^[^_]+_[^_]+_[^_]+_[^_]+_[^_]+_[^_]+$`)
	grEndpointRe  = regexp.MustCompile(`^\d{9}$`)
)

func billInvoiceRulesGR() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(supplierCountryIs(l10n.GR),
			// GR-R-001: Greek invoice ID structure.
			rules.Assert("GR-R-001-1", "Greek invoice ID must have 6 segments (GR-R-001-1)",
				is.Func("gr id segments", grIDSixSegments),
			),
			rules.Assert("GR-R-001-2", "Greek invoice ID segment 1 must be the supplier TIN (GR-R-001-2)",
				is.Func("gr id tin", grIDFirstSegmentTIN),
			),
			rules.Assert("GR-R-001-3", "Greek invoice ID segment 2 must match issue date YYYYMMDD (GR-R-001-3)",
				is.Func("gr id date", grIDSecondSegmentDate),
			),
			rules.Assert("GR-R-001-4", "Greek invoice ID segment 3 must be a positive integer (GR-R-001-4)",
				is.Func("gr id seq", grIDThirdSegmentPositive),
			),
			rules.Assert("GR-R-001-5", "Greek invoice ID segment 4 must be a valid Greek document type (GR-R-001-5)",
				is.Func("gr id type", grIDFourthSegmentType),
			),
			rules.Assert("GR-R-001-6", "Greek invoice ID segment 5 must not be empty (GR-R-001-6)",
				is.Func("gr id seg5", grIDFifthSegmentPresent),
			),
			rules.Assert("GR-R-001-7", "Greek invoice ID segment 6 must not be empty (GR-R-001-7)",
				is.Func("gr id seg6", grIDSixthSegmentPresent),
			),
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
			// GR-R-002: supplier name required (covered at schema level, but double-check).
			rules.Field("supplier",
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
			// GR-R-005: customer name required.
			rules.Field("customer",
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

// grFullInvoiceID assembles Series + Code into a single identifier for rule checks.
func grFullInvoiceID(inv *bill.Invoice) string {
	if inv == nil {
		return ""
	}
	if inv.Series != "" {
		return inv.Series.String() + inv.Code.String()
	}
	return inv.Code.String()
}

func grIDSixSegments(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	id := grFullInvoiceID(inv)
	if id == "" {
		return true
	}
	return grInvoiceIDRe.MatchString(id)
}

func grIDSegments(inv *bill.Invoice) []string {
	return strings.Split(grFullInvoiceID(inv), "_")
}

func grIDFirstSegmentTIN(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	segs := grIDSegments(inv)
	if len(segs) < 1 {
		return true
	}
	if !grTINRe.MatchString(segs[0]) {
		return false
	}
	if inv.Supplier != nil && inv.Supplier.TaxID != nil {
		return inv.Supplier.TaxID.Code.String() == segs[0]
	}
	return true
}

func grIDSecondSegmentDate(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	segs := grIDSegments(inv)
	if len(segs) < 2 {
		return true
	}
	if len(segs[1]) != 8 || !onlyDigits(segs[1]) {
		return false
	}
	if inv.IssueDate.IsZero() {
		return true
	}
	expected := strings.ReplaceAll(inv.IssueDate.String(), "-", "")
	return segs[1] == expected
}

func grIDThirdSegmentPositive(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	segs := grIDSegments(inv)
	if len(segs) < 3 {
		return true
	}
	n, err := strconv.Atoi(segs[2])
	return err == nil && n > 0
}

func grIDFourthSegmentType(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	segs := grIDSegments(inv)
	if len(segs) < 4 {
		return true
	}
	for _, t := range grDocumentTypes {
		if segs[3] == t {
			return true
		}
	}
	return false
}

func grIDFifthSegmentPresent(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	segs := grIDSegments(inv)
	return len(segs) < 5 || segs[4] != ""
}

func grIDSixthSegmentPresent(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	segs := grIDSegments(inv)
	return len(segs) < 6 || segs[5] != ""
}

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
	code := p.TaxID.Code.String()
	// Either the code is prefixed with EL or stored separately; the TaxID
	// model holds the bare code and the country prefix indicates EL for Greece.
	return grTINRe.MatchString(code)
}

func grSupplierInboxValid(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	if len(p.Inboxes) == 0 {
		return true // R010 handles presence
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

// avoid unused import warnings if strconv goes unreferenced during edits
var _ = cbc.CodeEmpty
