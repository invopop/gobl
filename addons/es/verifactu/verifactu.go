// Package verifactu provides the Verifactu addon
package verifactu

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

const (
	// V1 for Verifactu versions 1.x
	V1 cbc.Key = "es-verifactu-v1"
)

// Official stamps or codes validated by government agencies
const (
	// StampQR contains the URL included in the QR code.
	StampQR cbc.Key = "verifactu-qr"
)

func init() {
	tax.RegisterAddonDef(newAddon())
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
		Validator:   validate,
		Scenarios:   scenarios,
		Normalizer:  normalize,
		Corrections: invoiceCorrectionDefinitions,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	case *tax.Combo:
		normalizeTaxCombo(obj)
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Combo:
		return validateTaxCombo(obj)
	}
	return nil
}
