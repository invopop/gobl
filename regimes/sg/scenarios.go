// Package sg provides tax scenarios specific to Singapore GST regulations.
package sg

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// Reference: https://www.iras.gov.sg/media/docs/default-source/e-tax/etaxguide_gst_gst-general-guide-for-businesses(1).pdf?sfvrsn=8a66716d_97 (pg 26-27)

// Invoice type tags
const (
	TagInvoiceReceipt cbc.Key = "receipt"
)

func invoiceTags() *tax.TagSet {
	return &tax.TagSet{
		Schema: bill.ShortSchemaInvoice,
		List: []*cbc.Definition{
			{
				Key: TagInvoiceReceipt,
				Name: i18n.String{
					i18n.EN: "Receipt",
				},
				Desc: i18n.String{
					i18n.EN: `This type of invoice can be issued to a non-GST registered customer:

							1. Supplier name and GST registration number;
							2. The date of issue of the invoice.
							3. The total amount payable including tax.
							4. The word 'Price Payable includes GST'.`,
				},
			},
		},
	}
}

func invoiceScenarios() *tax.ScenarioSet {
	return &tax.ScenarioSet{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			// Reverse Charge - https://www.iras.gov.sg/media/docs/default-source/e-tax/gst-taxing-imported-services-by-way-of-reverse-charge-(2nd-edition).pdf?sfvrsn=64fc1f4a_23
			{
				Tags: []cbc.Key{tax.TagReverseCharge},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  tax.TagReverseCharge,
					Text: "Reverse Charge",
				},
				Desc: i18n.String{
					i18n.EN: `Reverse charge in Singapore applies mainly when a GST-registered business in Singapore imports services or digital services from overseas suppliers. Some examples include:

							- Advertising services from overseas platforms.
							- Consultancy, professional, or technical services from foreign providers.
							- Digital services such as software or streaming services supplied from overseas.`,
				},
			},
			// Simplified Tax Invoice
			{
				Tags: []cbc.Key{tax.TagSimplified},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  tax.TagSimplified,
					Text: "Price Payable includes GST",
				},
				Desc: i18n.String{
					i18n.EN: `This invoice can only be used when the total amount (inclusive of GST) is less than 1000 SGD and must include:

							1. Supplier name, address and GST registration number;
							2. An identifying number, e.g. invoice number.
							3. The date of issue of the invoice.
							4. Description of the goods or services supplied.
							5. The total amount payable including tax.
							6. The word "Price Payable includes GST".`,
				},
			},
			// Receipt
			{
				Tags: []cbc.Key{TagInvoiceReceipt},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  TagInvoiceReceipt,
					Text: "Price Payable includes GST",
				},
				Desc: i18n.String{
					i18n.EN: `This type of invoice can be issued to a non-GST registered customer:

							1. Supplier name and GST registration number;
							2. The date of issue of the invoice.
							3. The total amount payable including tax.
							4. The word "Price Payable includes GST".`,
				},
			},
			// Self-billed
			{
				Tags: []cbc.Key{tax.TagSelfBilled},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  tax.TagSelfBilled,
					Text: "Self-billed",
				},
				Desc: i18n.String{
					i18n.EN: "Self-billing is a billing arrangement between a GST-registered supplier and a GST-registered customer, where the customer, instead of the supplier, prepares the supplier's tax invoice/ customer accounting tax invoice and sends a copy to the supplier.",
				},
			},
		},
	}
}
