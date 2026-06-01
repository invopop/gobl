package nz

import (
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Reference: https://www.nzbn.govt.nz/whats-an-nzbn/about/
// Reference: https://www.gs1nz.org/services/glns/

const (
	// IdentityTypeNZBN represents the New Zealand Business Number (NZBN), a 13-digit
	// GS1 Global Location Number used as the primary business identifier and
	// Peppol participant ID.
	IdentityTypeNZBN cbc.Code = "NZBN"
)

const nzbnPattern = `^\d{13}$`

var identityDefs = []*cbc.Definition{
	{
		Code: IdentityTypeNZBN,
		Name: i18n.String{
			i18n.EN: "New Zealand Business Number (NZBN)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The New Zealand Business Number (NZBN) is a unique 13-digit identifier
				issued to all businesses registered in New Zealand. It is based on the
				GS1 Global Location Number (GLN) standard and is used as the Peppol
				participant identifier for e-invoicing. Mandatory for businesses with
				annual turnover of NZD 60,000 or more.
			`),
		},
	},
}

func normalizeOrgIdentity(id *org.Identity) {
	if id == nil || id.Type != IdentityTypeNZBN {
		return
	}
	code := strings.ToUpper(id.Code.String())
	code = tax.IdentityCodeBadCharsRegexp.ReplaceAllString(code, "")
	id.Code = cbc.Code(strings.TrimPrefix(code, string(l10n.NZ)))
}

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				org.IdentityTypeIn(IdentityTypeNZBN),
				rules.Field("code",
					rules.Assert("01", "NZBN must be 13 digits",
						is.Matches(nzbnPattern),
					),
					rules.Assert("02", "NZBN check digit mismatch",
						is.Func("checksum", nzbnChecksumValid),
					),
				),
			),
		),
	)
}

// nzbnChecksumValid verifies the GS1 Modulo-10 check digit of a 13-digit NZBN.
func nzbnChecksumValid(value any) bool {
	code, ok := value.(cbc.Code)
	// Guard: wrong length is caught by the format rule; don't emit a misleading checksum error.
	if !ok || len(code) != 13 {
		return true
	}
	val := code.String()

	sum := 0
	for i := 0; i < 12; i++ {
		d := int(val[i] - '0')
		if i%2 == 0 {
			sum += d * 1
		} else {
			sum += d * 3
		}
	}
	check := (10 - sum%10) % 10
	return check == int(val[12]-'0')
}
