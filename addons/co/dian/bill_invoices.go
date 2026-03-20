package dian

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

var invoiceCorrectionDefinitions = []*tax.CorrectionDefinition{
	{
		Schema: bill.ShortSchemaInvoice,
		Types: []cbc.Key{
			bill.InvoiceTypeCreditNote,
			bill.InvoiceTypeDebitNote,
		},
		Extensions: []cbc.Key{
			ExtKeyCreditCode,
			ExtKeyDebitCode,
		},
		ReasonRequired: true,
		Stamps: []cbc.Key{
			StampCUDE,
		},
	},
}

func normalizeInvoice(inv *bill.Invoice) {
	if inv == nil {
		return
	}
	normalizeInvoiceParty(inv.Supplier)
	normalizeInvoiceParty(inv.Customer)
}

func normalizeInvoiceParty(p *org.Party) {
	if p == nil || !isColombian(p.TaxID) {
		return
	}
	def := tax.Extensions{ExtKeyFiscalResponsibility: "R-99-PN"}
	p.Ext = def.Merge(p.Ext)
}

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		// Code 01: invoice type restriction
		rules.Assert("01", "invoice type must be one of standard, credit-note, debit-note, or proforma",
			bill.InvoiceTypeIn(
				bill.InvoiceTypeStandard,
				bill.InvoiceTypeCreditNote,
				bill.InvoiceTypeDebitNote,
				bill.InvoiceTypeProforma,
			),
		),
		// Supplier validation
		// Code 02: Colombian supplier must have at least one address
		// Code 03: Colombian supplier with tax ID code must have municipality ext
		// Code 04: Colombian supplier must have fiscal responsibility ext
		rules.Field("supplier",
			rules.When(
				is.Func("colombian supplier", supplierIsColombian),
				rules.Field("addresses",
					rules.Assert("02", "at least one address is required for Colombian suppliers", is.Length(1, 0)),
				),
			),
			rules.When(
				is.Func("colombian supplier with code", supplierMunicipalityRequired),
				rules.Field("ext",
					rules.Assert("03", fmt.Sprintf("extension '%s' is required", ExtKeyMunicipality),
						tax.ExtensionsRequire(ExtKeyMunicipality),
					),
				),
			),
			rules.When(
				is.Func("colombian supplier", supplierIsColombian),
				rules.Field("ext",
					rules.Assert("04", fmt.Sprintf("extension '%s' is required", ExtKeyFiscalResponsibility),
						tax.ExtensionsRequire(ExtKeyFiscalResponsibility),
					),
				),
			),
		),
		// Customer validation
		// Code 05: customer tax_id required when not simplified
		// Code 06: Colombian customer must have at least one address
		// Code 07: Colombian customer with tax ID code must have municipality ext
		// Code 08: Colombian customer must have fiscal responsibility ext
		rules.Field("customer",
			rules.When(
				is.Func("not simplified", invoiceNotSimplified),
				rules.Field("tax_id",
					rules.Assert("05", "customer tax ID is required", is.Present),
				),
			),
			rules.When(
				is.Func("colombian customer", customerIsColombian),
				rules.Field("addresses",
					rules.Assert("06", "at least one address is required for Colombian customers", is.Length(1, 0)),
				),
			),
			rules.When(
				is.Func("colombian customer with code", customerMunicipalityRequired),
				rules.Field("ext",
					rules.Assert("07", fmt.Sprintf("extension '%s' is required", ExtKeyMunicipality),
						tax.ExtensionsRequire(ExtKeyMunicipality),
					),
				),
			),
			rules.When(
				is.Func("colombian customer", customerIsColombian),
				rules.Field("ext",
					rules.Assert("08", fmt.Sprintf("extension '%s' is required", ExtKeyFiscalResponsibility),
						tax.ExtensionsRequire(ExtKeyFiscalResponsibility),
					),
				),
			),
		),
		// Preceding validation
		// Code 09: preceding required for credit/debit notes
		rules.When(
			bill.InvoiceTypeIn(bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote),
			rules.Field("preceding",
				rules.Assert("09", "preceding documents are required for credit and debit notes", is.Present),
			),
		),
		// Code 10: preceding credit code required for credit notes
		rules.When(
			bill.InvoiceTypeIn(bill.InvoiceTypeCreditNote),
			rules.Field("preceding",
				rules.Each(
					rules.Field("ext",
						rules.Assert("10", fmt.Sprintf("extension '%s' is required for credit notes", ExtKeyCreditCode),
							tax.ExtensionsRequire(ExtKeyCreditCode),
						),
					),
					rules.Field("reason",
						rules.Assert("12", "preceding reason is required", is.Present),
					),
				),
			),
		),
		// Code 11: preceding debit code required for debit notes
		rules.When(
			bill.InvoiceTypeIn(bill.InvoiceTypeDebitNote),
			rules.Field("preceding",
				rules.Each(
					rules.Field("ext",
						rules.Assert("11", fmt.Sprintf("extension '%s' is required for debit notes", ExtKeyDebitCode),
							tax.ExtensionsRequire(ExtKeyDebitCode),
						),
					),
					rules.Field("reason",
						rules.Assert("12", "preceding reason is required", is.Present),
					),
				),
			),
		),
	)
}

func isColombian(tID *tax.Identity) bool {
	return tID != nil && tID.Country.In("CO")
}

// municipalityCodeRequired checks if the municipality code is required for the given tax
// identity by checking to see if the customer is a Colombian company.
func municipalityCodeRequired(tID *tax.Identity) bool {
	if tID == nil {
		return false
	}
	if !tID.Country.In("CO") {
		return false
	}
	return tID.Code != ""
}

func supplierIsColombian(val any) bool {
	p, ok := val.(*org.Party)
	return ok && p != nil && isColombian(p.TaxID)
}

func supplierMunicipalityRequired(val any) bool {
	p, ok := val.(*org.Party)
	return ok && p != nil && municipalityCodeRequired(p.TaxID)
}

func customerIsColombian(val any) bool {
	p, ok := val.(*org.Party)
	return ok && p != nil && isColombian(p.TaxID)
}

func customerMunicipalityRequired(val any) bool {
	p, ok := val.(*org.Party)
	return ok && p != nil && municipalityCodeRequired(p.TaxID)
}

func invoiceNotSimplified(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	return !tax.TagSimplified.In(inv.GetTags()...)
}
