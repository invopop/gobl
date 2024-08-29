package hu

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// invoiceValidator adds validation checks to invoices which are relevant
// for the region.
type invoiceValidator struct {
	inv *bill.Invoice
}

func validateInvoice(inv *bill.Invoice) error {
	v := &invoiceValidator{inv: inv}
	return v.validate()
}

func (v *invoiceValidator) validate() error {
	inv := v.inv
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Type,
			validation.In(
				bill.InvoiceTypeStandard,
			),
			validation.Skip,
		),
		validation.Field(&inv.IssueDate,
			cal.DateAfter(*cal.NewDate(2010, 1, 1)),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.Required,
			validation.By(v.supplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.By(v.customer),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) supplier(value interface{}) error {
	obj, _ := value.(*org.Party)
	if obj == nil {
		return nil
	}

	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.By(v.localTaxID),
			validation.Skip,
		),
		validation.Field(&obj.Name,
			validation.Required,
			validation.Skip,
		),
		validation.Field(&obj.Addresses,
			validation.Length(1, 0),
			validation.Skip,
		),
		validation.Field(&obj.Identities,
			validation.When(
				isGroupVatID(obj.TaxID),
				validation.Required,
				validation.Length(1, 1),
				validation.By(v.memberGroupID),
			),
			validation.Skip),
	)
}

// The supplier must be a Hungarian entity.
func (v *invoiceValidator) localTaxID(value interface{}) error {
	obj, _ := value.(*tax.Identity)
	if obj == nil {
		return nil
	}

	return validation.ValidateStruct(obj,
		validation.Field(&obj.Country,
			validation.In(l10n.HU.Tax()),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) customer(value interface{}) error {
	obj, _ := value.(*org.Party)
	if obj == nil {
		return nil
	}

	return validation.ValidateStruct(obj,
		validation.Field(&obj.Identities,
			validation.When(
				isGroupVatID(obj.TaxID),
				validation.Required,
				validation.Length(1, 1),
				validation.By(v.memberGroupID),
			),
			validation.Skip),
	)
}

func isGroupVatID(taxID *tax.Identity) bool {
	if taxID == nil {
		return false
	}
	return len(taxID.Code) == 11 && taxID.Code.String()[8] == '5'
}

func (v *invoiceValidator) memberGroupID(value interface{}) error {
	obj, _ := value.([]*org.Identity)
	if obj == nil {
		return nil
	}

	return validation.ValidateStruct(obj[0],
		validation.Field(&obj[0].Code,
			validation.Length(11, 11),
			isGroupMember,
			validation.Skip,
		),
	)
}

var isGroupMember = validation.NewStringRule(checkGroupMemberDigit, "must be a group member ID")

func checkGroupMemberDigit(code string) bool {
	return len(code) == 11 && code[8:9] == "4"
}
