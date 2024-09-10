package bill

import (
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// ScenarioSummary determines a summary of the tax scenario for the invoice based on
// the document type and tax tags.
//
// Deprecated: tax regimes should be updated to automatically apply all the required
// extensions and meta-data to the invoice itself. This method will still be needed
// until all regimes have transitioned to the new approach.
func (inv *Invoice) ScenarioSummary() *tax.ScenarioSummary {
	r := inv.TaxRegime()
	if r == nil {
		return nil
	}
	return inv.scenarioSummary(r)
}

func (inv *Invoice) scenarioSummary(r *tax.Regime) *tax.ScenarioSummary {
	if r == nil {
		return nil
	}
	ss := r.ScenarioSet(ShortSchemaInvoice)
	if ss == nil {
		return nil
	}
	exts := make([]tax.Extensions, 0)
	tags := []cbc.Key{}

	if inv.Tax != nil {
		tags = inv.Tax.Tags
		if len(inv.Tax.Ext) > 0 {
			exts = append(exts, inv.Tax.Ext)
		}
	}
	for _, cat := range inv.Totals.Taxes.Categories {
		for _, rate := range cat.Rates {
			exts = append(exts, rate.Ext)
		}
	}

	return ss.SummaryFor(inv.Type, tags, exts)
}

func (inv *Invoice) prepareTags(r *tax.Regime) error {
	if r == nil {
		return nil
	}
	if inv.Tax == nil {
		return nil
	}

	// Check the tags are all valid and identified by the tax regime
	// as acceptable for invoices.
	for _, k := range inv.Tax.Tags {
		if t := r.Tag(k); t == nil {
			return validation.Errors{
				"tax": validation.Errors{
					"tags": fmt.Errorf("invalid tag '%v'", k),
				},
			}
		}
	}

	return nil
}

func (inv *Invoice) prepareScenarios(r *tax.Regime) error {
	// Use the scenario summary to add any notes to the invoice
	ss := inv.scenarioSummary(r)
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
		if inv.Tax.Ext[k] == "" {
			// only override if not set
			inv.Tax.Ext[k] = v
		}
	}

	return nil
}
