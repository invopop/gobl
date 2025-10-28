// Package ticket handles the validation rules in order to use
// GOBL with the Italian Agenzia delle Entrate format.
package ticket

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// Key to identify the AdE ticket addon
const (
	// V1 for AdE format
	V1 cbc.Key = "it-ticket-v1"
)

// Official stamps or codes validated by government agencies
const (
	// StampRef is the key to identify the reference provided by the AdE once the ticket is accepted
	// This code needs to be printed on the ticket
	StampRef cbc.Key = "ade-ref"
	// StampVoidRef is the key to identify the reference provided by the AdE once the ticket is voided
	// This code needs to be printed on the ticket representing the void action
	StampVoidRef cbc.Key = "ade-void-ref"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

// This validation follows the rules of the Italian Agenzia delle Entrate
// This addon will then be used to create documents using the following services
// https://www.agenziaentrate.gov.it/portale/schede/comunicazioni/fatture-e-corrispettivi
func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "Italy AdE ticket v1.x",
		},
		Sources: []*cbc.Source{
			{
				Title:       i18n.NewString("Italian AdE Cassetto Fiscale"),
				URL:         "https://www.agenziaentrate.gov.it/portale/schede/comunicazioni/fatture-e-corrispettivi",
				ContentType: "application/pdf",
			},
			{
				Title:       i18n.NewString("Italian AdE Fattura e Corrispettivi"),
				URL:         "https://www.agenziaentrate.gov.it/portale/schede/comunicazioni/fatture-e-corrispettivi",
				ContentType: "application/pdf",
			},
		},
		Extensions:  extensions,
		Validator:   validate,
		Normalizer:  normalize,
		Corrections: invoiceCorrectionDefinitions,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	case *org.Item:
		normalizeOrgItem(obj)
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
