package favat

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
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

func isExemptionNote(n *org.Note) bool {
	return n.Key == org.NoteKeyLegal && n.Src == ExtKeyExemption
}

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Field("type",
			rules.Assert("01", "invoice type must be standard or credit-note",
				is.In(bill.InvoiceTypeStandard, bill.InvoiceTypeCreditNote),
			),
		),
		// All preceding documents need issue_date and code
		rules.Field("preceding",
			rules.Each(
				rules.Field("issue_date",
					rules.Assert("02", "preceding issue date is required", is.Present),
				),
				rules.Field("code",
					rules.Assert("03", "preceding code is required", is.Present),
				),
			),
		),
		// Supplier validation
		rules.Field("supplier",
			rules.Assert("04", "supplier is required", is.Present),
			rules.Field("tax_id",
				rules.Assert("05", "supplier tax ID is required", is.Present),
				rules.Field("code",
					rules.Assert("06", "supplier tax ID code required", is.Present),
				),
			),
			rules.Field("name",
				rules.Assert("07", "supplier name is required", is.Present),
			),
			rules.Field("addresses",
				rules.Assert("08", "supplier addresses are required", is.Present),
				rules.Assert("09", "supplier first address must have a country",
					is.Func("first address country", firstAddressHasCountry),
				),
				rules.Assert("10", "supplier first address must have a street",
					is.Func("first address street", firstAddressHasStreet),
				),
			),
		),
		// Customer validation (required unless simplified)
		rules.When(is.Func("not simplified", invoiceNotSimplified),
			rules.Field("customer",
				rules.Assert("11", "customer is required", is.Present),
				rules.Field("tax_id",
					rules.Assert("12", "customer tax ID is required", is.Present),
				),
			),
		),
		// Customer JST identity check (invoice-level, needs both customer.ext and customer.identities)
		rules.Assert("13",
			fmt.Sprintf("customer requires identity with role '%s' and code for JST", cbc.Code("8")),
			is.Func("JST identity", invoiceCustomerJSTIdentityValid),
		),
		// Customer GroupVAT identity check (invoice-level)
		rules.Assert("14",
			fmt.Sprintf("customer requires identity with role '%s' and code for GroupVAT", cbc.Code("10")),
			is.Func("GroupVAT identity", invoiceCustomerGroupVATIdentityValid),
		),
		// Exemption note validation (object-level)
		rules.Assert("15", "missing exemption note",
			is.Func("exemption note present", invoiceExemptionNotePresent),
		),
		rules.Assert("16", "too many exemption notes",
			is.Func("single exemption note", invoiceSingleExemptionNote),
		),
	)
}

func firstAddressHasCountry(val any) bool {
	addresses, ok := val.([]*org.Address)
	if !ok || len(addresses) == 0 {
		return true // empty addresses caught by "is required"
	}
	return addresses[0].Country != ""
}

func firstAddressHasStreet(val any) bool {
	addresses, ok := val.([]*org.Address)
	if !ok || len(addresses) == 0 {
		return true // empty addresses caught by "is required"
	}
	return addresses[0].Street != ""
}

func invoiceNotSimplified(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	return !inv.HasTags(tax.TagSimplified)
}

func invoiceCustomerJSTIdentityValid(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Customer == nil {
		return true
	}
	if inv.Customer.Ext.Get(ExtKeyJST) != "1" {
		return true // JST not enabled, skip
	}
	return hasIdentityWithRole(inv.Customer.Identities, "8")
}

func invoiceCustomerGroupVATIdentityValid(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Customer == nil {
		return true
	}
	if inv.Customer.Ext.Get(ExtKeyGroupVAT) != "1" {
		return true // GroupVAT not enabled, skip
	}
	return hasIdentityWithRole(inv.Customer.Identities, "10")
}

func hasIdentityWithRole(identities []*org.Identity, roleCode cbc.Code) bool {
	for _, identity := range identities {
		if identity.Ext.Get(ExtKeyThirdPartyRole) == roleCode && identity.Code != "" {
			return true
		}
	}
	return false
}

func invoiceExemptionNotePresent(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	exemptionCode := inv.Tax.Ext.Get(ExtKeyExemption)
	if exemptionCode == "" {
		return true
	}
	for _, note := range inv.Notes {
		if isExemptionNote(note) {
			return true
		}
	}
	return false
}

func invoiceSingleExemptionNote(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	count := 0
	for _, note := range inv.Notes {
		if isExemptionNote(note) {
			count++
		}
	}
	return count <= 1
}
