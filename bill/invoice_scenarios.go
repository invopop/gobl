package bill

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

var invoiceScenarios = &tax.ScenarioSet{
	Schema: ShortSchemaInvoice,
	List: []*tax.Scenario{
		// Reverse Charges
		{
			Tags: []cbc.Key{tax.TagReverseCharge},
			Note: &tax.Note{
				Category: tax.CategoryVAT,
				Key:      tax.KeyReverseCharge,
				Text:     "Reverse charge: Customer to account for VAT to the relevant tax authority.",
			},
		},
	},
}

// InvoiceScenarios provides a standard set of scenarios to either be extended
// or overridden by a regime or addon.
func InvoiceScenarios() *tax.ScenarioSet {
	return invoiceScenarios
}

// GetType provides the invoice type as part of the tax.ScenarioDocument interface.
func (inv *Invoice) GetType() cbc.Key {
	return inv.Type
}

// GetExtensions goes through the invoice and grabs all the extensions that are in
// use and expected to be used as part of a scenario.
func (inv *Invoice) GetExtensions() []tax.Extensions {
	exts := make([]tax.Extensions, 0)
	if inv.Tax != nil {
		if len(inv.Tax.Ext) > 0 {
			exts = append(exts, inv.Tax.Ext)
		}
	}
	if inv.Totals != nil && inv.Totals.Taxes != nil {
		for _, cat := range inv.Totals.Taxes.Categories {
			for _, rate := range cat.Rates {
				exts = append(exts, rate.Ext)
			}
		}
	}
	return exts
}

// ScenarioSummary determines a summary of the tax scenario for the invoice based on
// the document type and tax tags.
//
// Deprecated: tax regimes should be updated to automatically apply all the required
// extensions and meta-data to the invoice itself. This method will still be needed
// until all regimes have transitioned to the new approach.
func (inv *Invoice) ScenarioSummary() *tax.ScenarioSummary {
	return inv.scenarioSummary()
}

func (inv *Invoice) scenarioSummary() *tax.ScenarioSummary {
	ss := tax.NewScenarioSet(ShortSchemaInvoice)

	if r := inv.RegimeDef(); r != nil {
		ss.Merge(r.Scenarios)
	}
	for _, a := range inv.AddonDefs() {
		ss.Merge(a.Scenarios)
	}

	return ss.SummaryFor(inv)
}

func (inv *Invoice) prepareScenarios() error {
	// Use the scenario summary to add any notes to the invoice
	ss := inv.scenarioSummary()
	if ss == nil {
		return nil
	}

	normalizers := tax.ExtractNormalizers(inv)

	for _, sn := range ss.Notes {
		// make sure we don't already have the same note
		found := false
		if inv.Tax != nil {
			for _, n := range inv.Tax.Notes {
				if sn.SameAs(n) {
					found = true
					break
				}
			}
		}
		if !found {
			// Normalize the note so addons can enrich it (e.g. en16931
			// adds UNTDID tax category extensions).
			sn.Normalize(normalizers)
			inv.Tax = inv.Tax.MergeNotes(sn)
		}
	}

	// Apply extensions at the document level
	if len(ss.Ext) > 0 {
		inv.Tax = inv.Tax.MergeExtensions(ss.Ext)
	}

	return nil
}
