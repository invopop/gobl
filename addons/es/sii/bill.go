package sii

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

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
	// SII implying that scenarios cannot be used.
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
			ExtKeyThirdPartyIssuer: "S",
		})
	}
	if inv.Ordering != nil && inv.Ordering.Issuer != nil {
		inv.Tax = inv.Tax.MergeExtensions(tax.Extensions{
			ExtKeyThirdPartyIssuer: "S",
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

func normalizeBillLine(line *bill.Line) {
	if line == nil || line.Item == nil {
		return
	}
	vt := line.Taxes.Get(tax.CategoryVAT)
	if vt == nil {
		return
	}
	switch line.Item.Key {
	case org.ItemKeyGoods:
		vt.Ext = vt.Ext.Set(ExtKeyProduct, ExtCodeProductGoods)
	case org.ItemKeyServices:
		vt.Ext = vt.Ext.Set(ExtKeyProduct, ExtCodeProductServices)
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
		validation.Field(&inv.Customer,
			validation.When(
				!inv.Tax.GetExt(ExtKeyDocType).In("F2", "R5"), // not simplified
				validation.Required,
			),
			validation.By(validateInvoiceCustomer),
			validation.Skip,
		),
		validation.Field(&inv.Tax,
			validation.Required,
			validation.By(validateInvoiceTax(inv.Type)),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.By(validateInvoiceLines),
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
			// SII requires all Tax IDs to have a code. Sales into
			// countries without a specific Tax ID code will have to enter
			// something here regardless, or issue simplified invoices.
			tax.RequireIdentityCode,
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
			validation.Length(0, 500),
			validation.Skip,
		),
	)
}

func validateInvoiceLines(val any) error {
	lines, _ := val.([]*bill.Line)
	if lines == nil {
		return nil
	}

	var ref *tax.Combo
	for _, l := range lines {
		for _, tc := range l.Taxes {
			if !tc.Category.In(tax.CategoryVAT, es.TaxCategoryIGIC) {
				continue
			}
			if ref == nil {
				// first tax combo becomes the reference
				ref = tc
				continue
			}
			if ref.Ext.Has(ExtKeyProduct) != tc.Ext.Has(ExtKeyProduct) {
				return fmt.Errorf("`%s` must be present in all tax combos or none", ExtKeyProduct)
			}
			if ref.Ext.Get(ExtKeyRegime) != tc.Ext.Get(ExtKeyRegime) {
				return fmt.Errorf("`%s` must be the same in all tax combos", ExtKeyRegime)
			}
		}
	}

	return nil
}
