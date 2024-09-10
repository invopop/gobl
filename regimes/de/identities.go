package de

import (
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

const (
	// IdentityKeyTaxNumber represents the German tax number (Steuernummer) issued to
	// people that can be included on invoices inside Germany. For international
	// sales, the registered VAT number (Umsatzsteueridentifikationsnummer) should
	// be used instead.
	IdentityKeyTaxNumber cbc.Key = "de-tax-number"
)

var identityKeyDefinitions = []*cbc.KeyDefinition{
	{
		Key: IdentityKeyTaxNumber,
		Name: i18n.String{
			i18n.EN: "Tax Number",
			i18n.DE: "Steuernummer",
		},
	},
}

func normalizeIdentity(id *org.Identity) error {
	if id == nil || id.Key != IdentityKeyTaxNumber {
		return nil
	}
	code := strings.ToUpper(id.Code.String())
	code = tax.IdentityCodeBadCharsRegexp.ReplaceAllString(code, "")
	id.Code = cbc.Code(code)
	return nil
}
