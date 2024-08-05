package hu

import (
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegime(New())
}

func New() *tax.Regime {
	return &tax.Regime{
		Country:  "HU",
		Currency: currency.HUF,
		Name: i18n.String{
			i18n.EN: "Hungary",
			i18n.HU: "Magyarország",
		},
		TimeZone:   "Europe/Budapest",
		Extensions: extensionKeys,
		Categories: taxCategories,
		Tags:       invoiceTags,
	}
}
