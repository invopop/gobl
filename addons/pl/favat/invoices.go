package favat

import (
	"fmt"

	"github.com/invopop/gobl/bill"
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
		// Only commercial and simplified supported at this time for Poland.
		// Rectification state determined by Preceding value.
		validation.Field(&inv.Type, validation.In(
			bill.InvoiceTypeStandard,
			bill.InvoiceTypeCreditNote,
		)),
		validation.Field(&inv.Preceding,
			validation.When(
				inv.Type.In(bill.InvoiceTypeCreditNote),
				validation.Required,
			),
			validation.Each(validation.By(v.preceding)),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.Required,
			validation.By(v.supplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.When(
				!inv.HasTags(tax.TagSimplified),
				validation.Required,
				validation.By(v.commercialCustomer),
			),
			validation.Skip,
		),
		validation.Field(&inv.Notes,
			validation.By(validateExemptionNote(inv)),
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
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) commercialCustomer(value interface{}) error {
	obj, _ := value.(*org.Party)
	if obj == nil {
		return nil
	}
	if obj.TaxID == nil {
		return nil // validation already handled, this prevents panics
	}
	// Customers must have a tax ID if a Polish entity
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.Required,
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) preceding(value interface{}) error {
	obj, ok := value.(*org.DocumentRef)
	if !ok || obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Ext,
			tax.ExtensionsRequire(ExtKeyEffectiveDate),
		),
		validation.Field(&obj.Reason, validation.Required),
	)
}

func normalizeInvoice(inv *bill.Invoice) {
	if inv.HasTags(tax.TagSelfBilled) {
		inv.Tax = inv.Tax.MergeExtensions(tax.Extensions{
			ExtKeySelfBilling: "1",
		})
	}

	if inv.HasTags(tax.TagReverseCharge) {
		inv.Tax = inv.Tax.MergeExtensions(tax.Extensions{
			ExtKeyReverseCharge: "1",
		})
	}

	// Even if we know that the invoice is exempt (has tag tax.KeyExempt), we cannot autogenerate values
	// under key ExtKeyExemption here, as there are multiple possible values for this extension.
}

func isExemptionNote(n *org.Note) bool {
	return n.Key == org.NoteKeyLegal && n.Src == ExtKeyExemption
}

// validateExemptionNote validates that when the invoice is marked as tax exempt, appropriate the note is added
func validateExemptionNote(inv *bill.Invoice) validation.RuleFunc {
	return func(val any) error {
		notes, ok := val.([]*org.Note)
		if !ok {
			return fmt.Errorf("expected []*org.Note, got %T", val)
		}

		exemptionCode := inv.Tax.Ext.Get(ExtKeyExemption)

		count := 0
		for index, note := range notes {
			if isExemptionNote(note) {
				if exemptionCode == "" {
					// Tax extension for exemption is not set, but an exemption note is present
					return fmt.Errorf("(%d: unexpected exemption note)", index)
				}
				if count > 0 {
					// More than one exemption note is present
					return fmt.Errorf("(%d: too many exemption notes)", index)
				}
				if exemptionCode != note.Code {
					// Code given in the note and in the extension are different
					return fmt.Errorf("(%d: note code %s must match extension %s)", index, note.Code, exemptionCode)
				}
				count++
			}
		}

		if exemptionCode != "" && count == 0 {
			// Exemption code is set, but no exemption note is present
			return fmt.Errorf("missing exemption note for code %s", exemptionCode)
		}

		return nil
	}
}
