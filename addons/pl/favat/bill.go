package favat

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Based on the keys, the extension should be set

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

func validateBillInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Type, validation.In(
			bill.InvoiceTypeStandard,
			bill.InvoiceTypeCreditNote,
		)),
		validation.Field(&inv.Preceding,
			validation.When(
				inv.Type.In(bill.InvoiceTypeCreditNote),
				validation.Required,
			),
			validation.Each(validation.By(validateBillInvoicePreceding)),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.Required,
			validation.By(validateBillInvoiceSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.When(
				!inv.HasTags(tax.TagSimplified),
				validation.Required,
				validation.By(validateBillInvoiceCustomer),
			),
			validation.Skip,
		),
		validation.Field(&inv.Notes,
			validation.By(validateExemptionNote(inv)),
			validation.Skip,
		),
	)
}

func validateBillInvoiceSupplier(value any) error {
	obj, ok := value.(*org.Party)
	if !ok || obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.Name,
			validation.Required,
			validation.Skip,
		),
		validation.Field(&obj.Addresses,
			validation.Required,
			validation.By(validateAddresses),
			validation.Skip,
		),
	)
}

func validateAddresses(val any) error {
	addresses, ok := val.([]*org.Address)
	if !ok || len(addresses) == 0 {
		return nil
	}

	return validation.ValidateStruct(addresses[0],
		validation.Field(&addresses[0].Country,
			validation.Required,
			validation.Skip,
		),
		validation.Field(&addresses[0].Street,
			validation.Required,
			validation.Skip,
		),
	)
}

func validateBillInvoiceCustomer(value any) error {
	obj, ok := value.(*org.Party)
	if !ok || obj == nil {
		return nil
	}

	// Customers must have a tax ID if a Polish entity
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
		validation.Field(&obj.Identities,
			validation.When(
				obj.Ext.Get(ExtKeyJST) == "1",                // Yes
				validation.By(validateIdentityWithRole("8")), // Local Government Unit (LGU) - recipient
			),
			validation.When(
				obj.Ext.Get(ExtKeyGroupVAT) == "1",            // Yes
				validation.By(validateIdentityWithRole("10")), // GV member - recipient
			),
			validation.Skip,
		),
	)
}

// validateIdentityWithRole checks that there is an identity with the given role extension
// and that the identity has a code.
func validateIdentityWithRole(roleCode cbc.Code) validation.RuleFunc {
	return func(value any) error {
		identities, ok := value.([]*org.Identity)
		if !ok {
			return nil
		}

		for _, identity := range identities {
			if identity.Ext.Get(ExtKeyThirdPartyRole) == roleCode && identity.Code != "" {
				return nil
			}
		}

		return fmt.Errorf("missing identity with role '%s' and code", roleCode)
	}
}

func validateBillInvoicePreceding(value any) error {
	obj, ok := value.(*org.DocumentRef)
	if !ok || obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.IssueDate, validation.Required),
		validation.Field(&obj.Code, validation.Required),
	)
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
				if count > 0 {
					// More than one exemption note is present
					return fmt.Errorf("(%d: too many exemption notes)", index)
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
