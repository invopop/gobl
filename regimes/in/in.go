// Package in provides models for dealing with India.
package in

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax region definition for India.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "IN",
		Currency:  currency.INR,
		TaxScheme: tax.CategoryGST,
		Name: i18n.String{
			i18n.EN: "India",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				India follows a dual GST (Goods and Services Tax) model, where both the Central
				and State Governments levy taxes. CGST (Central GST) is levied by the Central
				Government, SGST/UTGST (State/Union Territory GST) by state governments, and
				IGST (Integrated GST) on interstate supplies and imports.

				For intrastate supplies, CGST and SGST/UTGST apply in equal proportions. For
				interstate supplies and imports, IGST applies at a rate equivalent to CGST plus
				SGST. A Compensation Cess may apply on luxury and sin goods. Due to the dual
				model, tax rate allocations between central and state must be managed at the
				application level.

				GST rates vary by goods and services: 0.25%-3% for precious metals, 5% for
				basic goods, 12%-18% for standard goods and services, and 28% for luxury
				items. Exports and supplies to Special Economic Zones are zero-rated. Exempt
				supplies include fresh fruits and vegetables, educational services, and public
				road tolls.

				Businesses are identified by their GSTIN (Goods and Services Tax Identification
				Number), a unique 15-digit identifier with format and checksum validation.
				Items on invoices must include HSN (Harmonized System of Nomenclature) codes
				for goods classification. India supports both credit notes and debit notes for
				invoice corrections.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("GST Portal"),
				URL:   "https://www.gst.gov.in/",
			},
		},
		TimeZone: "Asia/Kolkata",
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios,
		},
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
			},
		},
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: taxCategories,
	}
}

// Validate function assesses the document type to determine if validation is required.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *org.Identity:
		return validateOrgIdentity(obj)
	case *org.Item:
		return validateOrgItem(obj)
	}
	return nil
}

// Normalize attempts to clean up the object passed to it.
func Normalize(doc interface{}) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	case *org.Identity:
		normalizeOrgIdentity(obj)
	}
}
