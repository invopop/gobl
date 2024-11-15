package bill

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

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

	inv.removePreviousScenarios(ss)
	return ss.SummaryFor(inv)
}

func (inv *Invoice) removePreviousScenarios(ss *tax.ScenarioSet) {
	if inv.Tax != nil && len(inv.Tax.Ext) > 0 {
		for _, ek := range ss.ExtensionKeys() {
			delete(inv.Tax.Ext, ek)
		}
	}
	for _, n := range ss.Notes() {
		for i, n2 := range inv.Notes {
			if n.SameAs(n2) {
				// remove from array
				inv.Notes = append(inv.Notes[:i], inv.Notes[i+1:]...)
			}
		}
	}
}

func (inv *Invoice) prepareScenarios() error {
	// Use the scenario summary to add any notes to the invoice
	ss := inv.scenarioSummary()
	if ss == nil {
		return nil
	}
	for _, n := range ss.Notes {
		// make sure we don't already have the same note in the invoice
		for _, n2 := range inv.Notes {
			if n.SameAs(n2) {
				n = nil
				break
			}
		}
		if n != nil {
			inv.Notes = append(inv.Notes, n)
		}
	}
	// Apply extensions at the document level
	for k, v := range ss.Ext {
		if inv.Tax == nil {
			inv.Tax = new(Tax)
		}
		if inv.Tax.Ext == nil {
			inv.Tax.Ext = make(tax.Extensions)
		}
		// Always override
		inv.Tax.Ext[k] = v
	}

	return nil
}
