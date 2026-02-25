// Package ae provides the tax region definition for United Arab Emirates.
package ae

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax region definition for AE.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "AE",
		Currency:  currency.AED,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "United Arab Emirates",
			i18n.AR: "الإمارات العربية المتحدة",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				The UAE tax system is administered by the Federal Tax Authority (FTA). VAT was
				introduced on January 1, 2018, with a standard rate of 5% applying to most goods
				and services.

				Businesses must register for VAT if taxable supplies and imports exceed AED 375,000
				in a 12-month period, with voluntary registration available above AED 187,500.
				Registered businesses receive a Tax Registration Number (TRN) which must be included
				on all tax invoices.

				VAT rates include 5% standard rate for most goods and services, 0% for certain
				essential goods, exports, and specific services, and exempt supplies covering some
				financial services and residential real estate.

				Simplified VAT invoices may be used when the recipient is not VAT registered, or
				when the transaction value does not exceed AED 10,000 for VAT-registered recipients.
				Credit notes are supported for correcting invoices.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Federal Tax Authority - VAT"),
				URL:   "https://tax.gov.ae/en/taxes/Vat/vat.topics/registration.for.vat.aspx",
			},
		},
		TimeZone: "Asia/Dubai",
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios,
		},
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
			},
		},
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: taxCategories,
	}
}

// Validate function assesses the document type to determine if validation is required.
// Note that, under the AE tax regime, validation of the supplier's tax ID is not
// necessary if the business does not meet the mandatory registration threshold.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Normalize attempts to clean up the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)

	}
}
