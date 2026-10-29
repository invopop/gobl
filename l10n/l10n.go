// Package l10n provides localization models and data.
package l10n

import (
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("l10n"),
		Code(""),
		ISOCountryCode(""),
		TaxCountryCode(""),
	)
	rules.Register(
		"l10n",
		rules.GOBL.Add("L10N"),
		isoCountryCodeRules(),
		taxCountryCodeRules(),
	)
}
