// Package verifactu provides the Verifactu addon
package verifactu

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// Key identifies the Verifactu addon family. Individual versions append a
	// suffix; the family key is used as the fault-code namespace so that
	// rules that carry across versions keep stable codes.
	Key cbc.Key = "es-verifactu"

	// V1 for Verifactu versions 1.x
	V1 cbc.Key = Key + "-v1"
)

// Official stamps or codes validated by government agencies
const (
	// StampQR contains the URL included in the QR code.
	StampQR cbc.Key = "verifactu-qr"
)

func init() {
	tax.RegisterAddonDef(newAddon())
	rules.RegisterWithGuard(
		Key.String(),
		rules.GOBL.Add("ES-VERIFACTU"),
		is.InContext(tax.AddonIn(V1)),
		billInvoiceRules(),
		taxComboRules(),
	)
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "Spain VERI*FACTU V1",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("VERI*FACTU error response code list"),
				URL:   "https://prewww2.aeat.es/static_files/common/internet/dep/aplicaciones/es/aeat/tikeV1.0/cont/ws/errores.properties",
			},
		},
		Tags: []*tax.TagSet{
			{
				Schema: bill.ShortSchemaInvoice,
				List: []*cbc.Definition{
					{
						Key: tax.TagReplacement,
						Name: i18n.String{
							i18n.EN: "Replacement Invoice",
							i18n.ES: "Factura de Sustitución",
						},
						Desc: i18n.NewString(here.Doc(`
							Used under special circumstances to indicate that this invoice replaces a previously
							issued simplified invoice. The previous document was correct, but the replacement is
							necessary to provide tax details of the customer.
						`)),
					},
				},
			},
		},
		Extensions:  extensions,
		Scenarios:   scenarios,
		Normalizer:  normalize,
		Corrections: invoiceCorrectionDefinitions,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeBillInvoice(obj)
	case *tax.Combo:
		normalizeTaxCombo(obj)
	}
}
