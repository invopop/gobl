package bill

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

var defaultInvoiceTags = &tax.TagSet{
	Schema: ShortSchemaInvoice,
	List: []*cbc.Definition{
		// Simplified invoices are issued when the complete fiscal details of
		// a customer are not available.
		{
			Key: tax.TagSimplified,
			Name: i18n.String{
				i18n.EN: "Simplified Invoice",
			},
			Desc: i18n.String{
				i18n.EN: here.Doc(`
					Used for B2C transactions when the client details are not available, check with
					local authorities for limits.
				`),
			},
		},

		// Reverse Charge mechanism is used when the supplier is not
		// required to charge VAT on the invoice and the customer is
		// responsible for paying the VAT to the tax authorities.
		{
			Key: tax.TagReverseCharge,
			Name: i18n.String{
				i18n.EN: "Reverse Charge",
			},
			Desc: i18n.String{
				i18n.EN: here.Doc(`
					Applied when the *customer* is responsible for paying taxes to the tax authorities. Often used
					when the supplier is not registered for tax in the customer's country, or for special cases
					inside the same country when the seller is unlikely to be able to collect the tax.
				`),
			},
		},

		// Self-billed invoices are issued by the customer instead of the
		// supplier. This is usually done when the customer is a large
		// company and the supplier is a small company.
		{
			Key: tax.TagSelfBilled,
			Name: i18n.String{
				i18n.EN: "Self-billed",
			},
			Desc: i18n.String{
				i18n.EN: here.Doc(`
					Used when the customer or third party issues the invoice on behalf of the supplier.
				`),
			},
		},

		// Customer rates (mainly for digital goods inside EU)
		{
			Key: tax.TagCustomerRates,
			Name: i18n.String{
				i18n.EN: "Customer rates",
			},
			Desc: i18n.String{
				i18n.EN: here.Doc(`
					When set, implies that taxes rates should be determined from the customer's location
					as opposed to the supplier's. This is typically used for digital goods and services.
				`),
			},
		},

		// Partial invoice document, implying that this is only a first part
		// and a final invoice for the remaining amount will be made later.
		// A few regimes use this tag to classify invoices, notably Italy.
		{
			Key: tax.TagPartial,
			Name: i18n.String{
				i18n.EN: "Partial",
			},
			Desc: i18n.String{
				i18n.EN: here.Doc(`
					Indicates that this invoice is a partial document, meaning it is not the final invoice
					for the transaction. This is often used in construction or large projects where multiple
					invoices are issued for different stages of the work.
				`),
			},
		},

		// Special bypass tag used to skip calculations on the document.
		{
			Key: tax.TagBypass,
			Name: i18n.String{
				i18n.EN: "Bypass",
			},
			Desc: i18n.String{
				i18n.EN: here.Doc(`
					The bypass tag is used for special circumstances where calculations on billing documents should
					be skipped. Normalization and validation will still occur, but no automatic tax or total calculations
					will be performed. This is useful for correcting documents or when importing historical data where
					the original calculations need to be preserved.
				`),
			},
		},
	},
}

// DefaultInvoiceTags is a convenience function to get the default invoice tags.
func DefaultInvoiceTags() *tax.TagSet {
	return defaultInvoiceTags
}
