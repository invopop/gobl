package tbai

import (
	"fmt"
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
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
			ExtKeyCorrection,
		},
	},
}

func normalizeInvoice(inv *bill.Invoice) {
	if inv == nil {
		return
	}
	normalizeInvoiceTax(inv)
	normalizeInvoiceCustomerRates(inv)
	normalizeInvoiceRegime(inv)
	normalizeInvoicePartyIdentity(inv.Customer)
}

// normalizeInvoiceCustomerRates flags VAT/IGIC combos as "no sujeta" when the
// invoice carries the ~customer-rates~ tag with a non-Spanish customer.
func normalizeInvoiceCustomerRates(inv *bill.Invoice) {
	if !inv.HasTags(tax.TagCustomerRates) {
		return
	}
	if inv.Customer == nil || inv.Customer.TaxID == nil {
		return
	}
	country := inv.Customer.TaxID.Country
	if country == "" || country == l10n.ES.Tax() {
		return
	}
	apply := func(tc *tax.Combo) {
		if tc == nil || !tc.Category.In(tax.CategoryVAT, es.TaxCategoryIGIC) {
			return
		}
		tc.Ext = tc.Ext.SetOneOf(ExtKeyExempt, "RL", "IE")
	}
	for _, line := range inv.Lines {
		if line == nil {
			continue
		}
		for _, tc := range line.Taxes {
			apply(tc)
		}
	}
	for _, ch := range inv.Charges {
		if ch == nil {
			continue
		}
		for _, tc := range ch.Taxes {
			apply(tc)
		}
	}
	for _, d := range inv.Discounts {
		if d == nil {
			continue
		}
		for _, tc := range d.Taxes {
			apply(tc)
		}
	}
}

// normalizeInvoiceRegime fills in the regime code on every VAT/IGIC combo:
// ~52~ under the simplified-scheme tag, ~01~ otherwise.
func normalizeInvoiceRegime(inv *bill.Invoice) {
	simplified := inv.HasTags(es.TagSimplifiedScheme)
	apply := func(tc *tax.Combo) {
		if tc == nil || !tc.Category.In(tax.CategoryVAT, es.TaxCategoryIGIC) {
			return
		}
		if simplified {
			tc.Ext = tc.Ext.SetIfEmpty(ExtKeyRegime, "52")
		}
		tc.Ext = tc.Ext.SetIfEmpty(ExtKeyRegime, "01")
	}
	for _, line := range inv.Lines {
		if line == nil {
			continue
		}
		for _, tc := range line.Taxes {
			apply(tc)
		}
	}
	for _, ch := range inv.Charges {
		if ch == nil {
			continue
		}
		for _, tc := range ch.Taxes {
			apply(tc)
		}
	}
	for _, d := range inv.Discounts {
		if d == nil {
			continue
		}
		for _, tc := range d.Taxes {
			apply(tc)
		}
	}
}

// normalizeInvoicePartyIdentity maps the customer's first identity key onto the
// es-tbai-identity-type extension when no Spanish NIF is present.
func normalizeInvoicePartyIdentity(cus *org.Party) {
	if cus == nil {
		return
	}
	if cus.TaxID != nil && cus.TaxID.Country == "ES" && cus.TaxID.Code != "" {
		return
	}
	if len(cus.Identities) == 0 {
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
		id.Ext = id.Ext.Merge(tax.ExtensionsOf(cbc.CodeMap{
			ExtKeyIdentityType: code,
		}))
	}
}

func normalizeInvoiceTax(inv *bill.Invoice) {
	tx := inv.Tax
	if tx == nil {
		tx = &bill.Tax{}
	}
	if tx.Ext.IsZero() {
		tx.Ext = tax.MakeExtensions()
	}
	if tx.Ext.Has(ExtKeyRegion) {
		return
	}
	if inv.Supplier == nil || len(inv.Supplier.Addresses) == 0 {
		return
	}
	addr := inv.Supplier.Addresses[0]
	// Take a set of different names for the same region and attempt
	// to use them to set the region code automatically.
	switch strings.ToLower(addr.Region) {
	case "alava", "álava", "araba", "vi":
		tx.Ext = tx.Ext.Set(ExtKeyRegion, ExtValueRegionVI)
	case "bizkaia", "vizcaya", "bi":
		tx.Ext = tx.Ext.Set(ExtKeyRegion, ExtValueRegionBI)
	case "gipuzkoa", "guipuzcoa", "guipúzcoa", "ss":
		tx.Ext = tx.Ext.Set(ExtKeyRegion, ExtValueRegionSS)
	default:
		return
	}
	if tx.Ext.Len() > 0 {
		inv.Tax = tx
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
		vt.Ext = vt.Ext.SetOneOf(ExtKeyProduct, "goods", "resale")
	case org.ItemKeyServices, cbc.KeyEmpty:
		vt.Ext = vt.Ext.Set(ExtKeyProduct, "services")
	}
}

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Assert("09", "invoice must be in EUR or provide exchange rate for conversion", currency.CanConvertTo(currency.EUR)),
		// Tax
		// Code 01: tax required
		// Code 02: region required in tax ext
		rules.Field("tax",
			rules.Assert("01", "tax is required", is.Present),
			rules.Field("ext",
				rules.Assert("02", fmt.Sprintf("extension '%s' is required", ExtKeyRegion),
					tax.ExtensionsRequire(ExtKeyRegion),
				),
			),
		),
		// Customer
		// Code 03: customer required for non-simplified invoices
		// Code 08: customer must have tax_id or an identity carrying the
		//          es-tbai-identity-type extension (the L7 IDType code).
		rules.When(
			is.Func("non-simplified", func(val any) bool {
				inv, ok := val.(*bill.Invoice)
				return ok && inv != nil && !inv.HasTags(tax.TagSimplified)
			}),
			rules.Field("customer",
				rules.Assert("03", "customer is required for non-simplified invoices", is.Present),
				rules.Assert("08", "customer must have a tax_id or an identity with ext 'es-tbai-identity-type'",
					is.Func("has tax_id or identity-type identity", customerHasTaxIDOrIdentity),
				),
			),
		),
		// Preceding
		// Code 04: preceding required for correction types
		rules.When(
			bill.InvoiceTypeIn(es.InvoiceCorrectionTypes...),
			rules.Field("preceding",
				rules.Assert("04", fmt.Sprintf("preceding documents are required for %s invoices",
					strings.Join(cbc.KeyStrings(es.InvoiceCorrectionTypes), ", ")),
					is.Present,
				),
			),
		),
		// Code 05: each preceding issue_date required
		// Code 06: each preceding ext correction required
		rules.Field("preceding",
			rules.Each(
				rules.Field("issue_date",
					rules.Assert("05", "preceding issue date is required", is.Present),
				),
				rules.Field("ext",
					rules.Assert("06", fmt.Sprintf("preceding ext '%s' is required", ExtKeyCorrection),
						tax.ExtensionsRequire(ExtKeyCorrection),
					),
				),
			),
		),
		// Notes
		// Code 07: must have a general note
		rules.Field("notes",
			rules.Assert("07", "with key 'general' missing",
				is.Func("has general note", notesHasGeneralKey),
			),
		),
		// Supplier
		// Code 10: activity ext required for Bizkaia individuals (Modelo 140 LROE)
		// Code 11: activity ext, when present, must be a valid epígrafe code
		rules.When(
			is.Func("Bizkaia individual", isBizkaiaIndividual),
			rules.Field("supplier",
				rules.Field("ext",
					rules.Assert("10", fmt.Sprintf("extension '%s' is required for Bizkaia individuals", ExtKeyBIActivity),
						tax.ExtensionsRequire(ExtKeyBIActivity),
					),
				),
			),
		),
		rules.Field("supplier",
			rules.Field("ext",
				rules.Assert("11", fmt.Sprintf("extension '%s' must be a valid Bizkaia activity code (epígrafe)", ExtKeyBIActivity),
					tax.ExtensionHasValidCode(ExtKeyBIActivity),
				),
			),
		),
	)
}

func notesHasGeneralKey(val any) bool {
	notes, ok := val.([]*org.Note)
	if !ok {
		return false
	}
	for _, n := range notes {
		if n.Key.In(org.NoteKeyGeneral) {
			return true
		}
	}
	return false
}

func isBizkaiaIndividual(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Tax == nil || inv.Supplier == nil {
		return false
	}
	return inv.Tax.Ext.Get(ExtKeyRegion) == ExtValueRegionBI &&
		es.TaxIdentityKey(inv.Supplier.TaxID) != es.TaxIdentityOrg
}

func customerHasTaxIDOrIdentity(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true // nil customer handled by the presence check above
	}
	return p.TaxID != nil || org.IdentityForExtKey(p.Identities, ExtKeyIdentityType) != nil
}
