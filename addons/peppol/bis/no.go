package bis

import (
	"regexp"
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// noVATRe matches the Norwegian VAT identifier format per NO-R-001:
// "NO" + 9 digits + "MVA".
var noVATRe = regexp.MustCompile(`^NO\d{9}MVA$`)

func orgPartyRulesNO() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(supplierCountryIs(l10n.NO),
			rules.Field("supplier",
				rules.Assert("NO-01", "Norwegian VAT must be NO+9 digits+MVA (NO-R-001)",
					is.Func("no vat format", norwegianVATFormat),
				),
			),
		),
	)
}

func norwegianVATFormat(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil || p.TaxID == nil {
		return true
	}
	if p.TaxID.Country.Code() != l10n.NO {
		return true
	}
	// The TaxID code stores the bare number; to validate the full Peppol-visible
	// form, we reconstruct the expected string.
	code := p.TaxID.Code.String()
	if code == "" {
		return true
	}
	// Accept either bare 9 digits (as typically stored in GOBL) or the full NOxxxMVA form.
	if onlyDigits(code) && len(code) == 9 {
		return true
	}
	return noVATRe.MatchString(strings.ToUpper(code))
}
