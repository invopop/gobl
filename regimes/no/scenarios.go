package no

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// Norwegian invoice tags.
const (
	// TagForetaksregisteret indicates the supplier is registered in the
	// Register of Business Enterprises (Foretaksregisteret), mandatory for
	// AS/ASA/NUF entities per foretaksregisterloven section 10-2.
	TagForetaksregisteret cbc.Key = "foretaksregisteret"
)

var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*cbc.Definition{
		{
			Key: TagForetaksregisteret,
			Name: i18n.String{
				i18n.EN: "Registered in the Register of Business Enterprises",
				i18n.NB: "Registrert i Foretaksregisteret",
			},
		},
	},
}

// invoiceScenarios returns Norwegian-specific invoice scenarios that inject
// legal notes based on tags. A custom set is used instead of
// bill.InvoiceScenarios() to provide bilingual (EN/NB) legal text.
func invoiceScenarios() *tax.ScenarioSet {
	return &tax.ScenarioSet{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			// Reverse charge
			{
				Tags: []cbc.Key{tax.TagReverseCharge},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  tax.TagReverseCharge,
					Text: "Reverse charge / Omvendt avgiftsplikt – Merverdiavgift ikke beregnet.",
				},
			},
			// Foretaksregisteret label for AS/ASA/NUF entities
			{
				Tags: []cbc.Key{TagForetaksregisteret},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  TagForetaksregisteret,
					Text: "Foretaksregisteret",
				},
			},
		},
	}
}
