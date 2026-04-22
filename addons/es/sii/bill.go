package sii

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
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

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Assert("15", "invoice must be in EUR or provide exchange rate for conversion", currency.CanConvertTo(currency.EUR)),
		// Preceding documents
		// Code 01: preceding required when corrective
		rules.When(
			bill.InvoiceTypeIn(bill.InvoiceTypeCorrective),
			rules.Field("preceding",
				rules.Assert("01", "preceding documents are required for corrective invoices", is.Present),
			),
		),
		// Code 02: each preceding issue date required
		rules.Field("preceding",
			rules.Each(
				rules.Field("issue_date",
					rules.Assert("02", "preceding issue date is required", is.Present),
				),
			),
		),
		// Code 03: each preceding tax required when corrective
		rules.When(
			bill.InvoiceTypeIn(bill.InvoiceTypeCorrective),
			rules.Field("preceding",
				rules.Each(
					rules.Field("tax",
						rules.Assert("03", "preceding invoice tax data is required for corrective invoices", is.Present),
					),
				),
			),
		),
		// Customer - standard (non-simplified) invoices
		// Code 04: customer required
		// Code 05: customer must have tax_id or identity with ext
		// Code 06: customer tax_id must have code
		rules.When(
			is.Func("not simplified", isNotSimplifiedInvoice),
			rules.Field("customer",
				rules.Assert("04", "customer is required", is.Present),
				rules.Assert("05", fmt.Sprintf("must have a tax_id, or an identity with ext '%s'", ExtKeyIdentityType),
					is.Func("has tax_id or identity", customerHasTaxIDOrIdentity),
				),
				rules.Field("tax_id",
					rules.Field("code",
						rules.Assert("06", "customer tax ID must have a code", is.Present),
					),
				),
			),
		),
		// Invoice tax extensions
		// Code 07: tax required
		// Code 08: doc_type required
		// Code 09: correction_type required when credit/debit doc type
		rules.Field("tax",
			rules.Assert("07", "tax object is required with extensions", is.Present),
			rules.Field("ext",
				rules.Assert("08", fmt.Sprintf("tax extension '%s' is required", ExtKeyDocType),
					tax.ExtensionsRequire(ExtKeyDocType),
				),
				rules.When(
					tax.ExtensionsHasCodes(ExtKeyDocType, "R1", "R2", "R3", "R4", "R5"),
					rules.Assert("09", fmt.Sprintf("tax extension '%s' is required", ExtKeyCorrectionType),
						tax.ExtensionsRequire(ExtKeyCorrectionType),
					),
				),
			),
		),
		// Code 10: standard invoice doc type must be F1, F2, or F3
		rules.When(
			bill.InvoiceTypeIn(bill.InvoiceTypeStandard),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("10", fmt.Sprintf("tax extension '%s' for standard invoices must be F1, F2, or F3", ExtKeyDocType),
						tax.ExtensionsHasCodes(ExtKeyDocType, "F1", "F2", "F3"),
					),
				),
			),
		),
		// Code 11: corrective invoice doc type must be R1-R5
		rules.When(
			bill.InvoiceTypeIn(bill.InvoiceTypeCorrective, bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("11", "doc type extension for corrective invoices must be R1, R2, R3, R4, or R5",
						tax.ExtensionsHasCodes(ExtKeyDocType, "R1", "R2", "R3", "R4", "R5"),
					),
				),
			),
		),
		// Notes
		// Code 12: general note text max 500 characters
		rules.Field("notes",
			rules.Each(
				rules.When(
					is.Func("general note", isGeneralNote),
					rules.Field("text",
						rules.Assert("12", "general note text must be 500 characters or less", is.Length(0, 500)),
					),
				),
			),
		),
		// Lines
		// Code 13: all VAT/IGIC tax combos must have consistent product presence
		// Code 14: all VAT/IGIC tax combos must have the same regime
		rules.Field("lines",
			rules.Assert("13", fmt.Sprintf("'%s' must be present in all tax combos or none", ExtKeyProduct),
				is.Func("consistent product", invoiceLinesHaveConsistentProduct),
			),
			rules.Assert("14", fmt.Sprintf("'%s' must be the same in all tax combos", ExtKeyRegime),
				is.Func("same regime", invoiceLinesHaveSameRegime),
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

func customerHasTaxIDOrIdentity(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true // nil customer handled by Required check
	}
	return p.TaxID != nil || org.IdentityForExtKey(p.Identities, ExtKeyIdentityType) != nil
}

func isGeneralNote(val any) bool {
	note, ok := val.(*org.Note)
	return ok && note != nil && note.Key == org.NoteKeyGeneral
}

func invoiceLinesHaveConsistentProduct(val any) bool {
	lines, ok := val.([]*bill.Line)
	if !ok {
		return true
	}
	var ref *tax.Combo
	for _, l := range lines {
		for _, tc := range l.Taxes {
			if !tc.Category.In(tax.CategoryVAT, es.TaxCategoryIGIC) {
				continue
			}
			if ref == nil {
				ref = tc
				continue
			}
			if ref.Ext.Has(ExtKeyProduct) != tc.Ext.Has(ExtKeyProduct) {
				return false
			}
		}
	}
	return true
}

func invoiceLinesHaveSameRegime(val any) bool {
	lines, ok := val.([]*bill.Line)
	if !ok {
		return true
	}
	var ref *tax.Combo
	for _, l := range lines {
		for _, tc := range l.Taxes {
			if !tc.Category.In(tax.CategoryVAT, es.TaxCategoryIGIC) {
				continue
			}
			if ref == nil {
				ref = tc
				continue
			}
			if ref.Ext.Get(ExtKeyRegime) != tc.Ext.Get(ExtKeyRegime) {
				return false
			}
		}
	}
	return true
}
