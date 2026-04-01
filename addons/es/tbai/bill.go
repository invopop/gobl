package tbai

import (
	"fmt"
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
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
}

func normalizeInvoiceTax(inv *bill.Invoice) {
	tx := inv.Tax
	if tx == nil {
		tx = &bill.Tax{}
	}
	if tx.Ext == nil {
		tx.Ext = make(tax.Extensions)
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
		tx.Ext[ExtKeyRegion] = "VI"
	case "bizkaia", "vizcaya", "bi":
		tx.Ext[ExtKeyRegion] = "BI"
	case "gipuzkoa", "guipuzcoa", "guipúzcoa", "ss":
		tx.Ext[ExtKeyRegion] = "SS"
	default:
		return
	}
	if len(tx.Ext) > 0 {
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
		rules.When(
			is.Func("non-simplified", func(val any) bool {
				inv, ok := val.(*bill.Invoice)
				return ok && inv != nil && !inv.HasTags(tax.TagSimplified)
			}),
			rules.Field("customer",
				rules.Assert("03", "customer is required for non-simplified invoices", is.Present),
				rules.Field("tax_id",
					rules.Assert("03a", "customer tax ID is required", is.Present),
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
