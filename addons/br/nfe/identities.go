package nfe

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
)

// Identity keys
const (
	IdentityKeyTaxable = "taxable"
)

var identities = []*cbc.Definition{
	{
		Key: org.IdentityKeyGTIN,
		Name: i18n.String{
			i18n.EN: "Product's GTIN (Global Trade Item Number)",
			i18n.PT: "GTIN (Global Trade Item Number) do produto",
		},
	},
	{
		Key: org.IdentityKeyGTIN.With(IdentityKeyTaxable),
		Name: i18n.String{
			i18n.EN: "Taxable item's GTIN (Global Trade Item Number)",
			i18n.PT: "GTIN (Global Trade Item Number) da unidade tributável",
		},
	},
	{
		Key: org.IdentityKeyNCM,
		Name: i18n.String{
			i18n.EN: "Product's NCM (Mercosur Common Nomenclature) code",
			i18n.PT: "Código NCM (Nomenclatura Comum do Mercosul) do produto",
		},
	},
}
