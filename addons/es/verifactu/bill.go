package verifactu

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
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

func normalizeBillInvoice(inv *bill.Invoice) {
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

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		// Preceding documents
		// Code 01: preceding required when corrective
		rules.When(
			rules.By("corrective", invoiceTypeIsCorrective),
			rules.Field("preceding",
				rules.Assert("01", "preceding documents are required for corrective invoices", rules.Required),
			),
		),
		// Code 02: each preceding issue date required
		rules.Field("preceding",
			rules.Each(
				rules.When(
					rules.By("not nil", precedingDocIsNotNil),
					rules.Field("issue_date",
						rules.Assert("02", "issue date is required", rules.Required),
					),
				),
			),
		),
		// Code 03: each preceding tax required when corrective
		rules.When(
			rules.By("corrective", invoiceTypeIsCorrective),
			rules.Field("preceding",
				rules.Each(
					rules.When(
						rules.By("not nil", precedingDocIsNotNil),
						rules.Field("tax",
							rules.Assert("03", "preceding invoice tax data is required for corrective invoices", rules.Required),
						),
					),
				),
			),
		),
		// Customer - simplified invoices (F2 or R5)
		// Code 04: no tax_id on simplified customer
		// Code 05: no identity type ext on simplified customer
		rules.When(
			rules.By("simplified", isSimplifiedInvoice),
			rules.Field("customer",
				rules.Assert("04", "customer tax ID must not be set for simplified invoices",
					rules.By("no tax_id", simplifiedCustomerHasNoTaxID),
				),
				rules.Assert("05", "customer identity type extension not allowed for simplified invoices",
					rules.By("no identity type ext", simplifiedCustomerHasNoIdentityType),
				),
			),
		),
		// Customer - standard invoices
		// Code 06: customer required
		// Code 07: customer must have tax_id or identity
		// Code 08: customer tax_id must have code
		rules.When(
			rules.By("not simplified", isNotSimplifiedInvoice),
			rules.Field("customer",
				rules.Assert("06", "customer is required", rules.Required),
				rules.Assert("07", "must have a tax_id or an identity with ext 'es-verifactu-v1-identity-type'",
					rules.By("has tax_id or identity", customerHasTaxIDOrIdentity),
				),
				rules.Assert("08", "tax ID must have a code",
					rules.By("tax_id has code", customerTaxIDHasCode),
				),
			),
		),
		// Invoice tax extensions
		// Code 09: tax required
		// Code 10: doc_type required
		// Code 13: correction_type required when credit/debit doc type
		rules.Field("tax",
			rules.Assert("09", "tax is required", rules.Required),
			rules.Field("ext",
				rules.Assert("10", "doc type is required",
					rules.By("has doc type", extHasDocType),
				),
				rules.When(
					rules.By("credit/debit doc type", extDocTypeIsCreditDebit),
					rules.Assert("13", "correction type is required",
						rules.By("has correction type", extHasCorrectionType),
					),
				),
			),
		),
		// Code 11: standard invoice doc type must be F1, F2, or F3
		rules.When(
			rules.By("standard invoice type", invoiceTypeIsStandard),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("11", "doc type for standard invoices must be F1, F2, or F3",
						rules.By("standard doc type", extDocTypeInStandard),
					),
				),
			),
		),
		// Code 12: corrective invoice doc type must be R1-R5
		rules.When(
			rules.By("corrective invoice type", invoiceTypeIsCorrectiveType),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("12", "doc type for corrective invoices must be R1, R2, R3, R4, or R5",
						rules.By("corrective doc type", extDocTypeInCreditDebit),
					),
				),
			),
		),
		// Notes
		// Code 14: general note text max 500 characters
		rules.Field("notes",
			rules.Each(
				rules.When(
					rules.By("general note", isGeneralNote),
					rules.Field("text",
						rules.Assert("14", "general note text must be 500 characters or less", rules.Length(0, 500)),
					),
				),
			),
		),
		// Lines
		// Code 15: each line must have at least one of VAT, IGIC, or IPSI
		rules.Field("lines",
			rules.Each(
				rules.Field("taxes",
					rules.Assert("15", "must include at least one of VAT, IGIC, or IPSI",
						rules.By("has required tax category", lineHasRequiredTax),
					),
				),
			),
		),
	)
}

func isSimplifiedInvoice(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && inv.Tax.GetExt(ExtKeyDocType).In("F2", "R5")
}

func isNotSimplifiedInvoice(val any) bool {
	return !isSimplifiedInvoice(val)
}

func invoiceTypeIsCorrective(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && inv.Type == bill.InvoiceTypeCorrective
}

func invoiceTypeIsStandard(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && inv.Type == bill.InvoiceTypeStandard
}

func invoiceTypeIsCorrectiveType(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && inv.Type.In(bill.InvoiceTypeCorrective, bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote)
}

func precedingDocIsNotNil(val any) bool {
	ref, ok := val.(*org.DocumentRef)
	return ok && ref != nil
}

func extHasDocType(val any) bool {
	ext, ok := val.(tax.Extensions)
	return ok && ext.Has(ExtKeyDocType)
}

func extDocTypeInStandard(val any) bool {
	ext, ok := val.(tax.Extensions)
	return ok && ext.Get(ExtKeyDocType).In("F1", "F2", "F3")
}

func extDocTypeInCreditDebit(val any) bool {
	ext, ok := val.(tax.Extensions)
	return ok && ext.Get(ExtKeyDocType).In("R1", "R2", "R3", "R4", "R5")
}

func extDocTypeIsCreditDebit(val any) bool {
	return extDocTypeInCreditDebit(val)
}

func extHasCorrectionType(val any) bool {
	ext, ok := val.(tax.Extensions)
	return ok && ext.Has(ExtKeyCorrectionType)
}

func simplifiedCustomerHasNoTaxID(val any) bool {
	p, ok := val.(*org.Party)
	return !ok || p == nil || p.TaxID == nil
}

func simplifiedCustomerHasNoIdentityType(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	return org.IdentityForExtKey(p.Identities, ExtKeyIdentityType) == nil
}

func customerHasTaxIDOrIdentity(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true // nil customer handled by Required check
	}
	return p.TaxID != nil || org.IdentityForExtKey(p.Identities, ExtKeyIdentityType) != nil
}

func customerTaxIDHasCode(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	return p.TaxID == nil || !p.TaxID.Code.IsEmpty()
}

func isGeneralNote(val any) bool {
	note, ok := val.(*org.Note)
	return ok && note != nil && note.Key == org.NoteKeyGeneral
}

func lineHasRequiredTax(val any) bool {
	ts, ok := val.(tax.Set)
	if !ok {
		return false
	}
	for _, c := range ts {
		if c != nil && c.Category.In(tax.CategoryVAT, es.TaxCategoryIGIC, es.TaxCategoryIPSI) {
			return true
		}
	}
	return false
}
