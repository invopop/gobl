// Package untdid defines the UN/EDIFACT data elements contained in the UNTDID (United Nations Trade Data Interchange Directory).
package untdid

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterCatalogueDef(newCatalogue())
}

func newCatalogue() *tax.CatalogueDef {
	return &tax.CatalogueDef{
		Key:  "untdid",
		Name: i18n.NewString("UN/EDIFACT Data Elements"),
		Extensions: []*cbc.KeyDefinition{
			extDocumentTypes, // 1001
			extPaymentMeans,  // 4461
			extAllowance,     // 5189
			extTaxCategory,   // 5305
			extCharge,        // 7161
		},
	}
}
