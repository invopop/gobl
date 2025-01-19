package sg

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Reference: https://www.iras.gov.sg/media/docs/default-source/e-tax/etaxguide_gst_gst-general-guide-for-businesses(1).pdf?sfvrsn=8a66716d_97 (pg 26-27)

// Invoice type tags
const (
	TagInvoiceReceipt cbc.Key = "receipt"
)

var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*cbc.Definition{
		{
			Key: TagInvoiceReceipt,
			Name: i18n.String{
				i18n.EN: "Receipt",
			},
		},
	},
}

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.When(
				inv.HasTags(TagInvoiceReceipt),
				validation.By(validateRecieptSupplier),
				validation.Skip,
			).Else(
				validation.By(validateInvoiceSupplier),
				validation.Skip,
			),
		),
		validation.Field(&inv.Customer,
			validation.When(
				inv.HasTags(TagInvoiceReceipt) || inv.HasTags(tax.TagSimplified),
				validation.Skip,
			).Else(validation.Required),
		),
	)
}

func validateInvoiceSupplier(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
		validation.Field(&p.Name,
			validation.Required,
		),
		validation.Field(&p.Addresses,
			validation.Required,
		),
	)
}

func validateRecieptSupplier(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
	)
}
