// Package sii provides the SII addon
package sii

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
	// Key identifies the SII addon family. Individual versions append a
	// suffix; the family key is used as the fault-code namespace so that
	// rules that carry across versions keep stable codes.
	Key cbc.Key = "es-sii"

	// V1 for SII versions 1.x
	V1 cbc.Key = Key + "-v1"
)

func init() {
	tax.RegisterAddonDef(newAddon())
	rules.RegisterWithGuard(
		Key.String(),
		rules.GOBL.Add("ES-SII"),
		is.InContext(tax.AddonIn(V1)),
		billInvoiceRules(),
		taxComboRules(),
	)
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "Spain SII V1.x",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Suministro Inmediato de Información (SII) at AEAT portal"),
				URL:   "https://sede.agenciatributaria.gob.es/Sede/iva/suministro-inmediato-informacion.html",
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
		normalizeInvoice(obj)
	case *bill.Line:
		normalizeBillLine(obj)
	case *tax.Combo:
		normalizeTaxCombo(obj)
	}
}
