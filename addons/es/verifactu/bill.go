package verifactu

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// forbiddenChars contains characters that are not allowed in certain string fields
var forbiddenChars = []rune{'<', '>', '"', '\'', '='}

var invoiceCorrectionDefinitions = tax.CorrectionSet{
	{
		Schema: bill.ShortSchemaInvoice,
		Extensions: []cbc.Key{
			ExtKeyDocType,
		},
		CopyTax: true,
	},
}

func normalizeInvoice(inv *bill.Invoice) {
	// Try to move any preceding choices to the document level
	for _, row := range inv.Preceding {
		if row == nil || len(row.Ext) == 0 {
			continue
		}
		found := false
		if row.Ext.Has(ExtKeyDocType) {
			if inv.Tax == nil || !found {
				inv.Tax = inv.Tax.MergeExtensions(tax.Extensions{
					ExtKeyDocType: row.Ext[ExtKeyDocType],
				})
				found = true // only assign first one
			}
			delete(row.Ext, ExtKeyDocType)
		}
	}

	// Try to normalize the correction type, which is especially complex for
	// Verifactu implying that scenarios cannot be used.
	switch inv.Type {
	case bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote:
		inv.Tax = inv.Tax.MergeExtensions(tax.Extensions{
			ExtKeyCorrectionType: "I",
		})
	case bill.InvoiceTypeCorrective:
		inv.Tax = inv.Tax.MergeExtensions(tax.Extensions{
			ExtKeyCorrectionType: "S",
		})
	}

	// Set default correction type, unless already provided.
	switch inv.Type {
	case bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote, bill.InvoiceTypeCorrective:
		// Don't try to override a previously set document type.
		// This is non-deterministic. May be overwritten by user *or*
		// scenarios.
		if !inv.Tax.Ext.Get(ExtKeyDocType).In("R2", "R3", "R4", "R5") {
			inv.Tax.Ext[ExtKeyDocType] = "R1"
		}
	}

	// Normalize the third party details
	if inv.HasTags(tax.TagSelfBilled) {
		inv.Tax = inv.Tax.MergeExtensions(tax.Extensions{
			ExtKeyIssuerType: ExtCodeIssuerTypeCustomer,
		})
	}
	if inv.Ordering != nil && inv.Ordering.Issuer != nil {
		inv.Tax = inv.Tax.MergeExtensions(tax.Extensions{
			ExtKeyIssuerType: ExtCodeIssuerTypeThirdParty,
		})
	}

	normalizeInvoicePartyIdentity(inv.Customer)
}

func normalizeInvoicePartyIdentity(cus *org.Party) {
	if cus == nil {
		return
	}
	if cus.TaxID != nil && cus.TaxID.Country == "ES" && cus.TaxID.Code != "" {
		// Spanish NIFs are already handled
		return
	}
	if len(cus.Identities) == 0 {
		// nothing to do if no identities
		return
	}
	id := cus.Identities[0]
	var code cbc.Code
	switch id.Key {
	case org.IdentityKeyPassport:
		code = ExtCodeIdentityTypePassport
	case org.IdentityKeyForeign:
		code = ExtCodeIdentityTypeForeign
	case org.IdentityKeyResident:
		code = ExtCodeIdentityTypeResident
	case org.IdentityKeyOther:
		code = ExtCodeIdentityTypeOther
	}
	if !code.IsEmpty() {
		id.Ext = id.Ext.Merge(tax.Extensions{
			ExtKeyIdentityType: code,
		})
	}
}

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Preceding,
			validation.When(
				// When an invoice is a "sustitutiva", the taxes from the preceding invoice must be included.
				inv.Type.In(bill.InvoiceTypeCorrective),
				validation.Required,
			),
			validation.Each(
				validation.By(validateInvoicePreceding(inv)),
				validation.Skip,
			),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.By(validateInvoiceSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.When(
				!inv.Tax.GetExt(ExtKeyDocType).In("F2", "R5"), // not simplified
				validation.Required,
			),
			validation.By(validateInvoiceCustomer),
			validation.Skip,
		),
		validation.Field(&inv.Ordering,
			validation.By(validateInvoiceOrdering),
			validation.Skip,
		),
		validation.Field(&inv.Tax,
			validation.Required,
			validation.By(validateInvoiceTax(inv.Type)),
			validation.Skip,
		),
		validation.Field(&inv.Notes,
			validation.Each(
				validation.By(validateNote),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}

func validateInvoiceSupplier(val any) error {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.Name,
			validation.By(validateNoForbiddenChars),
			validation.Skip,
		),
	)
}

func validateInvoiceCustomer(val any) error {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	if p.TaxID == nil && org.IdentityForExtKey(p.Identities, ExtKeyIdentityType) == nil {
		return fmt.Errorf("must have a tax_id, or an identity with ext '%s'", ExtKeyIdentityType)
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.TaxID,
			// VERI*FACTU requires all Tax IDs to have a code. Sales into
			// countries without a specific Tax ID code will have to enter
			// something here regardless, or issue simplified invoices.
			tax.RequireIdentityCode,
			validation.Skip,
		),
		validation.Field(&p.Name,
			validation.By(validateNoForbiddenChars),
			validation.Skip,
		),
	)
}

func validateInvoiceOrdering(val any) error {
	o, ok := val.(*bill.Ordering)
	if !ok || o == nil {
		return nil
	}
	return validation.ValidateStruct(o,
		validation.Field(&o.Issuer,
			validation.By(validateInvoiceOrderingIssuer),
			validation.Skip,
		),
	)
}

func validateInvoiceOrderingIssuer(val any) error {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.Name,
			validation.By(validateNoForbiddenChars),
			validation.Skip,
		),
	)
}

var docTypesStandard = []cbc.Code{ // Standard invoices
	"F1", "F2", "F3",
}
var docTypesCreditDebit = []cbc.Code{ // Credit or Debit notes
	"R1", "R2", "R3", "R4", "R5",
}

func validateInvoiceTax(it cbc.Key) validation.RuleFunc {
	return func(val any) error {
		obj := val.(*bill.Tax)
		return validation.ValidateStruct(obj,
			validation.Field(&obj.Ext,
				tax.ExtensionsRequire(ExtKeyDocType),
				validation.When(
					it.In(bill.InvoiceTypeStandard),
					tax.ExtensionsHasCodes(
						ExtKeyDocType,
						docTypesStandard...,
					),
				),
				validation.When(
					it.In(bill.InvoiceTypeCorrective, bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote),
					tax.ExtensionsHasCodes(
						ExtKeyDocType,
						docTypesCreditDebit...,
					),
				),
				validation.When(
					obj.Ext.Get(ExtKeyDocType).In(docTypesCreditDebit...),
					tax.ExtensionsRequire(ExtKeyCorrectionType),
				),
				validation.Skip,
			),
		)
	}
}

func validateInvoicePreceding(inv *bill.Invoice) validation.RuleFunc {
	return func(val any) error {
		p, ok := val.(*org.DocumentRef)
		if !ok || p == nil {
			return nil
		}
		return validation.ValidateStruct(p,
			validation.Field(&p.IssueDate,
				validation.Required,
				validation.Skip,
			),
			validation.Field(&p.Tax,
				// Tax data of previous invoices is required for substitutions
				validation.When(
					inv.Type.In(bill.InvoiceTypeCorrective),
					validation.Required,
				),
				validation.Skip,
			),
		)
	}
}

func validateNote(val any) error {
	note, ok := val.(*org.Note)
	if !ok || note == nil || note.Key != org.NoteKeyGeneral {
		return nil
	}
	return validation.ValidateStruct(note,
		validation.Field(&note.Text,
			validation.By(validateNoForbiddenChars),
			validation.Length(0, 500),
			validation.Skip,
		),
	)
}

// validateNoForbiddenChars validates that a string doesn't contain any of the forbidden characters: < > " ' =
func validateNoForbiddenChars(val any) error {
	str, _ := val.(string)

	for _, char := range str {
		for _, forbidden := range forbiddenChars {
			if char == forbidden {
				return fmt.Errorf("contains forbidden character: %c", char)
			}
		}
	}

	return nil
}
