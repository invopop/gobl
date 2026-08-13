// Package i18n provides internationalization models.
package i18n

import (
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("i18n"), String{})
	rules.Register(
		"i18n",
		rules.GOBL.Add("I18N"),
		langRules(),
	)
}

func langRules() *rules.Set {
	return rules.For(Lang(""),
		rules.Assert("01", "must be a valid language code", is.In(validLangValues()...)),
	)
}
