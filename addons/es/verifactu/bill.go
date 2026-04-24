package verifactu

import (
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
		rules.Assert("16", "invoice must be in EUR or provide exchange rate for conversion", currency.CanConvertTo(currency.EUR)),
		// Code 20: series and code combined max 60 chars
		rules.Assert("20", "invoice series and code combined must be 60 characters or less",
			is.Func("series and code fit", invoiceSeriesCodeFits),
		),
		// Preceding documents
		// Code 01: preceding required when corrective
		rules.When(
			bill.InvoiceTypeIn(bill.InvoiceTypeCorrective),
			rules.Field("preceding",
				rules.Assert("01", "preceding documents are required for corrective invoices", is.Present),
			),
		),
		// Code 02: each preceding issue date required
		// Code 21: each preceding series and code combined max 60 chars
		rules.Field("preceding",
			rules.Each(
				rules.When(
					is.Func("not nil", precedingDocIsNotNil),
					rules.Field("issue_date",
						rules.Assert("02", "issue date is required", is.Present),
					),
					rules.AssertIfPresent("21", "preceding series and code combined must be 60 characters or less",
						is.Func("series and code fit", precedingSeriesCodeFits),
					),
				),
			),
		),
		// Code 03: each preceding tax required when corrective
		rules.When(
			bill.InvoiceTypeIn(bill.InvoiceTypeCorrective),
			rules.Field("preceding",
				rules.Each(
					rules.When(
						is.Func("not nil", precedingDocIsNotNil),
						rules.Field("tax",
							rules.Assert("03", "preceding invoice tax data is required for corrective invoices", is.Present),
						),
					),
				),
			),
		),
		// Supplier
		// Code 18: supplier name max 120 chars
		rules.Field("supplier",
			rules.Field("name",
				rules.AssertIfPresent("18", "supplier name must be 120 characters or less", is.Length(0, 120)),
			),
		),
		// Customer - universal
		// Code 19: customer name max 120 chars
		// Code 22: non-ES customer tax ID code max 18 chars
		rules.Field("customer",
			rules.Field("name",
				rules.AssertIfPresent("19", "customer name must be 120 characters or less", is.Length(0, 120)),
			),
			rules.Field("tax_id",
				rules.When(
					is.Func("non-ES tax ID", taxIDIsNonES),
					rules.Field("code",
						rules.AssertIfPresent("22", "non-Spanish customer tax ID code must be 18 characters or less", is.Length(0, 18)),
					),
				),
			),
		),
		// Customer - simplified invoices (F2 or R5)
		// Code 04: no tax_id on simplified customer
		// Code 05: no identity type ext on simplified customer
		rules.When(
			is.Func("simplified", isSimplifiedInvoice),
			rules.Field("customer",
				rules.Field("tax_id",
					rules.Assert("04", "customer tax ID must not be set for simplified invoices",
						is.Nil,
					),
				),
				rules.Assert("05", "customer identity type extension not allowed for simplified invoices",
					is.Func("no identity type ext", simplifiedCustomerHasNoIdentityType),
				),
			),
		),
		// Customer - standard invoices
		// Code 06: customer required
		// Code 07: customer must have tax_id or identity
		// Code 08: customer tax_id must have code
		// Code 17: identity country required when identity type is not NIF-VAT
		rules.When(
			is.Func("not simplified", isNotSimplifiedInvoice),
			rules.Field("customer",
				rules.Assert("06", "customer is required", is.Present),
				rules.Assert("07", "must have a tax_id or an identity with ext 'es-verifactu-identity-type'",
					is.Func("has tax_id or identity", customerHasTaxIDOrIdentity),
				),
				rules.Field("tax_id",
					rules.Field("code",
						rules.Assert("08", "tax ID must have a code", is.Present),
					),
				),
				rules.Field("identities",
					rules.Each(
						rules.When(
							is.Func("has non-VAT identity type", identityHasNonVATType),
							rules.Field("country",
								rules.Assert("17", "country is required when ext 'es-verifactu-identity-type' is not 02 (NIF-VAT)", is.Present),
							),
						),
					),
				),
			),
		),
		// Invoice tax extensions
		// Code 09: tax required
		// Code 10: doc_type required
		// Code 13: correction_type required when credit/debit doc type
		rules.Field("tax",
			rules.Assert("09", "tax is required", is.Present),
			rules.Field("ext",
				rules.Assert("10", "doc type is required",
					tax.ExtensionsRequire(ExtKeyDocType),
				),
				rules.When(
					tax.ExtensionsHasCodes(ExtKeyDocType, "R1", "R2", "R3", "R4", "R5"),
					rules.Assert("13", "correction type extension is required",
						tax.ExtensionsRequire(ExtKeyCorrectionType),
					),
				),
			),
		),
		// Code 11: standard invoice doc type must be F1, F2, or F3
		rules.When(
			bill.InvoiceTypeIn(bill.InvoiceTypeStandard),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("11", "doc type extension for standard invoices must be F1, F2, or F3",
						tax.ExtensionsHasCodes(ExtKeyDocType, "F1", "F2", "F3"),
					),
				),
			),
		),
		// Code 12: corrective invoice doc type must be R1-R5
		rules.When(
			bill.InvoiceTypeIn(bill.InvoiceTypeCorrective, bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("12", "doc type extension for corrective invoices must be R1, R2, R3, R4, or R5",
						tax.ExtensionsHasCodes(ExtKeyDocType, "R1", "R2", "R3", "R4", "R5"),
					),
				),
			),
		),
		// Notes
		// Code 14: general note text max 500 characters
		rules.Field("notes",
			rules.Each(
				rules.When(
					is.Func("general note", isGeneralNote),
					rules.Field("text",
						rules.Assert("14", "general note text must be 500 characters or less", is.Length(0, 500)),
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
						tax.SetHasOneOf(tax.CategoryVAT, es.TaxCategoryIGIC, es.TaxCategoryIPSI),
					),
				),
			),
		),
		// Totals - tax categories
		// Code 23: non-retained tax rates cannot exceed 12
		rules.Field("totals",
			rules.Field("taxes",
				rules.Field("categories",
					rules.Assert("23", "non-retained tax rates cannot exceed 12",
						is.Func("non-retained rate count within limit", taxRatesWithinLimit),
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

func precedingDocIsNotNil(val any) bool {
	ref, ok := val.(*org.DocumentRef)
	return ok && ref != nil
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

func identityHasNonVATType(val any) bool {
	id, ok := val.(*org.Identity)
	return ok && id != nil && id.Ext.Has(ExtKeyIdentityType) && !id.Ext.Get(ExtKeyIdentityType).In(ExtCodeIdentityTypeVAT)
}

func isGeneralNote(val any) bool {
	note, ok := val.(*org.Note)
	return ok && note != nil && note.Key == org.NoteKeyGeneral
}

func invoiceSeriesCodeFits(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	return seriesCodeFits(inv.Series, inv.Code)
}

func precedingSeriesCodeFits(val any) bool {
	ref, ok := val.(*org.DocumentRef)
	if !ok || ref == nil {
		return true
	}
	return seriesCodeFits(ref.Series, ref.Code)
}

func seriesCodeFits(series, code cbc.Code) bool {
	joined := series.Join(code)
	return len(joined) <= 60
}

func taxIDIsNonES(val any) bool {
	tid, ok := val.(*tax.Identity)
	return ok && tid != nil && tid.Country != "" && tid.Country != "ES"
}

func taxRatesWithinLimit(val any) bool {
	cats, ok := val.([]*tax.CategoryTotal)
	if !ok {
		return true
	}
	count := 0
	for _, c := range cats {
		if c == nil || c.Retained {
			continue
		}
		count += len(c.Rates)
	}
	return count <= 12
}
