package ad

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

var scenarios = []*tax.ScenarioSet{
	invoiceScenarios,
}

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		{
			// Diplomatic Exemption — Form 980-C, Art. 15 Llei 11/2012
			Tags: []cbc.Key{"diplomatic"},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  "diplomatic",
				Text: "Exempt d'IGI - Art. 15 Llei 11/2012 (Règim diplomàtic)",
			},
		},
		{
			// Traveler / Export refund — Form 980-A
			// IGI must be shown separately at 1% and 4.5% with a Declaration of Non-Residency.
			Tags: []cbc.Key{tax.TagExport},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  tax.TagExport,
				Text: "Exempt d'IGI - Exportació / Devolució a viatgers (Form 980-A)",
			},
		},
		{
			// Non-Resident B2B — Form 980-B
			// Requires a Fiscal Representative in Andorra (name + NRT).
			Tags: []cbc.Key{"non-resident-b2b"},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  "non-resident-b2b",
				Text: "Operació amb no-resident - Representant fiscal requerit (Form 980-B)",
			},
		},
	},
}
