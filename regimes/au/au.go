package au

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New instantiates a new Australian regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   l10n.AU.Tax(),
		Currency:  currency.AUD,
		TaxScheme: tax.CategoryGST,
		Name: i18n.String{
			i18n.EN: "Australia",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Australia's indirect tax system is administered by the Australian
				Taxation Office (ATO). Goods and Services Tax (GST) applies at a
				standard rate of 10%, with zero-rated exports and exempt supplies for
				specific transactions such as certain financial services and residential
				rent.

				Businesses are identified by an Australian Business Number (ABN), an
				11-digit identifier used for GST registration and invoicing. Australian
				tax invoices must show the supplier's details and ABN, and invoices of
				AUD 1,000 or more must also identify the customer. Electronic invoicing
				is aligned with the Peppol PINT A-NZ specification.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("ATO - GST"),
				URL:   "https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/gst",
			},
			{
				Title: i18n.NewString("ATO - Tax invoices"),
				URL:   "https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/gst/tax-invoices",
			},
			{
				Title: i18n.NewString("Peppol PINT A-NZ BIS"),
				URL:   "https://docs.peppol.eu/poac/aunz/pint-aunz/bis/",
			},
			{
				Title: i18n.NewString("ATO - eInvoicing"),
				URL:   "https://www.ato.gov.au/businesses-and-organisations/invoicing-and-using-accounting-software/einvoicing",
			},
		},
		TimeZone: "Australia/Sydney",
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
		Categories: taxCategories(),
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateBillInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
