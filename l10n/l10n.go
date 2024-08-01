// Package l10n provides localization models and data.
package l10n

import (
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("l10n"),
		Code(""),
		ISOCountryCode(""),
		TaxCountryCode(""),
	)
}
